package lexer

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type item struct {
	typ itemType
	pos int
	val string
}

type itemType int

const (
	itemError itemType = iota
	itemEOF
	itemLeftBracket
	itemRightBracket
	itemEqualSign
	itemNewLine
	itemSection
	itemKey
	itemValue
)

const (
	leftBracket  string = "["
	rightBracket string = "]"
	equalSign    string = "="
	newLine      string = "\n"
	spaceChars   string = " \t\r\n"
)

const eof = -1

type stateFn func(*lexer) stateFn

type lexer struct {
	input string
	start int
	pos   int
	width int
	state stateFn
	items chan item
}

func (l *lexer) next() rune {
	if l.isEOF() {
		l.width = 0
		return eof
	}

	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = w
	l.pos += l.width
	return r
}

func (l *lexer) backup() {
	l.pos -= l.width
}

func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *lexer) ignore() {
	l.start = l.pos
}

func (l *lexer) ignoreSpace() {
	for strings.ContainsRune(spaceChars, l.next()) {
	}

	l.backup()
	l.ignore()
}

func (l *lexer) isEOF() bool {
	return l.pos >= len(l.input)
}

func (l *lexer) contains(c string) bool {
	return strings.Index(l.input[l.pos:], c) >= 0
}

func (l *lexer) acceptUntil(c string) {
	for !strings.ContainsRune(c, l.next()) {
	}
	l.backup()
}

func (l *lexer) emit(t itemType) {
	l.items <- item{t, l.pos, l.input[l.start:l.pos]}
	l.start = l.pos
}

func (l *lexer) run() {
	for l.state = lexLoop; l.state != nil; {
		l.state = l.state(l)
	}
	close(l.items)
}

func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- item{itemError, l.start, fmt.Sprintf(format, args...)}
	return nil
}

func lex(input string) *lexer {
	l := &lexer{
		input: input,
		items: make(chan item),
	}

	go l.run()
	return l
}

func lexLoop(l *lexer) stateFn {
	l.ignoreSpace()

	if string(l.peek()) == leftBracket {
		return lexLeftBracket
	} else if !l.isEOF() {
		return lexKey
	}

	l.emit(itemEOF)
	return nil
}

func lexLeftBracket(l *lexer) stateFn {
	l.pos += len(leftBracket)
	l.emit(itemLeftBracket)
	return lexSection
}

func lexSection(l *lexer) stateFn {
	if l.contains(rightBracket) {
		l.acceptUntil(rightBracket)
		l.emit(itemSection)
		return lexRightBracket
	}
	return l.errorf("no right bracket found")
}

func lexRightBracket(l *lexer) stateFn {
	l.pos += len(rightBracket)
	l.emit(itemRightBracket)
	return lexLoop
}

func lexKey(l *lexer) stateFn {
	if l.contains(equalSign) {
		l.acceptUntil(equalSign)
		l.emit(itemKey)
		return lexEqualSign

	}
	return l.errorf("no equal sign found")
}

func lexEqualSign(l *lexer) stateFn {
	l.pos += len(equalSign)
	l.emit(itemEqualSign)
	return lexValue
}

func lexValue(l *lexer) stateFn {
	if l.contains(newLine) {
		l.acceptUntil(newLine)
		l.emit(itemValue)
		return lexLoop
	}
	return l.errorf("no new line found")
}
