package content

import (
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"strings"
	"time"

	"github.com/PsionicAlch/course-platform/internal/database"
	"github.com/PsionicAlch/course-platform/internal/database/models"
	"github.com/PsionicAlch/course-platform/internal/utils"
	"github.com/adrg/frontmatter"
)

//go:embed tutorials/*.md
var tutorialsFS embed.FS

// TutorialMatter is a struct representation of the metadata found in each tutorial markdown file.
type TutorialMatter struct {
	Title        string   `yaml:"title"`
	Description  string   `yaml:"description"`
	ThumbnailURL string   `yaml:"thumbnail_url"`
	BannerURL    string   `yaml:"banner_url"`
	Keywords     []string `yaml:"keywords"`
	Key          string   `yaml:"key"`
}

// TutorialData is a struct representation for the information contained in a tutorial.
type TutorialData struct {
	TutorialMatter
	Content string
}

func (content *Content) RegisterTutorialsContent(db database.Database) {
	db.PrepareBulkTutorials()

	timerStart := time.Now()

	tutorials, err := db.GetAllTutorials("", nil)
	if err != nil {
		content.ErrorLog.Fatalf("Failed to read all tutorials from the database: %s\n", err)
	}

	files, err := tutorialsFS.ReadDir("tutorials")
	if err != nil {
		content.ErrorLog.Printf("Failed to read tutorials embedded file system: %s\n", err)
	}

	content.InfoLog.Printf("Parsing %d tutorials!\n", len(files))

	for _, file := range files {
		// Skip over any directories.
		if file.IsDir() {
			content.WarningLog.Fatalf("Found a directory in tutorials/%s. Directories are not supported.\n", file.Name())
			continue
		}

		// Read the contents of the tutorial file.
		output, err := tutorialsFS.ReadFile("tutorials/" + file.Name())
		if err != nil {
			content.ErrorLog.Fatalf("Failed to read \"%s\" from tutorials embedded file system: %s", file.Name(), err)
		}

		matter := new(TutorialMatter)
		tutData := new(TutorialData)

		// Parse the tutorial file and separate out the frontmatter from the tutorial contents.
		data, err := frontmatter.Parse(strings.NewReader(string(output)), matter)
		if err != nil {
			content.ErrorLog.Fatalf("Failed to parse markdown from \"%s\": %s", "tutorials/"+file.Name(), err)
		}

		// Create the file checksum to be able to see if the file data has changed at all.
		hasher := sha256.New()
		hasher.Write(output)
		fileChecksum := hex.EncodeToString(hasher.Sum(nil))

		// Slugs need to be unique so we can use it as a way to find individual tutorials.
		fileKeyIndex, fileKeyFound := utils.InSliceFunc(matter.Key, tutorials, func(fileKey string, tutorial *models.TutorialModel) bool {
			return fileKey == tutorial.FileKey
		})

		checksumMatch := false
		if fileKeyFound {
			checksumMatch = tutorials[fileKeyIndex].FileChecksum == string(fileChecksum)
		}

		// Skip the tutorial if it's already in the database and hasn't changed yet.
		if fileKeyFound && checksumMatch {
			continue
		}

		tutData.TutorialMatter = *matter
		tutData.Content = string(MarkdownToHTML(data))

		// This tutorial is new.
		if !fileKeyFound {
			db.InsertTutorial(tutData.Title, TitleToSlug(tutData.Title), tutData.Description, tutData.ThumbnailURL, tutData.BannerURL, tutData.Content, fileChecksum, tutData.Key, tutData.Keywords)

			continue
		}

		// This tutorial has been updated.
		if !checksumMatch {
			db.UpdateTutorial(tutorials[fileKeyIndex].ID, tutData.Title, TitleToSlug(tutData.Title), tutData.Description, tutData.ThumbnailURL, tutData.BannerURL, tutData.Content, string(fileChecksum), tutData.Key, tutData.Keywords, tutorials[fileKeyIndex].AuthorID)

			continue
		}
	}

	if err := db.RunBulkTutorials(); err != nil {
		content.ErrorLog.Fatalln(err)
	}

	timerEnd := time.Since(timerStart)

	content.InfoLog.Printf("Parsed %d tutorials in %s\n", len(files), timerEnd)
}
