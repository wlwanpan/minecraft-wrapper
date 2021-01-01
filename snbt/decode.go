package snbt

import (
	"github.com/mitchellh/mapstructure"
)

// Decode parses SNBT data to the provided value of a pointer s.
func Decode(r []byte, s interface{}) error {
	lxr := newLexer(r)
	if err := lxr.tokenize(); err != nil {
		return err
	}
	psr := newParser(lxr.tokens())
	dmap, err := psr.parse()
	if err != nil && err != ErrExhaustedAllTokens {
		return err
	}
	return mapstructure.Decode(dmap, s)
}
