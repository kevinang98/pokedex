package controllers

import (
	mMocks "pokedex/src/mocks/middleware"
	mocks "pokedex/src/mocks/services"
)

type serviceMock struct {
	pokedexSvc *mocks.PokedexService
	userSvc    *mocks.UserService
	mid        *mMocks.Middleware
}

func initServiceMock() serviceMock {
	pokedexSvc := new(mocks.PokedexService)
	userSvc := new(mocks.UserService)
	mid := new(mMocks.Middleware)

	return serviceMock{
		pokedexSvc: pokedexSvc,
		userSvc:    userSvc,
		mid:        mid,
	}
}

func initControllerMock(s serviceMock) *Controller {
	return NewController(s.pokedexSvc, s.userSvc, s.mid)
}
