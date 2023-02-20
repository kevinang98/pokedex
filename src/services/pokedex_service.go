package services

import (
	"pokedex/src/models"
	Repository "pokedex/src/repositories"
	"strconv"
	"strings"
)

type PokedexService interface {
	GetPokemon(id int, name, sort, option string, pokeType []string, offset, limit, order int64) ([]models.Pokedex, error)
	AddPokemon(pokeList []models.Pokedex) error
	UpdatePokemon(pokeList []models.Pokedex) (string, error)
	DeletePokemon(pidList []int) (string, error)
}

type PokedexServiceImpl struct {
	pr Repository.PokedexRepository
	ur Repository.UserRepository
}

func NewPokedexService(pr Repository.PokedexRepository, ur Repository.UserRepository) PokedexService {
	return &PokedexServiceImpl{
		pr: pr,
		ur: ur,
	}
}

func (p *PokedexServiceImpl) GetPokemon(id int, name, sort, option string, pokeType []string, offset, limit, order int64) ([]models.Pokedex, error) {
	var (
		listPID []int
		err     error
	)

	param := make(map[string]interface{}, 0)
	param["id"] = id

	user, err := p.ur.GetUserData(param)
	if err != nil {
		return nil, err
	}

	listID := strings.Split(user.CapturedPokemon, ",")
	mapListID := make(map[int]string, 0)

	for _, id := range listID {
		pid := 0

		pid, err := strconv.Atoi(id)
		if err != nil {
			return nil, err
		}

		if pid != 0 {
			mapListID[pid] = id
			listPID = append(listPID, pid)
		}
	}

	pokeList, err := p.pr.GetPokemon(name, sort, option, listPID, pokeType, offset, limit, order)
	if err != nil {
		return nil, err
	}

	for i, poke := range pokeList {
		if _, ok := mapListID[poke.PID]; ok {
			pokeList[i].Captured = true
		} else {
			pokeList[i].Captured = false
		}
	}

	return pokeList, nil
}

func (p *PokedexServiceImpl) AddPokemon(pokeList []models.Pokedex) error {
	if err := p.pr.AddPokemon(pokeList); err != nil {
		return err
	}

	return nil
}

func (p *PokedexServiceImpl) UpdatePokemon(pokeList []models.Pokedex) (string, error) {
	var (
		updatedID string
		err       error
	)

	for _, pl := range pokeList {
		if err = p.pr.UpdatePokemon(pl); err != nil {
			return updatedID, err
		}
		updatedID += updatedID + "," + strconv.Itoa(pl.PID)
	}

	return "", nil
}

func (p *PokedexServiceImpl) DeletePokemon(pidList []int) (string, error) {
	var (
		deletedID string
		err       error
	)

	for _, pl := range pidList {
		if err = p.pr.DeletePokemon(pl); err != nil {
			return deletedID, err
		}
		deletedID += deletedID + "," + strconv.Itoa(pl)
	}

	return "", nil
}
