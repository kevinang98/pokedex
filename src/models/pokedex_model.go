package models

type Pokedex struct {
	PID         int                `json:"pid" bson:"pid"`
	Name        string             `json:"name" bson:"name"`
	Race        string             `json:"race" bson:"race"`
	Type        []string           `json:"type" bson:"type"`
	Description PokemonDescription `json:"description" bson:"description"`
	Stats       PokemonStats       `json:"stats" bson:"stats"`
	Image       string             `json:"image" bson:"image"`
	Captured    bool               `json:"captured"`
}

type PokemonDescription struct {
	Detail string `json:"detail" bson:"detail"`
	Weight string `json:"weight" bson:"weight"`
	Height string `json:"height" bson:"height"`
}

type PokemonStats struct {
	HP     int `json:"hp" bson:"hp"`
	Attack int `json:"attack" bson:"attack"`
	Def    int `json:"def" bson:"def"`
	Speed  int `json:"speed" bson:"spped"`
}

type PokedexRequest struct {
	Name     string `query:"name"`
	Sort     string `query:"sort"`
	PokeType string `query:"type"`
	Offset   int64  `query:"offset"`
	Limit    int64  `query:"limit"`
	Order    string `query:"order"`
	Option   string `query:"option"`
}

type DeletePokedex struct {
	PID []int `json:"pid" bson:"pid"`
}
