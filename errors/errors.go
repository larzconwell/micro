package errors

import (
	"fmt"
	"strings"

	"github.com/larzconwell/micro/token"
)

// StepError indicates that an entire step has failed after a
// certain number of errors have occurred.
type StepError string

func (stepError StepError) Error() string {
	return fmt.Sprintf("error: Encountered errors during %s", string(stepError))
}

// ScanError is a scanner error that has occurred when trying to retrieve valid tokens.
type ScanError rune

func (scanError ScanError) Error() string {
	return fmt.Sprintf("scanner: Invalid token found '%c'", scanError)
}

// ParseError is a parser error that has occurred when encountering an unexpected token.
type ParseError struct {
	Expected []token.Token
	Actual   *token.Value
}

func (parseError *ParseError) Error() string {
	var expected []string
	for _, token := range parseError.Expected {
		expected = append(expected, token.String())
	}

	return fmt.Sprintf("parser: Expected '%s' but found '%s'", strings.Join(expected, ", "), parseError.Actual.Value)
}
