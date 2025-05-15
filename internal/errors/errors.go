package interrors

import "errors"

var (
	ErrBadCred  = errors.New("bad creds")
	ErrBusyLogin  = errors.New("login is busy")
)