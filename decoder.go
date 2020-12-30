package wrapper

import (
	"errors"
	"fmt"
	"strconv"
	"unicode"

	"github.com/mitchellh/mapstructure"
)

var (
	ErrIllegalChar = errors.New("illegal char")

	ErrExhaustedAllTokens = errors.New("exhausted all tokens")
)

type TokenType int

const (
	TTIllegal TokenType = iota
	TTStr
	TTInt
	TTFloat
	TTArrStart
	TTArrEnd
	TTMapStart
	TTMapEnd
)

var displayTokenType = []string{
	"ILL",
	"STR",
	"INT",
	"FLO",
	"ARRS",
	"ARRE",
	"MAPS",
	"MAPE",
}

var endSigChars = []byte{
	':',
	',',
	'b',
	'd',
	'f',
	's',
}

var specialEndSigChars = []byte{
	']',
	'}',
}

var emptyChars = []byte{
	' ',
	',',
	'\n',
	'\t',
}

func charIn(arr []byte, char byte) bool {
	for _, b := range arr {
		if b == char {
			return true
		}
	}
	return false
}

type Token struct {
	t TokenType
	v interface{}
}

func (t Token) String() string {
	if t.v == nil {
		return displayTokenType[t.t]
	}
	return fmt.Sprintf("%s:%v", displayTokenType[t.t], t.v)
}

type Lexer struct {
	r      []byte
	cur    int
	tokens []Token
}

func NewLexer(r []byte) *Lexer {
	return &Lexer{
		r:      r,
		cur:    -1,
		tokens: []Token{},
	}
}

func (lxr *Lexer) Next() bool {
	if lxr.cur+1 < len(lxr.r) {
		lxr.cur++
		return true
	}
	return false
}

func (lxr *Lexer) Char() byte {
	return lxr.r[lxr.cur]
}

func (lxr *Lexer) Tokens() []Token {
	return lxr.tokens
}

func (lxr *Lexer) Tokenize() error {
	for lxr.Next() {
		char := lxr.Char()
		switch char {
		case '{':
			lxr.appendToken(Token{t: TTMapStart})
		case '}':
			lxr.appendToken(Token{t: TTMapEnd})
		case '[':
			lxr.appendToken(Token{t: TTArrStart})
		case ']':
			lxr.appendToken(Token{t: TTArrEnd})
		case '"':
			s, err := lxr.buildStr("", '"')
			if err != nil {
				return err
			}
			lxr.appendToken(Token{t: TTStr, v: s})
		case '\'':
			s, err := lxr.buildStr("", '\'')
			if err != nil {
				return err
			}
			lxr.appendToken(Token{t: TTStr, v: s})
		case '-':
			n, t, err := lxr.buildNum("", true)
			if err != nil {
				return err
			}
			lxr.appendToken(Token{t: t, v: n})
		default:
			if charIn(emptyChars, char) {
				continue
			}
			if char >= '0' && char <= '9' {
				n, t, err := lxr.buildNum(string(char), false)
				if err != nil {
					return err
				}
				lxr.appendToken(Token{t: t, v: n})
				continue
			}
			if char <= unicode.MaxASCII {
				s, err := lxr.buildStr(string(char), ':')
				if err != nil {
					return err
				}
				lxr.appendToken(Token{t: TTStr, v: s})
				continue
			}
			return ErrIllegalChar
		}
	}
	return nil
}

func (lxr *Lexer) appendToken(t Token) {
	lxr.tokens = append(lxr.tokens, t)
}

func (lxr *Lexer) buildStr(start string, endChar byte) (string, error) {
	str := []byte(start)
	for lxr.Next() {
		char := lxr.Char()
		// The raw data for the player UUID is as follows: { UUID: [I; INT, INT, INT, INT] }
		// From the doc: "UUID of owner, stored as four ints.". Why the header 'I;' in the log???
		if char == endChar || char == ';' {
			break
		}
		str = append(str, char)
	}
	return string(str), nil
}

func (lxr *Lexer) buildNum(start string, isSigned bool) (string, TokenType, error) {
	num := []byte(start)
	if isSigned {
		num = append([]byte{'-'}, num...)
	}

	t := TTInt
	for lxr.Next() {
		char := lxr.Char()
		if charIn(specialEndSigChars, char) {
			lxr.cur--
			break
		}
		if charIn(endSigChars, char) {
			break
		}
		if char == '.' {
			if t == TTFloat {
				t = TTIllegal
			} else {
				t = TTFloat
			}
		}
		num = append(num, char)
	}
	return string(num), t, nil
}

type Parser struct {
	t   []Token
	cur int
}

func NewParser(t []Token) *Parser {
	return &Parser{
		t:   t,
		cur: -1,
	}
}

func (psr *Parser) Parse() (interface{}, error) {
	if !psr.Next() {
		return nil, ErrExhaustedAllTokens
	}
	token := psr.Token()
	switch token.t {
	case TTMapStart:
		return psr.buildMap()
	case TTArrStart:
		return psr.buildArr()
	case TTMapEnd:
		return nil, fmt.Errorf("invalid token at posn %d", psr.cur)
	case TTArrEnd:
		return nil, fmt.Errorf("invalid token at posn %d", psr.cur)
	case TTInt:
		return strconv.Atoi(token.v.(string))
	case TTFloat:
		return strconv.ParseFloat(token.v.(string), 64)
	default:
		return token.v, nil
	}
}

func (psr *Parser) Next() bool {
	if psr.cur+1 < len(psr.t) {
		psr.cur++
		return true
	}
	return false
}

func (psr *Parser) Token() Token {
	return psr.t[psr.cur]
}

func (psr *Parser) buildMap() (map[string]interface{}, error) {
	m := map[string]interface{}{}
	for psr.Next() {
		token := psr.Token()
		if token.t == TTMapEnd {
			break
		}
		key := token.v.(string)
		val, err := psr.Parse()
		if err != nil {
			return nil, err
		}
		m[key] = val
	}
	return m, nil
}

func (psr *Parser) buildArr() ([]interface{}, error) {
	arr := []interface{}{}
	for psr.Next() {
		token := psr.Token()
		if token.t == TTArrEnd {
			break
		}
		psr.cur--
		val, err := psr.Parse()
		if err != nil {
			return nil, err
		}
		arr = append(arr, val)
	}
	return arr, nil
}

func DecodeSNBT(r []byte, s interface{}) error {
	lxr := NewLexer(r)
	if err := lxr.Tokenize(); err != nil {
		return err
	}
	psr := NewParser(lxr.Tokens())
	dmap, err := psr.Parse()
	if err != nil && err != ErrExhaustedAllTokens {
		return err
	}
	return mapstructure.Decode(dmap, s)
}
