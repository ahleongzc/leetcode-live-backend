package scenario

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/google/uuid"
)

type AuthScenario interface {
	GenerateRandomToken() string
}

func NewAuthScenario() AuthScenario {
	return &AuthScenarioImpl{}
}

type AuthScenarioImpl struct{}

func (a *AuthScenarioImpl) GenerateRandomToken() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return uuid.NewString()
	}
	return base64.RawURLEncoding.EncodeToString(b)
}
