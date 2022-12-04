package goson

import "fmt"

type invalidPatternError struct {
	pattern string
	err     error
}

func (invalidPatternError *invalidPatternError) Error() string {
	return fmt.Sprintf("invalid pattern: %s, reason: %v", invalidPatternError.pattern, invalidPatternError.err)
}

func throwInvalidPatternError(pattern string, err error) *invalidPatternError {
	return &invalidPatternError{
		pattern: pattern,
		err:     err,
	}
}
