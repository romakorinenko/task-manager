package constant

const UserSessionKey = "user"
const AdminSessionKey = "admin"

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

const (
	ActiveUser  = true
	BlockedUser = false
)
