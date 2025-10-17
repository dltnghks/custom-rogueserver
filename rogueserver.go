/*
   Copyright (C) 2024  Pagefault Games

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package main

import (
	"context"
	"encoding/gob"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/pagefaultgames/rogueserver/api"
	"github.com/pagefaultgames/rogueserver/api/account"
	"github.com/pagefaultgames/rogueserver/db"
	"github.com/pagefaultgames/rogueserver/dbcount"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	// OTLP gRPC 익스포터
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

// -------- OpenTelemetry TracerProvider 초기화 --------

func initTracerProvider() (*sdktrace.TracerProvider, error) {
	ctx := context.Background()

	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT")
	if endpoint == "" {
		endpoint = "jaeger:4317"
	}
	log.Printf("Initializing OTLP gRPC exporter with endpoint: %s", endpoint)

	exp, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(endpoint),
	)
	if err != nil {
		return nil, err
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String("pokerogue-api-server"),
		),
	)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(res),
	)
	return tp, nil
}

// -------- 기능 컨텍스트 + 집계 + 요약 미들웨어 --------

func featureContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		feature := r.Header.Get("X-Pokerogue-Feature")
		action := r.Header.Get("X-Pokerogue-Action")

		ctx := r.Context()
		ctx = db.WithFeatureTags(ctx, feature, action)
		ctx = db.WithFeatureAgg(ctx, &db.FeatureAgg{})

		// otelhttp가 parent 요청 스팬을 만들어둠 → 태그 부착
		span := trace.SpanFromContext(ctx)
		if feature != "" {
			span.SetAttributes(
				attribute.String("feature", feature),
				attribute.String("action", action),
			)
		}

		// ★ 고루틴-로컬 컨텍스트 push/pop (콜사이트 무수정 핵심)
		db.PushCtxForMiddleware(ctx)
		defer db.PopCtxForMiddleware()

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)

		// 응답 직전 요약 속성 기록
		totalMs := time.Since(start).Milliseconds()
		agg := db.GetAgg(ctx)
		if agg != nil {
			ioMs := agg.DBDurationMs + agg.RedisDurationMs
			if ioMs < 0 {
				ioMs = 0
			}
			cpuMs := totalMs - ioMs
			if cpuMs < 0 {
				cpuMs = 0
			}

			bound := "mixed"
			if totalMs > 0 {
				ratio := float64(ioMs) / float64(totalMs)
				if ratio >= 0.6 {
					bound = "io-bound"
				} else if ratio <= 0.4 {
					bound = "cpu-bound"
				}
			}

			span.SetAttributes(
				attribute.Int64("io.db.read.count", agg.DBReadCount),
				attribute.Int64("io.db.write.count", agg.DBWriteCount),
				attribute.Int64("io.db.read.rows", agg.DBReadRows),
				attribute.Int64("io.db.write.rows", agg.DBWriteRows),
				attribute.Int64("io.db.duration.ms", agg.DBDurationMs),

				attribute.Int64("io.redis.read.count", agg.RedisReadCount),
				attribute.Int64("io.redis.write.count", agg.RedisWriteCount),
				attribute.Int64("io.redis.duration.ms", agg.RedisDurationMs),

				attribute.Int64("total.duration.ms", totalMs),
				attribute.Int64("cpu.duration.ms", cpuMs),
				attribute.String("workload.bound", bound),

				attribute.Int64("db.n_plus_1.count", agg.NPlus1),
			)
		}
	})
}

// --------------------------- main ---------------------------

func main() {
	tp, err := initTracerProvider()
	if err != nil {
		log.Fatalf("failed to initialize tracer provider: %s", err)
	}
	otel.SetTracerProvider(tp)
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	debug, _ := strconv.ParseBool(os.Getenv("debug"))
	dbcount.LoadCSVFile("/app/csv/credentials.csv")
	proto := getEnv("proto", "tcp")
	addr := getEnv("addr", "0.0.0.0:8001")
	tlscert := getEnv("tlscert", "")
	tlskey := getEnv("tlskey", "")
	dbuser := getEnv("dbuser", "pokerogue")
	dbpass := getEnv("dbpass", "pokerogue")
	dbproto := getEnv("dbproto", "tcp")
	dbaddr := getEnv("dbaddr", "localhost")
	dbname := getEnv("dbname", "pokeroguedb")
	discordclientid := getEnv("discordclientid", "")
	discordsecretid := getEnv("discordsecretid", "")
	googleclientid := getEnv("googleclientid", "")
	googlesecretid := getEnv("googlesecretid", "")
	callbackurl := getEnv("callbackurl", "http://localhost:8001/")
	gameurl := getEnv("gameurl", "https://pokerogue.net")
	discordbottoken := getEnv("discordbottoken", "")
	discordguildid := getEnv("discordguildid", "")

	account.GameURL = gameurl
	account.DiscordClientID = discordclientid
	account.DiscordClientSecret = discordsecretid
	account.DiscordCallbackURL = callbackurl + "/auth/discord/callback"
	account.GoogleClientID = googleclientid
	account.GoogleClientSecret = googlesecretid
	account.GoogleCallbackURL = callbackurl + "/auth/google/callback"
	account.DiscordSession, _ = discordgo.New("Bot " + discordbottoken)
	account.DiscordGuildID = discordguildid
	gob.Register([]interface{}{})
	gob.Register(map[string]interface{}{})

	// DB 연결
	err = db.Init(dbuser, dbpass, dbproto, dbaddr, dbname)
	if err != nil {
		log.Fatalf("failed to initialize database: %s", err)
	}

	// Listener
	listener, err := createListener(proto, addr)
	if err != nil {
		log.Fatalf("failed to create net listener: %s", err)
	}

	// 라우터
	mux := http.NewServeMux()

	// API 초기화
	if err := api.Init(mux); err != nil {
		log.Fatal(err)
	}

	// 핸들러 체인
	var base http.Handler = prodHandler(mux, gameurl)
	if debug {
		base = debugHandler(mux)
	}

	// 체인: otelhttp → featureContextMiddleware → base
	finalHandler := otelhttp.NewHandler(featureContextMiddleware(base), "http-server")

	// 서버 시작
	if tlscert == "" {
		err = http.Serve(listener, finalHandler)
	} else {
		err = http.ServeTLS(listener, finalHandler, tlscert, tlskey)
	}
	if err != nil {
		log.Fatalf("failed to create http server or server errored: %s", err)
	}
}

// ------------------------ helpers ------------------------

func createListener(proto, addr string) (net.Listener, error) {
	if proto == "unix" {
		_ = os.Remove(addr)
	}
	listener, err := net.Listen(proto, addr)
	if err != nil {
		return nil, err
	}
	if proto == "unix" {
		if err := os.Chmod(addr, 0777); err != nil {
			listener.Close()
			return nil, err
		}
	}
	return listener, nil
}

func prodHandler(router *http.ServeMux, clienturl string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Headers",
			"Authorization, Content-Type, X-Pokerogue-Feature, X-Pokerogue-Action")
		w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST")
		w.Header().Set("Access-Control-Allow-Origin", clienturl)

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		router.ServeHTTP(w, r)
	})
}

func debugHandler(router *http.ServeMux) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		router.ServeHTTP(w, r)
	})
}

func getEnv(key string, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}
