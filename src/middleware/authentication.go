package middleware

import (
	"pokedex/src/models"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type Middleware interface {
	GenerateToken(id int, username, role string) (string, error)
	ParseJwtToken(c echo.Context) (int, string, string)
	ValidatePokedexBody(poke []models.Pokedex) error
	ValidateIntPokedexBody(poke []int) error
}

type MiddlewareImpl struct {
}

func NewMiddleware() Middleware {
	return &MiddlewareImpl{}
}

func (m *MiddlewareImpl) GenerateToken(id int, username, role string) (string, error) {
	claims := &models.JwtCustomClaims{
		ID:       id,
		Username: username,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(models.JwtKey))
	if err != nil {
		return "", err
	}

	return t, nil
}

func (m *MiddlewareImpl) ParseJwtToken(c echo.Context) (int, string, string) {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*models.JwtCustomClaims)
	id := claims.ID
	username := claims.Username
	role := claims.Role

	return id, username, role
}
