package token

// ReservedIdents is a collection of identifiers that are reserved
// and the tokens that represent them.
var ReservedIdents = map[string]Token{
	"begin": Begin,
	"end":   End,
	"read":  Read,
	"write": Write,
}
