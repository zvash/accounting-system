// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: user.sql

package sql

import (
	"context"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (username, name, email, password)
VALUES ($1, $2, $3, $4)
RETURNING username, name, email, password, password_changed_at, created_at
`

type CreateUserParams struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, createUser,
		arg.Username,
		arg.Name,
		arg.Email,
		arg.Password,
	)
	var i User
	err := row.Scan(
		&i.Username,
		&i.Name,
		&i.Email,
		&i.Password,
		&i.PasswordChangedAt,
		&i.CreatedAt,
	)
	return i, err
}

const getUserByUserName = `-- name: GetUserByUserName :one
SELECT username, name, email, password, password_changed_at, created_at FROM users WHERE username = $1 LIMIT 1
`

func (q *Queries) GetUserByUserName(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRow(ctx, getUserByUserName, username)
	var i User
	err := row.Scan(
		&i.Username,
		&i.Name,
		&i.Email,
		&i.Password,
		&i.PasswordChangedAt,
		&i.CreatedAt,
	)
	return i, err
}