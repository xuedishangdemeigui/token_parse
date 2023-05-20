package main

import (
	"fmt"

	"github.com/pkg/errors"
)

var (
	ErrInvalidToken = errors.New("invalid token")
)

const (
	Apostrophe      = '`'
	DoubleQuotation = '"'
	SingleQuotation = '\''
	Delimiter       = '|'
)

type Keyword struct {
	Text    string
	Pos     int
	CurChar byte
	Len     int
}

func NewKeyword(text string) *Keyword {
	if len(text) <= 0 {
		panic(ErrInvalidToken)
	}

	k := &Keyword{
		Text:    text,
		Pos:     0,
		CurChar: text[0],
		Len:     len(text),
	}

	return k
}

func (k *Keyword) Advance() (isEnd bool) {
	for k.Pos < k.Len-1 {
		k.Pos += 1
		k.CurChar = k.Text[k.Pos]
		if k.CurChar != ' ' {
			break
		}
	}

	return k.Pos == k.Len-1
}

func (k *Keyword) GetNextToken() (token string, isEnd bool, err error) {
	var (
		begin     = k.Pos
		end       = 0
		quotation byte
	)

	switch k.CurChar {
	case Apostrophe:
		quotation = Apostrophe
	case DoubleQuotation:
		quotation = DoubleQuotation
	case SingleQuotation:
		quotation = SingleQuotation
	default:
		return "", isEnd, errors.WithStack(ErrInvalidToken)
	}

	for isEnd = k.Advance(); k.Pos < k.Len; isEnd = k.Advance() {
		if k.CurChar == quotation {
			end = k.Pos
			break
		}
	}

	if !isEnd {
		isEnd = k.Advance()
		if !isEnd && k.CurChar != Delimiter {
			return "", isEnd, errors.WithStack(ErrInvalidToken)
		}
		isEnd = k.Advance()
	}

	if end == 0 {
		return "", isEnd, errors.WithStack(ErrInvalidToken)
	}

	return k.Text[begin+1 : end], isEnd, nil
}

func ParseKeyword(text string) (tokens []string, err error) {
	k := NewKeyword(text)

	for {
		token, isEnd, err := k.GetNextToken()
		if err != nil {
			return nil, err
		}

		fmt.Printf("[INFO] tokem: %s\n", token)

		tokens = append(tokens, token)
		if isEnd {
			break
		}
	}

	return
}

func main() {
	tokens, err := ParseKeyword(`"hello" | '"world"' | "'gol|a'ng"`)
	if err != nil {
		fmt.Printf("[ERROR] parse failed, error: %+v\n", err)
		return
	}

	fmt.Printf("[INFO] tokens: %v\n", tokens)
}
