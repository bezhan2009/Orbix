package src

func SignOutUtil(user string, isWorking *bool) string {
	username, isSuccess := CheckUser(user)

	if isSuccess {
		*isWorking = false
	}

	return username
}
