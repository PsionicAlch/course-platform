package content

import (
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"strings"
	"sync"
	"time"

	"github.com/PsionicAlch/psionicalch-home/internal/database"
	"github.com/adrg/frontmatter"
)

//go:embed courses/*
var coursesFS embed.FS

type CourseMatter struct {
	Title        string   `yaml:"title"`
	Description  string   `yaml:"description"`
	ThumbnailURL string   `yaml:"thumbnail_url`
	BannerURL    string   `yaml:"banner_url"`
	Keywords     []string `yaml:"keywords"`
	Directory    string   `yaml:"directory"`
	Key          string   `yaml:"key"`
}

type CourseData struct {
	CourseMatter
	Content string
}

func (content *Content) RegisterCourseContent(waitGroup *sync.WaitGroup, db database.Database) {
	defer waitGroup.Done()

	timerStart := time.Now()

	files, err := coursesFS.ReadDir("courses")
	if err != nil {
		content.ErrorLog.Fatalf("Failed to read courses embedded file system: %s\n", err)
	}

	content.InfoLog.Printf("Parsing %d courses!\n", len(files)/2)

	for _, file := range files {
		if file.IsDir() {
			content.InfoLog.Printf("Directory found: courses/%s\n", file.Name())
			continue
		}

		output, err := coursesFS.ReadFile("courses/" + file.Name())
		if err != nil {
			content.ErrorLog.Fatalf("Failed to read \"%s\" from courses embedded file system: %s\n", file.Name(), err)
		}

		matter := new(CourseMatter)
		courseData := new(CourseData)

		data, err := frontmatter.Parse(strings.NewReader(string(output)), matter)
		if err != nil {
			content.ErrorLog.Fatalf("Failed to parse markdown from \"%s\": %s", "courses/"+file.Name(), err)
		}

		// Create the file checksum to be able to see if the file data has changed at all.
		hasher := sha256.New()
		hasher.Write(output)
		fileChecksum := hex.EncodeToString(hasher.Sum(nil))

		courseData.CourseMatter = *matter
		courseData.Content = string(MarkdownToHTML(data))

		content.InfoLog.Printf("Finished parsing \"%s\"\n. Data: %#v\n. Checksum: %s\n", "courses/"+file.Name(), courseData, fileChecksum)
	}

	timerEnd := time.Since(timerStart)

	content.InfoLog.Printf("Parsed %d courses in %s\n", len(files)/2, timerEnd)
}
