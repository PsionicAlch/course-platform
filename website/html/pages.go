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
	Author       *models.UserModel
	LenTutorials uint
	Tutorials    *TutorialsListComponent
}

type AuthorsCoursesPage struct {
	BasePage
}

type CertificatePage struct {
	Certificate *models.CertificateModel
	User        *models.UserModel
	Course      *models.CourseModel
}

type CoursesPage struct {
	BasePage
	Courses *CoursesListComponent
}

type CoursesCoursePage struct {
	BasePage
	Course   *models.CourseModel
	Author   *models.UserModel
	Chapters int
}

type CoursesPurchasesPage struct {
	BasePage
	Course             *models.CourseModel
	Author             *models.UserModel
	CoursePurchaseForm *CoursePurchaseFormComponent
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
	User                       *models.UserModel
	NumTutorialsBookmarked     uint
	HasAffiliateHistory        bool
	Courses                    []*models.CourseModel
	HasMoreCourses             bool
	TutorialsBookmarked        []*models.TutorialModel
	HasMoreTutorialsBookmarked bool
	TutorialsLiked             []*models.TutorialModel
	HasMoreTutorialsLiked      bool
}

type ProfileAffiliateHistoryPage struct {
	BasePage
	User             *models.UserModel
	AffiliateHistory *AffiliateHistoryListComponent
}

type ProfileCertificate struct {
	BasePage
	Certificate *models.CertificateModel
	Course      *models.CourseModel
	User        *models.UserModel
	Author      *models.UserModel
	Chapters    []*models.ChapterModel
}

type ProfileCourses struct {
	BasePage
	Courses *CoursesListComponent
}

type ProfileCourse struct {
	BasePage
	Course             *models.CourseModel
	Chapter            *models.ChapterModel
	Chapters           []*models.ChapterModel
	LastChapter        bool
	Completed          map[string]bool
	HasCompletedCourse bool
}

type ProfileTutorialsBookmarksPage struct {
	BasePage
	Tutorials *TutorialsListComponent
}

type ProfileTutorialsLikedPage struct {
	BasePage
	Tutorials *TutorialsListComponent
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
	Author             *models.UserModel
	Course             *models.CourseModel
	TutorialLiked      bool
	TutorialBookmarked bool
	Comments           *CommentsListComponent
}

type LoadingScreenPage struct {
	Title   string
	PingURL string
}
