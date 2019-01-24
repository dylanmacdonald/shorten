package api

// ShortenRequest represents a request to shorten a URL
type ShortenRequest struct {
	URL string `query:"url"`
}
