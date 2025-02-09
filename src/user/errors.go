package user

type NotFound struct {
	Err error
}

func (e NotFound) Error() string {
	return "User not found"
}

type InternalServerError struct {
	Err error
}

func (e InternalServerError) Error() string {
	return "Internal server error"
}
