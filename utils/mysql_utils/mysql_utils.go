package mysql_utils

import (
	"main/utils/errors"
	"strings"

	"github.com/go-sql-driver/mysql"
)

const (
	errorNoRows              = "no rows in result set"
	alreadyExistsErrorNumber = 1062
)

func ParseError(err error) *errors.RestErr {
	sqlErr, ok := err.(*mysql.MySQLError)
	if !ok {
		if strings.Contains(err.Error(), errorNoRows) {
			return errors.NewNotFoundError("no record matching given id")
		}
		return errors.NewInternalServerError("error parsing database response")
	}
	switch sqlErr.Number {
	case alreadyExistsErrorNumber:
		return errors.NewBadRequestError("invalid data")
	}
	return errors.NewInternalServerError("error parsing request")
}
