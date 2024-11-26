package database

//go:generate enumer -type=AuthorizationLevel -json
type AuthorizationLevel int

const (
	All AuthorizationLevel = iota
	User
	Admin
	Author
)
