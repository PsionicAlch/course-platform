package authentication

const (
	AuthenticationToken = "authentication"
	EmailToken          = "email"
)

func NewToken() (string, error) {
	tokenBytes, err := RandomBytes(32)
	if err != nil {
		return "", err
	}

	return BytesToString(tokenBytes), nil
}
