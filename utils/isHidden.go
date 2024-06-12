package utils

import "goCmd/safe/block/bun/component"

func IsHidden() bool {
	return component.IsBanned()
}
