package errors

type UserAlreadyExists struct {
	msg string
}

func CreateUserAlreadyExists() UserAlreadyExists {
	return UserAlreadyExists{
		msg: "a user with this email address already exists",
	}
}

func (err UserAlreadyExists) Error() string {
	return err.msg
}
