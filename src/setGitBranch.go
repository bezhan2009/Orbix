package src

import "goCmd/system"

func SetGitBranch() {
	var errGitBranch error
	system.GitBranch, errGitBranch = GetCurrentGitBranch()
	if errGitBranch != nil {
		system.GitBranch = ""
	}
}
