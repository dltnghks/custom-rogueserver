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

	"github.com/bwmarrin/discordgo"
	"github.com/pagefaultgames/rogueserver/api"
	"github.com/pagefaultgames/rogueserver/api/account"
	"github.com/pagefaultgames/rogueserver/db"
	"github.com/pagefaultgames/rogueserver/dbcount"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	// CHANGED: HTTP 익스포터 대신 gRPC 익스포터를 사용합니다.
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
	// ADDED: gRPC 옵션을 위해 추가합니다.
)

// initTracerProvider 함수를 수정하여 환경 변수를 사용하도록 변경합니다.
func initTracerProvider() (*sdktrace.TracerProvider, error) {
	ctx := context.Background()

	// CHANGED: 하드코딩된 주소 대신 환경 변수에서 Jaeger 엔드포인트를 읽어옵니다.
	// docker-compose.yml에 설정된 OTEL_EXPORTER_OTLP_TRACES_ENDPOINT 값을 사용합니다.
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT")
	if endpoint == "" {
		endpoint = "jaeger:4317" // 환경 변수가 없을 경우의 기본값
	}

	log.Printf("Initializing OTLP gRPC exporter with endpoint: %s", endpoint)

	// CHANGED: OTLP/HTTP 익스포터를 OTLP/gRPC 익스포터로 변경합니다.
	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithInsecure(),         // Docker 내부 통신이므로 TLS 없이 연결
		otlptracegrpc.WithEndpoint(endpoint), // 환경 변수에서 읽어온 주소 사용
	)
	if err != nil {
		return nil, err
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String("pokerogue-api-server"), // Jaeger UI에 표시될 서비스 이름
		),
	)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
	return tp, nil
}

// featureContextMiddleware 함수는 변경 없음
func featureContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		feature := r.Header.Get("X-Pokerogue-Feature")
		if feature != "" {
			span := trace.SpanFromContext(r.Context())
			span.SetAttributes(attribute.String("pokerogue.feature", feature))
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	// OpenTelemetry Tracer Provider 초기화
	tp, err := initTracerProvider()
	if err != nil {
		log.Fatalf("failed to initialize tracer provider: %s", err)
	}
	otel.SetTracerProvider(tp)

	// 애플리케이션 종료 시 Tracer Provider를 안전하게 종료합니다.
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	// ... (환경 변수 설정 부분은 변경 없음) ...
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

	// get database connection
	err = db.Init(dbuser, dbpass, dbproto, dbaddr, dbname)
	if err != nil {
		log.Fatalf("failed to initialize database: %s", err)
	}

	// create listener
	listener, err := createListener(proto, addr)
	if err != nil {
		log.Fatalf("failed to create net listener: %s", err)
	}

	mux := http.NewServeMux()

	// init api
	if err := api.Init(mux); err != nil {
		log.Fatal(err)
	}

	// start web server
	handler := prodHandler(mux, gameurl)
	if debug {
		handler = debugHandler(mux)
	}

	// 미들웨어 체인 적용
	finalHandler := otelhttp.NewHandler(featureContextMiddleware(handler), "http-server")

	// 중복된 서버 시작 코드를 제거하고 하나로 합칩니다.
	if tlscert == "" {
		err = http.Serve(listener, finalHandler)
	} else {
		err = http.ServeTLS(listener, finalHandler, tlscert, tlskey)
	}
	if err != nil {
		log.Fatalf("failed to create http server or server errored: %s", err)
	}
}

// createListener 함수는 변경 없음
func createListener(proto, addr string) (net.Listener, error) {
	if proto == "unix" {
		os.Remove(addr)
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
		// CORS 허용 헤더에 'X-Pokerogue-Feature' 추가
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, X-Pokerogue-Feature")
		w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST")
		w.Header().Set("Access-Control-Allow-Origin", clienturl)

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		router.ServeHTTP(w, r)
	})
}

// debugHandler 함수는 변경 없음
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

// getEnv 함수는 변경 없음
func getEnv(key string, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}
