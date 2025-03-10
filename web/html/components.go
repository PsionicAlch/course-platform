package html

import "github.com/PsionicAlch/course-platform/internal/database/models"

type NavbarComponent struct {
	User *models.UserModel
}

func NewNavbarComponent(user *models.UserModel) *NavbarComponent {
	return &NavbarComponent{
		User: user,
	}
}

/*
FormControlComponent is a data type to represent the properties required
to render a form-control component.

Form Control Props Required:

  - Label (a string for the label)

  - Type (a string to control the input type)

  - Name (a string to be used as the input name, id and label for)

  - Errors (a string array of errors)

  - ValidationURL (a string pointing to the URL that the input can use to validate)

  - Value (a string to be used as input value)
*/
type FormControlComponent struct {
	Label         string
	Type          string
	Name          string
	Errors        []string
	ValidationURL string
	Value         string
}

/*
PasswordControlComponent is a data type to represent the properties required
to render a form-control component.

Form Control Props Required:

  - Label (a string for the label)

  - Name (a string to be used as the input name, id and label for)

  - Errors (a string array of errors)

  - ValidationURL (a string pointing to the URL that the input can use to validate)

  - Value (a string to be used as input value)
*/
type PasswordControlComponent struct {
	Label         string
	Name          string
	Errors        []string
	ValidationURL string
	Value         string
}

type LoginFormComponent struct {
	EmailInput    *FormControlComponent
	PasswordInput *PasswordControlComponent
	ErrorMessage  string
}

type SignupFormComponent struct {
	FirstNameInput       *FormControlComponent
	LastNameInput        *FormControlComponent
	EmailInput           *FormControlComponent
	PasswordInput        *PasswordControlComponent
	ConfirmPasswordInput *PasswordControlComponent
	ErrorMessage         string
}

type ForgotPasswordFormComponent struct {
	EmailInput   *FormControlComponent
	ErrorMessage string
}

type ResetPasswordFormComponent struct {
	EmailToken           string
	PasswordInput        *PasswordControlComponent
	ConfirmPasswordInput *PasswordControlComponent
	ErrorMessage         string
}

type NewDiscountFormComponent struct {
	TitleInput       *FormControlComponent
	DescriptionInput *FormControlComponent
	UsesInput        *FormControlComponent
	AmountInput      *FormControlComponent
	ErrorMessage     string
}

type CoursePurchaseFormComponent struct {
	AffiliateCodeInput      *FormControlComponent
	AffiliatePointsInput    *FormControlComponent
	DiscountCodeInput       *FormControlComponent
	AffiliateCodeDiscount   uint
	AffiliatePointsDiscount uint
	DiscountCodeDiscount    uint
	Total                   float64
	CourseSlug              string
	ErrorMessage            string
}

type ChangeFirstNameFormComponent struct {
	FirstNameInput *FormControlComponent
}

type ChangeLastNameFormComponent struct {
	LastNameInput *FormControlComponent
}

type ChangeEmailFormComponent struct {
	EmailInput *FormControlComponent
}

type ChangePasswordFormComponent struct {
	PreviousPasswordInput *PasswordControlComponent
	NewPasswordInput      *PasswordControlComponent
}

type TutorialsListComponent struct {
	Tutorials    []*models.TutorialModel
	LastTutorial *models.TutorialModel
	QueryURL     string
	ErrorMessage string
}

type CommentsListComponent struct {
	Comments     []*models.CommentModel
	LastComment  *models.CommentModel
	Users        map[string]*models.UserModel
	QueryURL     string
	ErrorMessage string
}

type CoursesListComponent struct {
	Courses      []*models.CourseModel
	LastCourse   *models.CourseModel
	QueryURL     string
	ErrorMessage string
}

type AdminUsersListComponent struct {
	Users               []*models.UserModel
	LastUser            *models.UserModel
	TutorialsLiked      map[string]uint
	TutorialsBookmarked map[string]uint
	CoursesBought       map[string]uint
	TutorialsWritten    map[string]uint
	CoursesWritten      map[string]uint
	BaseURL             string
	URLQuery            string
	ErrorMessage        string
}

type SelectComponent struct {
	Name         string
	Options      map[string]string
	Selected     string
	URL          string
	ErrorMessage string
}

type AdminTutorialsListComponent struct {
	Tutorials    []*models.TutorialModel
	LastTutorial *models.TutorialModel
	Authors      map[string]*models.UserModel
	Keywords     map[string][]string
	Comments     map[string]uint
	Likes        map[string]uint
	Bookmarks    map[string]uint
	BaseURL      string
	URLQuery     string
	ErrorMessage string
}

type AdminDiscountsListComponent struct {
	Discounts    []*models.DiscountModel
	LastDiscount *models.DiscountModel
	DiscountUsed map[string]uint
	BaseURL      string
	URLQuery     string
	ErrorMessage string
}

type AdminCoursesListComponent struct {
	Courses      []*models.CourseModel
	LastCourse   *models.CourseModel
	Authors      map[string]*models.UserModel
	Keywords     map[string][]string
	Purchases    map[string]uint
	BaseURL      string
	URLQuery     string
	ErrorMessage string
}

type AdminCommentsListComponent struct {
	Comments     []*models.CommentModel
	LastComment  *models.CommentModel
	Users        map[string]*models.UserModel
	Tutorials    map[string]*models.TutorialModel
	BaseURL      string
	URLQuery     string
	ErrorMessage string
}

type AffiliateHistoryListComponent struct {
	AffiliateHistory     []*models.AffiliatePointsHistoryModel
	LastAffiliateHistory *models.AffiliatePointsHistoryModel
	QueryURL             string
	ErrorMessage         string
}

type AdminCoursePurchaseListComponent struct {
	Purchases    []*models.CoursePurchaseModel
	LastPurchase *models.CoursePurchaseModel
	Users        map[string]*models.UserModel
	Courses      map[string]*models.CourseModel
	BaseURL      string
	URLQuery     string
	ErrorMessage string
}

type AdminRefundsListComponent struct {
	Refunds      []*models.RefundModel
	LastRefund   *models.RefundModel
	Users        map[string]*models.UserModel
	Courses      map[string]*models.CourseModel
	BaseURL      string
	URLQuery     string
	ErrorMessage string
}
