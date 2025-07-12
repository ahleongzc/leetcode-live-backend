package handler

import (
	"net/http"

	"github.com/ahleongzc/leetcode-live-backend/internal/config"
	"github.com/ahleongzc/leetcode-live-backend/internal/service"
	"github.com/ahleongzc/leetcode-live-backend/internal/util"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(
	authService service.AuthService,
) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (a *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	request := &struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	err := ReadJSONHTTPReq(w, r, request)
	if err != nil {
		HandleErrorResponseHTTP(w, err)
		return
	}

	sessionToken, err := a.authService.Login(ctx, request.Email, request.Password)
	if err != nil {
		HandleErrorResponseHTTP(w, err)
		return
	}

	headers := http.Header{}
	headers.Set(config.SESSION_TOKEN_HEADER_KEY, sessionToken)

	WriteJSONHTTP(w, nil, http.StatusOK, headers)
}

func (a *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sessionToken, err := util.GetSessionToken(ctx)
	if err != nil {
		HandleErrorResponseHTTP(w, err)
		return
	}

	if err := a.authService.Logout(ctx, sessionToken); err != nil {
		HandleErrorResponseHTTP(w, err)
		return
	}

	WriteJSONHTTP(w, nil, http.StatusOK, nil)
}
