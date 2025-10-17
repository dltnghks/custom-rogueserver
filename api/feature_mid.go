package api

import (
	"context"
	"net/http"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const (
	HdrFeature = "X-Pokerogue-Feature"
	HdrAction  = "X-Pokerogue-Action" // 선택
)

type featureKey struct{}

type featureCtx struct {
	Feature string
	Action  string
	start   time.Time
}

func withFeature(ctx context.Context, f featureCtx) context.Context {
	return context.WithValue(ctx, featureKey{}, f)
}

// 다른 패키지(예: db)에서 읽어갈 수 있게 공개 함수로
func FeatureFrom(ctx context.Context) (feature, action string, ok bool) {
	if v, ok := ctx.Value(featureKey{}).(featureCtx); ok {
		return v.Feature, v.Action, true
	}
	return "", "", false
}

var featTracer = otel.Tracer("pokerogue/feature")

// after: 요청 총소요시간(ms)을 받아 최종 속성 찍는 콜백
func FeatureMiddleware(after func(ctx context.Context, totalMs int64)) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			feat := strings.TrimSpace(r.Header.Get(HdrFeature))
			if feat == "" {
				feat = "Unknown"
			}
			act := strings.TrimSpace(r.Header.Get(HdrAction))

			ctx, span := featTracer.Start(r.Context(), "Feature."+feat)
			defer span.End()
			span.SetAttributes(
				attribute.String("feature", feat),
				attribute.String("action", act),
			)

			f := featureCtx{Feature: feat, Action: act, start: time.Now()}
			ctx = withFeature(ctx, f)

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)

			if after != nil {
				totalMs := time.Since(f.start).Milliseconds()
				after(ctx, totalMs)
			}
		})
	}
}

// 활성 스팬 꺼내고 싶을 때(선택)
func CurrentSpan(ctx context.Context) trace.Span {
	return trace.SpanFromContext(ctx)
}
