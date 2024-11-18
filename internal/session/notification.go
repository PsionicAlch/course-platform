package session

import (
	"context"
)

const (
	InfoMessageKey    = "info-flash-messages"
	WarningMessageKey = "warning-flash-messages"
	ErrorMessageKey   = "error-flash-messages"
)

func (s *Session) SetMessage(ctx context.Context, key, msg string) {
	messages, ok := s.session.Get(ctx, key).([]string)
	if !ok {
		messages = []string{}
	}

	messages = append(messages, msg)

	s.session.Put(ctx, key, messages)
}

func (s *Session) GetMessages(ctx context.Context, key string) []string {
	messages, ok := s.session.Pop(ctx, key).([]string)
	if !ok {
		messages = []string{}
	}

	return messages
}

func (s *Session) SetInfoMessage(ctx context.Context, msg string) {
	s.SetMessage(ctx, InfoMessageKey, msg)
}

func (s *Session) GetInfoMessages(ctx context.Context) []string {
	return s.GetMessages(ctx, InfoMessageKey)
}

func (s *Session) SetWarningMessage(ctx context.Context, msg string) {
	s.SetMessage(ctx, WarningMessageKey, msg)
}

func (s *Session) GetWarningMessages(ctx context.Context) []string {
	return s.GetMessages(ctx, WarningMessageKey)
}

func (s *Session) SetErrorMessage(ctx context.Context, msg string) {
	s.SetMessage(ctx, ErrorMessageKey, msg)
}

func (s *Session) GetErrorMessages(ctx context.Context) []string {
	return s.GetMessages(ctx, ErrorMessageKey)
}
