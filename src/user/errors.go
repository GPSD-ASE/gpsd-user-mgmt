package user

type BadRequest struct {
	Err error
}

func (e BadRequest) Error() string {
	return "User not found"
}

type Unauthorized struct {
	Err error
}

func (e Unauthorized) Error() string {
	return "User does not have required access"
}

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
