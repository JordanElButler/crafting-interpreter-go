package main

import (
	"fmt"
)

type ScanError struct {
	Code ErrorCode
	Expr string
}

func (e *ScanError) Error() string {
	return "error scanning input: " + e.Code.String() + ": `" + e.Expr + "`"
}

type ErrorCode string

const (
	ErrUnspecifiedError ErrorCode = "Unknown error"
	ErrUnrecognizedToken = "Unrecognized token"
	ErrString = "String error"
	ErrNumber = "Number error"
	ErrIdentifier = "Identifer error"
)
func (e ErrorCode) String() string {
	return string(e)
}

type TokenType string

const (
	Undefined TokenType = "Undefined"

	Left_Paren = "("
	Right_Paren = ")"
	Left_Brace = "{"
	Right_Brace = "}"
	Comma = ","
	Dot = "."
	Minus = "-"
	Plus = "+"
	Semicolon = ";"
	Slash = "/"
	Star = "*"

	Bang = "!"
	BangEqual = "!="
	Equal = "="
	Equal_Equal = "=="
	Greater = ">"
	Greater_Equal = ">="
	Less = "<"
	Less_Equal = "<="

	Identifier = "Identifier"
	String = "String"
	Number = "Number"

	And = "and"
	Class = "class"
	Else = "else"
	False = "false"
	Fun = "fun"
	For = "for"
	If = "if"
	Nil = "nil"
	Or = "or"
	Print = "print"
	Return = "return"
	Super = "super"
	This = "this"
	True = "true"
	Var = "var"
	While = "while"

	EOF = "EOF"
)

var ReservedMap = []TokenType {
	And,
	Class,
	Else,
	False,
	Fun,
	For,
	If,
	Nil,
	Or,
	Print,
	Return,
	Super,
	This,
	True,
	Var,
	While,
}

var PunctuationOpMap = []TokenType {
	Left_Paren,
	Right_Paren,
	Left_Brace,
	Right_Brace,
	Comma,
	Dot,
	Minus,
	Plus,
	Semicolon,
	Slash,
	Star,

	Bang,
	BangEqual,
	Equal,
	Equal_Equal,
	Greater,
	Greater_Equal,
	Less,
	Less_Equal,	
}
func (t TokenType) String() string {
	return string(t)
}

func (t TokenType) Runes() []rune {
	return []rune(t)
}

type Token struct {
	ttype  TokenType
	val    []rune
	indexBegin  int
	indexEnd int
	line   int
	offset int
}

func (t Token) String() string {
	return fmt.Sprintf("%v: (%v) start: %v ending: %v", t.ttype, t.val, t.indexBegin, t.indexEnd)
}

type TokenScanner struct {
	tokens    []Token
	hasError  bool
	scanError ScanError
	source    []rune
	sourceLen int
	index     int
	line      int
	offset    int
}

func newTokenScanner(source string) TokenScanner {
	return TokenScanner{
		tokens:    make([]Token, 0),
		hasError: false,
		scanError: ScanError{Code: ErrUnspecifiedError, Expr: "" },
		source:    []rune(source),
		sourceLen: len([]rune(source)),
		index:     0,
	}
}

func (ts *TokenScanner) setError(m_error ScanError) {
	ts.hasError = true
	ts.scanError = m_error
}

func (ts *TokenScanner) addToken(ttype TokenType, val []rune, indexBegin int, indexEnd int) {

	t := Token{
		ttype,
		val,
		indexBegin,
		indexEnd,
		ts.line,
		ts.offset,
	}
	ts.tokens = append(ts.tokens, t)
}
func (ts *TokenScanner) isAtEnd() bool {
	return ts.index >= ts.sourceLen
}

func (ts *TokenScanner) atEnd(index int) bool {
	return ts.sourceLen <= index
}

func (ts *TokenScanner) peak() rune {
	if !ts.isAtEnd() {
		return ts.source[ts.index]
	} else {
		return rune(0)
	}
}

func (ts *TokenScanner) peakAhead() rune {
	if !ts.atEnd(ts.index + 1) {
		return ts.source[ts.index + 1]
	} else {
		return rune(0)
	}
}

func (ts *TokenScanner) advance(l int) {
	ts.index += l
	ts.offset += l
}

