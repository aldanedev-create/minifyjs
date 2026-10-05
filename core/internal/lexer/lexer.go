package lexer

import "fmt"

type TokenKind uint8

const (
	Illegal TokenKind = iota
	String
)

type Token struct {
	Kind  TokenKind
	Value string
}

type cursor struct {
	src string
	pos int
}

func (c *cursor) advance() byte {
	value := c.src[c.pos]
	c.pos++
	return value
}

func (c *cursor) peek() byte {
	return c.src[c.pos]
}

func (c *cursor) eof() bool {
	return c.pos >= len(c.src)
}

type Lexer struct {
	cur    cursor
	errors []error
}

func (l *Lexer) emit(kind TokenKind, value string, _, _, _ int, _ bool) Token {
	return Token{Kind: kind, Value: value}
}

func errorf(line, col int, message string) error {
	return fmt.Errorf("%d:%d: %s", line, col, message)
}