package user

import (
	"context"
	"errors"
	"gpsd-user-mgmt/src/db"
	"log"

	"github.com/jackc/pgx/v5"
)

type User struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	DevID string `json:"devID"`
	Role  string `json:"role"`
}

const (
	GET_USER    = "SELECT id, name, deviceID, role FROM users WHERE id = $1"
	GET_USERS   = "SELECT id, name, deviceID, role FROM users LIMIT $1 OFFSET $2"
	ADD_USER    = "INSERT INTO users (name, role, deviceID, createdAt, updatedAt) VALUES ($1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) RETURNING id"
	UPDATE_USER = "UPDATE users set name = $1, role = $2, updatedAt = CURRENT_TIMESTAMP WHERE id = $3"
	DELETE_USER = "DELETE FROM users WHERE id = $1"
)

func GetUser(id string) (User, error) {
	var result User
	row := db.Pool.QueryRow(context.Background(), GET_USER, id)

	err := row.Scan(
		&result.Id,
		&result.Name,
		&result.DevID,
		&result.Role,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return result, NotFound{}
		}
		return result, err
	}

	return result, nil
}

func GetUsers(limit, offset string) ([]User, error) {
	var result []User
	rows, err := db.Pool.Query(context.Background(), GET_USERS, limit, offset)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var user User
		err = rows.Scan(
			&user.Id,
			&user.Name,
			&user.DevID,
			&user.Role,
		)
		if err != nil {
			log.Printf("Scan error: %s\n", err.Error())
			return nil, err
		}
		result = append(result, user)
	}

	return result, nil
}

func AddUser(user User) (int, error) {
	row := db.Pool.QueryRow(
		context.Background(),
		ADD_USER,
		user.Name,
		user.Role,
		user.DevID,
	)

	var userId int
	err := row.Scan(
		&userId,
	)

	if err != nil {
		return userId, err
	}

	return userId, nil
}

func UpdateUser(userId int, user User) error {
	_ = db.Pool.QueryRow(
		context.Background(),
		UPDATE_USER,
		user.Name,
		user.Role,
		userId,
	)

	return nil
}

func DeleteUser(userId int) error {
	_ = db.Pool.QueryRow(
		context.Background(),
		DELETE_USER,
		userId,
	)

	return nil
}
