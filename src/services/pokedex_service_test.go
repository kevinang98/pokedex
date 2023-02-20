package services

import (
	"database/sql"
	"errors"
	"pokedex/src/models"
	Repository "pokedex/src/repositories"
	"reflect"
	"testing"
)

func TestPokedexServiceImpl_GetPokemon(t *testing.T) {
	mockListPID := []int{1}
	type fields struct {
		pr Repository.PokedexRepository
		ur Repository.UserRepository
	}
	type args struct {
		id       int
		name     string
		sort     string
		option   string
		pokeType []string
		offset   int64
		limit    int64
		order    int64
	}
	tests := []struct {
		name                   string
		fields                 fields
		args                   args
		want                   []models.Pokedex
		wantErr                bool
		wantErrMsg             string
		expectReturnGetUser    models.User
		expectErrGetUser       error
		expectReturnGetPokemon []models.Pokedex
		expectErrGetPokemon    error
	}{
		{
			name: "success_get_pokemon",
			args: args{
				id:     1,
				name:   "ar",
				sort:   "pid",
				option: "all",
				pokeType: []string{
					"Dragon",
				},
				offset: 0,
				limit:  5,
				order:  1,
			},
			want: []models.Pokedex{
				{
					PID:  1,
					Name: "Charizard",
					Race: "Fire Dragon Monster",
					Type: []string{
						"Dragon",
						"Fire",
					},
					Description: models.PokemonDescription{
						Detail: "",
						Weight: "110kg",
						Height: "10m",
					},
					Stats: models.PokemonStats{
						HP:     100,
						Attack: 30,
						Def:    25,
						Speed:  15,
					},
					Image:    "charizard.jpg",
					Captured: true,
				},
				{
					PID:  2,
					Name: "Charmander",
					Race: "Fire Monster",
					Type: []string{
						"Fire",
					},
					Description: models.PokemonDescription{
						Detail: "",
						Weight: "50kg",
						Height: "4m",
					},
					Stats: models.PokemonStats{
						HP:     30,
						Attack: 10,
						Def:    8,
						Speed:  7,
					},
					Image:    "charmander.jpg",
					Captured: false,
				},
			},
			wantErr:    false,
			wantErrMsg: "",
			expectReturnGetUser: models.User{
				ID:              1,
				Username:        "admin",
				Password:        "admin",
				Role:            "admin",
				CapturedPokemon: "1",
				CreatedAt: sql.NullString{
					String: "2023-02-20 10:33:41",
					Valid:  true,
				},
				UpdatedAt: sql.NullString{
					String: "2023-02-20 12:31:23",
					Valid:  true,
				},
			},
			expectErrGetUser: nil,
			expectReturnGetPokemon: []models.Pokedex{
				{
					PID:  1,
					Name: "Charizard",
					Race: "Fire Dragon Monster",
					Type: []string{
						"Dragon",
						"Fire",
					},
					Description: models.PokemonDescription{
						Detail: "",
						Weight: "110kg",
						Height: "10m",
					},
					Stats: models.PokemonStats{
						HP:     100,
						Attack: 30,
						Def:    25,
						Speed:  15,
					},
					Image:    "charizard.jpg",
					Captured: true,
				},
				{
					PID:  2,
					Name: "Charmander",
					Race: "Fire Monster",
					Type: []string{
						"Fire",
					},
					Description: models.PokemonDescription{
						Detail: "",
						Weight: "50kg",
						Height: "4m",
					},
					Stats: models.PokemonStats{
						HP:     30,
						Attack: 10,
						Def:    8,
						Speed:  7,
					},
					Image:    "charmander.jpg",
					Captured: false,
				},
			},
			expectErrGetPokemon: nil,
		},
		{
			name: "error_get_user_data",
			args: args{
				id:     1,
				name:   "ar",
				sort:   "pid",
				option: "catched",
				pokeType: []string{
					"Dragon",
				},
				offset: 0,
				limit:  5,
				order:  1,
			},
			want:                   nil,
			wantErr:                true,
			wantErrMsg:             "GetUser",
			expectReturnGetUser:    models.User{},
			expectErrGetUser:       errors.New("GetUser"),
			expectReturnGetPokemon: nil,
			expectErrGetPokemon:    nil,
		},
		{
			name: "error_get_pokemon",
			args: args{
				id:     1,
				name:   "ar",
				sort:   "pid",
				option: "catched",
				pokeType: []string{
					"Dragon",
				},
				offset: 0,
				limit:  5,
				order:  1,
			},
			want:       nil,
			wantErr:    true,
			wantErrMsg: "GetPokemon",
			expectReturnGetUser: models.User{
				ID:              1,
				Username:        "admin",
				Password:        "admin",
				Role:            "admin",
				CapturedPokemon: "1",
				CreatedAt: sql.NullString{
					String: "2023-02-20 10:33:41",
					Valid:  true,
				},
				UpdatedAt: sql.NullString{
					String: "2023-02-20 12:31:23",
					Valid:  true,
				},
			},
			expectErrGetUser:       nil,
			expectReturnGetPokemon: nil,
			expectErrGetPokemon:    errors.New("GetPokemon"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoMock := initRepositoryMock()
			s := initServiceMock(repoMock)

			param := make(map[string]interface{}, 0)
			param["id"] = tt.args.id

			repoMock.userRepo.On("GetUserData", param).
				Return(tt.expectReturnGetUser, tt.expectErrGetUser)

			repoMock.pokedexRepo.On("GetPokemon", tt.args.name, tt.args.sort, tt.args.option, mockListPID, tt.args.pokeType,
				tt.args.offset, tt.args.limit, tt.args.order).
				Return(tt.expectReturnGetPokemon, tt.expectErrGetPokemon)

			got, err := s.P.GetPokemon(tt.args.id, tt.args.name, tt.args.sort, tt.args.option, tt.args.pokeType, tt.args.offset, tt.args.limit, tt.args.order)
			if (err != nil) != tt.wantErr {
				t.Errorf("PokedexServiceImpl.GetPokemon() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				if !reflect.DeepEqual(err.Error(), tt.wantErrMsg) {
					t.Errorf("UserServiceImpl.GetPokemon() error = %v, want %v", err.Error(), tt.wantErrMsg)
				}
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PokedexServiceImpl.GetPokemon() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPokedexServiceImpl_AddPokemon(t *testing.T) {
	type fields struct {
		pr Repository.PokedexRepository
		ur Repository.UserRepository
	}
	type args struct {
		pokeList []models.Pokedex
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantErr    bool
		wantErrMsg string
		expectErr  error
	}{
		{
			name: "success_add_pokemon",
			args: args{
				[]models.Pokedex{
					{
						PID:  1,
						Name: "Charizard",
						Race: "Fire Dragon Monster",
						Type: []string{
							"Dragon",
							"Fire",
						},
						Description: models.PokemonDescription{
							Detail: "",
							Weight: "110kg",
							Height: "10m",
						},
						Stats: models.PokemonStats{
							HP:     100,
							Attack: 30,
							Def:    25,
							Speed:  15,
						},
						Image:    "charizard.jpg",
						Captured: true,
					},
				},
			},
			wantErr:    false,
			wantErrMsg: "",
			expectErr:  nil,
		},
		{
			name: "error_add_pokemon",
			args: args{
				[]models.Pokedex{
					{
						PID:  1,
						Name: "Charizard",
						Race: "Fire Dragon Monster",
						Type: []string{
							"Dragon",
							"Fire",
						},
						Description: models.PokemonDescription{
							Detail: "",
							Weight: "110kg",
							Height: "10m",
						},
						Stats: models.PokemonStats{
							HP:     100,
							Attack: 30,
							Def:    25,
							Speed:  15,
						},
						Image:    "charizard.jpg",
						Captured: true,
					},
				},
			},
			wantErr:    true,
			wantErrMsg: "AddPokemon",
			expectErr:  errors.New("AddPokemon"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoMock := initRepositoryMock()
			s := initServiceMock(repoMock)

			repoMock.pokedexRepo.On("AddPokemon", tt.args.pokeList).
				Return(tt.expectErr)

			err := s.P.AddPokemon(tt.args.pokeList)
			if (err != nil) != tt.wantErr {
				t.Errorf("PokedexServiceImpl.AddPokemon() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				if !reflect.DeepEqual(err.Error(), tt.wantErrMsg) {
					t.Errorf("UserServiceImpl.AddPokemon() error = %v, want %v", err.Error(), tt.wantErrMsg)
				}
			}
		})
	}
}

func TestPokedexServiceImpl_UpdatePokemon(t *testing.T) {
	type fields struct {
		pr Repository.PokedexRepository
		ur Repository.UserRepository
	}
	type args struct {
		pokeList []models.Pokedex
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		want       string
		wantErr    bool
		wantErrMsg string
		expectErr  error
	}{
		{
			name: "success_update_pokemon",
			args: args{
				[]models.Pokedex{
					{
						PID:  1,
						Name: "Charizard",
						Race: "Fire Dragon Monster",
						Type: []string{
							"Dragon",
							"Fire",
						},
						Description: models.PokemonDescription{
							Detail: "",
							Weight: "110kg",
							Height: "10m",
						},
						Stats: models.PokemonStats{
							HP:     110,
							Attack: 30,
							Def:    25,
							Speed:  15,
						},
						Image:    "charizard.jpg",
						Captured: true,
					},
				},
			},
			want:       "",
			wantErr:    false,
			wantErrMsg: "",
			expectErr:  nil,
		},
		{
			name: "error_update_pokemon",
			args: args{
				[]models.Pokedex{
					{
						PID:  1,
						Name: "Charizard",
						Race: "Fire Dragon Monster",
						Type: []string{
							"Dragon",
							"Fire",
						},
						Description: models.PokemonDescription{
							Detail: "",
							Weight: "110kg",
							Height: "10m",
						},
						Stats: models.PokemonStats{
							HP:     110,
							Attack: 30,
							Def:    25,
							Speed:  15,
						},
						Image:    "charizard.jpg",
						Captured: true,
					},
				},
			},
			want:       "",
			wantErr:    true,
			wantErrMsg: "UpdatePokemon",
			expectErr:  errors.New("UpdatePokemon"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoMock := initRepositoryMock()
			s := initServiceMock(repoMock)

			repoMock.pokedexRepo.On("UpdatePokemon", tt.args.pokeList[0]).
				Return(tt.expectErr)

			got, err := s.P.UpdatePokemon(tt.args.pokeList)
			if (err != nil) != tt.wantErr {
				t.Errorf("PokedexServiceImpl.UpdatePokemon() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				if !reflect.DeepEqual(err.Error(), tt.wantErrMsg) {
					t.Errorf("UserServiceImpl.UpdatePokemon() error = %v, want %v", err.Error(), tt.wantErrMsg)
				}
			}
			if got != tt.want {
				t.Errorf("PokedexServiceImpl.UpdatePokemon() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPokedexServiceImpl_DeletePokemon(t *testing.T) {
	type fields struct {
		pr Repository.PokedexRepository
		ur Repository.UserRepository
	}
	type args struct {
		pidList []int
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		want       string
		wantErr    bool
		wantErrMsg string
		expectErr  error
	}{
		{
			name: "success_delete_pokemon",
			args: args{
				[]int{1},
			},
			want:       "",
			wantErr:    false,
			wantErrMsg: "",
			expectErr:  nil,
		},
		{
			name: "error_delete_pokemon",
			args: args{
				[]int{1},
			},
			want:       "",
			wantErr:    true,
			wantErrMsg: "DeletePokemon",
			expectErr:  errors.New("DeletePokemon"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoMock := initRepositoryMock()
			s := initServiceMock(repoMock)

			repoMock.pokedexRepo.On("DeletePokemon", tt.args.pidList[0]).
				Return(tt.expectErr)

			got, err := s.P.DeletePokemon(tt.args.pidList)
			if (err != nil) != tt.wantErr {
				t.Errorf("PokedexServiceImpl.DeletePokemon() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				if !reflect.DeepEqual(err.Error(), tt.wantErrMsg) {
					t.Errorf("UserServiceImpl.DeletePokemon() error = %v, want %v", err.Error(), tt.wantErrMsg)
				}
			}
			if got != tt.want {
				t.Errorf("PokedexServiceImpl.DeletePokemon() = %v, want %v", got, tt.want)
			}
		})
	}
}
