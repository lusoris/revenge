package matcher

import "errors"

var (
	// ErrNoMatches indicates no matching content was found
	ErrNoMatches = errors.New("no matching content found")

	// ErrNoTitle indicates the scan result has no parsed title
	ErrNoTitle = errors.New("no title to search")

	// ErrSearchFailed indicates the external search failed
	ErrSearchFailed = errors.New("external search failed")

	// ErrCreateFailed indicates content creation failed
	ErrCreateFailed = errors.New("content creation failed")

	// ErrLowConfidence indicates the match confidence is too low
	ErrLowConfidence = errors.New("match confidence too low")
)
