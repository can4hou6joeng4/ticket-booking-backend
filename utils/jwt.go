package utils

import "github.com/golang-jwt/jwt/v5"

func GenerateJWT(claims jwt.MapClaims, method jwt.SigningMethod, jwtSecret string) (string, error) {
	return jwt.NewWithClaims(method, claims).SignedString([]byte(jwtSecret))
}
