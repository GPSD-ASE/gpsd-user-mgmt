package user

import (
	"context"
	"errors"
	"gpsd-user-mgmt/src/db"
	"log"
	"log/slog"
	"strings"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UserId       int    `json:"id"`
	UserName     string `json:"username"`
	Email        string `json:"email"`
	PasswordHash string `json:"password,omitempty"`
	DeviceID     string `json:"deviceID"`
	Role         string `json:"role"`
}

const (
	REPORTER = 2
	ADMIN    = 1
)

const (
	get_user = `SELECT user_id, user_name, email, device_id, role_name
				FROM public.user
				JOIN public.user_role ON public.user.role_id = public.user_role.role_id
				WHERE user_id = $1`

	get_user_from_name = `SELECT user_id, user_name, email, password_hash, device_id, role_name
				FROM public.user
				JOIN public.user_role ON public.user.role_id = public.user_role.role_id
				WHERE user_name = $1`

	get_users = `SELECT user_id, user_name, email, device_id, role_name
				FROM public.user
				JOIN public.user_role ON public.user.role_id = public.user_role.role_id
				LIMIT $1
				OFFSET $2`

	add_user = `INSERT INTO public.user (user_name, email, password_hash, device_id, role_id)
				VALUES ($1, $2, $3, $4, $5)
				RETURNING user_id`

	update_user = `UPDATE public.user
				   SET user_name = $2, email = $3, device_id = $4, role_id = $5
				   WHERE user_id = $1`

	delete_user = `DELETE FROM public.user 
				   WHERE user_id = $1`
)

func GetUser(id int) (User, error) {
	var result User
	row := db.Pool.QueryRow(context.Background(), get_user, id)

	err := row.Scan(
		&result.UserId,
		&result.UserName,
		&result.Email,
		&result.DeviceID,
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

func GetUserFromName(name string) (User, error) {
	var result User
	row := db.Pool.QueryRow(context.Background(), get_user_from_name, name)

	err := row.Scan(
		&result.UserId,
		&result.UserName,
		&result.Email,
		&result.PasswordHash,
		&result.DeviceID,
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
		slog.Error(err.Error())
		return nil, err
	}

	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.UserId,
			&user.UserName,
			&user.Email,
			&user.DeviceID,
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

func (user *User) AddUser() error {
	roleID := getRoleID(user.Role)

	hash, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hash)

	row := db.Pool.QueryRow(
		context.Background(),
		add_user,
		user.UserName,
		user.Email,
		user.PasswordHash,
		user.DeviceID,
		roleID,
	)

	var userId int
	err = row.Scan(
		&userId,
	)
	user.UserId = userId
	user.Role = getRoleString(roleID)
	user.PasswordHash = ""

	if err != nil {
		return err
	}

	return nil
}

func getRoleID(role string) int {
	var roleID int
	switch strings.ToLower(role) {
	case "2":
		roleID = REPORTER
	case "1":
		roleID = ADMIN
	default:
		roleID = REPORTER
	}
	return roleID
}

func getRoleString(id int) string {
	var role string
	switch id {
	case 2:
		role = "reporter"
	case 1:
		role = "admin"
	default:
		role = "reporter"
	}
	return role
}

func UpdateUser(userId int, user User) error {
	savedUser, err := GetUser(userId)
	if err != nil {
		return err
	}
	roleID := getRoleID(user.Role)

	if user.Email == "" {
		user.Email = savedUser.Email
	}

	if user.UserName == "" {
		user.UserName = savedUser.UserName
	}

	if user.DeviceID == "" {
		user.DeviceID = savedUser.DeviceID
	}

	_, err = db.Pool.Query(
		context.Background(),
		update_user,
		userId,
		user.UserName,
		user.Email,
		user.DeviceID,
		roleID,
	)

	return err
}

func DeleteUser(userId int) error {
	_, err := GetUser(userId)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	_ = db.Pool.QueryRow(
		context.Background(),
		delete_user,
		userId,
	)

	return nil
}
