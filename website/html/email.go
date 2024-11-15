package html

type BaseEmail struct {
	Title string
}

func NewBaseEmail(title string) BaseEmail {
	return BaseEmail{
		Title: title,
	}
}

type GreetingEmail struct {
	BaseEmail
	FirstName     string
	LatestCourses []struct {
		Name        string
		Description string
		Slug        string
	}
	AffiliateCode string
}

func NewGreetingEmail(firstName, affiliateCode string) *GreetingEmail {
	return &GreetingEmail{
		BaseEmail: NewBaseEmail("Welcome to PsionicAlch"),
		FirstName: firstName,
		LatestCourses: []struct {
			Name        string
			Description string
			Slug        string
		}{},
		AffiliateCode: affiliateCode,
	}
}
