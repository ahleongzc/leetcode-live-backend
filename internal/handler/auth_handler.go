package handler

import (
	"net/http"

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

	sessionID, err := a.authService.Login(ctx, request.Email, request.Password)
	if err != nil {
		HandleErrorResponseHTTP(w, err)
		return
	}

	payload := util.NewJSONPayload()
	payload.Add("session_id", sessionID)

	WriteJSONHTTP(w, payload, http.StatusOK, nil)
}

func (a *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}
