package validation

import (
	"fmt"
	"net/url"
	"strings"
	"unicode/utf8"
)

const (
	MaxTweetLength       = 280
	MaxTweetMedia        = 4
	MaxNameLength        = 50
	MaxDescriptionLength = 160
	MaxWebsiteLength     = 100
)

// Error is a user-facing validation failure. Handlers turn it into a 400.
type Error struct {
	Field   string
	Message string
}

func (e Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

func newError(field, format string, args ...any) error {
	return Error{Field: field, Message: fmt.Sprintf(format, args...)}
}

// TweetContent trims the text and checks it against the length limits. Empty
// text is allowed only when the tweet carries at least one attachment.
func TweetContent(content string, mediaCount int) (string, error) {
	trimmed := strings.TrimSpace(content)

	if trimmed == "" && mediaCount == 0 {
		return "", newError("tweet_data", "must not be empty")
	}

	if utf8.RuneCountInString(trimmed) > MaxTweetLength {
		return "", newError("tweet_data", "must be at most %d characters", MaxTweetLength)
	}

	if mediaCount > MaxTweetMedia {
		return "", newError("tweet_media_ids", "must contain at most %d items", MaxTweetMedia)
	}

	return trimmed, nil
}

// Name checks the display name of a profile.
func Name(name string) (string, error) {
	trimmed := strings.TrimSpace(name)

	if trimmed == "" {
		return "", newError("name", "must not be empty")
	}

	if utf8.RuneCountInString(trimmed) > MaxNameLength {
		return "", newError("name", "must be at most %d characters", MaxNameLength)
	}

	return trimmed, nil
}

// Description checks the profile bio. An empty bio is allowed.
func Description(description string) (string, error) {
	trimmed := strings.TrimSpace(description)

	if utf8.RuneCountInString(trimmed) > MaxDescriptionLength {
		return "", newError("description", "must be at most %d characters", MaxDescriptionLength)
	}

	return trimmed, nil
}

// Website checks the profile link. An empty link is allowed; a non-empty one
// must be an absolute http(s) URL.
func Website(website string) (string, error) {
	trimmed := strings.TrimSpace(website)

	if trimmed == "" {
		return "", nil
	}

	if utf8.RuneCountInString(trimmed) > MaxWebsiteLength {
		return "", newError("website", "must be at most %d characters", MaxWebsiteLength)
	}

	parsed, err := url.Parse(trimmed)
	if err != nil {
		return "", newError("website", "must be a valid URL")
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return "", newError("website", "must start with http:// or https://")
	}

	if parsed.Host == "" {
		return "", newError("website", "must contain a host")
	}

	return trimmed, nil
}

