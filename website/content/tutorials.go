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

//go:embed tutorials/*.md
var tutorialsFS embed.FS

// TutorialMatter is a struct representation of the metadata found in each tutorial markdown file.
type TutorialMatter struct {
	Title        string   `yaml:"title"`
	Description  string   `yaml:"description"`
	ThumbnailURL string   `yaml:"thumbnail_url"`
	BannerURL    string   `yaml:"banner_url"`
	Keywords     []string `yaml:"keywords"`
}

// TutorialData is a struct representation for the information contained in a tutorial.
type TutorialData struct {
	TutorialMatter
	Content string
}

func (content *Content) RegisterTutorialsContent(waitGroup *sync.WaitGroup, db database.Database) {
	defer waitGroup.Done()

	var parsedTutorials []*models.TutorialModel

	timerStart := time.Now()

	tutorials, err := db.GetAllTutorials()
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
		slugIndex, slugFound := utils.InSliceFunc(TitleToSlug(matter.Title), tutorials, func(slug string, tutorial *models.TutorialModel) bool {
			return slug == tutorial.Slug
		})

		checksumMatch := false
		if slugFound {
			checksumMatch = tutorials[slugIndex].FileChecksum == string(fileChecksum)
		}

		// Skip the tutorial if it's already in the database and hasn't changed yet.
		if slugFound && checksumMatch {
			continue
		}

		tutData.TutorialMatter = *matter
		tutData.Content = string(MarkdownToHTML(data))

		// This tutorial is new.
		if !slugFound {
			var keywordModels []*models.KeywordModel
			for _, keyword := range tutData.Keywords {
				keywordModel := new(models.KeywordModel)
				keywordModel.Keyword = keyword

				keywordModels = append(keywordModels, keywordModel)
			}

			tutorialToAdd := new(models.TutorialModel)

			tutorialToAdd.Title = tutData.Title
			tutorialToAdd.Slug = TitleToSlug(tutData.Title)
			tutorialToAdd.Description = tutData.Description
			tutorialToAdd.ThumbnailURL = tutData.ThumbnailURL
			tutorialToAdd.BannerURL = tutData.BannerURL
			tutorialToAdd.Content = tutData.Content
			tutorialToAdd.FileChecksum = string(fileChecksum)
			tutorialToAdd.Keywords = keywordModels

			parsedTutorials = append(parsedTutorials, tutorialToAdd)

			continue
		}

		// This tutorial has been updated.
		if !checksumMatch {
			var keywordModels []*models.KeywordModel
			for _, keyword := range tutData.Keywords {
				keywordModel := new(models.KeywordModel)
				keywordModel.Keyword = keyword

				keywordModels = append(keywordModels, keywordModel)
			}

			tutorialToUpdate := new(models.TutorialModel)
			tutorialToUpdate.ID = tutorials[slugIndex].ID
			tutorialToUpdate.Title = tutData.Title
			tutorialToUpdate.Slug = TitleToSlug(tutData.Title)
			tutorialToUpdate.Description = tutData.Description
			tutorialToUpdate.ThumbnailURL = tutData.ThumbnailURL
			tutorialToUpdate.BannerURL = tutData.BannerURL
			tutorialToUpdate.Content = tutData.Content
			tutorialToUpdate.Published = false
			tutorialToUpdate.AuthorID = tutorials[slugIndex].AuthorID
			tutorialToUpdate.FileChecksum = string(fileChecksum)
			tutorialToUpdate.Keywords = keywordModels

			parsedTutorials = append(parsedTutorials, tutorialToUpdate)

			continue
		}
	}

	if len(parsedTutorials) > 0 {
		if err := db.BulkAddTutorials(parsedTutorials); err != nil {
			content.ErrorLog.Printf("Failed to bulk insert tutorials to database: %s\n", err)
		}
	}

	timerEnd := time.Since(timerStart)

	content.InfoLog.Printf("Parsed %d tutorials in %s\n", len(files), timerEnd)
}
