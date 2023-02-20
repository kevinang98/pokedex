package controllers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"pokedex/src/middleware"
	"pokedex/src/models"
	Service "pokedex/src/services"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestUserControllerImpl_RegisterUser(t *testing.T) {
	type fields struct {
		us Service.UserService
	}
	type args struct {
		username string
		password string
		role     string
	}
	tests := []struct {
		name                  string
		fields                fields
		args                  args
		want                  string
		wantStatus            int
		wantErr               bool
		expectErrRegisterUser error
	}{
		{
			name: "success_register_user",
			args: args{
				username: "admin",
				password: "admin",
				role:     "admin",
			},
			want:                  `{"status":"success","data":"username admin success registered"}`,
			wantStatus:            http.StatusOK,
			wantErr:               false,
			expectErrRegisterUser: nil,
		},
		{
			name: "error_username_password_empty",
			args: args{
				username: "",
				password: "",
				role:     "admin",
			},
			want:                  `{"status":"fail","data":"` + models.ErrorUserPassEmpty + `"}`,
			wantStatus:            http.StatusBadRequest,
			wantErr:               true,
			expectErrRegisterUser: nil,
		},
		{
			name: "error_username_password_too_long",
			args: args{
				username: "adminadminadminadminadminadminadminadminadminadminadmin",
				password: "adminadminadminadminadminadminadminadminadminadminadmin",
				role:     "admin",
			},
			want:                  `{"status":"fail","data":"` + models.ErrorUserPassTooLong + `"}`,
			wantStatus:            http.StatusBadRequest,
			wantErr:               true,
			expectErrRegisterUser: nil,
		},
		{
			name: "error_invalid_role",
			args: args{
				username: "admin",
				password: "admin",
				role:     "another",
			},
			want:                  `{"status":"fail","data":"` + models.ErrorInvalidRole + `"}`,
			wantStatus:            http.StatusBadRequest,
			wantErr:               true,
			expectErrRegisterUser: nil,
		},
		{
			name: "error_register_user",
			args: args{
				username: "admin",
				password: "admin",
				role:     "admin",
			},
			want:                  `{"status":"error","data":"RegisterUser"}`,
			wantStatus:            http.StatusInternalServerError,
			wantErr:               true,
			expectErrRegisterUser: errors.New("RegisterUser"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			svcMock := initServiceMock()
			c := initControllerMock(svcMock)

			data := url.Values{}
			data.Set("username", tt.args.username)
			data.Set("password", tt.args.password)
			data.Set("role", tt.args.role)

			req, err := http.NewRequest(http.MethodPost, "/user/register", strings.NewReader(data.Encode()))
			if err != nil {
				t.Errorf("UserControllerImpl.RegisterUser() error = %v", err)
			}
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			svcMock.userSvc.On("RegisterUser", tt.args.username, tt.args.password, tt.args.role).
				Return(tt.expectErrRegisterUser)

			if assert.NoError(t, c.U.RegisterUser(ctx)) {
				assert.Equal(t, tt.wantStatus, rec.Code)
				assert.Equal(t, tt.want, strings.ReplaceAll(rec.Body.String(), "\n", ""))
			}
		})
	}
}

func TestUserControllerImpl_LoginUser(t *testing.T) {
	type fields struct {
		us Service.UserService
	}
	type args struct {
		username string
		password string
	}
	tests := []struct {
		name                      string
		fields                    fields
		args                      args
		want                      string
		wantStatus                int
		wantErr                   bool
		expectReturnLoginUser     *models.User
		expectErrLoginUser        error
		expectReturnGenerateToken string
		expectErrGenerateToken    error
	}{
		{
			name: "success_login_user",
			args: args{
				username: "admin",
				password: "admin",
			},
			want:       `{"status":"success","data":"token"}`,
			wantStatus: http.StatusOK,
			wantErr:    false,
			expectReturnLoginUser: &models.User{
				ID:              1,
				Username:        "admin",
				Password:        "admin",
				Role:            "admin",
				CapturedPokemon: "1,2",
				CreatedAt: sql.NullString{
					String: "2023-02-20 10:33:41",
					Valid:  true,
				},
				UpdatedAt: sql.NullString{
					String: "2023-02-20 12:31:23",
					Valid:  true,
				},
			},
			expectErrLoginUser:        nil,
			expectReturnGenerateToken: "token",
			expectErrGenerateToken:    nil,
		},
		{
			name: "error_username_password_empty",
			args: args{
				username: "",
				password: "",
			},
			want:                      `{"status":"fail","data":"` + models.ErrorUserPassEmpty + `"}`,
			wantStatus:                http.StatusBadRequest,
			wantErr:                   true,
			expectReturnLoginUser:     &models.User{},
			expectErrLoginUser:        nil,
			expectReturnGenerateToken: "",
			expectErrGenerateToken:    nil,
		},
		{
			name: "error_username_password_too_long",
			args: args{
				username: "adminadminadminadminadminadminadminadminadminadminadminadmin",
				password: "adminadminadminadminadminadminadminadminadminadminadminadmin",
			},
			want:                      `{"status":"fail","data":"` + models.ErrorUserPassTooLong + `"}`,
			wantStatus:                http.StatusBadRequest,
			wantErr:                   true,
			expectReturnLoginUser:     &models.User{},
			expectErrLoginUser:        nil,
			expectReturnGenerateToken: "",
			expectErrGenerateToken:    nil,
		},
		{
			name: "error_login_user",
			args: args{
				username: "admin",
				password: "admin",
			},
			want:                      `{"status":"error","data":"LoginUser"}`,
			wantStatus:                http.StatusInternalServerError,
			wantErr:                   true,
			expectReturnLoginUser:     &models.User{},
			expectErrLoginUser:        errors.New("LoginUser"),
			expectReturnGenerateToken: "",
			expectErrGenerateToken:    nil,
		},
		{
			name: "error_generate_token",
			args: args{
				username: "admin",
				password: "admin",
			},
			want:       `{"status":"error","data":"GenerateToken"}`,
			wantStatus: http.StatusInternalServerError,
			wantErr:    true,
			expectReturnLoginUser: &models.User{
				ID:              1,
				Username:        "admin",
				Password:        "admin",
				Role:            "admin",
				CapturedPokemon: "1,2",
				CreatedAt: sql.NullString{
					String: "2023-02-20 10:33:41",
					Valid:  true,
				},
				UpdatedAt: sql.NullString{
					String: "2023-02-20 12:31:23",
					Valid:  true,
				},
			},
			expectErrLoginUser:        nil,
			expectReturnGenerateToken: "",
			expectErrGenerateToken:    errors.New("GenerateToken"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			svcMock := initServiceMock()
			c := initControllerMock(svcMock)

			data := url.Values{}
			data.Set("username", tt.args.username)
			data.Set("password", tt.args.password)

			req, err := http.NewRequest(http.MethodPost, "/user/login", strings.NewReader(data.Encode()))
			if err != nil {
				t.Errorf("UserControllerImpl.LoginUser() error = %v", err)
			}
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			svcMock.userSvc.On("LoginUser", tt.args.username, tt.args.password).
				Return(tt.expectReturnLoginUser, tt.expectErrLoginUser)

			svcMock.mid.On("GenerateToken", tt.expectReturnLoginUser.ID, tt.expectReturnLoginUser.Username, tt.expectReturnLoginUser.Role).
				Return(tt.expectReturnGenerateToken, tt.expectErrGenerateToken)

			if assert.NoError(t, c.U.LoginUser(ctx)) {
				assert.Equal(t, tt.wantStatus, rec.Code)
				assert.Equal(t, tt.want, strings.ReplaceAll(rec.Body.String(), "\n", ""))
			}
		})
	}
}

