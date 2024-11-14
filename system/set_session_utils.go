package system

func SetGitBranch(sd *Session) {
	var errGitBranch error
	sd.GitBranch, errGitBranch = GetCurrentGitBranch()
	if errGitBranch != nil {
		sd.GitBranch = ""
	}
}

func SetPath(sd *Session) {
	sd.Path = Getwd()
}
