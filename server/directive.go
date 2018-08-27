package main

import (
	"context"

	"github.com/exfly/manageme/model"

	"github.com/99designs/gqlgen/graphql"
)

func Logined(ctx context.Context, next graphql.Resolver) (res interface{}, err error) {
	user := getUser(ctx)
	if user != nil {
		return next(ctx)
	} else {
		return nil, ErrNotLogined
	}
}
func RequirePermission(ctx context.Context, next graphql.Resolver, permission model.Permission, meta *map[string]interface{}) (res interface{}, err error) {
	result := model.CheckPermission(ctx, permission, meta)
	if result {
		return next(ctx)
	} else {
		return nil, nil
	}
}
