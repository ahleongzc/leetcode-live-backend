package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
	"github.com/ahleongzc/leetcode-live-backend/internal/handler"
	"github.com/ahleongzc/leetcode-live-backend/internal/util"

	"github.com/rs/zerolog"
)

type Middleware struct {
	logger *zerolog.Logger
}

func NewMiddleware(logger *zerolog.Logger) *Middleware {
	return &Middleware{
		logger: logger,
	}
}

func (m *Middleware) CORS(next http.Handler) http.Handler {
	var trustedOrigins map[string]struct{}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		if util.IsProdEnv() {
			if _, exists := trustedOrigins[origin]; !exists {
				handler.HandleErrorResponseHTTP(w, common.ErrForbidden)
				return
			}
		}

		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", fmt.Sprintf("%s, %s, %s, %s, %s",
			http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions))
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, Content-Length, X-Session-Id")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) RecordRequestTimestampMS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := util.SetStartRequestTimestampMS(r.Context())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *Middleware) Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.logger.Info().
			Str("method", r.Method).
			Str("origin", r.Header.Get("Origin")).
			Str("url", r.URL.String()).
			Str("remote_addr", r.RemoteAddr).
			Str("user_agent", r.UserAgent()).
			Msg("")

		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				stackTrace := debug.Stack()
				m.logger.Error().
					Interface("panic", err).
					Bytes("stackTrace", stackTrace).
					Msg("panic recovered in request")
				handler.HandleErrorResponseHTTP(w, fmt.Errorf("%w", common.ErrInternalServerError))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
