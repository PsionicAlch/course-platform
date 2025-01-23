package payments

import "errors"

var (
	// ErrUserAlreadyOwnsCourse represents the user already owns the course that they're trying to buy.
	ErrUserAlreadyOwnsCourse = errors.New("user already owns this course")

	// ErrCantUseOwnAffiliateCode represents the user is trying to get a discount by using their own affiliate code.
	ErrCantUseOwnAffiliateCode = errors.New("user can't use their own affiliate code")

	// ErrUserDoesNotExist represents that the user trying to buy the course could not be found in the database.
	ErrUserDoesNotExist = errors.New("user does not exist")

	// ErrInvalidAffiliateCode represents the affiliate code provided doesn't belong to any user.
	ErrInvalidAffiliateCode = errors.New("invalid affiliate code provided")

	// ErrInvalidDiscountCode represents the discount code provided either doesn't exist, isn't active, or has been
	// used too much.
	ErrInvalidDiscountCode = errors.New("invalid discount code provided")

	// ErrInsufficientAffiliatePoints represents the user has less affiliate points available than what they're trying
	// to use.
	ErrInsufficientAffiliatePoints = errors.New("user doesn't have enough affiliate points")

	// ErrUserHasNotBoughtCourse represents the user has yet to purchase this course.
	ErrUserHasNotBoughtCourse = errors.New("user hasn't purchased the course")
)
