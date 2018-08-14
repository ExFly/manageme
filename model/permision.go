package model

import (
	"context"

	"github.com/spf13/viper"
)

func CheckPermission(ctx context.Context, permission Permission, payloads ...interface{}) bool {
	debug := viper.GetBool("server.debug")
	if debug && permission == PermissionDebug {
		return true
	}
	return false
}
