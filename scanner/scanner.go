package scanner

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"unicode"

	"github.com/larzconwell/micro/errors"
	"github.com/larzconwell/micro/token"
)

// Scanner represents a scanner reading a buffered input.
type Scanner struct {
	Reader    *bufio.Reader
	MaxErrors int
}

// New creates a scanner from the given reader and max error count.
func New(reader *bufio.Reader, maxErrors int) *Scanner {
	return &Scanner{
		Reader:    reader,
		MaxErrors: maxErrors,
	}
}

// Scan scans all the tokens from the given input and may encounter
// a set number of errors before failing completely.
func (scanner *Scanner) Scan() (tokens []*token.Value, err error) {
	var errorCount int

	for {
		if errorCount >= scanner.MaxErrors {
			break
		}

		token, err := scanner.Next()
		if err == io.EOF {
			break
		}

		switch err.(type) {
		case nil:
			tokens = append(tokens, token)
			continue
		case errors.ScanError:
			errorCount++

			fmt.Fprintf(os.Stderr, "error #%d: %s\n", errorCount, err)
		default:
			return tokens, err
		}
	}

	if errorCount > 0 {
		err = errors.StepError("token scanning")
	}

	return tokens, err
}

// Next reads the scanners reader until it reads a token,
// or an error if an invalid token has been encountered. io.EOF is
// returned once the end of the input has been read, even if a token.End
// has been encountered such that it can be detected if an invalid program
// has been provided that does not include the token.End token.
func (scanner *Scanner) Next() (*token.Value, error) {
	for {
		char, _, err := scanner.Reader.ReadRune()
		if err != nil {
			return nil, err
		}

		switch {
		case unicode.IsSpace(char):
			continue
		case char == '(':
			return &token.Value{Token: token.LeftParen, Value: "("}, nil
		case char == ')':
			return &token.Value{Token: token.RightParen, Value: ")"}, nil
		case char == ';':
			return &token.Value{Token: token.Semicolon, Value: ";"}, nil
		case char == ',':
			return &token.Value{Token: token.Comma, Value: ","}, nil
		case char == '+':
			return &token.Value{Token: token.PlusOP, Value: "+"}, nil
		case unicode.IsLetter(char):
			return scanner.readIdentToken(char)
		case unicode.IsDigit(char):
			return scanner.readIntLiteralToken(char)
		case char == ':':
			return scanner.readAssignOPToken()
		case char == '-':
			token, err := scanner.readMinusOPOrCommentToken()
			if err != nil {
				return nil, err
			}

			if token != nil {
				return token, nil
			}
		default:
			return nil, errors.ScanError(char)
		}
	}
}

// token.readIdent reads an identifier token from the scanners reader,
// initial is the character that triggered the identifier scanning.
func (scanner *Scanner) readIdentToken(initial rune) (*token.Value, error) {
	chars := []rune{initial}

	for {
		char, _, err := scanner.Reader.ReadRune()
		if err != nil {
			return nil, err
		}

		chars = append(chars, char)

		// If the next character is not a letter, digit, or an underscore,
		// then we've reached the end of the identifier and we need to unread
		// the last character we read.
		if !(unicode.IsLetter(char) || unicode.IsDigit(char) || char == '_') {
			value := string(chars[:len(chars)-1])

			tok, ok := token.ReservedIdents[value]
			if !ok {
				tok = token.Ident
			}

			return &token.Value{
				Token: tok,
				Value: value,
			}, scanner.Reader.UnreadRune()
		}
	}
}

// token.readIntLiteral reads an integer literal token from the scanners reader,
// initial is the character that triggered the identifier scanning.
// IntLiteral ::= Digit | IntLiteral Digit
func (scanner *Scanner) readIntLiteralToken(initial rune) (*token.Value, error) {
	chars := []rune{initial}

	for {
		char, _, err := scanner.Reader.ReadRune()
		if err != nil {
			return nil, err
		}

		chars = append(chars, char)

		// If the next character is not a digit, then we've reached the end
		// of the literal and we need to unread the last character we read.
		if !unicode.IsDigit(char) {
			return &token.Value{
				Token: token.IntLiteral,
				Value: string(chars[:len(chars)-1]),
			}, scanner.Reader.UnreadRune()
		}
	}
}

// token.readAssignOP reads an assignment operator token from the scanners reader,
// it is assumed the colon has been consumed already.
func (scanner *Scanner) readAssignOPToken() (*token.Value, error) {
	char, _, err := scanner.Reader.ReadRune()
	if err != nil {
		return nil, err
	}

	if char != '=' {
		err = scanner.Reader.UnreadRune()
		if err != nil {
			return nil, err
		}

		// The invalid character is the colon, not what's after it.
		return nil, errors.ScanError(':')
	}

	return &token.Value{Token: token.AssignOP, Value: ":="}, nil
}

// token.readMinusOPOrComment reads a minus operator token or a comment token from
// the scanners reader, it is assumed the initial dash has been consumer already.
func (scanner *Scanner) readMinusOPOrCommentToken() (*token.Value, error) {
	char, _, err := scanner.Reader.ReadRune()
	if err != nil {
		return nil, err
	}

	if char != '-' {
		return &token.Value{Token: token.MinusOP, Value: "-"}, scanner.Reader.UnreadRune()
	}

	for {
		char, _, err = scanner.Reader.ReadRune()
		if err != nil {
			return nil, err
		}

		if char == '\n' {
			break
		}
	}

	return nil, nil
}
