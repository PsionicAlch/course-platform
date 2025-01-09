package generators

import (
	"bytes"
	"time"

	"github.com/PsionicAlch/psionicalch-home/internal/database"
	"github.com/PsionicAlch/psionicalch-home/internal/database/models"
	"github.com/PsionicAlch/psionicalch-home/internal/render"
	"github.com/PsionicAlch/psionicalch-home/internal/utils"
	"github.com/PsionicAlch/psionicalch-home/web/html"
)

func RSSFeed(loggers utils.Loggers, db database.Database, renderer render.Renderer) func() (string, error) {
	return func() (string, error) {
		rssData := &html.GeneralRSS{
			LastBuildTime: time.Now(),
		}

		published := true
		feed := new(bytes.Buffer)

		tutorials, err := db.GetAllTutorials("", &published)
		if err != nil {
			loggers.ErrorLog.Printf("Failed to get all tutorials: %s\n", err)

			if err := renderer.Render(feed, nil, "errors-500-rss", nil); err != nil {
				loggers.ErrorLog.Println(err)
			}

			return feed.String(), err
		}

		rssData.Tutorials = tutorials

		courses, err := db.GetAllCourses("", &published)
		if err != nil {
			loggers.ErrorLog.Printf("Failed to get all courses: %s\n", err)

			if err := renderer.Render(feed, nil, "errors-500-rss", nil); err != nil {
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

					if err := renderer.Render(feed, nil, "errors-500-rss", nil); err != nil {
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

					if err := renderer.Render(feed, nil, "errors-500-rss", nil); err != nil {
						loggers.ErrorLog.Println(err)
					}

					return feed.String(), err
				}

				authorsCache[course.AuthorID.String] = author
			}

			authors[course.ID] = author
		}

		rssData.Authors = authors

		if err := renderer.Render(feed, nil, "general-rss", rssData); err != nil {
			loggers.ErrorLog.Println(err)
		}

		return feed.String(), nil
	}
}

func TutorialsRSSFeed(loggers utils.Loggers, db database.Database, renderer render.Renderer) func() (string, error) {
	return func() (string, error) {
		rssData := &html.TutorialsRSS{
			LastBuildTime: time.Now(),
		}

		published := true

		feed := new(bytes.Buffer)

		tutorials, err := db.GetAllTutorials("", &published)
		if err != nil {
			loggers.ErrorLog.Printf("Failed to get all tutorials: %s\n", err)

			if err := renderer.Render(feed, nil, "errors-500-rss", nil); err != nil {
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

					if err := renderer.Render(feed, nil, "errors-500-rss", nil); err != nil {
						loggers.ErrorLog.Println(err)
					}

					return feed.String(), err
				}

				authorsCache[tutorial.AuthorID.String] = author
			}

			authors[tutorial.ID] = author
		}

		rssData.Authors = authors

		if err := renderer.Render(feed, nil, "tutorials-rss", rssData); err != nil {
			loggers.ErrorLog.Println(err)
		}

		return feed.String(), nil
	}
}

func CoursesRSSFeed(loggers utils.Loggers, db database.Database, renderer render.Renderer) func() (string, error) {
	return func() (string, error) {
		rssData := &html.CoursesRSS{
			LastBuildTime: time.Now(),
		}

		published := true

		feed := new(bytes.Buffer)

		courses, err := db.GetAllCourses("", &published)
		if err != nil {
			loggers.ErrorLog.Printf("Failed to get all courses: %s\n", err)

			if err := renderer.Render(feed, nil, "errors-500-rss", nil); err != nil {
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

					if err := renderer.Render(feed, nil, "errors-500-rss", nil); err != nil {
						loggers.ErrorLog.Println(err)
					}

					return feed.String(), err
				}

				authorsCache[course.AuthorID.String] = author
			}

			authors[course.ID] = author
		}

		rssData.Authors = authors

		if err := renderer.Render(feed, nil, "courses-rss", rssData); err != nil {
			loggers.ErrorLog.Println(err)
		}

		return feed.String(), nil
	}
}

