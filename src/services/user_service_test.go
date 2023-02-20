package services

import (
	"database/sql"
	"errors"
	"pokedex/src/models"
	Repository "pokedex/src/repositories"
	"reflect"
	"testing"
)

func TestUserServiceImpl_RegisterUser(t *testing.T) {
	type fields struct {
		ur Repository.UserRepository
	}
	type args struct {
		username string
		password string
		role     string
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
			name: "success_register_user",
			args: args{
				username: "admin",
				password: "admin",
				role:     "admin",
			},
			wantErr:    false,
			wantErrMsg: "",
			expectErr:  nil,
		},
		{
			name: "error_register_user",
			args: args{
				username: "admin",
				password: "admin",
				role:     "admin",
			},
			wantErr:    true,
			wantErrMsg: "RegisterUser",
			expectErr:  errors.New("RegisterUser"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoMock := initRepositoryMock()
			s := initServiceMock(repoMock)

			repoMock.userRepo.On("RegisterUser", tt.args.username, tt.args.password, tt.args.role).
				Return(tt.expectErr)

			err := s.U.RegisterUser(tt.args.username, tt.args.password, tt.args.role)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserServiceImpl.RegisterUser() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				if !reflect.DeepEqual(err.Error(), tt.wantErrMsg) {
					t.Errorf("UserServiceImpl.RegisterUser() error = %v, want %v", err.Error(), tt.wantErrMsg)
				}
			}
		})
	}
}

func TestUserServiceImpl_LoginUser(t *testing.T) {
	type fields struct {
		r *Repository.Repository
	}
	type args struct {
		username string
		password string
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		want         *models.User
		wantErr      bool
		wantErrMsg   string
		expectReturn models.User
		expectErr    error
	}{
		{
			name: "success_login_user",
			args: args{
				username: "admin",
				password: "admin",
			},
			want: &models.User{
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
			wantErr:    false,
			wantErrMsg: "",
			expectReturn: models.User{
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
			expectErr: nil,
		},
		{
			name: "error_username_not_exist",
			args: args{
				username: "admin123",
				password: "admin",
			},
			want:         nil,
			wantErr:      true,
			wantErrMsg:   models.ErrorUsernameNotExist,
			expectReturn: models.User{},
			expectErr:    nil,
		},
		{
			name: "error_invalid_password",
			args: args{
				username: "admin",
				password: "admin123",
			},
			want:       nil,
			wantErr:    true,
			wantErrMsg: models.ErrorInvalidPassword,
			expectReturn: models.User{
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
			expectErr: nil,
		},
		{
			name: "error_get_user_data",
			args: args{
				username: "admin",
				password: "admin",
			},
			want:         nil,
			wantErr:      true,
			wantErrMsg:   "GetUser",
			expectReturn: models.User{},
			expectErr:    errors.New("GetUser"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoMock := initRepositoryMock()
			s := initServiceMock(repoMock)

			param := make(map[string]interface{}, 0)
			param["username"] = tt.args.username

			repoMock.userRepo.On("GetUserData", param).
				Return(tt.expectReturn, tt.expectErr)

			got, err := s.U.LoginUser(tt.args.username, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserServiceImpl.LoginUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				if !reflect.DeepEqual(err.Error(), tt.wantErrMsg) {
					t.Errorf("UserServiceImpl.LoginUser() error = %v, want %v", err.Error(), tt.wantErrMsg)
				}
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserServiceImpl.LoginUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserServiceImpl_UpdateCapturedPokemonUser(t *testing.T) {
	mockId := "1,2,3"
	type fields struct {
		ur Repository.UserRepository
	}
	type args struct {
		id  int
		pid string
	}
	tests := []struct {
		name                    string
		fields                  fields
		args                    args
		wantErr                 bool
		wantErrMsg              string
		expectReturnGetUser     models.User
		expectErrGetUser        error
		expectErrUpdateCaptured error
	}{
		{
			name: "success_update_captured_pokemon_user",
			args: args{
				id:  1,
				pid: "1,2",
			},
			wantErr:    false,
			wantErrMsg: "",
			expectReturnGetUser: models.User{
				ID:              1,
				Username:        "admin",
				Password:        "admin",
				Role:            "admin",
				CapturedPokemon: "3",
				CreatedAt: sql.NullString{
					String: "2023-02-20 10:33:41",
					Valid:  true,
				},
				UpdatedAt: sql.NullString{
					String: "2023-02-20 12:31:23",
					Valid:  true,
				},
			},
			expectErrGetUser:        nil,
			expectErrUpdateCaptured: nil,
		},
		{
			name: "error_get_user_data",
			args: args{
				id:  1,
				pid: "1,2",
			},
			wantErr:                 true,
			wantErrMsg:              "GetUser",
			expectReturnGetUser:     models.User{},
			expectErrGetUser:        errors.New("GetUser"),
			expectErrUpdateCaptured: nil,
		},
		{
			name: "error_update_captured_pokemon_user",
			args: args{
				id:  1,
				pid: "1,2",
			},
			wantErr:    true,
			wantErrMsg: "UpdateCaptured",
			expectReturnGetUser: models.User{
				ID:              1,
				Username:        "admin",
				Password:        "admin",
				Role:            "admin",
				CapturedPokemon: "3",
				CreatedAt: sql.NullString{
					String: "2023-02-20 10:33:41",
					Valid:  true,
				},
				UpdatedAt: sql.NullString{
					String: "2023-02-20 12:31:23",
					Valid:  true,
				},
			},
			expectErrGetUser:        nil,
			expectErrUpdateCaptured: errors.New("UpdateCaptured"),
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

			repoMock.userRepo.On("UpdateCapturedPokemonUser", tt.args.id, mockId).
				Return(tt.expectErrUpdateCaptured)

			err := s.U.UpdateCapturedPokemonUser(tt.args.id, tt.args.pid)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserServiceImpl.UpdateCapturedPokemonUser() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				if !reflect.DeepEqual(err.Error(), tt.wantErrMsg) {
					t.Errorf("UserServiceImpl.UpdateCapturedPokemonUser() error = %v, want %v", err.Error(), tt.wantErrMsg)
				}
			}
		})
	}
}
