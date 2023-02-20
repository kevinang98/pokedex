package repositories

import (
	"context"
	"errors"
	"pokedex/src/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PokedexRepository interface {
	GetPokemon(name, sort string, option string, listPID []int, pokeType []string, offset, limit, order int64) ([]models.Pokedex, error)
	AddPokemon(pokelist []models.Pokedex) error
	UpdatePokemon(poke models.Pokedex) error
	DeletePokemon(pid int) error
}

type PokedexRepositoryImpl struct {
	db *mongo.Client
}

func NewPokedexRepository(db *mongo.Client) PokedexRepository {
	return &PokedexRepositoryImpl{
		db: db,
	}
}

func (p *PokedexRepositoryImpl) GetPokemon(name, sort, option string, listPID []int, pokeType []string, offset, limit, order int64) ([]models.Pokedex, error) {
	var (
		po  *models.Pokedex
		pos []models.Pokedex
	)

	ctx := context.Background()
	filter := bson.M{}
	options := options.Find()

	pokeCol := p.db.Database(models.PokemonDB).Collection(models.PokemonCol)

	if name != "" {
		filter["name"] = primitive.Regex{Pattern: ".*" + name + ".*"}
	}

	if len(pokeType) > 0 {
		filter["type"] = bson.M{
			"$in": pokeType,
		}
	}

	if option == "catched" {
		filter["pid"] = bson.M{
			"$in": listPID,
		}
	} else if option == "uncatched" {
		filter["pid"] = bson.M{
			"$nin": listPID,
		}
	}

	if sort != "" {
		options.SetSort(bson.M{sort: order})
	}

	if limit > 0 {
		options.SetSkip(offset)
		options.SetLimit(limit)
	}

	cursor, err := pokeCol.Find(ctx, filter, options)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("no result")
	} else if err != nil {
		return nil, err
	}

	for cursor.Next(ctx) {
		po = nil
		err = cursor.Decode(&po)
		if err != nil {
			return nil, err
		}
		pos = append(pos, *po)
	}

	return pos, nil
}

func (p *PokedexRepositoryImpl) AddPokemon(pokelist []models.Pokedex) error {
	var (
		l    interface{}
		list []interface{}
	)
	ctx := context.Background()

	for _, poke := range pokelist {
		l = nil
		l = poke
		list = append(list, l)
	}

	pokeCol := p.db.Database(models.PokemonDB).Collection(models.PokemonCol)

	_, err := pokeCol.InsertMany(ctx, list)
	if err != nil {
		return err
	}

	return nil
}

func (p *PokedexRepositoryImpl) UpdatePokemon(poke models.Pokedex) error {
	var (
		pBson bson.M
	)
	ctx := context.Background()

	filter := bson.M{"pid": poke.PID}

	pByte, err := bson.Marshal(poke)
	if err != nil {
		return err
	}

	err = bson.Unmarshal(pByte, &pBson)
	if err != nil {
		return err
	}

	pokeCol := p.db.Database(models.PokemonDB).Collection(models.PokemonCol)

	_, err = pokeCol.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: pBson}})
	if err != nil {
		return err
	}

	return nil
}

func (p *PokedexRepositoryImpl) DeletePokemon(pid int) error {
	ctx := context.Background()

	filter := bson.M{"pid": pid}

	pokeCol := p.db.Database(models.PokemonDB).Collection(models.PokemonCol)

	_, err := pokeCol.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}
