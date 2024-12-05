package html

import "github.com/PsionicAlch/psionicalch-home/internal/database/models"

type BasePage struct {
	Navbar *NavbarComponent
}

func NewBasePage(user *models.UserModel) BasePage {
	return BasePage{
		Navbar: NewNavbarComponent(user),
	}
}

type AccountsLoginPage struct {
	BasePage
	LoginForm *LoginFormComponent
}

type AccountsSignupPage struct {
	BasePage
	SignupForm *SignupFormComponent
}

type AccountsForgotPasswordPage struct {
	BasePage
	ForgotPasswordForm *ForgotPasswordFormComponent
}

type AccountsResetPasswordPage struct {
	BasePage
	ResetPasswordForm *ResetPasswordFormComponent
}

type AdminCommentsPage struct {
	BasePage
	NumComments uint
	URLQuery    string
	Tutorials   []*models.TutorialModel
	Users       []*models.UserModel
	Comments    *AdminCommentsListComponent
}

type AdminCoursesPage struct {
	BasePage
	NumCourses    uint
	URLQuery      string
	PublishStatus []string
	Authors       []*models.UserModel
	Keywords      []string
	Courses       *AdminCoursesListComponent
}

type AdminDiscountsPage struct {
	BasePage
	NumDiscounts    uint
	URLQuery        string
	DiscountStatus  []string
	NewDiscountForm *NewDiscountFormComponent
	Discounts       *AdminDiscountsListComponent
}

type AdminPurchasesPage struct {
	BasePage
}

type AdminRefundsPage struct {
	BasePage
}

type AdminTutorialsPage struct {
	BasePage
	NumTutorials  uint
	URLQuery      string
	PublishStatus []string
	Authors       []*models.UserModel
	Keywords      []string
	Tutorials     *AdminTutorialsListComponent
}

type AdminUsersPage struct {
	BasePage
	NumUsers            uint
	URLQuery            string
	AuthorizationLevels []string
	Users               *AdminUsersListComponent
}

type AuthorsTutorialsPage struct {
	BasePage
}

type AuthorsCoursesPage struct {
	BasePage
}

type CoursesPage struct {
	BasePage
	Courses *CoursesListComponent
}

type CoursesCoursePage struct {
	BasePage
	Course   *models.CourseModel
	Author   *models.AuthorModel
	Chapters int
}

type CoursesPurchasesPage struct {
	BasePage
}

type Errors404Page struct {
	BasePage
}

type Errors500Page struct {
	BasePage
}

type GeneralHomePage struct {
	BasePage
}

type GeneralAffiliateProgramPage struct {
	BasePage
}

type GeneralPrivacyPolicyPage struct {
	BasePage
}

type GeneralRefundPolicyPage struct {
	BasePage
}

type ProfilePage struct {
	BasePage
}

type ProfileAffiliateHistoryPage struct {
	BasePage
}

type ProfileCourses struct {
	BasePage
}

type ProfileCourse struct {
	BasePage
}

type ProfileTutorialsBookmarksPage struct {
	BasePage
}

type ProfileTutorialsLikedPage struct {
	BasePage
}

type SettingsPage struct {
	BasePage
}

type TutorialsPage struct {
	BasePage
	Tutorials *TutorialsListComponent
}

type TutorialsTutorialPage struct {
	BasePage
	User               *models.UserModel
	Tutorial           *models.TutorialModel
	Keywords           []string
	Author             *models.AuthorModel
	Course             *models.CourseModel
	TutorialLiked      bool
	TutorialBookmarked bool
	Comments           *CommentsListComponent
}
