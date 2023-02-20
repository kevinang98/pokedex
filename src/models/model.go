package models

import (
	"github.com/golang-jwt/jwt/v4"
)

type Request struct {
	Data interface{} `json:"data"`
}

type Response struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

type JwtCustomClaims struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}