func (ts *TokenScanner) scanToken() {

	if ts.string() {
		
	} else if ts.number() {
		
	} else if ts.identifier() {
		
	} else if ts.punctuationOrOp() {

	} else if ts.whitespace() {
		
	} else {
		ts.setError(ScanError {
			Code: ErrUnrecognizedToken,
			Expr: "",
		})
		ts.advance(1)
	}
}
func (ts *TokenScanner) string() bool {
	if ts.source[ts.index] == '"' {
		res := ts.scanString()
		if !res {
			ts.setError(ScanError {
				Code: ErrString,
				Expr: "",
			})
		}
		return true
	} else {
		return false
	}
}
func (ts *TokenScanner) scanString() bool {
	start := ts.index
	index := start + 1
	closedFlag := false
	stringval := make([]rune, 0)
	
	// do escape characters
	for !ts.atEnd(index) && !closedFlag {
		var c = ts.source[index]
		if c == 0x5c {
			index += 1
			c = ts.source[index]
			if c == 0x22  || c == 0x5c {
				stringval = append(stringval, c)
			} else if c == '\n' {
				stringval = append(stringval, '\n')
			} else if c == 'r' {
				stringval = append(stringval, '\r')
			} else if c == 't' {
				stringval = append(stringval, '\t')
			} else {
				stringval = append(stringval, 0x5c)
				stringval = append(stringval, c)
			}
		} else if c == '"' {
			closedFlag = true
		} else {
			stringval = append(stringval, c)
		}
		index += 1
	}
	
	ts.advance(index - start)
	
	if closedFlag {
		ts.addToken(String, stringval, start, index)
		return true
	} else {
		return false	
	}
}
func (ts *TokenScanner) number() bool {
	if isDigit(ts.source[ts.index]) {
		res := ts.scanNumber()
		if !res {
			ts.setError(ScanError {
				Code: ErrNumber,
				Expr: "",
			})
		}
		return true
	} else {
		return false
	}
}

func (ts *TokenScanner) scanNumber() bool {
	start := ts.index
	index := start
	
	decimalFlag := false
	numberval := make([]rune, 0)
	
	
	for !ts.atEnd(index) {
		c := ts.source[index]
		if isDigit(c) {
			numberval = append(numberval, c)
		} else if c == '.' && !decimalFlag {
			decimalFlag = true
			numberval = append(numberval, '.')
		} else {
			break
		}
		index += 1
	}
	
	ts.advance(index - start)
	
	if len(numberval) > 0 && numberval[len(numberval)-1] != '.' {
		ts.addToken(Number, numberval, start, index)
		return true
	} else {
		return false
	}
}

func (ts *TokenScanner) identifier() bool {
	c := ts.source[ts.index]
	if isAlphabetical(c) || c == '_' {
		res := ts.scanIdentifier()
		if !res {
			ts.setError(ScanError {
				Code: ErrIdentifier,
				Expr: "",
			})
		}
		return true
	} else {
		return false
	}
}

func (ts *TokenScanner) scanIdentifier() bool {
	start := ts.index
	index := start
	endFlag := false
	identifierval := make([]rune, 0)
	
	for !ts.atEnd(index) && !endFlag {
		c := ts.source[index]
		
		if isDigit(c) || isAlphabetical(c) || c == '_' {
			identifierval = append(identifierval, c)
		} else {
			endFlag = true
		}
		index += 1
	}
	
	ts.advance(index - start)
	
	if true {
		isReserved, reserved := checkReserved(identifierval)
		
		if isReserved {
			ts.addToken(reserved, identifierval, start, index)
		} else {
			ts.addToken(Identifier, identifierval, start, index)
		}
		return true
	} else {
		return false
	}
}

func checkReserved(ident []rune) (bool, TokenType) {
	
	for _, v := range ReservedMap {
		if string(ident) == string(v) {
			return true, v
		}
	}
	
	return false, Undefined
}

func (ts *TokenScanner) punctuationOrOp() bool {
	res := ts.scanPunctuationOrOp()
	if res {
		return true
	} else {
		return false
	}	
}

func (ts *TokenScanner) scanPunctuationOrOp() bool {
	start := ts.index
	foundFlag := false
	foundLen := -1
	foundToken := Undefined
	
	for _, v := range PunctuationOpMap {
		if ts.sourceLen >= start + len(v) && string(ts.source[start : start + len(v)]) == string(v) && foundLen < len(v) {
			foundFlag = true
			foundLen = len(v)
			foundToken = v
		}
	}
	
	if foundFlag {
		ts.advance(foundLen)
		ts.addToken(foundToken, []rune(foundToken), start, start + len(foundToken))
		return true
	} else {
		return false
	}
}

func (ts * TokenScanner) whitespace() bool {
	res := ts.scanWhitespace()
	if res {
		return true
	} else {
		return false
	}
}
func (ts *TokenScanner) scanWhitespace() bool {
	start := ts.index
	index := start
	lines := 0
	c := ts.source[index]
	switch c {
		case ' ', '\t': {
			index += 1
		}
		case '\n', '\r': {
			lines += 1
			index += 1
		}
		default: {
			break
		}
	}
	
	ts.advance(1)
	if lines != 0 {
		ts.line += lines
		ts.offset = 0
	}

	if start != index {
		return true
	} else {
		return false
	}
}

func scanTokens(source string) ([]Token, bool, ScanError) {

	ts := newTokenScanner(source)

	for !ts.isAtEnd() && !ts.hasError {
		ts.scanToken()
	}

	if ts.isAtEnd() {
		ts.addToken(EOF, make([]rune, 0), ts.sourceLen, ts.sourceLen)
	}
	
	return ts.tokens, ts.hasError, ts.scanError
}

func isAlphabetical(r rune) bool {
	var t = false
	t = t || ( r >= 'a' && r <= 'z' )
	t = t || ( r >= 'A' && r <= 'Z' )
	
	return t
}

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}
