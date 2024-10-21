package tutorials

import (
	"embed"
	"sort"
	"strings"
	"time"

	"github.com/PsionicAlch/psionicalch-home/internal/utils"
	"github.com/adrg/frontmatter"
)

//go:embed content
var content embed.FS

// TutorialMatter is a struct representation of the metadata found in each tutorial markdown file.
type TutorialMatter struct {
	Title       string    `yaml:"title"`
	Slug        string    `yaml:"slug"`
	Description string    `yaml:"description"`
	Date        time.Time `yaml:"date"`
	Image       string    `yaml:"image"`
	Keywords    []string  `yaml:"keywords"`
	Published   bool      `yaml:"published"`
}

// TutorialData is a struct representation for the information contained in a tutorial.
type TutorialData struct {
	TutorialMatter
	Content string
}

// TutorialCache is a type alias to represent the processed tutorials.
type TutorialsCache map[string]*TutorialData

// Tutorials is a struct to hold the required data for associated functions.
type Tutorials struct {
	utils.Loggers
	cache TutorialsCache
}

func SetupTutorials() *Tutorials {
	// Construct loggers.
	loggers := utils.CreateLoggers("TUTORIALS")

	// Read all files from embedded file systems.
	files, err := content.ReadDir("content")
	if err != nil {
		loggers.ErrorLog.Fatal("Failed to read files from embedded content folder: ", err)
	}

	// Create a cache variable that is long enough to contain all the files.
	cache := make(map[string]*TutorialData, len(files))

	// Read and parse each file.
	for _, file := range files {
		// Just in case there are any folders in the folder skip over them.
		if file.IsDir() {
			loggers.WarningLog.Printf("Found a folder (%s) in the embedded content folder. Please remove it.", file.Name())
			continue
		}

		// Read the content of the file into a byte slice.
		output, err := content.ReadFile("content/" + file.Name())
		if err != nil {
			loggers.ErrorLog.Fatalf("Failed to read \"%s\" from embedded content folder: %s\n", file.Name(), err)
		}

		matter := new(TutorialMatter)

		// Parse the frontmatter to get the file "metadata".
		data, err := frontmatter.Parse(strings.NewReader(string(output)), matter)
		if err != nil {
			loggers.ErrorLog.Fatalf("Failed to parse markdown from \"%s\": %s\n", "content/"+file.Name(), err)
		}

		// If the tutorial has been published then parse the markdown into HTML and add it to the cache.
		if matter.Published {
			cache[matter.Slug] = &TutorialData{
				TutorialMatter: *matter,
				Content:        string(utils.MarkdownToHTML(data)),
			}
		}
	}

	loggers.InfoLog.Printf("Generated tutorials cache with %d tutorials.", len(cache))

	return &Tutorials{
		Loggers: loggers,
		cache:   cache,
	}
}

// GetAllTutorials returns a slice of *TutorialData from the internal cache.
func (t *Tutorials) GetAllTutorials() []*TutorialData {
	tuts := make([]*TutorialData, 0, len(t.cache))

	for _, tut := range t.cache {
		tuts = append(tuts, tut)
	}

	sort.SliceStable(tuts, func(i, j int) bool {
		return tuts[i].Date.After(tuts[j].Date)
	})

	return tuts
}
