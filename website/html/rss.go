package html

import (
	"github.com/PsionicAlch/psionicalch-home/internal/database/models"
)

type GeneralRSS struct {
	Tutorials []*models.TutorialModel
	Courses   []*models.CourseModel
	Authors   map[string]*models.UserModel
}
