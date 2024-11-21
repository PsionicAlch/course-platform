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

type AdminAdminsPage struct {
	BasePage
}

type AdminAuthorsPage struct {
	BasePage
}

type AdminCommentsPage struct {
	BasePage
}

type AdminCoursesPage struct {
	BasePage
}

type AdminDiscountsPage struct {
	BasePage
}

type AdminPurchasesPage struct {
	BasePage
}

type AdminRefundsPage struct {
	BasePage
}

type AdminTutorialsPage struct {
	BasePage
}

type AdminUsersPage struct {
	BasePage
}

type AuthorsTutorialsPage struct {
	BasePage
}

type AuthorsCoursesPage struct {
	BasePage
}

type CoursesPage struct {
	BasePage
}

type CoursesCoursePage struct {
	BasePage
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
