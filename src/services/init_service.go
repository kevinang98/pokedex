package services

import (
	mocks "pokedex/src/mocks/repositories"
)

type repositoryMock struct {
	pokedexRepo *mocks.PokedexRepository
	userRepo    *mocks.UserRepository
}

func initRepositoryMock() repositoryMock {
	pokedexRepo := new(mocks.PokedexRepository)
	userRepo := new(mocks.UserRepository)

	return repositoryMock{
		pokedexRepo: pokedexRepo,
		userRepo:    userRepo,
	}
}

func initServiceMock(r repositoryMock) *Service {
	return NewService(r.pokedexRepo, r.userRepo)
}
