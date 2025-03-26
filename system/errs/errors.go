package errs

import (
	"errors"
	"fmt"
	"goCmd/system"
)

func TypeError(necessaryType string, gettedType string) (err error) {
	return errors.New(fmt.Sprintf("Type error: mismatched types %s and %s", necessaryType, gettedType))
}

var (
	UserAlreadyExists      = errors.New("user Already Exists")
	UserNotFound           = errors.New("user not found")
	PasswordIncorrect      = errors.New("password incorrect")
	UserNotFoundOrPassword = errors.New("user not found")
	SessionIsNil           = errors.New("session is Nil")
	SessionNotExists       = errors.New("session Not Exists")
	SessionExpired         = errors.New("session Expired")
	CommandArgsNotFound    = errors.New("command args Not Found")
	ExactMatchNotFound     = errors.New(fmt.Sprintf("exact match not found in %s", system.OrbixRunningUsersFileName))
	ValidationError        = errors.New("validation error: do not use special symbols")
	VariableDoesNotExist   = errors.New("variable Does Not Exist")
)
