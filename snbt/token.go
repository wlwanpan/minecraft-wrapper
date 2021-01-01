package snbt

import (
	"fmt"
)

type tokenType int

const (
	tkIllegal tokenType = iota
	tkStr
	tkInt
	tkFloat
	tkArrStart
	tkArrEnd
	tkMapStart
	tkMapEnd
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

var endSigBytes = []byte{
	':',
	',',
	'b',
	'd',
	'f',
	's',
}

var specialEndSigBytes = []byte{
	']',
	'}',
}

var emptyBytes = []byte{
	' ',
	',',
	'\n',
	'\t',
}

func bytesContain(bs []byte, tb byte) bool {
	for _, b := range bs {
		if b == tb {
			return true
		}
	}
	return false
}

type token struct {
	t tokenType
	v interface{}
}

func (t token) string() string {
	if t.v == nil {
		return displayTokenType[t.t]
	}
	return fmt.Sprintf("%s:%v", displayTokenType[t.t], t.v)
}
