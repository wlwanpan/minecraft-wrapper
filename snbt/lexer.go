package snbt

import (
	"errors"
	"fmt"
	"unicode"
)

var (
	ErrIllegalChar = errors.New("illegal character")
)

type lexer struct {
	r   []byte
	cur int
	tks []token
}

func newLexer(r []byte) *lexer {
	return &lexer{
		r:   r,
		cur: -1,
		tks: []token{},
	}
}

func (lxr *lexer) next() bool {
	if lxr.cur+1 < len(lxr.r) {
		lxr.cur++
		return true
	}
	return false
}

func (lxr *lexer) curr() byte {
	return lxr.r[lxr.cur]
}

func (lxr *lexer) tokens() []token {
	return lxr.tks
}

func (lxr *lexer) tokenize() error {
	for lxr.next() {
		b := lxr.curr()
		switch b {
		case '{':
			lxr.appendToken(token{t: tkMapStart})
		case '}':
			lxr.appendToken(token{t: tkMapEnd})
		case '[':
			lxr.appendToken(token{t: tkArrStart})
		case ']':
			lxr.appendToken(token{t: tkArrEnd})
		case '"':
			s, err := lxr.buildStr("", '"')
			if err != nil {
				return err
			}
			lxr.appendToken(token{t: tkStr, v: s})
		case '\'':
			s, err := lxr.buildStr("", '\'')
			if err != nil {
				return err
			}
			lxr.appendToken(token{t: tkStr, v: s})
		case '-':
			n, t, err := lxr.buildNum("", true)
			if err != nil {
				return err
			}
			lxr.appendToken(token{t: t, v: n})
		default:
			if bytesContain(emptyBytes, b) {
				continue
			}
			if b >= '0' && b <= '9' {
				n, t, err := lxr.buildNum(string(b), false)
				if err != nil {
					return err
				}
				lxr.appendToken(token{t: t, v: n})
				continue
			}
			if b <= unicode.MaxASCII {
				s, err := lxr.buildStr(string(b), ':')
				if err != nil {
					return err
				}
				lxr.appendToken(token{t: tkStr, v: s})
				continue
			}
			return fmt.Errorf("%s: %b", ErrIllegalChar, b)
		}
	}
	return nil
}

func (lxr *lexer) appendToken(t token) {
	lxr.tks = append(lxr.tks, t)
}

func (lxr *lexer) buildStr(start string, end byte) (string, error) {
	str := []byte(start)
	for lxr.next() {
		b := lxr.curr()
		// The raw data for the player UUID is as follows: { UUID: [I; INT, INT, INT, INT] }
		// From the doc: "UUID of owner, stored as four ints.". NBT, why???
		if b == end || b == ';' {
			break
		}
		str = append(str, b)
	}
	return string(str), nil
}

func (lxr *lexer) buildNum(start string, isSigned bool) (string, tokenType, error) {
	num := []byte(start)
	if isSigned {
		num = append([]byte{'-'}, num...)
	}

	t := tkInt
	for lxr.next() {
		b := lxr.curr()
		if bytesContain(specialEndSigBytes, b) {
			lxr.cur--
			break
		}
		if bytesContain(endSigBytes, b) {
			break
		}
		if b == '.' {
			if t == tkFloat {
				t = tkIllegal
			} else {
				t = tkFloat
			}
		}
		num = append(num, b)
	}
	return string(num), t, nil
}