func AuthorTutorialsRSSFeed(loggers utils.Loggers, db database.Database, renderer render.Renderer) func(authorSlug string) (string, error) {
	return func(authorSlug string) (string, error) {
		rssData := &html.TutorialAuthorRSS{
			LastBuildTime: time.Now(),
		}

		published := true

		feed := new(bytes.Buffer)

		author, err := db.GetUserBySlug(authorSlug, database.Author)
		if err != nil {
			loggers.ErrorLog.Printf("Failed to get author by slug (\"%s\"): %s\n", authorSlug, err)

			if err := renderer.Render(feed, nil, "errors-500-rss", nil); err != nil {
				loggers.ErrorLog.Println(err)
			}

			return feed.String(), err
		}

		if author == nil {
			if err := renderer.Render(feed, nil, "errors-404-rss", nil); err != nil {
				loggers.ErrorLog.Println(err)
			}

			return feed.String(), err
		}

		rssData.Author = author

		tutorials, err := db.GetAllTutorials(author.ID, &published)
		if err != nil {
			loggers.ErrorLog.Printf("Failed to get all tutorials: %s\n", err)

			if err := renderer.Render(feed, nil, "errors-500-rss", nil); err != nil {
				loggers.ErrorLog.Println(err)
			}

			return feed.String(), err
		}

		if len(tutorials) == 0 {
			if err := renderer.Render(feed, nil, "errors-404-rss", nil); err != nil {
				loggers.ErrorLog.Println(err)
			}

			return feed.String(), err
		}

		rssData.Tutorials = tutorials

		if err := renderer.Render(feed, nil, "author-tutorials-rss", rssData); err != nil {
			loggers.ErrorLog.Println(err)
		}

		return feed.String(), err
	}
}

func AuthorCoursesRSSFeed(loggers utils.Loggers, db database.Database, renderer render.Renderer) func(authorSlug string) (string, error) {
	return func(authorSlug string) (string, error) {
		rssData := &html.CourseAuthorRSS{
			LastBuildTime: time.Now(),
		}

		published := true

		feed := new(bytes.Buffer)

		author, err := db.GetUserBySlug(authorSlug, database.Author)
		if err != nil {
			loggers.ErrorLog.Printf("Failed to get author by slug (\"%s\"): %s\n", authorSlug, err)

			if err := renderer.Render(feed, nil, "errors-500-rss", nil); err != nil {
				loggers.ErrorLog.Println(err)
			}

			return feed.String(), err
		}

		if author == nil {
			if err := renderer.Render(feed, nil, "errors-404-rss", nil); err != nil {
				loggers.ErrorLog.Println(err)
			}

			return feed.String(), err
		}

		rssData.Author = author

		courses, err := db.GetAllCourses(author.ID, &published)
		if err != nil {
			loggers.ErrorLog.Printf("Failed to get all courses: %s\n", err)

			if err := renderer.Render(feed, nil, "errors-500-rss", nil); err != nil {
				loggers.ErrorLog.Println(err)
			}

			return feed.String(), err
		}

		if len(courses) == 0 {
			if err := renderer.Render(feed, nil, "errors-404-rss", nil); err != nil {
				loggers.ErrorLog.Println(err)
			}

			return feed.String(), err
		}

		rssData.Courses = courses

		if err := renderer.Render(feed, nil, "author-courses-rss", rssData); err != nil {
			loggers.ErrorLog.Println(err)
		}

		return feed.String(), err
	}
}

func TutorialRSSFeed(loggers utils.Loggers, db database.Database, renderer render.Renderer) func(tutorialSlug string) (string, error) {
	return func(tutorialSlug string) (string, error) {
		rssData := &html.TutorialRSS{
			LastBuildTime: time.Now(),
		}

		feed := new(bytes.Buffer)

		tutorial, err := db.GetTutorialBySlug(tutorialSlug)
		if err != nil {
			loggers.ErrorLog.Printf("Failed to get tutorial by slug (\"%s\"): %s\n", tutorialSlug, err)

			if err := renderer.Render(feed, nil, "errors-500-rss", nil); err != nil {
				loggers.ErrorLog.Println(err)
			}

			return feed.String(), err
		}

		if !tutorial.Published || !tutorial.AuthorID.Valid {
			if err := renderer.Render(feed, nil, "errors-404-rss", nil); err != nil {
				loggers.ErrorLog.Println(err)
			}

			return feed.String(), err
		}

		rssData.Tutorial = tutorial

		author, err := db.GetUserByID(tutorial.AuthorID.String, database.Author)
		if err != nil {
			loggers.ErrorLog.Printf("Failed to get author by ID (\"%s\"): %s\n", tutorial.AuthorID.String, err)

			if err := renderer.Render(feed, nil, "errors-500-rss", nil); err != nil {
				loggers.ErrorLog.Println(err)
			}

			return feed.String(), err
		}

		rssData.Author = author

		if err := renderer.Render(feed, nil, "tutorial-rss", rssData); err != nil {
			loggers.ErrorLog.Println(err)
		}

		return feed.String(), err
	}
}
