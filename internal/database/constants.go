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
	Disputed
)

// String converts PaymentStatus to a string.
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
	case Disputed:
		return "Disputed"
	default:
		return ""
	}
}

// PaymentStatusFromString converts a string to a PaymentStatus.
func PaymentStatusFromString(s string) PaymentStatus {
	switch s {
	case Pending.String():
		return Pending
	case RequiresAction.String():
		return RequiresAction
	case Processing.String():
		return Processing
	case Succeeded.String():
		return Succeeded
	case Failed.String():
		return Failed
	case Cancelled.String():
		return Cancelled
	case Refunded.String():
		return Refunded
	case Disputed.String():
		return Disputed
	default:
		return Pending
	}
}

type RefundStatus int

const (
	RefundPending RefundStatus = iota
	RefundRequiresAction
	RefundSucceeded
	RefundFailed
	RefundCancelled
	DisputeWarningNeedsResponse
	DisputeWarningUnderReview
	DisputeWarningClosed
	DisputeNeedsResponse
	DisputeUnderReview
	DisputeWon
	DisputeLost
)

// String converts a RefundStatus to a string.
func (r RefundStatus) String() string {
	switch r {
	case RefundPending:
		return "Refund Pending"
	case RefundRequiresAction:
		return "Refund Requires Action"
	case RefundSucceeded:
		return "Refund Succeeded"
	case RefundFailed:
		return "Refund Failed"
	case RefundCancelled:
		return "Refund Cancelled"
	case DisputeWarningNeedsResponse:
		return "Dispute Warning Needs Response"
	case DisputeWarningUnderReview:
		return "Dispute Warning Under Review"
	case DisputeWarningClosed:
		return "Dispute Warning Closed"
	case DisputeNeedsResponse:
		return "Dispute Needs Response"
	case DisputeUnderReview:
		return "Dispute Under Review"
	case DisputeWon:
		return "Dispute Won"
	case DisputeLost:
		return "Dispute Lost"
	default:
		return ""
	}
}

// RefundStatusFromString converts a string to a RefundStatus.
func RefundStatusFromString(s string) RefundStatus {
	switch s {
	case RefundRequiresAction.String():
		return RefundRequiresAction
	case RefundSucceeded.String():
		return RefundSucceeded
	case RefundFailed.String():
		return RefundFailed
	case RefundCancelled.String():
		return RefundCancelled
	case DisputeWarningNeedsResponse.String():
		return DisputeWarningNeedsResponse
	case DisputeWarningUnderReview.String():
		return DisputeWarningUnderReview
	case DisputeWarningClosed.String():
		return DisputeWarningClosed
	case DisputeNeedsResponse.String():
		return DisputeNeedsResponse
	case DisputeUnderReview.String():
		return DisputeUnderReview
	case DisputeWon.String():
		return DisputeWon
	case DisputeLost.String():
		return DisputeLost
	default:
		return RefundPending
	}
}
