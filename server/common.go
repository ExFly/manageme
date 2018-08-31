package main

import (
	"context"
	"errors"

	"github.com/exfly/manageme/model"
)

var (
	// ErrNotLogined like the name
	ErrNotLogined = errors.New("Not Logined")
	// ErrNoPermission like the name
	ErrNoPermission = errors.New("No Permission")
	// ErrBadRequest like the name
	ErrBadRequest    = errors.New("Bad Request")
	ErrInvalidID     = errors.New("Invalid ID")
	ErrNotFound      = errors.New("Not found")
	ErrInternalError = errors.New("Internal Error")
)

func getUser(ctx context.Context) *model.User {
	user, ok := ctx.Value("user").(*model.User)
	if !ok {
		return nil
	}
	return user
}

func mergeError(errs []error) error {
	if errs == nil {
		return nil
	}
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}
