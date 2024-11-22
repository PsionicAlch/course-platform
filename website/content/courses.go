package content

import (
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"strings"
	"sync"
	"time"

	"github.com/PsionicAlch/psionicalch-home/internal/database"
	"github.com/PsionicAlch/psionicalch-home/internal/database/models"
	"github.com/PsionicAlch/psionicalch-home/internal/utils"
	"github.com/adrg/frontmatter"
)

//go:embed courses/**/*.md courses/*.md
var coursesFS embed.FS

type CourseMatter struct {
	Title        string   `yaml:"title"`
	Description  string   `yaml:"description"`
	ThumbnailURL string   `yaml:"thumbnail_url"`
	BannerURL    string   `yaml:"banner_url"`
	Keywords     []string `yaml:"keywords"`
	Directory    string   `yaml:"directory"`
	Key          string   `yaml:"key"`
}

type CourseData struct {
	CourseMatter
	Content string
}

type ChapterMatter struct {
	Title     string `yaml:"title"`
	Chapter   int    `yaml:"chapter"`
	CourseKey string `yaml:"course_key"`
	Key       string `yaml:"key"`
}

type ChapterData struct {
	ChapterMatter
	Content string
}

func (content *Content) RegisterCourseContent(waitGroup *sync.WaitGroup, db database.Database) {
	defer waitGroup.Done()

	timerStart := time.Now()

	courses, err := db.GetAllCourses()
	if err != nil {
		content.ErrorLog.Fatalf("Failed to get all courses: %s\n", err)
	}

	chapters, err := db.GetAllChapters()
	if err != nil {
		content.ErrorLog.Fatalf("Failed to get all chapters: %s\n", err)
	}

	files, err := coursesFS.ReadDir("courses")
	if err != nil {
		content.ErrorLog.Fatalf("Failed to read courses embedded file system: %s\n", err)
	}

	db.PrepareBulkCourses()

	content.InfoLog.Printf("Parsing %d courses!\n", len(files)/2)

	for _, file := range files {
		if file.IsDir() {
			chapterFiles, err := coursesFS.ReadDir("courses/" + file.Name())
			if err != nil {
				content.ErrorLog.Fatalf("Failed to read chapter (\"%s\") files in courses embedded file system: %s\n", "courses/"+file.Name(), err)
			}

			for _, chapterFile := range chapterFiles {
				filePath := "courses/" + file.Name() + "/" + chapterFile.Name()
				content.ParseChapterFile(filePath, db, chapters)
			}
		} else {
			filePath := "courses/" + file.Name()
			content.ParseCourseFile(filePath, db, courses)
		}
	}

	if err := db.RunBulkCourses(); err != nil {
		content.ErrorLog.Fatalf("Failed to bulk parse courses: %s\n", err)
	}

	timerEnd := time.Since(timerStart)

	content.InfoLog.Printf("Parsed %d courses in %s\n", len(files)/2, timerEnd)
}

func (content *Content) ParseChapterFile(filePath string, db database.Database, chapters []*models.ChapterModel) {
	output, err := coursesFS.ReadFile(filePath)
	if err != nil {
		content.ErrorLog.Printf("Failed to read chapter file (\"%s\") in courses embedded file system: %s\n", filePath, err)
	}

	chapterMatter := new(ChapterMatter)
	chapterData := new(ChapterData)

	data, err := frontmatter.Parse(strings.NewReader(string(output)), chapterMatter)
	if err != nil {
		content.ErrorLog.Fatalf("Failed to parse frontmatter from \"%s\": %s\n", filePath, err)
	}

	// Create the file checksum to be able to see if the file data has changed at all.
	hasher := sha256.New()
	hasher.Write(output)
	fileChecksum := hex.EncodeToString(hasher.Sum(nil))

	chapterData.ChapterMatter = *chapterMatter
	chapterData.Content = string(MarkdownToHTML(data))

	fileKeyIndex, fileKeyFound := utils.InSliceFunc(chapterMatter.Key, chapters, func(fileKey string, chapter *models.ChapterModel) bool {
		return fileKey == chapter.FileKey
	})

	checksumMatch := false
	if fileKeyFound {
		checksumMatch = chapters[fileKeyIndex].FileChecksum == fileChecksum
	}

	// The chapter already exists and hasn't been updated.
	if fileKeyFound && checksumMatch {
		return
	}

	// The chapter does not yet exist.
	if !fileKeyFound {
		db.InsertChapter(chapterData.Title, chapterData.Chapter, chapterData.Content, fileChecksum, chapterData.Key, chapterData.CourseKey)
		return
	}

	// The chapter has been updated.
	if !checksumMatch {
		db.UpdateChapter(chapters[fileKeyIndex].ID, chapterData.Title, chapterData.Chapter, chapterData.Content, fileChecksum, chapterData.Key, chapterData.CourseKey)
		return
	}
}

func (content *Content) ParseCourseFile(filePath string, db database.Database, courses []*models.CourseModel) {
	output, err := coursesFS.ReadFile(filePath)
	if err != nil {
		content.ErrorLog.Printf("Failed to read course file (\"%s\") in courses embedded file system: %s\n", filePath, err)
	}

	courseMatter := new(CourseMatter)
	courseData := new(CourseData)

	data, err := frontmatter.Parse(strings.NewReader(string(output)), courseMatter)
	if err != nil {
		content.ErrorLog.Fatalf("Failed to parse frontmatter from \"%s\": %s\n", filePath, err)
	}

	// Create the file checksum to be able to see if the file data has changed at all.
	hasher := sha256.New()
	hasher.Write(output)
	fileChecksum := hex.EncodeToString(hasher.Sum(nil))

	courseData.CourseMatter = *courseMatter
	courseData.Content = string(MarkdownToHTML(data))

	fileKeyIndex, fileKeyFound := utils.InSliceFunc(courseMatter.Key, courses, func(fileKey string, course *models.CourseModel) bool {
		return fileKey == course.FileKey
	})

	checksumMatch := false
	if fileKeyFound {
		checksumMatch = courses[fileKeyIndex].FileChecksum == fileChecksum
	}

	// The chapter already exists and hasn't been updated.
	if fileKeyFound && checksumMatch {
		return
	}

	// The chapter does not yet exist.
	if !fileKeyFound {
		db.InsertCourse(courseData.Title, TitleToSlug(courseData.Title), courseData.Description, courseData.ThumbnailURL, courseData.BannerURL, courseData.Content, fileChecksum, courseData.Key, courseData.Keywords)
		return
	}

	// The chapter has been updated.
	if !checksumMatch {
		db.UpdateCourse(courses[fileKeyIndex].ID, courseData.Title, TitleToSlug(courseData.Title), courseData.Description, courseData.ThumbnailURL, courseData.BannerURL, courseData.Content, fileChecksum, courseData.Key, courseData.Keywords, courses[fileKeyIndex].AuthorID)
		return
	}
}
