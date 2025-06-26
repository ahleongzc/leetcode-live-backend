package handler

import (
	"net/http"

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

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (a *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	request := &struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	err := ReadJSON(w, r, request)
	if err != nil {
		HandleErrorResponse(w, err)
		return
	}

	sessionID, err := a.authService.Login(ctx, request.Email, request.Password)
	if err != nil {
		HandleErrorResponse(w, err)
		return
	}

	payload := NewJSONPayload()
	payload.Add("session_id", sessionID)

	WriteJSON(w, payload, http.StatusOK, nil)
}
