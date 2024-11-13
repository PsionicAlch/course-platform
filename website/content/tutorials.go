package content

import (
	"embed"
	"sync"

	"github.com/PsionicAlch/psionicalch-home/internal/database"
)

//go:embed tutorials/*.md
var tutorialsFS embed.FS

// TutorialMatter is a struct representation of the metadata found in each tutorial markdown file.
type TutorialMatter struct {
	Title        string   `yaml:"title"`
	Slug         string   `yaml:"slug"`
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
	// defer waitGroup.Done()

	// var tutorialsToAdd, tutorialsToUpdate []*models.TutorialModel

	// timerStart := time.Now()

	// // Get all tutorial file checksums from database.
	// tutorials, err := db.GetAllTutorials()
	// if err != nil {
	// 	content.ErrorLog.Fatalf("Failed to read all tutorials from the database: %s\n", err)
	// }

	// // Read all files from embedded file systems.
	// files, err := tutorialsFS.ReadDir("tutorials")
	// if err != nil {
	// 	content.ErrorLog.Printf("Failed to read tutorials embedded file system: %s\n", err)
	// }

	// content.InfoLog.Printf("Loading %d tutorials into the database!\n", len(files))

	// // Read and parse each file.
	// for _, file := range files {
	// 	// Just in case there are any folders in the folder skip over them.
	// 	if file.IsDir() {
	// 		content.WarningLog.Fatalf("Found a directory in tutorials/%s. Directories are not supported.\n", file.Name())
	// 		continue
	// 	}

	// 	// Read the content of the file into a byte slice.
	// 	output, err := tutorialsFS.ReadFile("tutorials/" + file.Name())
	// 	if err != nil {
	// 		content.ErrorLog.Fatalf("Failed to read \"%s\" from tutorials embedded file system: %s", file.Name(), err)
	// 	}

	// 	matter := new(TutorialMatter)
	// 	tutData := new(TutorialData)

	// 	// Parse the frontmatter to get the file "metadata".
	// 	data, err := frontmatter.Parse(strings.NewReader(string(output)), matter)
	// 	if err != nil {
	// 		content.ErrorLog.Fatalf("Failed to parse markdown from \"%s\": %s", "tutorials/"+file.Name(), err)
	// 	}

	// 	hasher := sha256.New()
	// 	hasher.Write(output)
	// 	fileChecksum := hasher.Sum(nil)

	// 	slugIndex, slugFound := utils.InSliceFunc(matter.Slug, tutorials, func(slug string, tutorial *models.TutorialModel) bool {
	// 		return tutorial.Slug == slug
	// 	})

	// 	_, checksumFound := utils.InSliceFunc(string(fileChecksum), tutorials, func(fileChecksum string, tutorial *models.TutorialModel) bool {
	// 		return tutorial.FileChecksum == fileChecksum
	// 	})

	// 	if slugFound && checksumFound {
	// 		continue
	// 	}

	// 	tutData.TutorialMatter = *matter
	// 	tutData.Content = string(MarkdownToHTML(data))

	// 	if !slugFound {
	// 		var keywordModels []*models.KeywordModel
	// 		for _, keyword := range tutData.Keywords {
	// 			keywordModel := new(models.KeywordModel)
	// 			keywordModel.Keyword = keyword

	// 			keywordModels = append(keywordModels, keywordModel)
	// 		}

	// 		tutorialToAdd := new(models.TutorialModel)

	// 		tutorialToAdd.Title = tutData.Title
	// 		tutorialToAdd.Slug = tutData.Slug
	// 		tutorialToAdd.Description = tutData.Description
	// 		tutorialToAdd.ThumbnailURL = tutData.ThumbnailURL
	// 		tutorialToAdd.BannerURL = tutData.BannerURL
	// 		tutorialToAdd.Content = tutData.Content
	// 		tutorialToAdd.FileChecksum = string(fileChecksum)
	// 		tutorialToAdd.Keywords = keywordModels

	// 		tutorialsToAdd = append(tutorialsToAdd, tutorialToAdd)

	// 		continue
	// 	}

	// 	if !checksumFound {
	// 		var keywordModels []*models.KeywordModel
	// 		for _, keyword := range tutData.Keywords {
	// 			keywordModel := new(models.KeywordModel)
	// 			keywordModel.Keyword = keyword

	// 			keywordModels = append(keywordModels, keywordModel)
	// 		}

	// 		tutorialToUpdate := new(models.TutorialModel)
	// 		tutorialToUpdate.ID = tutorials[slugIndex].ID
	// 		tutorialToUpdate.Title = tutData.Title
	// 		tutorialToUpdate.Slug = tutData.Slug
	// 		tutorialToUpdate.Description = tutData.Description
	// 		tutorialToUpdate.ThumbnailURL = tutData.ThumbnailURL
	// 		tutorialToUpdate.BannerURL = tutData.BannerURL
	// 		tutorialToUpdate.Content = tutData.Content
	// 		tutorialToUpdate.Published = false
	// 		tutorialToUpdate.AuthorID = tutorials[slugIndex].AuthorID
	// 		tutorialToUpdate.FileChecksum = string(fileChecksum)
	// 		tutorialToUpdate.Keywords = keywordModels

	// 		tutorialsToUpdate = append(tutorialsToUpdate, tutorialToUpdate)

	// 		continue
	// 	}
	// }

	// if len(tutorialsToAdd) > 0 {
	// 	if err := db.AddNewTutorialBulk(tutorialsToAdd); err != nil {
	// 		content.ErrorLog.Fatalf("Failed to bulk add new tutorials: %s", err)
	// 	}
	// }

	// if len(tutorialsToUpdate) > 0 {
	// 	if err := db.UpdateTutorialBulk(tutorialsToUpdate); err != nil {
	// 		content.ErrorLog.Fatalf("Failed to bulk update tutorials: %s", err)
	// 	}
	// }

	// timerEnd := time.Since(timerStart)

	// content.InfoLog.Printf("Parsed %d tutorials in %s\n", len(files), timerEnd)
}
