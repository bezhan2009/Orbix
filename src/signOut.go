package src

func SignOutUtil(username string) {
	removeUserFromRunningFile(username)
	Orbix("")
}
