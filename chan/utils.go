package _chan

import (
	"goCmd/structs"
	"strings"
)

func UpdateChan(chanName string) {
	chanName = strings.ToLower(strings.TrimSpace(chanName))

	// System -> User
	if chanName == "system__user" {
		UserStatusAuth = false
		UserNameStatus = false
	}

	// Orbix -> SRC(Restore)
	if chanName == "orbix__src_restore" {
		DirUser = ""
		UserName = ""
	}

	// Orbix -> SRC(Newprompt)
	if chanName == "orbix__src_prompt" {
		UseNewPrompt = false
		UseOldPrompt = false
		EnableSecure = false
	}

	// Scripts -> Orbix(Func)
	if chanName == "scripts__orbix_func" {
		SaveVarsFn = func() {

		}
		LoadConfigsFn = func() error {
			return nil
		}
	}

	// Scripts -> Orbix(LoopData)
	if chanName == "scripts__orbix_loop_data" {
		LoopData = &structs.OrbixLoopData{}
	}

	// Environment -> SRC(nikcame)
	if chanName == "environment__src_get_user_nickname" {
		GetUserNikcName = func() string {
			return ""
		}
	}
}
