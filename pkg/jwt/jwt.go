package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtService struct {
	SecretKey []byte
	TTL       int
}

type JwtClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func NewJwtService(secretKey string, ttl int) *JwtService {
	return &JwtService{
		SecretKey: []byte(secretKey),
		TTL:       ttl,
	}
}

func (s *JwtService) GenerateToken(userID, role string) (string, error) {
	// Implementation for signing the JWT
	now := time.Now()
	claims := JwtClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(s.TTL) * time.Second)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(s.SecretKey)
}

func (s *JwtService) Parse(tokenString string) (*JwtClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return s.SecretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JwtClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrTokenInvalidClaims
}
