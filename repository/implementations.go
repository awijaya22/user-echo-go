package repository

import (
	"context"
	"errors"
	"strconv"
	"strings"
)

func (r *Repository) GetUserByID(ctx context.Context, userID string) (user User, err error) {
	err = r.Db.QueryRowContext(ctx, "SELECT id, full_name, password, phone_number, sucessful_login_count FROM users WHERE id = $1", userID).Scan(&user.ID, &user.FullName, &user.Password, &user.PhoneNumber, &user.SuccessfulLoginCount)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			err = ErrUserNotFound
		}
		return
	}
	return

}

func (r *Repository) CreateUser(ctx context.Context, user User) (ID string, err error) {
	err = r.Db.QueryRowContext(ctx, "INSERT INTO users (full_name, password, phone_number) VALUES ($1, $2, $3) returning id", user.FullName, user.Password, user.PhoneNumber).Scan(&ID)
	if err != nil {
		return
	}
	return
}

func (r *Repository) GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (user User, err error) {
	err = r.Db.QueryRowContext(ctx, "SELECT id, full_name, password, phone_number, sucessful_login_count FROM users WHERE phone_number = $1", phoneNumber).Scan(&user.ID, &user.FullName, &user.Password, &user.PhoneNumber, &user.SuccessfulLoginCount)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			err = ErrUserNotFound
		}
		return
	}
	return
}

func (r *Repository) UpdateSuccessfulLoginCount(ctx context.Context, userID string) (err error) {
	_, err = r.Db.ExecContext(ctx, "UPDATE users SET sucessful_login_count = sucessful_login_count + 1 WHERE id = $1", userID)
	if err != nil {
		return
	}
	return
}

func (r *Repository) UpdateUser(ctx context.Context, updateReq UpdateUser) (err error) {
	var args []interface{}
	var setClauses []string

	if updateReq.FullName != nil {
		setClauses = append(setClauses, "full_name = $"+strconv.Itoa(len(args)+1))
		args = append(args, updateReq.FullName)
	}

	if updateReq.PhoneNumber != nil {
		setClauses = append(setClauses, "phone_number = $"+strconv.Itoa(len(args)+1))
		args = append(args, updateReq.PhoneNumber)
	}

	if len(setClauses) == 0 {
		return errors.New("no fields to update")
	}

	args = append(args, updateReq.ID)

	query := "UPDATE users SET " + strings.Join(setClauses, ", ") + " WHERE id = $ " + strconv.Itoa(len(args))
	_, err = r.Db.ExecContext(ctx, query, args...)
	if err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint \"users_phone_number_key\"" {
			return ErrPhoneNumberAlreadyExists
		}
		return  err
	}

	return
}
