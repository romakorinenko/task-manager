package errs

import "fmt"

type UserExistsErr struct{}

func (u UserExistsErr) Error() string {
	return fmt.Sprintf("user with the login already exists")
}

type BadReqErr struct{}

func (b BadReqErr) Error() string {
	return fmt.Sprintf("bad request")
}
