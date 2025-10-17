package db

import (
	"bytes"
	"context"
	"database/sql"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var (
	// SQL 분류 정규식 (가벼운 1차 파싱)
	reVerb = regexp.MustCompile(`(?i)^\s*([a-z]+)`)
	reTbl  = regexp.MustCompile(`(?i)\bfrom\s+([a-z0-9_]+)|\binto\s+([a-z0-9_]+)|\bupdate\s+([a-z0-9_]+)`)
	reNum  = regexp.MustCompile(`\b\d+\b`)
	reStr  = regexp.MustCompile(`'[^']*'|"[^"]*"`)
)

func classifySQL(sqlText string) (verb, rw, table string) {
	m := reVerb.FindStringSubmatch(sqlText)
	if len(m) >= 2 {
		verb = strings.ToUpper(m[1])
	} else {
		verb = "UNKNOWN"
	}
	switch verb {
	case "SELECT", "SHOW", "DESCRIBE", "EXPLAIN":
		rw = "read"
	default:
		rw = "write"
	}
	if m := reTbl.FindStringSubmatch(strings.ToLower(sqlText)); m != nil {
		for i := 1; i < len(m); i++ {
			if m[i] != "" {
				table = m[i]
				break
			}
		}
	}
	return
}

func normalizeSQL(sqlText string) string {
	s := reStr.ReplaceAllString(sqlText, "?")
	s = reNum.ReplaceAllString(s, "?")
	s = strings.Join(strings.Fields(s), " ")
	return s
}

// ---------- 기능 태그 컨텍스트 ----------

type featureKey struct{}

type FeatureTags struct {
	Feature string // 예: GameLoop / Battle / Retry ...
	Action  string // 선택: start / retry ...
}

func WithFeatureTags(ctx context.Context, feature, action string) context.Context {
	return context.WithValue(ctx, featureKey{}, FeatureTags{Feature: feature, Action: action})
}

func getFeatureTags(ctx context.Context) (feature, action string) {
	if v, ok := ctx.Value(featureKey{}).(FeatureTags); ok {
		return v.Feature, v.Action
	}
	return "", ""
}

// ---------- 요청별 집계기 ----------

type ctxKeyAgg struct{}

type FeatureAgg struct {
	DBReadCount  int64
	DBWriteCount int64
	DBReadRows   int64
	DBWriteRows  int64
	DBDurationMs int64 // MariaDB I/O 시간 합

	RedisReadCount  int64
	RedisWriteCount int64
	RedisDurationMs int64

	NPlus1    int64
	templates map[string]int
}

func WithFeatureAgg(ctx context.Context, agg *FeatureAgg) context.Context {
	if agg.templates == nil {
		agg.templates = make(map[string]int)
	}
	return context.WithValue(ctx, ctxKeyAgg{}, agg)
}

func getAgg(ctx context.Context) *FeatureAgg {
	if v := ctx.Value(ctxKeyAgg{}); v != nil {
		if agg, ok := v.(*FeatureAgg); ok {
			return agg
		}
	}
	return nil
}

// 외부(미들웨어 after 콜백)에서 요약을 읽을 때 사용
func GetAgg(ctx context.Context) *FeatureAgg { return getAgg(ctx) }

func recordTemplate(ctx context.Context, tpl string) {
	if agg := getAgg(ctx); agg != nil {
		agg.templates[tpl]++
		if agg.templates[tpl] == 5 {
			agg.NPlus1++
		}
	}
}

func addDBCost(ctx context.Context, rw string, rows int64, dur time.Duration) {
	if agg := getAgg(ctx); agg != nil {
		if rw == "read" {
			agg.DBReadCount++
			agg.DBReadRows += rows
		} else {
			agg.DBWriteCount++
			agg.DBWriteRows += rows
		}
		agg.DBDurationMs += dur.Milliseconds()
	}
}

// ---------- 고루틴-로컬 컨텍스트 (콜사이트 무수정용) ----------

var glsMu sync.RWMutex
var gls = map[uint64]context.Context{}

func pushCtx(ctx context.Context) {
	glsMu.Lock()
	gls[getGID()] = ctx
	glsMu.Unlock()
}
func popCtx() {
	glsMu.Lock()
	delete(gls, getGID())
	glsMu.Unlock()
}
func currentCtx() context.Context {
	glsMu.RLock()
	c := gls[getGID()]
	glsMu.RUnlock()
	if c != nil {
		return c
	}
	return context.Background()
}

// goroutine id 취득(서버 내부용)
func getGID() uint64 {
	var b [64]byte
	n := runtime.Stack(b[:], false)
	i := bytes.IndexByte(b[:n], ' ')
	j := bytes.IndexByte(b[i+1:n], ' ')
	id, _ := strconv.ParseUint(string(b[i+1:i+1+j]), 10, 64)
	return id
}

// 미들웨어에서 사용할 공개 함수
func PushCtxForMiddleware(ctx context.Context) { pushCtx(ctx) }
func PopCtxForMiddleware()                     { popCtx() }

// ---------- 트레이싱 래퍼: Exec / Query / ScanOne ----------

var dbTracer = otel.Tracer("pokerogue/db")

// 내부 공통: 실제 실행 + 스팬/집계
func execCore(ctx context.Context, execFn func(context.Context) (sql.Result, error), verb, rw, table, tpl string) (sql.Result, error) {
	feature, action := getFeatureTags(ctx)

	ctx, span := dbTracer.Start(ctx, "DB."+verb)
	span.SetAttributes(
		attribute.String("feature", feature),
		attribute.String("action", action),

		attribute.String("db.system", "mariadb"),
		attribute.String("db.operation", verb),
		attribute.String("db.rw", rw),
		attribute.String("db.sql.table", table),
		attribute.String("db.sql.template", tpl),
	)

	start := time.Now()
	res, err := execFn(ctx)
	dur := time.Since(start)
	if err != nil {
		span.RecordError(err)
		span.End()
		return nil, err
	}
	affected := int64(0)
	if n, e := res.RowsAffected(); e == nil {
		affected = n
		span.SetAttributes(attribute.Int64("db.rows_affected", n))
	}
	addDBCost(ctx, rw, affected, dur)
	span.End()
	return res, nil
}

func queryCore(ctx context.Context, queryFn func(context.Context) (*sql.Rows, error), verb, rw, table, tpl string) (*TracedRows, error) {
	feature, action := getFeatureTags(ctx)

	ctx, span := dbTracer.Start(ctx, "DB."+verb)
	span.SetAttributes(
		attribute.String("feature", feature),
		attribute.String("action", action),

		attribute.String("db.system", "mariadb"),
		attribute.String("db.operation", verb),
		attribute.String("db.rw", rw),
		attribute.String("db.sql.table", table),
		attribute.String("db.sql.template", tpl),
	)

	rows, err := queryFn(ctx)
	if err != nil {
		span.RecordError(err)
		span.End()
		return nil, err
	}
	return &TracedRows{
		ctx:   ctx,
		rows:  rows,
		span:  span,
		rw:    rw,
		start: time.Now(),
	}, nil
}

// 외부에서 직접 쓰는 래퍼(컨텍스트 있는 버전)
func Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	verb, rw, table := classifySQL(query)
	tpl := normalizeSQL(query)
	recordTemplate(ctx, tpl)
	return execCore(ctx, func(c context.Context) (sql.Result, error) {
		return raw.ExecContext(c, query, args...)
	}, verb, rw, table, tpl)
}

