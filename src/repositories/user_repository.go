package repositories

import (
	"errors"
	"fmt"
	"pokedex/src/models"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	RegisterUser(username, password, role string) error
	GetUserData(req map[string]interface{}) (models.User, error)
	UpdateCapturedPokemonUser(id int, captured string) error
}

type UserRepositoryImpl struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &UserRepositoryImpl{
		db: db,
	}
}

const (
	RegisterUser = `
	INSERT INTO ` + "`user`" + `
	(username, pass, role, captured_pokemon, created_at, updated_at)
	VALUES(?, AES_ENCRYPT(?,'` + models.UserKey + `'), ?, '', NOW(), '');
	`

	GetUserData = `
	SELECT id, username, AES_DECRYPT(pass, '` + models.UserKey + `') AS password, role, captured_pokemon, created_at, updated_at
	FROM ` + "`user`" + `
	<CONDITION>`

	UpdateUserCapturedPokemon = `
	UPDATE ` + "`user`" + `
	SET captured_pokemon = ?, updated_at = NOW()
	WHERE id = ?`
)

func (u *UserRepositoryImpl) RegisterUser(username, password, role string) error {
	_, err := u.db.Exec(RegisterUser, username, password, role)
	if err != nil {
		m, ok := err.(*mysql.MySQLError)
		if !ok {
			return err
		}
		if m.Number == 1062 {
			return errors.New("username already used")
		}
		return err
	}

	return nil
}

func (u *UserRepositoryImpl) GetUserData(req map[string]interface{}) (models.User, error) {
	var (
		user           models.User
		paramCondition []string
		paramValue     []interface{}
	)

	query := GetUserData

	if len(req) <= 0 {
		query = strings.ReplaceAll(query, "<CONDITION>", "")
	} else {
		query = strings.ReplaceAll(query, "<CONDITION>", "WHERE ")

		for key, value := range req {
			paramCondition = append(paramCondition, fmt.Sprintf("%s = ?", key))
			paramValue = append(paramValue, value)
		}
	}

	rows, err := u.db.Queryx(query+strings.Join(paramCondition, " AND "), paramValue...)
	if err != nil {
		return user, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.StructScan(&user)
		if err != nil {
			return user, err
		}
	}

	return user, err
}

func (u *UserRepositoryImpl) UpdateCapturedPokemonUser(id int, captured string) error {
	_, err := u.db.Exec(UpdateUserCapturedPokemon, captured, id)
	if err != nil {
		return err
	}

	return nil
}
