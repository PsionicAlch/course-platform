package html

import (
	"time"

	"github.com/PsionicAlch/psionicalch-home/internal/database/models"
)

type GeneralRSS struct {
	LastBuildTime time.Time
	Tutorials     []*models.TutorialModel
	Courses       []*models.CourseModel
	Authors       map[string]*models.UserModel
}

type TutorialsRSS struct {
	LastBuildTime time.Time
	Tutorials     []*models.TutorialModel
	Authors       map[string]*models.UserModel
}

type CoursesRSS struct {
	LastBuildTime time.Time
	Courses       []*models.CourseModel
	Authors       map[string]*models.UserModel
}

type CourseAuthorRSS struct {
	LastBuildTime time.Time
	Courses       []*models.CourseModel
	Author        *models.UserModel
}

type TutorialAuthorRSS struct {
	LastBuildTime time.Time
	Tutorials     []*models.TutorialModel
	Author        *models.UserModel
}

type TutorialRSS struct {
	LastBuildTime time.Time
	Tutorial      *models.TutorialModel
	Author        *models.UserModel
}
