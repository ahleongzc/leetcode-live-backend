package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
	"github.com/ahleongzc/leetcode-live-backend/internal/config"
	httphandler "github.com/ahleongzc/leetcode-live-backend/internal/handler/http_handler"
	"github.com/ahleongzc/leetcode-live-backend/internal/service"
	"github.com/ahleongzc/leetcode-live-backend/internal/util"

	"github.com/rs/zerolog"
)

type Middleware struct {
	authService service.AuthService
	logger      *zerolog.Logger
}

func NewMiddleware(
	authService service.AuthService,
	logger *zerolog.Logger,
) *Middleware {
	return &Middleware{
		logger:      logger,
		authService: authService,
	}
}

func (m *Middleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionToken := r.Header.Get(config.SESSION_TOKEN_HEADER_KEY)
		updatedSessionToken, err := m.authService.ValidateAndRefreshSessionToken(ctx, sessionToken)
		if err != nil {
			httphandler.HandleErrorResponseHTTP(w, err)
			return
		}

		ctx = util.SetSessionToken(ctx, updatedSessionToken)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) CORS(next http.Handler) http.Handler {
	var trustedOrigins map[string]struct{}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		if util.IsProdEnv() {
			if _, exists := trustedOrigins[origin]; !exists {
				httphandler.HandleErrorResponseHTTP(w, common.ErrForbidden)
				return
			}
		}

		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", fmt.Sprintf("%s, %s, %s, %s, %s",
			http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions))
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, Content-Length, X-Session-Token, X-Interview-Token")
		w.Header().Set("Access-Control-Expose-Headers", "X-Session-Token, X-Interview-Token")
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
				httphandler.HandleErrorResponseHTTP(w, fmt.Errorf("%w", common.ErrInternalServerError))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) SetUserID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionToken, err := util.GetSessionToken(ctx)
		if err != nil {
			httphandler.HandleErrorResponseHTTP(w, err)
			return
		}

		userID, err := m.authService.GetUserIDFromSessionToken(ctx, sessionToken)
		if err != nil {
			httphandler.HandleErrorResponseHTTP(w, err)
			return
		}

		ctx = util.SetUserID(ctx, userID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) SetSessionTokenInResponseHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionToken, err := util.GetSessionToken(ctx)
		if err != nil {
			httphandler.HandleErrorResponseHTTP(w, err)
			return
		}

		w.Header().Set(config.SESSION_TOKEN_HEADER_KEY, sessionToken)
		next.ServeHTTP(w, r)
	})
}
