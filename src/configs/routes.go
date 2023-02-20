package configs

import (
	"net/http"
	Controller "pokedex/src/controllers"
	"pokedex/src/models"

	"github.com/golang-jwt/jwt/v4"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

type Route struct {
	app    *echo.Echo
	router *echo.Group
	pc     Controller.PokedexController
	uc     Controller.UserController
}

func NewRoute(app *echo.Echo, pc Controller.PokedexController, uc Controller.UserController) *Route {
	return &Route{
		app:        app,
		router:     app.Group(""),
		pc: pc,
		uc: uc,
	}
}

func (r *Route) Init() {
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(models.JwtCustomClaims)
		},
		SigningKey: []byte(models.JwtKey),
	}

	user := r.router.Group("/user")
	user.POST("/register", r.uc.RegisterUser)
	user.POST("/login", r.uc.LoginUser)

	userJwt := r.router.Group("/user")
	userJwt.Use(echojwt.WithConfig(config))
	userJwt.PUT("/captured-pokemon", r.uc.UpdateCapturedPokemonUser)

	pokedex := r.router.Group("/pokedex")
	pokedex.Use(echojwt.WithConfig(config))
	pokedex.GET("", r.pc.GetPokemon)
	pokedex.POST("", r.pc.AddPokemon)
	pokedex.PUT("", r.pc.UpdatePokemon)
	pokedex.DELETE("", r.pc.DeletePokemon)

	r.RouteNotFound()
}

func (r *Route) RouteNotFound() {
	r.router.Any("*", func(e echo.Context) error {
		return e.JSON(http.StatusNotFound, models.Response{
			Status: "Fail",
			Data:   "Route not found",
		})
	})
}
