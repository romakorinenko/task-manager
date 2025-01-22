package constant

const UserSessionKey = "user"

const (
	Blocker = iota + 1
	High
	Medium
	Low
)

const (
	AdminRole = "ADMIN"
	UserRole  = "USER"
)

const (
	OpenTaskStatus       = "OPEN"
	InProgressTaskStatus = "IN_PROGRESS"
	DoneTaskStatus       = "DONE"
)

var TaskStatuses = []string{OpenTaskStatus, InProgressTaskStatus, DoneTaskStatus}

const (
	ActiveUser  = true
	BlockedUser = false
)
