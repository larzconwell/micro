package parser

import (
	"fmt"
	"os"

	"github.com/larzconwell/micro/errors"
	"github.com/larzconwell/micro/token"
)

// Parser represents a parser reading from a token set.
type Parser struct {
	TokenSet   *token.Set
	MaxErrors  int
	errorCount int
}

// New creates a parser from a given list of tokens
// and max number of errors that may occur.
func New(tokens []*token.Value, maxErrors int) *Parser {
	return &Parser{
		TokenSet:  token.NewSet(tokens),
		MaxErrors: maxErrors,
	}
}

// Parse parses the tokens handling up to the max number of errors.
func (parser *Parser) Parse() error {
	err := parser.parseProgram()
	if err != nil || parser.errorCount > 0 {
		return errors.StepError("parsing")
	}

	return nil
}

// parseProgram matches the entire program token stream.
func (parser *Parser) parseProgram() error {
	err := parser.match(token.Begin)
	if err != nil {
		return err
	}

	err = parser.parseStatementList()
	if err != nil {
		return err
	}

	return parser.match(token.End)
}

// parseStatementList matches one or more statements in a row.
func (parser *Parser) parseStatementList() error {
	err := parser.parseStatement()
	if err != nil {
		return err
	}

	for {
		next := parser.TokenSet.Peek()
		if next == nil {
			return nil
		}

		switch next.Token {
		case token.Ident, token.Read, token.Write:
			err = parser.parseStatement()
			if err != nil {
				return err
			}
		default:
			return nil
		}
	}
}

// parseStatement handles parsing any single valid statement.
func (parser *Parser) parseStatement() error {
	next := parser.TokenSet.Next()
	if next == nil {
		return parser.handleError(nil, token.Ident, token.Read, token.Write)
	}

	var err error
	switch next.Token {
	case token.Ident:
		err = parser.parseAssignment()
	case token.Read:
		err = parser.parseRead()
	case token.Write:
		err = parser.parseWrite()
	default:
		return parser.handleError(next, token.Ident, token.Read, token.Write)
	}

	if err != nil {
		return err
	}

	return parser.match(token.Semicolon)
}

// parseAssignment parses an assignment, and it is assumed the initial token.Ident has been parsed.
func (parser *Parser) parseAssignment() error {
	err := parser.match(token.AssignOP)
	if err != nil {
		return err
	}

	return parser.parseExpression()
}

// parseRead parses a read call, and it is assumed the initial token.Read has been parsed.
func (parser *Parser) parseRead() error {
	err := parser.match(token.LeftParen)
	if err != nil {
		return err
	}

	err = parser.parseIdentList()
	if err != nil {
		return err
	}

	return parser.match(token.RightParen)
}

// parseWrite parses a write call, and it is assumed the initial token.Write has been parsed.
func (parser *Parser) parseWrite() error {
	err := parser.match(token.LeftParen)
	if err != nil {
		return err
	}

	err = parser.parseExpressionList()
	if err != nil {
		return err
	}

	return parser.match(token.RightParen)
}

// parseIdentList parses a list of identifiers.
func (parser *Parser) parseIdentList() error {
	err := parser.match(token.Ident)
	if err != nil {
		return err
	}

	for {
		next := parser.TokenSet.Peek()
		if next == nil || next.Token != token.Comma {
			return nil
		}

		err = parser.match(token.Comma)
		if err != nil {
			return err
		}

		err = parser.match(token.Ident)
		if err != nil {
			return err
		}
	}
}

// parseExpression parses a single expression.
func (parser *Parser) parseExpression() error {
	err := parser.parsePrimary()
	if err != nil {
		return err
	}

	for {
		next := parser.TokenSet.Peek()
		if next == nil || !(next.Token == token.PlusOP || next.Token == token.MinusOP) {
			return nil
		}

		err = parser.parseArithmeticOP()
		if err != nil {
			return err
		}

		err = parser.parsePrimary()
		if err != nil {
			return err
		}
	}
}

// parseExpressionList parses a list of expressions.
func (parser *Parser) parseExpressionList() error {
	err := parser.parseExpression()
	if err != nil {
		return err
	}

	for {
		next := parser.TokenSet.Peek()
		if next == nil || next.Token != token.Comma {
			return nil
		}

		err = parser.match(token.Comma)
		if err != nil {
			return err
		}

		err = parser.parseExpression()
		if err != nil {
			return err
		}
	}
}

// parsePrimary parses a primary expression.
func (parser *Parser) parsePrimary() error {
	next := parser.TokenSet.Next()
	if next == nil {
		return parser.handleError(nil, token.LeftParen, token.Ident, token.IntLiteral)
	}

	switch next.Token {
	case token.LeftParen:
		err := parser.parseExpression()
		if err != nil {
			return err
		}

		return parser.match(token.RightParen)
	case token.Ident, token.IntLiteral:
		return nil
	default:
		return parser.handleError(next, token.LeftParen, token.Ident, token.IntLiteral)
	}
}

// parseArithmeticOP parses an arithmetic operator.
func (parser *Parser) parseArithmeticOP() error {
	return parser.match(token.PlusOP, token.MinusOP)
}

// match retrieves the next token and checks if it matches one of the given tokens.
func (parser *Parser) match(tokens ...token.Token) error {
	next := parser.TokenSet.Next()
	if next == nil {
		return parser.handleError(nil, tokens...)
	}

	for _, token := range tokens {
		if next.Token == token {
			return nil
		}
	}

	return parser.handleError(parser.TokenSet.Current, tokens...)
}

// handleError handles incrementing the error count and printing any errors that occur.
// After handling the error if the max error has been reached then it returns a step
// error that should be passed up.
func (parser *Parser) handleError(actual *token.Value, expected ...token.Token) error {
	parser.errorCount++

	if actual == nil {
		actual = &token.Value{Value: "EOF"}
	}

	err := &errors.ParseError{Expected: expected, Actual: actual}
	fmt.Fprintf(os.Stderr, "error #%d: %s\n", parser.errorCount, err)

	if parser.errorCount >= parser.MaxErrors {
		return errors.StepError("parsing")
	}

	return nil
}
