package middleware

import (
	"errors"
	"pokedex/src/models"
)

func (m *MiddlewareImpl) ValidatePokedexBody(poke []models.Pokedex) error {
	for _, p := range poke {
		if p.PID <= 0 {
			return errors.New(models.ErrorPIDZero)
		}
	}

	return nil
}

func (m *MiddlewareImpl) ValidateIntPokedexBody(poke []int) error {
	for _, p := range poke {
		if p <= 0 {
			return errors.New(models.ErrorPIDZero)
		}
	}

	return nil
}
