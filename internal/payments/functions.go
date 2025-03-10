package payments

import (
	"time"

	"math/rand"

	"github.com/PsionicAlch/course-platform/internal/database"
	"github.com/oklog/ulid/v2"
)

// CoursePrice is the price of a course.
const CoursePrice float64 = 200.0

// AffiliateCodeDiscount is the discount percentage from using an affiliate code.
const AffiliateCodeDiscount float64 = 10.0 / 100.0

// AffiliatePointDiscount is the discount percentage from using an affiliate point.
const AffiliatePointDiscount float64 = 1.0 / 100.0

// CalculatePrice will determine the cost of the course in cents. It takes into account the discounts when determining the price.
func (payment *Payments) CalculatePrice(userId, affiliateCode, discountCode string, affiliatePointsUsed uint) (int64, error) {
	affiliateCodeDiscount, err := payment.ValidateAffiliateCode(userId, affiliateCode)
	if err != nil {
		payment.ErrorLog.Printf("Failed to validate affiliate code: %s\n", err)
		return 0, ErrInvalidAffiliateCode
	}

	discountCodeDiscount, err := payment.ValidateDiscountCode(discountCode)
	if err != nil {
		payment.ErrorLog.Printf("Failed to validate discount code: %s\n", err)
		return 0, ErrInvalidDiscountCode
	}

	affiliatePointsDiscount, err := payment.ValidateAffiliatePointsUsed(userId, affiliatePointsUsed)
	if err != nil {
		payment.ErrorLog.Printf("Failed to validate affiliate points: %s\n", err)
		return 0, ErrInsufficientAffiliatePoints
	}

	// Ensure that the discount will never be more than 100%
	discount := affiliateCodeDiscount + discountCodeDiscount + affiliatePointsDiscount
	if discount > 1.0 {
		discount = 1.0
	}

	total := int64((CoursePrice - (CoursePrice * discount)) * 100)
	if total < 0 {
		total = 0
	}

	return total, nil
}

// ValidateAffiliateCode ensure that the provided affiliate code is allowed to be used.
func (payment *Payments) ValidateAffiliateCode(userId, affiliateCode string) (float64, error) {
	if affiliateCode == "" {
		return 0, nil
	}

	user, err := payment.Database.GetUserByID(userId, database.All)
	if err != nil {
		payment.ErrorLog.Printf("Failed to get user by ID (\"%s\"): %s\n", userId, err)
		return 0, err
	}

	if user == nil {
		return 0, ErrUserDoesNotExist
	}

	if user.AffiliateCode == affiliateCode {
		return 0, ErrCantUseOwnAffiliateCode
	}

	affiliateUser, err := payment.Database.GetUserByAffiliateCode(affiliateCode, database.All)
	if err != nil {
		payment.ErrorLog.Printf("Failed to get user by affiliate code (\"%s\"): %s\n", affiliateCode, err)
		return 0, err
	}

	if affiliateUser == nil {
		return 0, ErrInvalidAffiliateCode
	}

	return AffiliateCodeDiscount, nil
}

// ValidateDiscountCode ensures the provided discount code is allowed to be used.
func (payment *Payments) ValidateDiscountCode(discountCode string) (float64, error) {
	if discountCode == "" {
		return 0, nil
	}

	discount, err := payment.Database.GetDiscountByCode(discountCode)
	if err != nil {
		payment.ErrorLog.Printf("Failed to get discount from the database: %s\n", err)
		return 0, err
	}

	if discount == nil {
		return 0, ErrInvalidDiscountCode
	}

	if !discount.Active {
		return 0, ErrInvalidDiscountCode
	}

	discountUsed, err := payment.Database.CountCoursesWhereDiscountWasUsed(discountCode)
	if err != nil {
		payment.ErrorLog.Printf("Failed to count the number of courses bought with discount code (\"%s\"): %s\n", discountCode, err)
		return 0, err
	}

	if discountUsed >= discount.Uses {
		return 0, ErrInvalidDiscountCode
	}

	return float64(discount.Discount) / 100.0, nil
}

// ValidateAffiliatePointsUsed ensures that the user is allowed to use the provided amount of affiliate points.
func (payment *Payments) ValidateAffiliatePointsUsed(userId string, affiliatePointsUsed uint) (float64, error) {
	user, err := payment.Database.GetUserByID(userId, database.All)
	if err != nil {
		payment.ErrorLog.Printf("Failed to get user by ID (\"%s\"): %s\n", userId, err)
		return 0, err
	}

	if user.AffiliatePoints < int(affiliatePointsUsed) {
		return 0, ErrInsufficientAffiliatePoints
	}

	return float64(affiliatePointsUsed) * AffiliatePointDiscount, nil
}

// GeneratePaymentKey creates a new and unique payment key.
func GeneratePaymentKey() (string, error) {
	now := time.Now()
	entropy := rand.New(rand.NewSource(now.UnixNano()))
	ms := ulid.Timestamp(now)
	id, err := ulid.New(ms, entropy)

	return id.String(), err
}
