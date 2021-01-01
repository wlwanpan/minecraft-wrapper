package snbt

import (
	"errors"
	"fmt"
	"strconv"
)

var (
	ErrExhaustedAllTokens = errors.New("exhausted all tokens")

	ErrInvalidToken = errors.New("invalid token")
)

type parser struct {
	t   []token
	cur int
}

func newParser(t []token) *parser {
	return &parser{
		t:   t,
		cur: -1,
	}
}

func (psr *parser) parse() (interface{}, error) {
	if !psr.next() {
		return nil, ErrExhaustedAllTokens
	}
	token := psr.token()
	switch token.t {
	case tkMapStart:
		return psr.buildMap()
	case tkArrStart:
		return psr.buildArr()
	case tkMapEnd:
		return nil, fmt.Errorf("%w: at posn %d", ErrInvalidToken, psr.cur)
	case tkArrEnd:
		return nil, fmt.Errorf("%w: at posn %d", ErrInvalidToken, psr.cur)
	case tkInt:
		return strconv.Atoi(token.v.(string))
	case tkFloat:
		return strconv.ParseFloat(token.v.(string), 64)
	default:
		return token.v, nil
	}
}

func (psr *parser) next() bool {
	if psr.cur+1 < len(psr.t) {
		psr.cur++
		return true
	}
	return false
}

func (psr *parser) token() token {
	return psr.t[psr.cur]
}

func (psr *parser) buildMap() (map[string]interface{}, error) {
	m := map[string]interface{}{}
	for psr.next() {
		token := psr.token()
		if token.t == tkMapEnd {
			break
		}
		key := token.v.(string)
		val, err := psr.parse()
		if err != nil {
			return nil, err
		}
		m[key] = val
	}
	return m, nil
}

func (psr *parser) buildArr() ([]interface{}, error) {
	arr := []interface{}{}
	for psr.next() {
		token := psr.token()
		if token.t == tkArrEnd {
			break
		}
		psr.cur--
		val, err := psr.parse()
		if err != nil {
			return nil, err
		}
		arr = append(arr, val)
	}
	return arr, nil
}
