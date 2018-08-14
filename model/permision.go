package model

import (
	"context"

	"github.com/spf13/viper"
)

// CheckPermission l
func CheckPermission(ctx context.Context, permission Permission, payloads ...interface{}) bool {
	// TODO: foolish checker!!!
	debug := viper.GetBool("server.debug")
	if debug && permission == PermissionDebug {
		return true
	}
	return false
}
