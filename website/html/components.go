package html

type HeaderComponentData struct {
	Heading  string
	Text     string
	LinkHref string
	LinkText string
}

func CreateHeaderComponent(heading, text, linkhref, linktext string) *HeaderComponentData {
	return &HeaderComponentData{
		Heading:  heading,
		Text:     text,
		LinkHref: linkhref,
		LinkText: linktext,
	}
}
