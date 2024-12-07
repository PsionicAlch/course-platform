package database

//go:generate enumer -type=AuthorizationLevel -json
type AuthorizationLevel int

const (
	All AuthorizationLevel = iota
	User
	Admin
	Author
)

type PaymentStatus int

const (
	Pending PaymentStatus = iota
	RequiresAction
	Processing
	Succeeded
	Failed
	Cancelled
	Refunded
)

func (p PaymentStatus) String() string {
	switch p {
	case Pending:
		return "Pending"
	case RequiresAction:
		return "Requires Action"
	case Processing:
		return "Processing"
	case Succeeded:
		return "Succeeded"
	case Failed:
		return "Failed"
	case Cancelled:
		return "Cancelled"
	case Refunded:
		return "Refunded"
	}

	return ""
}
