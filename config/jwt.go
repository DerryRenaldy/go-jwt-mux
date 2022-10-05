package config

import "github.com/golang-jwt/jwt/v4"

var JWT_KEY = []byte("kjsdhfiuadj3i298398hf9a8hdnd")

type JWTClaim struct {
	Username string
	jwt.RegisteredClaims
}
