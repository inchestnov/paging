package repo

import (
	"context"
	"database/sql"
	"paging/internal/pkg/models"
	"strconv"
	"strings"
)

type UsersRepo struct {
	db *sql.DB
}

var nilUser = models.User{}

func NewUsersRepo(db *sql.DB) *UsersRepo {
	return &UsersRepo{db}
}

func (u *UsersRepo) Create(ctx context.Context, user models.User) (models.User, error) {
	var id int64
	err := u.db.QueryRowContext(ctx, "INSERT INTO users(name) VALUES ($1) RETURNING id", user.Name).Scan(&id)
	if err != nil {
		return nilUser, err
	}

	user.ID = id
	return user, nil
}

func (u *UsersRepo) CreateRandom(ctx context.Context, usersCount int64) error {
	var query strings.Builder
	query.WriteString("INSERT INTO users(name) VALUES ")
	for i := int64(1); i <= usersCount; i++ {
		query.WriteString("('User #" + strconv.FormatInt(i, 10) + "')")
		if i != usersCount {
			query.WriteString(", ")
		}
	}

	_, err := u.db.ExecContext(ctx, query.String())
	return err
}

func (u *UsersRepo) DeleteAll(ctx context.Context) error {
	_, err := u.db.ExecContext(ctx, "DELETE FROM users; ALTER SEQUENCE users_id_seq RESTART WITH 1")
	return err
}
