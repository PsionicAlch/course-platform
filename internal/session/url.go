package session

import "context"

// TODO: Write documentation

const RedirectURLKey = "redirect-url-key"

func (s *Session) SetRedirectURL(ctx context.Context, url string) {
	s.session.Put(ctx, RedirectURLKey, url)
}

func (s *Session) GetRedirectURL(ctx context.Context) string {
	return s.session.PopString(ctx, RedirectURLKey)
}
