package user

import (
	"errors"
	"fmt"
)

var errFirstnameRequired = errors.New("Firstname is required")
var errLastnameRequired = errors.New("Lastname is required")

type ErrNotFound struct {
	userID string
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("User '%s' doesn't exist", e.userID)
}
