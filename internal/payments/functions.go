package payments

import (
	"github.com/PsionicAlch/psionicalch-home/internal/database"
)

// The price of a course.
const CoursePrice float64 = 200.0
const AffiliateCodeDiscount float64 = 10.0 / 100.0
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

	total := int64((CoursePrice - (CoursePrice * (affiliateCodeDiscount + discountCodeDiscount + affiliatePointsDiscount))) * 100)

	return total, nil
}

func (payment *Payments) ValidateAffiliateCode(userId, affiliateCode string) (float64, error) {
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

	affiliateUser, err := payment.Database.GetUserByAffiliateCode(affiliateCode)
	if err != nil {
		payment.ErrorLog.Printf("Failed to get user by affiliate code (\"%s\"): %s\n", affiliateCode, err)
		return 0, err
	}

	if affiliateUser == nil {
		return 0, ErrInvalidAffiliateCode
	}

	return AffiliateCodeDiscount, nil
}

func (payment *Payments) ValidateDiscountCode(discountCode string) (float64, error) {
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

	if discountUsed > discount.Uses {
		return 0, ErrInvalidDiscountCode
	}

	return float64(discount.Discount) / 100.0, nil
}

func (payment *Payments) ValidateAffiliatePointsUsed(userId string, affiliatePointsUsed uint) (float64, error) {
	user, err := payment.Database.GetUserByID(userId, database.All)
	if err != nil {
		payment.ErrorLog.Printf("Failed to get user by ID (\"%s\"): %s\n", userId, err)
		return 0, err
	}

	if user.AffiliatePoints < affiliatePointsUsed {
		return 0, ErrInsufficientAffiliatePoints
	}

	return float64(affiliatePointsUsed) * AffiliatePointDiscount, nil
}
