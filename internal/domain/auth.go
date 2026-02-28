package domain

// AuthUseCase defines the business-logic contract for authentication.
type AuthUseCase interface {
	// GenerateToken validates credentials and returns a signed JWT.
	GenerateToken(username, password string) (string, error)
	// ValidateToken parses and validates a JWT, returning the subject claim.
	ValidateToken(token string) (string, error)
}

