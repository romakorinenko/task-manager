package errs

type UserExistsErr struct{}

func (u UserExistsErr) Error() string {
	return "user with the login already exists"
}

type BadReqErr struct{}

func (b BadReqErr) Error() string {
	return "bad request"
}
