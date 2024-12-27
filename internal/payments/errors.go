package payments

import "errors"

var (
	ErrUserAlreadyOwnsCourse       = errors.New("user already owns this course")
	ErrCantUseOwnAffiliateCode     = errors.New("user can't use their own affiliate code")
	ErrUserDoesNotExist            = errors.New("user does not exist")
	ErrInvalidAffiliateCode        = errors.New("invalid affiliate code provided")
	ErrInvalidDiscountCode         = errors.New("invalid discount code provided")
	ErrInsufficientAffiliatePoints = errors.New("user doesn't have enough affiliate points")
	ErrUserHasNotBoughtCourse      = errors.New("user hasn't purchased the course")
)
