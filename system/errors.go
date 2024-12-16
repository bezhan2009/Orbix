package system

import (
	"errors"
	"fmt"
)

var (
	UserAlreadyExists      = errors.New("user Already Exists")
	UserNotFound           = errors.New("user not found")
	PasswordIncorrect      = errors.New("password incorrect")
	UserNotFoundOrPassword = errors.New("user not found")
	SessionIsNil           = errors.New("session is Nil")
	SessionNotExists       = errors.New("session Not Exists")
	SessionExpired         = errors.New("session Expired")
	ExactMatchNotFound     = errors.New(fmt.Sprintf("exact match not found in %s", OrbixRunningUsersFileName))
)
