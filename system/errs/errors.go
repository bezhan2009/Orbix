package errs

import (
	"errors"
	"fmt"
	"goCmd/system"
)

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
)
