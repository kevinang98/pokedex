package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"pokedex/src/middleware"
	"pokedex/src/models"
	Service "pokedex/src/services"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestPokedexControllerImpl_GetPokemon(t *testing.T) {
	mockTypes := []string{"dragon"}
	mockOrder := int64(1)
	type fields struct {
		ps Service.PokedexService
		m  middleware.Middleware
	}
	type args struct {
		params models.PokedexRequest
	}
	tests := []struct {
		name                         string
		fields                       fields
		args                         args
		want                         string
		wantStatus                   int
		wantErr                      bool
		expectReturnGetPokemon       []models.Pokedex
		expectErrGetPokemon          error
		expectReturnParseJWTID       int
		expectReturnParseJWTUsername string
		expectReturnParseJWTRole     string
	}{
		{
			name: "success_get_pokemon",
			args: args{
				params: models.PokedexRequest{
					Name:     "ra",
					Sort:     "pid",
					PokeType: "dragon",
					Offset:   0,
					Limit:    5,
					Order:    "asc",
					Option:   "catched",
				},
			},
			want:       `{"status":"success","data":[{"pid":1,"name":"Charizard","race":"Fire Dragon Monster","type":["Dragon","Fire"],"description":{"detail":"","weight":"110kg","height":"10m"},"stats":{"hp":100,"attack":30,"def":25,"speed":15},"image":"charizard.jpg","captured":true}]}`,
			wantStatus: http.StatusOK,
			wantErr:    false,
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
			},
			expectErrGetPokemon:          nil,
			expectReturnParseJWTID:       1,
			expectReturnParseJWTUsername: "admin",
			expectReturnParseJWTRole:     "admin",
		},
		{
			name: "error_access_invalid",
			args: args{
				params: models.PokedexRequest{
					Name:     "ra",
					Sort:     "pid",
					PokeType: "dragon",
					Offset:   0,
					Limit:    5,
					Order:    "asc",
					Option:   "catched",
				},
			},
			want:       `{"status":"fail","data":"` + models.ErrorNoAccess + `"}`,
			wantStatus: http.StatusUnauthorized,
			wantErr:    true,
			expectReturnGetPokemon: []models.Pokedex{
				{},
			},
			expectErrGetPokemon:          nil,
			expectReturnParseJWTID:       1,
			expectReturnParseJWTUsername: "admin",
			expectReturnParseJWTRole:     "user",
		},
		{
			name: "error_offset_limit_invalid",
			args: args{
				params: models.PokedexRequest{
					Name:     "ra",
					Sort:     "pid",
					PokeType: "dragon",
					Offset:   -1,
					Limit:    -1,
					Order:    "asc",
					Option:   "catched",
				},
			},
			want:       `{"status":"fail","data":"` + models.ErrorOffsetLimitParam + `"}`,
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
			expectReturnGetPokemon: []models.Pokedex{
				{},
			},
			expectErrGetPokemon:          nil,
			expectReturnParseJWTID:       1,
			expectReturnParseJWTUsername: "admin",
			expectReturnParseJWTRole:     "admin",
		},
		{
			name: "error_sort_invalid",
			args: args{
				params: models.PokedexRequest{
					Name:     "ra",
					Sort:     "test",
					PokeType: "dragon",
					Offset:   0,
					Limit:    5,
					Order:    "asc",
					Option:   "catched",
				},
			},
			want:       `{"status":"fail","data":"` + models.ErrorSortParam + `"}`,
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
			expectReturnGetPokemon: []models.Pokedex{
				{},
			},
			expectErrGetPokemon:          nil,
			expectReturnParseJWTID:       1,
			expectReturnParseJWTUsername: "admin",
			expectReturnParseJWTRole:     "admin",
		},
		{
			name: "error_order_invalid",
			args: args{
				params: models.PokedexRequest{
					Name:     "ra",
					Sort:     "pid",
					PokeType: "dragon",
					Offset:   0,
					Limit:    5,
					Order:    "test",
					Option:   "catched",
				},
			},
			want:       `{"status":"fail","data":"` + models.ErrorOrderParam + `"}`,
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
			expectReturnGetPokemon: []models.Pokedex{
				{},
			},
			expectErrGetPokemon:          nil,
			expectReturnParseJWTID:       1,
			expectReturnParseJWTUsername: "admin",
			expectReturnParseJWTRole:     "admin",
		},
		{
			name: "error_get_pokemon",
			args: args{
				params: models.PokedexRequest{
					Name:     "ra",
					Sort:     "pid",
					PokeType: "dragon",
					Offset:   0,
					Limit:    5,
					Order:    "asc",
					Option:   "catched",
				},
			},
			want:                         `{"status":"error","data":"GetPokemon"}`,
			wantStatus:                   http.StatusInternalServerError,
			wantErr:                      true,
			expectReturnGetPokemon:       nil,
			expectErrGetPokemon:          errors.New("GetPokemon"),
			expectReturnParseJWTID:       1,
			expectReturnParseJWTUsername: "admin",
			expectReturnParseJWTRole:     "admin",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			svcMock := initServiceMock()
			c := initControllerMock(svcMock)

			claims := &models.JwtCustomClaims{
				ID:       tt.expectReturnParseJWTID,
				Username: tt.expectReturnParseJWTUsername,
				Role:     tt.expectReturnParseJWTRole,
				StandardClaims: jwt.StandardClaims{
					ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
				},
			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			tokenStr, err := token.SignedString([]byte(models.JwtKey))
			if err != nil {
				t.Errorf("PokedexControllerImpl.GetPokemon() error = %v", err)
			}

			req, err := http.NewRequest(http.MethodGet, "/pokedex", nil)
			if err != nil {
				t.Errorf("PokedexControllerImpl.GetPokemon() error = %v", err)
			}
			req.Header.Add("Authorization", "Bearer "+tokenStr)

			q := req.URL.Query()
			q.Add("name", tt.args.params.Name)
			q.Add("sort", tt.args.params.Sort)
			q.Add("type", tt.args.params.PokeType)
			q.Add("offset", strconv.Itoa(int(tt.args.params.Offset)))
			q.Add("limit", strconv.Itoa(int(tt.args.params.Limit)))
			q.Add("order", tt.args.params.Order)
			q.Add("option", tt.args.params.Option)
			req.URL.RawQuery = q.Encode()

			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			svcMock.mid.On("ParseJwtToken", ctx).
				Return(tt.expectReturnParseJWTID, tt.expectReturnParseJWTUsername, tt.expectReturnParseJWTRole)

			svcMock.pokedexSvc.On("GetPokemon", tt.expectReturnParseJWTID, tt.args.params.Name, tt.args.params.Sort,
				tt.args.params.Option, mockTypes, tt.args.params.Offset, tt.args.params.Limit, mockOrder).
				Return(tt.expectReturnGetPokemon, tt.expectErrGetPokemon)

			if assert.NoError(t, c.P.GetPokemon(ctx)) {
				assert.Equal(t, tt.wantStatus, rec.Code)
				assert.Equal(t, tt.want, strings.ReplaceAll(rec.Body.String(), "\n", ""))
			}
		})
	}
}

