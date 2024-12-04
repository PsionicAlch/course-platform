package html

import "github.com/PsionicAlch/psionicalch-home/internal/database/models"

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
}

type SignupFormComponent struct {
	FirstNameInput       *FormControlComponent
	LastNameInput        *FormControlComponent
	EmailInput           *FormControlComponent
	PasswordInput        *PasswordControlComponent
	ConfirmPasswordInput *PasswordControlComponent
}

type ForgotPasswordFormComponent struct {
	EmailInput *FormControlComponent
}

type ResetPasswordFormComponent struct {
	EmailToken           string
	PasswordInput        *PasswordControlComponent
	ConfirmPasswordInput *PasswordControlComponent
}

type NewDiscountFormComponent struct {
	TitleInput       *FormControlComponent
	DescriptionInput *FormControlComponent
	UsesInput        *FormControlComponent
	AmountInput      *FormControlComponent
	ErrorMessage     string
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
