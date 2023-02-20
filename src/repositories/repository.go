package repositories

import (
	"github.com/jmoiron/sqlx"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	P PokedexRepository
	U UserRepository
}

func NewRepository(mysql *sqlx.DB, mongo *mongo.Client) *Repository {
	return &Repository{
		P: NewPokedexRepository(mongo),
		U: NewUserRepository(mysql),
	}
}
