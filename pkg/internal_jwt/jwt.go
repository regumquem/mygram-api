package internal_jwt

import (
	"errors"

	"github.com/golang-jwt/jwt"
)

const secretKey = "rahasia"

func ParseToken(jwtToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(jwtToken, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token")
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, errors.New("invalid token")
	}

	return token, nil
}

func GenerateToken(jwtClaim jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaim)
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func ValidateToken(jwtToken string) (jwt.MapClaims, error) {
	invalidTokenMessage := "invalid token"

	token, err := ParseToken(jwtToken)
	if err != nil {
		return nil, errors.New(invalidTokenMessage)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New(invalidTokenMessage)
	}

	return claims, nil
}
