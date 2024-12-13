package api

import (
	"context"
	"fmt"
	"go-project-template/internal/domain"
	"go-project-template/internal/repository"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type api struct {
	logger *zap.Logger
	// statsd     *statsd.Client
	httpClient *http.Client

	userRepo domain.UserRepository
}

func NewAPI(ctx context.Context, logger *zap.Logger, pool *pgxpool.Pool) *api {
	userRepo := repository.NewUserRepository(pool)

	client := &http.Client{}

	return &api{
		logger:     logger,
		httpClient: client,

		userRepo: userRepo,
	}
}

func (a *api) Server(port int) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: a.Routes(),
	}
}

func (a *api) Routes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{os.Getenv("ALLOWED_DOMAIN")},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
	}))

	r.Use(a.loggingMiddleware)
	r.Use(a.requestIdMiddleware)
	r.Use(middleware.Recoverer)

	r.Route("/v1", func(r chi.Router) {
		// Health Check
		r.Get("/health", a.healthCheckHandler)

		r.Route("/users", func(r chi.Router) {
			// Users
			r.Post("/", a.upsertUserHandler)
			r.Delete("/{userid}", a.deleteUserHandler)
			r.Get("/{userid}", a.getByIdUserHandler)
			r.Get("/", a.getUserListHandler)

		})
	})

	return r
}

type LoggingResponseWriter struct {
	w          http.ResponseWriter
	statusCode int
	bytes      int
}

func (lrw *LoggingResponseWriter) Header() http.Header {
	return lrw.w.Header()
}

func (lrw *LoggingResponseWriter) Write(bb []byte) (int, error) {
	wb, err := lrw.w.Write(bb)
	lrw.bytes += wb
	return wb, err
}

func (lrw *LoggingResponseWriter) WriteHeader(statusCode int) {
	lrw.w.WriteHeader(statusCode)
	lrw.statusCode = statusCode
}

func (a *api) requestIdMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		w.Header().Set("X-Project-Request-Id", id)
		next.ServeHTTP(w, r)
	})
}

func (a *api) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip logging health checks
		if r.RequestURI == "/v1/health" {
			next.ServeHTTP(w, r)
			return
		}

		start := time.Now()
		lrw := &LoggingResponseWriter{w: w}

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(lrw, r)

		duration := time.Since(start).Milliseconds()

		remoteAddr := r.Header.Get("X-Forwarded-For")
		if remoteAddr == "" {
			if ip, _, err := net.SplitHostPort(r.RemoteAddr); err != nil {
				remoteAddr = "unknown"
			} else {
				remoteAddr = ip
			}
		}

		fields := []zap.Field{
			zap.Int64("duration", duration),
			zap.String("method", r.Method),
			zap.String("remote#addr", remoteAddr),
			zap.Int("response#bytes", lrw.bytes),
			zap.Int("status", lrw.statusCode),
			zap.String("uri", r.RequestURI),
			zap.String("request#id", lrw.Header().Get("X-Project-Request-Id")),
		}

		if lrw.statusCode == 200 || lrw.statusCode == 201 {
			a.logger.Info("", fields...)
		} else {
			err := lrw.Header().Get("X-Project-Error")
			a.logger.Error(err, fields...)
		}

		// tags := []string{fmt.Sprintf("status:%d", lrw.statusCode)}
		// _ = a.statsd.Histogram("api.latency", float64(duration), nil, 1.0)
		// _ = a.statsd.Incr("api.calls", tags, 1.0)
		// if lrw.statusCode >= 500 {
		// 	_ = a.statsd.Incr("api.errors", nil, 1.0)
		// }
	})
}
