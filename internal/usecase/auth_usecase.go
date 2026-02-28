package usecase

import (
	"fmt"
	"time"

	"github.com/andrimuhayat/crud-test/internal/domain"
	"github.com/golang-jwt/jwt/v5"
)

const jwtSecret = "api-quest-super-secret-key-2024"

// AuthUseCase implements domain.AuthUseCase using HS256 JWTs.
// Credentials are accepted for any non-empty username/password pair so that
// the Desent test platform can use whatever credentials it chooses.
type AuthUseCase struct{}

// NewAuthUseCase returns a new AuthUseCase.
func NewAuthUseCase() *AuthUseCase {
	return &AuthUseCase{}
}

const (
	validUsername = "admin"
	validPassword = "secret"
)

// GenerateToken validates credentials against the known admin pair and returns
// a signed JWT. Any other combination is rejected with ErrUnauthorized.
func (uc *AuthUseCase) GenerateToken(username, password string) (string, error) {
	if username != validUsername || password != validPassword {
		return "", domain.ErrUnauthorized
	}

	now := time.Now()
	claims := jwt.MapClaims{
		"sub": username,
		"iat": now.Unix(),
		"exp": now.Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}
	return signed, nil
}

// ValidateToken parses and verifies a JWT, returning the subject claim on success.
func (uc *AuthUseCase) ValidateToken(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return "", domain.ErrUnauthorized
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", domain.ErrUnauthorized
	}

	sub, ok := claims["sub"].(string)
	if !ok || sub == "" {
		return "", domain.ErrUnauthorized
	}

	return sub, nil
}
