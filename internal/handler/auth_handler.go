package handler

import (
	"net/http"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
	"github.com/ahleongzc/leetcode-live-backend/internal/service"
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
	headers.Set(common.SESSION_TOKEN_HEADER_KEY, sessionToken)

	WriteJSONHTTP(w, nil, http.StatusOK, headers)
}

func (a *AuthHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	WriteJSONHTTP(w, nil, http.StatusOK, nil)
}

func (a *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}
