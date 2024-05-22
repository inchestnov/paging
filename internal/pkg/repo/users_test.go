package repo

import (
	"context"
	"github.com/stretchr/testify/assert"
	"paging/internal/pkg/models"
	"paging/internal/pkg/tests"
	"reflect"
	"testing"
)

func TestMain(m *testing.M) {
	tests.TestWithPostgres(m)
}

func TestUsersRepo_Create(t *testing.T) {
	db, deferFunc := tests.ConnectToPostgres(t)
	defer deferFunc()

	repo := NewUsersRepo(db)

	type args struct {
		ctx  context.Context
		user models.User
	}
	tests := []struct {
		name    string
		args    args
		want    models.User
		wantErr bool
	}{
		{
			name: "all ok",
			args: args{
				user: models.User{Name: "NewUser"},
			},
			want: models.User{Name: "NewUser"},
		},
	}
	for _, ttt := range tests {
		tt := ttt
		t.Run(tt.name, func(t *testing.T) {
			tt.args.ctx = context.Background()

			got, err := repo.Create(tt.args.ctx, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			tt.want.ID = got.ID
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Create() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsersRepo_CreateRandom(t *testing.T) {
	db, deferFunc := tests.ConnectToPostgres(t)
	defer deferFunc()

	repo := NewUsersRepo(db)

	type args struct {
		ctx        context.Context
		usersCount int64
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "size = 100",
			args: args{
				usersCount: 100,
			},
		},
		{
			name: "size = 10_000",
			args: args{
				usersCount: 10_000,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			tt.args.ctx = ctx

			assert.NoError(t, repo.DeleteAll(ctx))
			assert.NoError(t, repo.CreateRandom(ctx, tt.args.usersCount))

			var actualCount int64
			assert.NoError(t, db.QueryRow("SELECT COUNT(1) FROM users").Scan(&actualCount))
			assert.Equal(t, tt.args.usersCount, actualCount)
		})
	}
}
