package generators

import (
	"bytes"
	"time"

	"github.com/PsionicAlch/psionicalch-home/internal/database"
	"github.com/PsionicAlch/psionicalch-home/internal/database/models"
	"github.com/PsionicAlch/psionicalch-home/internal/render"
	"github.com/PsionicAlch/psionicalch-home/internal/utils"
	"github.com/PsionicAlch/psionicalch-home/website/html"
)

func RSSFeed(loggers utils.Loggers, db database.Database, renderer render.Renderer) func() (string, error) {
	return func() (string, error) {
		rssData := html.GeneralRSS{}

		published := true
		feed := new(bytes.Buffer)

		tutorials, err := db.GetAllTutorials("", &published)
		if err != nil {
			loggers.ErrorLog.Printf("Failed to get all tutorials: %s\n", err)

			if err := renderer.Render(feed, nil, "errors-500", nil); err != nil {
				loggers.ErrorLog.Println(err)
			}

			return feed.String(), err
		}

		rssData.Tutorials = tutorials

		courses, err := db.GetAllCourses("", &published)
		if err != nil {
			loggers.ErrorLog.Printf("Failed to get all courses: %s\n", err)

			if err := renderer.Render(feed, nil, "errors-500", nil); err != nil {
				loggers.ErrorLog.Println(err)
			}

			return feed.String(), err
		}

		rssData.Courses = courses

		authors := make(map[string]*models.UserModel, len(tutorials)+len(courses))
		authorsCache := make(map[string]*models.UserModel)

		for _, tutorial := range tutorials {
			author, contains := authorsCache[tutorial.AuthorID.String]
			if !contains {
				author, err = db.GetUserByID(tutorial.AuthorID.String, database.Author)
				if err != nil {
					loggers.ErrorLog.Printf("Failed to get author by ID (\"%s\"): %s\n", tutorial.AuthorID.String, err)

					if err := renderer.Render(feed, nil, "errors-500", nil); err != nil {
						loggers.ErrorLog.Println(err)
					}

					return feed.String(), err
				}

				authorsCache[tutorial.AuthorID.String] = author
			}

			authors[tutorial.ID] = author
		}

		for _, course := range courses {
			author, contains := authorsCache[course.AuthorID.String]
			if !contains {
				loggers.InfoLog.Printf("Adding author to cache: %s\n", course.AuthorID.String)

				author, err = db.GetUserByID(course.AuthorID.String, database.Author)
				if err != nil {
					loggers.ErrorLog.Printf("Failed to get author by ID (\"%s\"): %s\n", course.AuthorID.String, err)

					if err := renderer.Render(feed, nil, "errors-500", nil); err != nil {
						loggers.ErrorLog.Println(err)
					}

					return feed.String(), err
				}

				authorsCache[course.AuthorID.String] = author
			}

			authors[course.ID] = author
		}

		rssData.Authors = authors

		rssData.LastBuildTime = time.Now()

		if err := renderer.Render(feed, nil, "general", rssData); err != nil {
			loggers.ErrorLog.Println(err)
		}

		return feed.String(), nil
	}
}

func TutorialsRSSFeed(loggers utils.Loggers, db database.Database, renderer render.Renderer) func() (string, error) {
	return func() (string, error) {
		rssData := html.TutorialsRSS{}

		published := true

		feed := new(bytes.Buffer)

		tutorials, err := db.GetAllTutorials("", &published)
		if err != nil {
			loggers.ErrorLog.Printf("Failed to get all tutorials: %s\n", err)

			if err := renderer.Render(feed, nil, "errors-500", nil); err != nil {
				loggers.ErrorLog.Println(err)
			}

			return feed.String(), err
		}

		rssData.Tutorials = tutorials

		authors := make(map[string]*models.UserModel, len(tutorials))
		authorsCache := make(map[string]*models.UserModel)

		for _, tutorial := range tutorials {
			author, contains := authorsCache[tutorial.AuthorID.String]
			if !contains {
				author, err = db.GetUserByID(tutorial.AuthorID.String, database.Author)
				if err != nil {
					loggers.ErrorLog.Printf("Failed to get author by ID (\"%s\"): %s\n", tutorial.AuthorID.String, err)

					if err := renderer.Render(feed, nil, "errors-500", nil); err != nil {
						loggers.ErrorLog.Println(err)
					}

					return feed.String(), err
				}

				authorsCache[tutorial.AuthorID.String] = author
			}

			authors[tutorial.ID] = author
		}

		rssData.Authors = authors

		rssData.LastBuildTime = time.Now()

		if err := renderer.Render(feed, nil, "tutorials", rssData); err != nil {
			loggers.ErrorLog.Println(err)
		}

		return feed.String(), nil
	}
}

