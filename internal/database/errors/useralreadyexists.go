package errors

// UserAlreadyExists represents that the system failed to add the new user to the database because
// and instance of that user already exists in the database.
type UserAlreadyExists struct{}

// CreateUserAlreadyExists creates a new instance of the UserAlreadyExists error.
func CreateUserAlreadyExists() UserAlreadyExists {
	return UserAlreadyExists{}
}

// Error returns an error message.
func (err UserAlreadyExists) Error() string {
	return "a user with that email already exists"
}
