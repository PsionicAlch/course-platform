package gatekeeper

import "time"

func validateAuthenticationToken(authToken string, token *Token) bool {
	// Check to make sure that a token was passed back.
	if token == nil {
		return false
	}

	// Make sure the authentication token isn't empty and that it's the same one we gave.
	if token.Token == "" || token.Token != authToken {
		return false
	}

	// Make sure that the authentication token has the correct type.
	if token.TokenType != authenticationTokenType {
		return false
	}

	// Make sure the token is hasn't expired yet.
	if time.Now().After(token.ValidUntil) {
		return false
	}

	return true
}