func TestPokedexControllerImpl_AddPokemon(t *testing.T) {
	type fields struct {
		ps Service.PokedexService
		m  middleware.Middleware
	}
	type args struct {
		pokedex []models.Pokedex
	}
	tests := []struct {
		name                         string
		fields                       fields
		args                         args
		want                         string
		wantStatus                   int
		wantErr                      bool
		expectReturnParseJWTID       int
		expectReturnParseJWTUsername string
		expectReturnParseJWTRole     string
		expectErrValidatePID         error
		expectErrAddPokemon          error
	}{
		{
			name: "success_add_pokemon",
			args: args{
				pokedex: []models.Pokedex{
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
			want:                         `{"status":"success","data":"success add pokemons"}`,
			wantStatus:                   http.StatusOK,
			wantErr:                      false,
			expectReturnParseJWTID:       1,
			expectReturnParseJWTUsername: "admin",
			expectReturnParseJWTRole:     "admin",
			expectErrValidatePID:         nil,
			expectErrAddPokemon:          nil,
		},
		{
			name: "error_access_invalid",
			args: args{
				pokedex: []models.Pokedex{
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
			want:                         `{"status":"fail","data":"` + models.ErrorNoAccess + `"}`,
			wantStatus:                   http.StatusUnauthorized,
			wantErr:                      true,
			expectReturnParseJWTID:       1,
			expectReturnParseJWTUsername: "admin",
			expectReturnParseJWTRole:     "user",
			expectErrValidatePID:         nil,
			expectErrAddPokemon:          nil,
		},
		{
			name: "error_validate_body",
			args: args{
				pokedex: []models.Pokedex{
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
			want:                         `{"status":"fail","data":"ValidateBody"}`,
			wantStatus:                   http.StatusBadRequest,
			wantErr:                      true,
			expectReturnParseJWTID:       1,
			expectReturnParseJWTUsername: "admin",
			expectReturnParseJWTRole:     "admin",
			expectErrValidatePID:         errors.New("ValidateBody"),
			expectErrAddPokemon:          nil,
		},
		{
			name: "error_add_pokemon",
			args: args{
				pokedex: []models.Pokedex{
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
			want:                         `{"status":"error","data":"AddPokemon"}`,
			wantStatus:                   http.StatusInternalServerError,
			wantErr:                      true,
			expectReturnParseJWTID:       1,
			expectReturnParseJWTUsername: "admin",
			expectReturnParseJWTRole:     "admin",
			expectErrValidatePID:         nil,
			expectErrAddPokemon:          errors.New("AddPokemon"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var reqBody models.Request

			e := echo.New()
			svcMock := initServiceMock()
			c := initControllerMock(svcMock)

			claims := &models.JwtCustomClaims{
				ID:       tt.expectReturnParseJWTID,
				Username: tt.expectReturnParseJWTUsername,
				Role:     tt.expectReturnParseJWTRole,
				StandardClaims: jwt.StandardClaims{
					ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
				},
			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			tokenStr, err := token.SignedString([]byte(models.JwtKey))
			if err != nil {
				t.Errorf("PokedexControllerImpl.AddPokemon() error = %v", err)
			}

			reqBody.Data = tt.args.pokedex
			body, err := json.Marshal(reqBody)
			if err != nil {
				t.Errorf("PokedexControllerImpl.AddPokemon() error = %v", err)
			}

			req, err := http.NewRequest(http.MethodPost, "/pokedex", bytes.NewBuffer(body))
			if err != nil {
				t.Errorf("PokedexControllerImpl.AddPokemon() error = %v", err)
			}
			req.Header.Add("Authorization", "Bearer "+tokenStr)
			req.Header.Add("Content-Type", "application/json")

			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			svcMock.mid.On("ParseJwtToken", ctx).
				Return(tt.expectReturnParseJWTID, tt.expectReturnParseJWTUsername, tt.expectReturnParseJWTRole)

			svcMock.mid.On("ValidatePokedexBody", tt.args.pokedex).
				Return(tt.expectErrValidatePID)

			svcMock.pokedexSvc.On("AddPokemon", tt.args.pokedex).
				Return(tt.expectErrAddPokemon)

			if assert.NoError(t, c.P.AddPokemon(ctx)) {
				assert.Equal(t, tt.wantStatus, rec.Code)
				assert.Equal(t, tt.want, strings.ReplaceAll(rec.Body.String(), "\n", ""))
			}
		})
	}
}

func TestPokedexControllerImpl_UpdatePokemon(t *testing.T) {
	type fields struct {
		ps Service.PokedexService
		m  middleware.Middleware
	}
	type args struct {
		pokedex []models.Pokedex
	}
	tests := []struct {
		name                         string
		fields                       fields
		args                         args
		want                         string
		wantStatus                   int
		wantErr                      bool
		expectReturnParseJWTID       int
		expectReturnParseJWTUsername string
		expectReturnParseJWTRole     string
		expectErrValidatePID         error
		expectReturnUpdatePokemon    string
		expectErrUpdatePokemon       error
	}{
		{
			name: "success_update_pokemon",
			args: args{
				pokedex: []models.Pokedex{
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
			want:                         `{"status":"success","data":"success update pokemons"}`,
			wantStatus:                   http.StatusOK,
			wantErr:                      false,
			expectReturnParseJWTID:       1,
			expectReturnParseJWTUsername: "admin",
			expectReturnParseJWTRole:     "admin",
			expectErrValidatePID:         nil,
			expectReturnUpdatePokemon:    "",
			expectErrUpdatePokemon:       nil,
		},
		{
			name: "error_access_invalid",
			args: args{
				pokedex: []models.Pokedex{
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
			want:                         `{"status":"fail","data":"` + models.ErrorNoAccess + `"}`,
			wantStatus:                   http.StatusUnauthorized,
			wantErr:                      true,
			expectReturnParseJWTID:       1,
			expectReturnParseJWTUsername: "admin",
			expectReturnParseJWTRole:     "user",
			expectErrValidatePID:         nil,
			expectReturnUpdatePokemon:    "",
			expectErrUpdatePokemon:       nil,
		},
		{
			name: "error_validate_body",
			args: args{
				pokedex: []models.Pokedex{
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
			want:                         `{"status":"fail","data":"ValidateBody"}`,
			wantStatus:                   http.StatusBadRequest,
			wantErr:                      true,
			expectReturnParseJWTID:       1,
			expectReturnParseJWTUsername: "admin",
			expectReturnParseJWTRole:     "admin",
			expectErrValidatePID:         errors.New("ValidateBody"),
			expectReturnUpdatePokemon:    "",
			expectErrUpdatePokemon:       nil,
		},
		{
			name: "error_update_pokemon",
			args: args{
				pokedex: []models.Pokedex{
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
			want:                         `{"status":"error","data":"UpdatePokemon"}`,
			wantStatus:                   http.StatusInternalServerError,
			wantErr:                      true,
			expectReturnParseJWTID:       1,
			expectReturnParseJWTUsername: "admin",
			expectReturnParseJWTRole:     "admin",
			expectErrValidatePID:         nil,
			expectReturnUpdatePokemon:    "",
			expectErrUpdatePokemon:       errors.New("UpdatePokemon"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var reqBody models.Request

			e := echo.New()
			svcMock := initServiceMock()
			c := initControllerMock(svcMock)

			claims := &models.JwtCustomClaims{
				ID:       tt.expectReturnParseJWTID,
				Username: tt.expectReturnParseJWTUsername,
				Role:     tt.expectReturnParseJWTRole,
				StandardClaims: jwt.StandardClaims{
					ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
				},
			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			tokenStr, err := token.SignedString([]byte(models.JwtKey))
			if err != nil {
				t.Errorf("PokedexControllerImpl.UpdatePokemon() error = %v", err)
			}

			reqBody.Data = tt.args.pokedex
			body, err := json.Marshal(reqBody)
			if err != nil {
				t.Errorf("PokedexControllerImpl.UpdatePokemon() error = %v", err)
			}

			req, err := http.NewRequest(http.MethodPut, "/pokedex", bytes.NewBuffer(body))
			if err != nil {
				t.Errorf("PokedexControllerImpl.UpdatePokemon() error = %v", err)
			}
			req.Header.Add("Authorization", "Bearer "+tokenStr)
			req.Header.Add("Content-Type", "application/json")

			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			svcMock.mid.On("ParseJwtToken", ctx).
				Return(tt.expectReturnParseJWTID, tt.expectReturnParseJWTUsername, tt.expectReturnParseJWTRole)

			svcMock.mid.On("ValidatePokedexBody", tt.args.pokedex).
				Return(tt.expectErrValidatePID)

			svcMock.pokedexSvc.On("UpdatePokemon", tt.args.pokedex).
				Return(tt.expectReturnUpdatePokemon, tt.expectErrUpdatePokemon)

			if assert.NoError(t, c.P.UpdatePokemon(ctx)) {
				assert.Equal(t, tt.wantStatus, rec.Code)
				assert.Equal(t, tt.want, strings.ReplaceAll(rec.Body.String(), "\n", ""))
			}
		})
	}
}

func TestPokedexControllerImpl_DeletePokemon(t *testing.T) {
	type fields struct {
		ps Service.PokedexService
		m  middleware.Middleware
	}
	type args struct {
		pokedex models.DeletePokedex
	}
	tests := []struct {
		name                         string
		fields                       fields
		args                         args
		want                         string
		wantStatus                   int
		wantErr                      bool
		expectReturnParseJWTID       int
		expectReturnParseJWTUsername string
		expectReturnParseJWTRole     string
		expectErrValidatePID         error
		expectReturnDeletePokemon    string
		expectErrDeletePokemon       error
	}{
		{
			name: "success_delete_pokemon",
			args: args{
				pokedex: models.DeletePokedex{
					PID: []int{
						1,
					},
				},
			},
			want:                         `{"status":"success","data":"success delete pokemons"}`,
			wantStatus:                   http.StatusOK,
			wantErr:                      false,
			expectReturnParseJWTID:       1,
			expectReturnParseJWTUsername: "admin",
			expectReturnParseJWTRole:     "admin",
			expectErrValidatePID:         nil,
			expectReturnDeletePokemon:    "",
			expectErrDeletePokemon:       nil,
		},
		{
			name: "error_access_invalid",
			args: args{
				pokedex: models.DeletePokedex{
					PID: []int{
						1,
					},
				},
			},
			want:                         `{"status":"fail","data":"` + models.ErrorNoAccess + `"}`,
			wantStatus:                   http.StatusUnauthorized,
			wantErr:                      true,
			expectReturnParseJWTID:       1,
			expectReturnParseJWTUsername: "admin",
			expectReturnParseJWTRole:     "user",
			expectErrValidatePID:         nil,
			expectReturnDeletePokemon:    "",
			expectErrDeletePokemon:       nil,
		},
		{
			name: "error_validate_body",
			args: args{
				pokedex: models.DeletePokedex{
					PID: []int{
						1,
					},
				},
			},
			want:                         `{"status":"fail","data":"ValidatePID"}`,
			wantStatus:                   http.StatusBadRequest,
			wantErr:                      true,
			expectReturnParseJWTID:       1,
			expectReturnParseJWTUsername: "admin",
			expectReturnParseJWTRole:     "admin",
			expectErrValidatePID:         errors.New("ValidatePID"),
			expectReturnDeletePokemon:    "",
			expectErrDeletePokemon:       nil,
		},
		{
			name: "error_delete_pokemon",
			args: args{
				pokedex: models.DeletePokedex{
					PID: []int{
						1,
					},
				},
			},
			want:                         `{"status":"error","data":"DeletePokemon"}`,
			wantStatus:                   http.StatusInternalServerError,
			wantErr:                      true,
			expectReturnParseJWTID:       1,
			expectReturnParseJWTUsername: "admin",
			expectReturnParseJWTRole:     "admin",
			expectErrValidatePID:         nil,
			expectReturnDeletePokemon:    "",
			expectErrDeletePokemon:       errors.New("DeletePokemon"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var reqBody models.Request

			e := echo.New()
			svcMock := initServiceMock()
			c := initControllerMock(svcMock)

			claims := &models.JwtCustomClaims{
				ID:       tt.expectReturnParseJWTID,
				Username: tt.expectReturnParseJWTUsername,
				Role:     tt.expectReturnParseJWTRole,
				StandardClaims: jwt.StandardClaims{
					ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
				},
			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			tokenStr, err := token.SignedString([]byte(models.JwtKey))
			if err != nil {
				t.Errorf("PokedexControllerImpl.DeletePokemon() error = %v", err)
			}

			reqBody.Data = tt.args.pokedex
			body, err := json.Marshal(reqBody)
			if err != nil {
				t.Errorf("PokedexControllerImpl.DeletePokemon() error = %v", err)
			}

			req, err := http.NewRequest(http.MethodDelete, "/pokedex", bytes.NewBuffer(body))
			if err != nil {
				t.Errorf("PokedexControllerImpl.DeletePokemon() error = %v", err)
			}
			req.Header.Add("Authorization", "Bearer "+tokenStr)
			req.Header.Add("Content-Type", "application/json")

			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			svcMock.mid.On("ParseJwtToken", ctx).
				Return(tt.expectReturnParseJWTID, tt.expectReturnParseJWTUsername, tt.expectReturnParseJWTRole)

			svcMock.mid.On("ValidateIntPokedexBody", tt.args.pokedex.PID).
				Return(tt.expectErrValidatePID)

			svcMock.pokedexSvc.On("DeletePokemon", tt.args.pokedex.PID).
				Return(tt.expectReturnDeletePokemon, tt.expectErrDeletePokemon)

			if assert.NoError(t, c.P.DeletePokemon(ctx)) {
				assert.Equal(t, tt.wantStatus, rec.Code)
				assert.Equal(t, tt.want, strings.ReplaceAll(rec.Body.String(), "\n", ""))
			}
		})
	}
}