func Query(ctx context.Context, query string, args ...any) (*TracedRows, error) {
	verb, rw, table := classifySQL(query)
	tpl := normalizeSQL(query)
	recordTemplate(ctx, tpl)
	return queryCore(ctx, func(c context.Context) (*sql.Rows, error) {
		return raw.QueryContext(c, query, args...)
	}, verb, rw, table, tpl)
}

func ScanOne(ctx context.Context, dest []any, query string, args ...any) error {
	rows, err := Query(ctx, query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()
	if !rows.Next() {
		if rows.Err() != nil {
			return rows.Err()
		}
		return sql.ErrNoRows
	}
	return rows.Scan(dest...)
}

// Rows 래퍼: Next 호출 횟수로 rows_returned 기록
type TracedRows struct {
	ctx   context.Context
	rows  *sql.Rows
	span  trace.Span
	rw    string
	start time.Time
	count int64
}

func (tr *TracedRows) Next() bool {
	ok := tr.rows.Next()
	if ok {
		tr.count++
	}
	return ok
}
func (tr *TracedRows) Scan(dest ...any) error { return tr.rows.Scan(dest...) }
func (tr *TracedRows) Err() error             { return tr.rows.Err() }
func (tr *TracedRows) Columns() ([]string, error) {
	return tr.rows.Columns()
}
func (tr *TracedRows) ColumnTypes() ([]*sql.ColumnType, error) {
	return tr.rows.ColumnTypes()
}
func (tr *TracedRows) NextResultSet() bool { return tr.rows.NextResultSet() }

func (tr *TracedRows) Close() error {
	err := tr.rows.Close()
	dur := time.Since(tr.start)
	tr.span.SetAttributes(attribute.Int64("db.rows_returned", tr.count))
	addDBCost(tr.ctx, tr.rw, tr.count, dur)
	tr.span.End()
	return err
}

// ---------- 단일 행 SELECT 계측: TracedRow ----------

type TracedRow struct {
	ctx   context.Context
	tx    *sql.Tx // 트랜잭션 경로면 사용
	query string
	args  []any

	rw    string
	table string
	tpl   string

	start time.Time
	span  trace.Span
}

// 공통 생성
func newTracedRow(ctx context.Context, tx *sql.Tx, verb, rw, table, tpl, query string, args []any) *TracedRow {
	ctx, span := dbTracer.Start(ctx, "DB."+verb)
	feature, action := getFeatureTags(ctx)
	span.SetAttributes(
		attribute.String("feature", feature),
		attribute.String("action", action),

		attribute.String("db.system", "mariadb"),
		attribute.String("db.operation", verb),
		attribute.String("db.rw", rw),
		attribute.String("db.sql.table", table),
		attribute.String("db.sql.template", tpl),
	)

	return &TracedRow{
		ctx:   ctx,
		tx:    tx,
		query: query,
		args:  args,
		rw:    rw,
		table: table,
		tpl:   tpl,
		start: time.Now(),
		span:  span,
	}
}

// Scan 시 실제 SELECT 수행 + 집계
func (tr *TracedRow) Scan(dest ...any) error {
	var rows *sql.Rows
	var err error

	if tr.tx != nil {
		rows, err = tr.tx.QueryContext(tr.ctx, tr.query, tr.args...)
	} else {
		rows, err = raw.QueryContext(tr.ctx, tr.query, tr.args...)
	}
	if err != nil {
		tr.span.RecordError(err)
		tr.span.End()
		return err
	}
	defer rows.Close()

	var count int64
	if rows.Next() {
		if err := rows.Scan(dest...); err != nil {
			tr.span.RecordError(err)
			tr.span.End()
			return err
		}
		count = 1
	}
	if err := rows.Err(); err != nil {
		tr.span.RecordError(err)
		tr.span.End()
		return err
	}

	dur := time.Since(tr.start)
	tr.span.SetAttributes(attribute.Int64("db.rows_returned", count))
	addDBCost(tr.ctx, tr.rw, count, dur)
	tr.span.End()

	if count == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// ---------- tracedDB / tracedTx (콜사이트 무수정용) ----------

type tracedDB struct{ inner *sql.DB }

// Context 없는 Exec/Query도 현재 고루틴 컨텍스트로 트레이싱
func (t *tracedDB) Exec(query string, args ...any) (sql.Result, error) {
	ctx := currentCtx()
	verb, rw, table := classifySQL(query)
	tpl := normalizeSQL(query)
	recordTemplate(ctx, tpl)
	return execCore(ctx, func(c context.Context) (sql.Result, error) {
		return t.inner.ExecContext(c, query, args...)
	}, verb, rw, table, tpl)
}
func (t *tracedDB) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return Exec(ctx, query, args...)
}
func (t *tracedDB) Query(query string, args ...any) (*TracedRows, error) {
	ctx := currentCtx()
	verb, rw, table := classifySQL(query)
	tpl := normalizeSQL(query)
	recordTemplate(ctx, tpl)
	return queryCore(ctx, func(c context.Context) (*sql.Rows, error) {
		return t.inner.QueryContext(c, query, args...)
	}, verb, rw, table, tpl)
}
func (t *tracedDB) QueryContext(ctx context.Context, query string, args ...any) (*TracedRows, error) {
	return Query(ctx, query, args...)
}

// 단일 행: *TracedRow 반환 (콜사이트는 .Scan(...)만 쓰므로 영향 없음)
func (t *tracedDB) QueryRow(query string, args ...any) *TracedRow {
	ctx := currentCtx()
	verb, rw, table := classifySQL(query)
	tpl := normalizeSQL(query)
	recordTemplate(ctx, tpl)
	return newTracedRow(ctx, nil, verb, rw, table, tpl, query, args)
}
func (t *tracedDB) QueryRowContext(ctx context.Context, query string, args ...any) *TracedRow {
	verb, rw, table := classifySQL(query)
	tpl := normalizeSQL(query)
	recordTemplate(ctx, tpl)
	return newTracedRow(ctx, nil, verb, rw, table, tpl, query, args)
}

func (t *tracedDB) Begin() (*tracedTx, error) {
	tx, err := t.inner.BeginTx(currentCtx(), nil)
	if err != nil {
		return nil, err
	}
	return &tracedTx{tx: tx}, nil
}
func (t *tracedDB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*tracedTx, error) {
	tx, err := t.inner.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return &tracedTx{tx: tx}, nil
}

// 트랜잭션 래퍼
type tracedTx struct{ tx *sql.Tx }

func (t *tracedTx) Commit() error   { return t.tx.Commit() }
func (t *tracedTx) Rollback() error { return t.tx.Rollback() }

func execWithTx(ctx context.Context, tx *sql.Tx, query string, args ...any) (sql.Result, error) {
	verb, rw, table := classifySQL(query)
	tpl := normalizeSQL(query)
	recordTemplate(ctx, tpl)
	return execCore(ctx, func(c context.Context) (sql.Result, error) {
		return tx.ExecContext(c, query, args...)
	}, verb, rw, table, tpl)
}
func queryWithTx(ctx context.Context, tx *sql.Tx, query string, args ...any) (*TracedRows, error) {
	verb, rw, table := classifySQL(query)
	tpl := normalizeSQL(query)
	recordTemplate(ctx, tpl)
	return queryCore(ctx, func(c context.Context) (*sql.Rows, error) {
		return tx.QueryContext(c, query, args...)
	}, verb, rw, table, tpl)
}

func (t *tracedTx) Exec(query string, args ...any) (sql.Result, error) {
	return execWithTx(currentCtx(), t.tx, query, args...)
}
func (t *tracedTx) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return execWithTx(ctx, t.tx, query, args...)
}
func (t *tracedTx) Query(query string, args ...any) (*TracedRows, error) {
	return queryWithTx(currentCtx(), t.tx, query, args...)
}
func (t *tracedTx) QueryContext(ctx context.Context, query string, args ...any) (*TracedRows, error) {
	return queryWithTx(ctx, t.tx, query, args...)
}

// 트랜잭션 단일 행도 *TracedRow 반환, Scan 시 tx.QueryContext 사용
func (t *tracedTx) QueryRow(query string, args ...any) *TracedRow {
	ctx := currentCtx()
	verb, rw, table := classifySQL(query)
	tpl := normalizeSQL(query)
	recordTemplate(ctx, tpl)
	return newTracedRow(ctx, t.tx, verb, rw, table, tpl, query, args)
}
func (t *tracedTx) QueryRowContext(ctx context.Context, query string, args ...any) *TracedRow {
	verb, rw, table := classifySQL(query)
	tpl := normalizeSQL(query)
	recordTemplate(ctx, tpl)
	return newTracedRow(ctx, t.tx, verb, rw, table, tpl, query, args)
}
