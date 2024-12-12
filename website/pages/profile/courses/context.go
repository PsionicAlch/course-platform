package courses

import (
	"context"
	"net/http"

	"github.com/PsionicAlch/psionicalch-home/internal/database/models"
)

type ContextKey string

const CourseContextKey ContextKey = "course-context-key"

func NewContextWithCourseModel(course *models.CourseModel, ctx context.Context) context.Context {
	return context.WithValue(ctx, CourseContextKey, course)
}

func GetCourseFromRequest(r *http.Request) *models.CourseModel {
	course, ok := r.Context().Value(CourseContextKey).(*models.CourseModel)
	if !ok {
		return nil
	}

	return course
}
