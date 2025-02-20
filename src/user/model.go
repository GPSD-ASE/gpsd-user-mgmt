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
	get_user    = "SELECT id, name, deviceID, role FROM users WHERE id = $1"
	get_users   = "SELECT id, name, deviceID, role FROM users LIMIT $1 OFFSET $2"
	add_user    = "INSERT INTO users (name, role, deviceID, createdAt, updatedAt) VALUES ($1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) RETURNING id"
	update_user = "UPDATE users set name = $1, role = $2, updatedAt = CURRENT_TIMESTAMP WHERE id = $3"
	delete_user = "DELETE FROM users WHERE id = $1"

	GET_USER_INCIDENTS = "SELECT incidents.* FROM incidents JOIN user_incidents ON incidents.id = user_incidents.incident_id WHERE user_id = $1"
)

func GetUser(id int) (User, error) {
	var result User
	row := db.Pool.QueryRow(context.Background(), get_user, id)

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

func GetUserIncidents(id string) (User, error) {
	var result User
	row := db.Pool.QueryRow(context.Background(),
		GET_USER_INCIDENTS,
		id,
	)

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

func GetUsers(limit, offset int) ([]User, error) {
	var result []User
	rows, err := db.Pool.Query(context.Background(), get_users, limit, offset)
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
		add_user,
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
	_, err := GetUser(userId)
	if err != nil {
		return err
	}

	_, err = db.Pool.Query(
		context.Background(),
		update_user,
		user.Name,
		user.Role,
		userId,
	)

	return err
}

func DeleteUser(userId int) error {
	_, err := GetUser(userId)
	if err != nil {
		return err
	}
	_ = db.Pool.QueryRow(
		context.Background(),
		delete_user,
		userId,
	)

	return nil
}
