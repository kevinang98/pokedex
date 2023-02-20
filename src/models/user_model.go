package models

import (
	"database/sql"
)

type User struct {
	ID              int            `json:"id" db:"id"`
	Username        string         `json:"username" db:"username"`
	Password        string         `db:"password"`
	Role            string         `json:"role" db:"role"`
	CapturedPokemon string         `json:"captured_pokemon" db:"captured_pokemon"`
	CreatedAt       sql.NullString `db:"created_at"`
	UpdatedAt       sql.NullString `db:"updated_at"`
}

type UserDto struct {
	ID              int    `json:"id"`
	Username        string `json:"username"`
	Role            string `json:"role"`
	CapturedPokemon string `json:"captured_pokemon"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

type PokemonCaptured struct {
	PID []int `json:"pid"`
}
