package _chan

import "strings"

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
		SaveVarsFn = nil
		LoadConfigsFn = nil
	}

	// Scripts -> Orbix(LoopData)
	if chanName == "scripts__orbix_loop_data" {
		LoopData = nil
	}
}
