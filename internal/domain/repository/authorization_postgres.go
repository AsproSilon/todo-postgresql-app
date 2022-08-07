package repository

import (
	"aspro/pkg/client/postgresql"
	"aspro/pkg/client/postgresql/model"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type AuthorizationPostgres struct {
	db *sqlx.DB
}

func NewAuthorizationPostgres(db *sqlx.DB) *AuthorizationPostgres {
	return &AuthorizationPostgres{db: db}
}

func (r *AuthorizationPostgres) CreateUser(user model.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (name, username, password_hash) values ($1, $2, $3) RETURNING id", postgresql.UsersTable)
	row := r.db.QueryRow(query, user.Name, user.Username, user.Password)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *AuthorizationPostgres) GetUser(username, password string) (model.User, error) {
	var user model.User
	query := fmt.Sprintf("SELECT id FROM %s username=$1 AND password_hash=$2", postgresql.UsersTable)
	err := r.db.Get(&user, query, username, password)

	return user, err
}
