package repositories

import (
	"database/sql/driver"
	"errors"
	"pokedex/src/models"
	"reflect"
	"regexp"
	"testing"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func TestUserRepositoryImpl_RegisterUser(t *testing.T) {
	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		username string
		password string
		role     string
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantErrMsg  string
		wantErr     bool
		expectedErr error
	}{
		{
			name: "success_register_user",
			args: args{
				username: "admin",
				password: "admin",
				role:     "admin",
			},
			wantErr:     false,
			expectedErr: nil,
			wantErrMsg:  "",
		},
		{
			name: "error_register_user",
			args: args{
				username: "admin",
				password: "admin",
				role:     "admin",
			},
			wantErr:     true,
			expectedErr: errors.New("RegisterUser"),
			wantErrMsg:  "RegisterUser",
		},
		{
			name: "error_username_exist",
			args: args{
				username: "admin",
				password: "admin",
				role:     "admin",
			},
			wantErr: true,
			expectedErr: &mysql.MySQLError{
				Number: 1062,
			},
			wantErrMsg: models.ErrorQueryUsernameUsed,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, mock := initRepositoryMock(t)

			mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `user`")).WithArgs(tt.args.username, tt.args.password, tt.args.role).
				WillReturnResult(driver.RowsAffected(1)).WillReturnError(tt.expectedErr)

			err := r.U.RegisterUser(tt.args.username, tt.args.password, tt.args.role)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserRepositoryImpl.RegisterUser() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				if !reflect.DeepEqual(err.Error(), tt.wantErrMsg) {
					t.Errorf("UserRepositoryImpl.RegisterUser() error = %v, want %v", err.Error(), tt.wantErrMsg)
				}
			}
		})
	}
}

func TestUserRepositoryImpl_GetUserData(t *testing.T) {
	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		req map[string]interface{}
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		want        models.User
		wantErrMsg  string
		wantErr     bool
		expectedErr error
		expectedID  int
	}{
		{
			name: "success_get_user_data",
			args: args{
				req: map[string]interface{}{"id": 1},
			},
			wantErr:     false,
			expectedErr: nil,
			wantErrMsg:  "",
			expectedID:  1,
		},
		{
			name: "error_get_user_data",
			args: args{
				req: map[string]interface{}{},
			},
			wantErr:     true,
			expectedErr: errors.New("RegisterUser"),
			wantErrMsg:  "RegisterUser",
			expectedID:  0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, mock := initRepositoryMock(t)

			if tt.expectedID > 0 {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT")).WithArgs(tt.expectedID).
					WillReturnRows(mock.NewRows([]string{"id"})).WillReturnError(tt.expectedErr)
			} else {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT")).WithArgs().
					WillReturnRows(mock.NewRows([]string{"id"})).WillReturnError(tt.expectedErr)
			}

			got, err := r.U.GetUserData(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserRepositoryImpl.GetUserData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				if !reflect.DeepEqual(err.Error(), tt.wantErrMsg) {
					t.Errorf("UserRepositoryImpl.GetUserData() error = %v, want %v", err.Error(), tt.wantErrMsg)
				}
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserRepositoryImpl.GetUserData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserRepositoryImpl_UpdateCapturedPokemonUser(t *testing.T) {
	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		id       int
		captured string
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantErrMsg  string
		wantErr     bool
		expectedErr error
	}{
		{
			name: "success_update_captured_pokemon_user",
			args: args{
				id:       1,
				captured: "1,2",
			},
			wantErr:     false,
			expectedErr: nil,
			wantErrMsg:  "",
		},
		{
			name: "error_update_captured_pokemon_user",
			args: args{
				id:       1,
				captured: "1,2",
			},
			wantErr:     true,
			expectedErr: errors.New("UpdateCaptured"),
			wantErrMsg:  "UpdateCaptured",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, mock := initRepositoryMock(t)

			mock.ExpectExec(regexp.QuoteMeta("UPDATE `user`")).WithArgs(tt.args.captured, tt.args.id).
				WillReturnResult(driver.RowsAffected(1)).WillReturnError(tt.expectedErr)

			err := r.U.UpdateCapturedPokemonUser(tt.args.id, tt.args.captured)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserRepositoryImpl.UpdateCapturedPokemonUser() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				if !reflect.DeepEqual(err.Error(), tt.wantErrMsg) {
					t.Errorf("UserRepositoryImpl.UpdateCapturedPokemonUser() error = %v, want %v", err.Error(), tt.wantErrMsg)
				}
			}
		})
	}
}
