package services

import (
	Repository "pokedex/src/repositories"
)

type Service struct {
	P PokedexService
	U UserService
}

func NewService(pr Repository.PokedexRepository, ur Repository.UserRepository) *Service {
	return &Service{
		P: NewPokedexService(pr, ur),
		U: NewUserService(ur),
	}
}
