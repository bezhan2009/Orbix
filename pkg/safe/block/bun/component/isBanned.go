package component

import "os"

func IsBanned() bool {
	_, err := os.Open("safe/block/bun/component/bunnedUser.json")

	if err != nil {
		return false
	} else {
		return true
	}
}
