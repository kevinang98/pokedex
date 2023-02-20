package controllers

import (
	"pokedex/src/middleware"
	Service "pokedex/src/services"
)

type Controller struct {
	P PokedexController
	U UserController
}

func NewController(ps Service.PokedexService, us Service.UserService, m middleware.Middleware) *Controller {
	return &Controller{
		P: NewPokedexController(ps, m),
		U: NewUserController(us, m),
	}
}
