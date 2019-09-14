package token

// Valid tokens.
const (
	Begin Token = iota
	End
	Read
	Write
	Ident
	IntLiteral
	LeftParen
	RightParen
	Semicolon
	Comma
	AssignOP
	PlusOP
	MinusOP
)

// Token represents a single valid token from input text.
type Token int

func (token Token) String() string {
	switch token {
	case Begin:
		return "begin"
	case End:
		return "end"
	case Read:
		return "read"
	case Write:
		return "write"
	case Ident:
		return "identifier"
	case IntLiteral:
		return "integer"
	case LeftParen:
		return "("
	case RightParen:
		return ")"
	case Semicolon:
		return ";"
	case Comma:
		return ","
	case AssignOP:
		return ":="
	case PlusOP:
		return "+"
	case MinusOP:
		return "-"
	default:
		return "unknown"
	}
}
