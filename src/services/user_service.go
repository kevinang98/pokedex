package services

import (
	"errors"
	"pokedex/src/models"
	Repository "pokedex/src/repositories"
	"strings"
)

type UserService interface {
	RegisterUser(username, password, role string) error
	LoginUser(username, password string) (*models.User, error)
	UpdateCapturedPokemonUser(id int, pid string) error
}

type UserServiceImpl struct {
	ur Repository.UserRepository
}

func NewUserService(ur Repository.UserRepository) UserService {
	return &UserServiceImpl{
		ur: ur,
	}
}

func (u *UserServiceImpl) RegisterUser(username, password, role string) error {
	if err := u.ur.RegisterUser(username, password, role); err != nil {
		return err
	}

	return nil
}

func (u *UserServiceImpl) LoginUser(username, password string) (*models.User, error) {
	var (
		user models.User
	)

	param := make(map[string]interface{}, 0)
	param["username"] = username

	user, err := u.ur.GetUserData(param)
	if err != nil {
		return nil, err
	}

	if (user == models.User{}) {
		return nil, errors.New(models.ErrorUsernameNotExist)
	} else if user.Password != password {
		return nil, errors.New(models.ErrorInvalidPassword)
	}

	return &user, nil
}

func (u *UserServiceImpl) UpdateCapturedPokemonUser(id int, pid string) error {
	list := make(map[string]bool)
	listID := []string{}

	param := make(map[string]interface{}, 0)
	param["id"] = id

	user, err := u.ur.GetUserData(param)
	if err != nil {
		return err
	}

	pids := append(strings.Split(pid, ","), strings.Split(user.CapturedPokemon, ",")...)

	for _, id := range pids {
		if _, value := list[id]; !value {
			if strings.ReplaceAll(id, " ", "") != "" {
				list[id] = true
				listID = append(listID, id)
			}
		}
	}

	capturedID := strings.Join(listID, ",")

	err = u.ur.UpdateCapturedPokemonUser(id, capturedID)
	if err != nil {
		return err
	}

	return nil
}
