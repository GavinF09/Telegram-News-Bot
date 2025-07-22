package commands

var (
	// fake database
	// echoUsers = make([]int, 0)
	// the val of the map is not important
	echoUsers = make(map[int64]bool)
)

func EchoMessage(chatID int64) bool {
	// check if echo is enabled for user
	_, ok := echoUsers[chatID]
	return ok
}

func ToggleEcho(userid int64) int {
	// see if user already enabled echo
	_, ok := echoUsers[userid]
	// user is already registered
	if ok {
		delete(echoUsers, userid)
		return 0
	} else {
		echoUsers[userid] = true
		return 1
	}
	// return -1, errors.New("something failed")
}
