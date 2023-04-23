package pgdb

import (
	"account-management-service/internal/entity"
	"account-management-service/pkg/postgres"
	"context"
	"errors"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestUserRepo_CreateUser(t *testing.T) {
	type args struct {
		ctx  context.Context
		user entity.User
	}

	type MockBehavior func(m pgxmock.PgxPoolIface, args args)

	testCases := []struct {
		name         string
		args         args
		mockBehavior MockBehavior
		want         int
		wantErr      bool
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				user: entity.User{
					Username: "test_user",
					Password: "Qwerty1!",
				},
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				rows := pgxmock.NewRows([]string{"id"}).
					AddRow(1)

				m.ExpectQuery("INSERT INTO users").
					WithArgs(args.user.Username, args.user.Password).
					WillReturnRows(rows)
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "user already exists",
			args: args{
				ctx: context.Background(),
				user: entity.User{
					Username: "test_user",
					Password: "Qwerty1!",
				},
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				m.ExpectQuery("INSERT INTO users").
					WithArgs(args.user.Username, args.user.Password).
					WillReturnError(&pgconn.PgError{
						Code: "23505",
					})
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "unexpected error",
			args: args{
				ctx: context.Background(),
				user: entity.User{
					Username: "test_user",
					Password: "Qwerty1!",
				},
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				m.ExpectQuery("INSERT INTO users").
					WithArgs(args.user.Username, args.user.Password).
					WillReturnError(errors.New("some error"))
			},
			want:    0,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			poolMock, _ := pgxmock.NewPool()
			defer poolMock.Close()
			tc.mockBehavior(poolMock, tc.args)

			postgresMock := &postgres.Postgres{
				Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
				Pool:    poolMock,
			}
			userRepoMock := NewUserRepo(postgresMock)

			got, err := userRepoMock.CreateUser(tc.args.ctx, tc.args.user)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.want, got)

			err = poolMock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func TestUserRepo_GetUserByUsernameAndPassword(t *testing.T) {
	type args struct {
		ctx      context.Context
		username string
		password string
	}

	type MockBehavior func(m pgxmock.PgxPoolIface, args args)

	testCases := []struct {
		name         string
		args         args
		mockBehavior MockBehavior
		want         entity.User
		wantErr      bool
	}{
		{
			name: "OK",
			args: args{
				ctx:      context.Background(),
				username: "test_user",
				password: "Qwerty1!",
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				rows := pgxmock.NewRows([]string{"id", "username", "password", "created_at"}).
					AddRow(1, args.username, args.password, time.UnixMilli(123456))

				m.ExpectQuery("SELECT id, username, password, created_at FROM users").
					WithArgs(args.username, args.password).
					WillReturnRows(rows)
			},
			want: entity.User{
				Id:        1,
				Username:  "test_user",
				Password:  "Qwerty1!",
				CreatedAt: time.UnixMilli(123456),
			},
			wantErr: false,
		},
		{
			name: "user not found",
			args: args{
				ctx:      context.Background(),
				username: "test_user",
				password: "Qwerty1!",
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				m.ExpectQuery("SELECT id, username, password, created_at FROM users").
					WithArgs(args.username, args.password).
					WillReturnError(pgx.ErrNoRows)
			},
			want:    entity.User{},
			wantErr: true,
		},
		{
			name: "unexpected error",
			args: args{
				ctx:      context.Background(),
				username: "test_user",
				password: "Qwerty1!",
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				m.ExpectQuery("SELECT id, username, password, created_at FROM users").
					WithArgs(args.username, args.password).
					WillReturnError(errors.New("some error"))
			},
			want:    entity.User{},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			poolMock, _ := pgxmock.NewPool()
			defer poolMock.Close()
			tc.mockBehavior(poolMock, tc.args)

			postgresMock := &postgres.Postgres{
				Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
				Pool:    poolMock,
			}
			userRepoMock := NewUserRepo(postgresMock)

			got, err := userRepoMock.GetUserByUsernameAndPassword(tc.args.ctx, tc.args.username, tc.args.password)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.want, got)

			err = poolMock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func TestUserRepo_GetUserById(t *testing.T) {
	type args struct {
		ctx context.Context
		id  int
	}

	type MockBehavior func(m pgxmock.PgxPoolIface, args args)

	testCases := []struct {
		name         string
		args         args
		mockBehavior MockBehavior
		want         entity.User
		wantErr      bool
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				rows := pgxmock.NewRows([]string{"id", "username", "password", "created_at"}).
					AddRow(args.id, "test_user", "Qwerty1!", time.UnixMilli(123456))

				m.ExpectQuery("SELECT id, username, password, created_at FROM users").
					WithArgs(args.id).
					WillReturnRows(rows)
			},
			want: entity.User{
				Id:        1,
				Username:  "test_user",
				Password:  "Qwerty1!",
				CreatedAt: time.UnixMilli(123456),
			},
			wantErr: false,
		},
		{
			name: "user not found",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				m.ExpectQuery("SELECT id, username, password, created_at FROM users").
					WithArgs(args.id).
					WillReturnError(pgx.ErrNoRows)
			},
			want:    entity.User{},
			wantErr: true,
		},
		{
			name: "unexpected error",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				m.ExpectQuery("SELECT id, username, password, created_at FROM users").
					WithArgs(args.id).
					WillReturnError(errors.New("some error"))
			},
			want:    entity.User{},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			poolMock, _ := pgxmock.NewPool()
			defer poolMock.Close()
			tc.mockBehavior(poolMock, tc.args)

			postgresMock := &postgres.Postgres{
				Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
				Pool:    poolMock,
			}
			userRepoMock := NewUserRepo(postgresMock)

			got, err := userRepoMock.GetUserById(tc.args.ctx, tc.args.id)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.want, got)

			err = poolMock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func TestUserRepo_GetUserByUsername(t *testing.T) {
	type args struct {
		ctx      context.Context
		username string
	}

	type MockBehavior func(m pgxmock.PgxPoolIface, args args)

	testCases := []struct {
		name         string
		args         args
		mockBehavior MockBehavior
		want         entity.User
		wantErr      bool
	}{
		{
			name: "OK",
			args: args{
				ctx:      context.Background(),
				username: "test_user",
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				rows := pgxmock.NewRows([]string{"id", "username", "password", "created_at"}).
					AddRow(1, args.username, "Qwerty1!", time.UnixMilli(123456))

				m.ExpectQuery("SELECT id, username, password, created_at FROM users").
					WithArgs(args.username).
					WillReturnRows(rows)
			},
			want: entity.User{
				Id:        1,
				Username:  "test_user",
				Password:  "Qwerty1!",
				CreatedAt: time.UnixMilli(123456),
			},
			wantErr: false,
		},
		{
			name: "user not found",
			args: args{
				ctx:      context.Background(),
				username: "test_user",
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				m.ExpectQuery("SELECT id, username, password, created_at FROM users").
					WithArgs(args.username).
					WillReturnError(pgx.ErrNoRows)
			},
			want:    entity.User{},
			wantErr: true,
		},
		{
			name: "unexpected error",
			args: args{
				ctx:      context.Background(),
				username: "test_user",
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				m.ExpectQuery("SELECT id, username, password, created_at FROM users").
					WithArgs(args.username).
					WillReturnError(errors.New("some error"))
			},
			want:    entity.User{},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			poolMock, _ := pgxmock.NewPool()
			defer poolMock.Close()
			tc.mockBehavior(poolMock, tc.args)

			postgresMock := &postgres.Postgres{
				Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
				Pool:    poolMock,
			}
			userRepoMock := NewUserRepo(postgresMock)

			got, err := userRepoMock.GetUserByUsername(tc.args.ctx, tc.args.username)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.want, got)

			err = poolMock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}