func TestUserControllerImpl_UpdateCapturedPokemonUser(t *testing.T) {
	mockListPID := "1"
	type fields struct {
		us Service.UserService
		m  middleware.Middleware
	}
	type args struct {
		cap models.PokemonCaptured
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
		expectErrUpdateCaptured      error
	}{
		{
			name: "success_update_captured_pokemon",
			args: args{
				cap: models.PokemonCaptured{
					PID: []int{
						1,
					},
				},
			},
			want:                         `{"status":"success","data":"pokemon captured has updated"}`,
			wantStatus:                   http.StatusOK,
			wantErr:                      false,
			expectReturnParseJWTID:       1,
			expectReturnParseJWTUsername: "admin",
			expectReturnParseJWTRole:     "admin",
			expectErrValidatePID:         nil,
			expectErrUpdateCaptured:      nil,
		},
		{
			name: "error_validate_pokedex_body",
			args: args{
				cap: models.PokemonCaptured{
					PID: []int{
						1,
					},
				},
			},
			want:                         `{"status":"fail","data":"ValidatePokedexBody"}`,
			wantStatus:                   http.StatusBadRequest,
			wantErr:                      true,
			expectReturnParseJWTID:       1,
			expectReturnParseJWTUsername: "admin",
			expectReturnParseJWTRole:     "admin",
			expectErrValidatePID:         errors.New("ValidatePokedexBody"),
			expectErrUpdateCaptured:      nil,
		},
		{
			name: "error_update_capture_pokemon",
			args: args{
				cap: models.PokemonCaptured{
					PID: []int{
						1,
					},
				},
			},
			want:                         `{"status":"error","data":"UpdateCapture"}`,
			wantStatus:                   http.StatusInternalServerError,
			wantErr:                      true,
			expectReturnParseJWTID:       1,
			expectReturnParseJWTUsername: "admin",
			expectReturnParseJWTRole:     "admin",
			expectErrValidatePID:         nil,
			expectErrUpdateCaptured:      errors.New("UpdateCapture"),
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
				t.Errorf("UserControllerImpl.UpdateCapturedPokemonUser() error = %v", err)
			}

			reqBody.Data = tt.args.cap
			body, err := json.Marshal(reqBody)
			if err != nil {
				t.Errorf("UserControllerImpl.UpdateCapturedPokemonUser() error = %v", err)
			}

			req, err := http.NewRequest(http.MethodPut, "/user/captured-pokemon", bytes.NewBuffer(body))
			if err != nil {
				t.Errorf("UserControllerImpl.UpdateCapturedPokemonUser() error = %v", err)
			}
			req.Header.Add("Authorization", "Bearer "+tokenStr)
			req.Header.Add("Content-Type", "application/json")

			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			svcMock.mid.On("ParseJwtToken", ctx).
				Return(tt.expectReturnParseJWTID, tt.expectReturnParseJWTUsername, tt.expectReturnParseJWTRole)

			svcMock.mid.On("ValidateIntPokedexBody", tt.args.cap.PID).
				Return(tt.expectErrValidatePID)

			svcMock.userSvc.On("UpdateCapturedPokemonUser", tt.expectReturnParseJWTID, mockListPID).
				Return(tt.expectErrUpdateCaptured)

			if assert.NoError(t, c.U.UpdateCapturedPokemonUser(ctx)) {
				assert.Equal(t, tt.wantStatus, rec.Code)
				assert.Equal(t, tt.want, strings.ReplaceAll(rec.Body.String(), "\n", ""))
			}
		})
	}
}
