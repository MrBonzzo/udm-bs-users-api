package users

import (
	"fmt"
	"main/datasources/mysql/users_db"
	"main/utils/date_utils"
	"main/utils/errors"
	"strings"
)

const (
	indexUniqueEmail = "email_UNIQUE"
	errorNoRows      = "no rows in result set"
	queryInsertUser  = "INSERT INTO users(first_name, last_name, email, date_created) VALUES(?, ?, ?, ?);"
	querySelectUser  = "SELECT id, first_name, last_name, email, date_created FROM users WHERE id = ?"
)

var (
	userDB = make(map[int64]*User)
)

func (user *User) Save() *errors.RestErr {

	stmt, err := users_db.Client.Prepare(queryInsertUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	user.DateCreated = date_utils.GetNowStirng()
	insertResult, err := stmt.Exec(user.FirstName, user.LastName, user.Email, user.DateCreated)
	if err != nil {
		if strings.Contains(err.Error(), indexUniqueEmail) {
			return errors.NewBadRequestError(
				fmt.Sprintf("email %s already exists", user.Email),
			)
		}
		return errors.NewInternalServerError(
			fmt.Sprintf("error when trying to save user: %s", err.Error()),
		)
	}
	userId, err := insertResult.LastInsertId()
	if err != nil {
		return errors.NewInternalServerError(fmt.Sprintf("error when trying to save user: %s", err.Error()))
	}

	user.Id = userId
	return nil
}

func (user *User) Get() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(querySelectUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	// result, err := stmt.Query(user.Id) call needs defer result.Close() operation,
	// otherwise result := stmt.QueryRow(user.Id) doesn't need err variable nor closing
	result := stmt.QueryRow(user.Id)
	if err := result.Scan(
		&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated,
	); err != nil {
		if strings.Contains(err.Error(), errorNoRows) {
			return errors.NewNotFoundError(fmt.Sprintf("user %d not found", user.Id))
		}
		return errors.NewInternalServerError(
			fmt.Sprintf("error when trying to save user %d: %s", user.Id, err.Error()),
		)
	}

	return nil
}
