package token

// Set is a collection of tokens that can be stepped through.
type Set struct {
	Current *Value
	Tokens  []*Value
}

// NewSet creates a Set from a list of tokens.
func NewSet(tokens []*Value) *Set {
	return &Set{Tokens: tokens}
}

// Next gets the next token in the token set or nil if there's no more tokens to process.
func (set *Set) Next() *Value {
	if len(set.Tokens) == 0 {
		set.Current = nil
		return nil
	}

	set.Current = set.Tokens[0]
	set.Tokens = set.Tokens[1:]
	return set.Current
}

// Peek peeks into the next token without advancing to it.
func (set *Set) Peek() *Value {
	if len(set.Tokens) == 0 {
		return nil
	}

	return set.Tokens[0]
}
