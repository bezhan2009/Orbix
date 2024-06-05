package utils

import "goCmd/security/block/bun/component"

func IsHidden() bool {
	return component.IsBanned()
}