func CoursesRSSFeed(loggers utils.Loggers, db database.Database, renderer render.Renderer) func() (string, error) {
	return func() (string, error) {
		rssData := html.CoursesRSS{}

		published := true

		feed := new(bytes.Buffer)

		courses, err := db.GetAllCourses("", &published)
		if err != nil {
			loggers.ErrorLog.Printf("Failed to get all courses: %s\n", err)

			if err := renderer.Render(feed, nil, "errors-500", nil); err != nil {
				loggers.ErrorLog.Println(err)
			}

			return feed.String(), err
		}

		rssData.Courses = courses

		authors := make(map[string]*models.UserModel, len(courses))
		authorsCache := make(map[string]*models.UserModel)

		for _, course := range courses {
			author, contains := authorsCache[course.AuthorID.String]
			if !contains {
				author, err = db.GetUserByID(course.AuthorID.String, database.Author)
				if err != nil {
					loggers.ErrorLog.Printf("Failed to get author by ID (\"%s\"): %s\n", course.AuthorID.String, err)

					if err := renderer.Render(feed, nil, "errors-500", nil); err != nil {
						loggers.ErrorLog.Println(err)
					}

					return feed.String(), err
				}

				authorsCache[course.AuthorID.String] = author
			}

			authors[course.ID] = author
		}

		rssData.Authors = authors

		rssData.LastBuildTime = time.Now()

		if err := renderer.Render(feed, nil, "courses", rssData); err != nil {
			loggers.ErrorLog.Println(err)
		}

		return feed.String(), nil
	}
}

func AuthorTutorialsRSSFeed(loggers utils.Loggers, db database.Database, renderer render.Renderer) func(authorSlug string) (string, error) {
	return func(authorSlug string) (string, error) {
		rssData := html.TutorialAuthorRSS{}

		published := true

		feed := new(bytes.Buffer)

		author, err := db.GetUserBySlug(authorSlug, database.Author)
		if err != nil {
			loggers.ErrorLog.Printf("Failed to get author by slug (\"%s\"): %s\n", authorSlug, err)

			if err := renderer.Render(feed, nil, "errors-500", nil); err != nil {
				loggers.ErrorLog.Println(err)
			}

			return feed.String(), err
		}

		if author == nil {
			if err := renderer.Render(feed, nil, "errors-404", nil); err != nil {
				loggers.ErrorLog.Println(err)
			}

			return feed.String(), err
		}

		rssData.Author = author

		tutorials, err := db.GetAllTutorials(author.ID, &published)
		if err != nil {
			loggers.ErrorLog.Printf("Failed to get all tutorials: %s\n", err)

			if err := renderer.Render(feed, nil, "errors-500", nil); err != nil {
				loggers.ErrorLog.Println(err)
			}

			return feed.String(), err
		}

		if len(tutorials) == 0 {
			if err := renderer.Render(feed, nil, "errors-404", nil); err != nil {
				loggers.ErrorLog.Println(err)
			}

			return feed.String(), err
		}

		rssData.Tutorials = tutorials
		rssData.LastBuildTime = time.Now()

		if err := renderer.Render(feed, nil, "author-tutorials", rssData); err != nil {
			loggers.ErrorLog.Println(err)
		}

		return feed.String(), err
	}
}

func AuthorCoursesRSSFeed(loggers utils.Loggers, db database.Database, renderer render.Renderer) func(authorSlug string) (string, error) {
	return func(authorSlug string) (string, error) {
		rssData := html.CourseAuthorRSS{}

		published := true

		feed := new(bytes.Buffer)

		author, err := db.GetUserBySlug(authorSlug, database.Author)
		if err != nil {
			loggers.ErrorLog.Printf("Failed to get author by slug (\"%s\"): %s\n", authorSlug, err)

			if err := renderer.Render(feed, nil, "errors-500", nil); err != nil {
				loggers.ErrorLog.Println(err)
			}

			return feed.String(), err
		}

		if author == nil {
			if err := renderer.Render(feed, nil, "errors-404", nil); err != nil {
				loggers.ErrorLog.Println(err)
			}

			return feed.String(), err
		}

		rssData.Author = author

		courses, err := db.GetAllCourses(author.ID, &published)
		if err != nil {
			loggers.ErrorLog.Printf("Failed to get all courses: %s\n", err)

			if err := renderer.Render(feed, nil, "errors-500", nil); err != nil {
				loggers.ErrorLog.Println(err)
			}

			return feed.String(), err
		}

		if len(courses) == 0 {
			if err := renderer.Render(feed, nil, "errors-404", nil); err != nil {
				loggers.ErrorLog.Println(err)
			}

			return feed.String(), err
		}

		rssData.Courses = courses

		if err := renderer.Render(feed, nil, "author-courses", rssData); err != nil {
			loggers.ErrorLog.Println(err)
		}

		return feed.String(), err
	}
}
