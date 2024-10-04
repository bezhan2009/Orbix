package src

import "goCmd/system"

func SetGitBranch(sd *system.Session) {
	var errGitBranch error
	sd.GitBranch, errGitBranch = GetCurrentGitBranch()
	if errGitBranch != nil {
		sd.GitBranch = ""
	}
}

func SetPath(sd *system.Session) {
	sd.Path = Getwd()
}
