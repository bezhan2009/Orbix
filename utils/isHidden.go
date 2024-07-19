package utils

import "goCmd/pkg/safe/block/bun/component"

func IsHidden() bool {
	return component.IsBanned()
}
