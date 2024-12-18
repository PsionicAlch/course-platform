package rss

import (
	"net/http"
	"time"

	"github.com/PsionicAlch/psionicalch-home/internal/database"
	"github.com/PsionicAlch/psionicalch-home/internal/database/models"
	"github.com/PsionicAlch/psionicalch-home/internal/render"
	"github.com/PsionicAlch/psionicalch-home/internal/utils"
	"github.com/PsionicAlch/psionicalch-home/website/html"
	"github.com/PsionicAlch/psionicalch-home/website/pages"
)

type Handlers struct {
	utils.Loggers
	Renderers *pages.Renderers
	Database  database.Database
}

func SetupHandlers(rssRenderer render.Renderer, db database.Database) *Handlers {
	loggers := utils.CreateLoggers("RSS HANDLERS")

	return &Handlers{
		Loggers:   loggers,
		Renderers: pages.CreateRenderers(nil, nil, rssRenderer),
		Database:  db,
	}
}

func (h *Handlers) RSSGet(w http.ResponseWriter, r *http.Request) {
	rssData := html.GeneralRSS{}

	published := true

	tutorials, err := h.Database.GetAllTutorials(&published)
	if err != nil {
		h.ErrorLog.Printf("Failed to get all tutorials: %s\n", err)

		if err := h.Renderers.RSS.RenderXML(w, "errors-500", nil); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	rssData.Tutorials = tutorials

	courses, err := h.Database.GetAllCourses(&published)
	if err != nil {
		h.ErrorLog.Printf("Failed to get all courses: %s\n", err)

		if err := h.Renderers.RSS.RenderXML(w, "errors-500", nil); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	rssData.Courses = courses

	authors := make(map[string]*models.UserModel, len(tutorials)+len(courses))
	authorsCache := make(map[string]*models.UserModel)

	for _, tutorial := range tutorials {
		author, contains := authorsCache[tutorial.AuthorID.String]
		if !contains {
			h.InfoLog.Printf("Adding author to cache: %s\n", tutorial.AuthorID.String)

			author, err = h.Database.GetUserByID(tutorial.AuthorID.String, database.Author)
			if err != nil {
				h.ErrorLog.Printf("Failed to get author by ID (\"%s\"): %s\n", tutorial.AuthorID.String, err)

				if err := h.Renderers.RSS.RenderXML(w, "errors-500", nil); err != nil {
					h.ErrorLog.Println(err)
				}

				return
			}

			authorsCache[tutorial.AuthorID.String] = author
		}

		authors[tutorial.ID] = author
	}

	for _, course := range courses {
		author, contains := authorsCache[course.AuthorID.String]
		if !contains {
			h.InfoLog.Printf("Adding author to cache: %s\n", course.AuthorID.String)

			author, err = h.Database.GetUserByID(course.AuthorID.String, database.Author)
			if err != nil {
				h.ErrorLog.Printf("Failed to get author by ID (\"%s\"): %s\n", course.AuthorID.String, err)

				if err := h.Renderers.RSS.RenderXML(w, "errors-500", nil); err != nil {
					h.ErrorLog.Println(err)
				}

				return
			}

			authorsCache[course.AuthorID.String] = author
		}

		authors[course.ID] = author
	}

	rssData.Authors = authors

	rssData.LastBuildTime = time.Now()

	if err := h.Renderers.RSS.RenderXML(w, "general", rssData); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) RSSTutorialsGet(w http.ResponseWriter, r *http.Request) {
	rssData := html.TutorialsRSS{}

	published := true

	tutorials, err := h.Database.GetAllTutorials(&published)
	if err != nil {
		h.ErrorLog.Printf("Failed to get all tutorials: %s\n", err)

		if err := h.Renderers.RSS.RenderXML(w, "errors-500", nil); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	rssData.Tutorials = tutorials

	authors := make(map[string]*models.UserModel, len(tutorials))
	authorsCache := make(map[string]*models.UserModel)

	for _, tutorial := range tutorials {
		author, contains := authorsCache[tutorial.AuthorID.String]
		if !contains {
			h.InfoLog.Printf("Adding author to cache: %s\n", tutorial.AuthorID.String)

			author, err = h.Database.GetUserByID(tutorial.AuthorID.String, database.Author)
			if err != nil {
				h.ErrorLog.Printf("Failed to get author by ID (\"%s\"): %s\n", tutorial.AuthorID.String, err)

				if err := h.Renderers.RSS.RenderXML(w, "errors-500", nil); err != nil {
					h.ErrorLog.Println(err)
				}

				return
			}

			authorsCache[tutorial.AuthorID.String] = author
		}

		authors[tutorial.ID] = author
	}

	rssData.Authors = authors

	rssData.LastBuildTime = time.Now()

	if err := h.Renderers.RSS.RenderXML(w, "tutorials", rssData); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) RSSCoursesGet(w http.ResponseWriter, r *http.Request) {
	rssData := html.CoursesRSS{}

	published := true

	courses, err := h.Database.GetAllCourses(&published)
	if err != nil {
		h.ErrorLog.Printf("Failed to get all courses: %s\n", err)

		if err := h.Renderers.RSS.RenderXML(w, "errors-500", nil); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	rssData.Courses = courses

	authors := make(map[string]*models.UserModel, len(courses))
	authorsCache := make(map[string]*models.UserModel)

	for _, course := range courses {
		author, contains := authorsCache[course.AuthorID.String]
		if !contains {
			h.InfoLog.Printf("Adding author to cache: %s\n", course.AuthorID.String)

			author, err = h.Database.GetUserByID(course.AuthorID.String, database.Author)
			if err != nil {
				h.ErrorLog.Printf("Failed to get author by ID (\"%s\"): %s\n", course.AuthorID.String, err)

				if err := h.Renderers.RSS.RenderXML(w, "errors-500", nil); err != nil {
					h.ErrorLog.Println(err)
				}

				return
			}

			authorsCache[course.AuthorID.String] = author
		}

		authors[course.ID] = author
	}

	rssData.Authors = authors

	rssData.LastBuildTime = time.Now()

	if err := h.Renderers.RSS.RenderXML(w, "courses", rssData); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) RSSTutorialAuthorGet(w http.ResponseWriter, r *http.Request) {
	rssData := html.TutorialAuthorRSS{}

	if err := h.Renderers.RSS.RenderXML(w, "author-tutorials", rssData); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) RSSCourseAuthorGet(w http.ResponseWriter, r *http.Request) {
	rssData := html.CourseAuthorRSS{}

	if err := h.Renderers.RSS.RenderXML(w, "author-courses", rssData); err != nil {
		h.ErrorLog.Println(err)
	}
}
