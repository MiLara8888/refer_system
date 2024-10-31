package refersrest

import (
	storage "refers_rest/pkg/storage/refersdb"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// генерация токенa
func GenerateTokenJwt(email string, secretKey string, exp time.Duration) (string, error) {
	claims := storage.JwtToken{
		CreateTime: int(time.Now().Unix()),
		ExpDate:    int(time.Now().Add(exp).Unix()),
		Email:      email,
		RegisteredClaims: jwt.RegisteredClaims{
			// A usual scenario is to set the expiration time relative to the current time
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(exp)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

