package jsonparser

type Lexer struct {
	input    string
	position int
	ch       byte
}

type TokenType string

const (
	UNKNOWN       TokenType = "UNKNOWN"
	SPACE         TokenType = "SPACE"
	NEWLINE       TokenType = "NEWLINE"
	TAB           TokenType = "TAB"
	CR            TokenType = "CR"
	LF            TokenType = "LF"
	FF            TokenType = "FF"
	VT            TokenType = "VT"
	NUMBER        TokenType = "NUMBER"
	STRING        TokenType = "STRING"
	BOOL          TokenType = "BOOL"
	NULL          TokenType = "NULL"
	COLON         TokenType = ":"
	COMMA         TokenType = ","
	OPEN_BRACE    TokenType = "{"
	OPEN_BRACKET  TokenType = "["
	CLOSE_BRACE   TokenType = "}"
	CLOSE_BRACKET TokenType = "]"
	QUOTE         TokenType = "\""
)

type Token struct {
	Type    TokenType
	Literal string
}

func NewLexer(input string) *Lexer {
	return &Lexer{input: input}
}

func (l *Lexer) Tokenize() []Token {
	tokens := []Token{}
	for l.position < len(l.input) {
		l.ch = l.input[l.position]
		if l.ch == '{' {
			tokens = append(tokens, Token{Type: OPEN_BRACE, Literal: "{"})
		} else if l.ch == '}' {
			tokens = append(tokens, Token{Type: CLOSE_BRACE, Literal: "}"})
		} else if l.ch == '[' {
			tokens = append(tokens, Token{Type: OPEN_BRACKET, Literal: "["})
		} else if l.ch == ']' {
			tokens = append(tokens, Token{Type: CLOSE_BRACKET, Literal: "]"})
		} else if l.ch == ':' {
			tokens = append(tokens, Token{Type: COLON, Literal: ":"})
		} else if l.ch == ',' {
			tokens = append(tokens, Token{Type: COMMA, Literal: ","})
		} else if l.ch == '"' {
			tokens = append(tokens, Token{Type: STRING, Literal: l.readString()})
		} else if l.isDigit() {
			tokens = append(tokens, Token{Type: NUMBER, Literal: l.readNumber()})
		} else if l.isBool() {
			tokens = append(tokens, Token{Type: BOOL, Literal: l.readBool()})
		} else if l.isNull() {
			tokens = append(tokens, Token{Type: NULL, Literal: l.readNull()})
		} else if l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' || l.ch == '\f' || l.ch == '\v' {
			l.position++
			continue
		} else {
			tokens = append(tokens, Token{Type: UNKNOWN, Literal: string(l.ch)})
		}
		l.position++
	}
	return tokens
}

func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.ch = l.input[position]
		if l.ch == '"' {
			break
		}
		position++
	}
	result := l.input[l.position+1 : position]
	l.position = position
	return result
}

func (l *Lexer) isDigit() bool {
	return (l.ch >= '0' && l.ch <= '9') ||
		(l.ch == '-' &&
			l.input[l.position+1] >= '0' &&
			l.input[l.position+1] <= '9') ||
		l.ch == '.'
}

func (l *Lexer) readNumber() string {
	position := l.position
	for {
		l.ch = l.input[position]
		if !l.isDigit() {
			break
		}
		position++
	}
	result := l.input[l.position:position]
	l.position = position - 1
	return result
}

func (l *Lexer) readBool() string {
	if l.input[l.position:l.position+4] == "true" {
		result := l.input[l.position : l.position+4]
		l.position = l.position + 3
		return result
	} else {
		result := l.input[l.position : l.position+5]
		l.position = l.position + 4
		return result
	}
}

func (l *Lexer) isBool() bool {
	if l.position+4 > len(l.input) || l.position+5 > len(l.input) {
		return false
	}
	return l.input[l.position:l.position+4] == "true" || l.input[l.position:l.position+5] == "false"
}

func (l *Lexer) readNull() string {
	result := l.input[l.position : l.position+4]
	l.position = l.position + 3
	return result
}

func (l *Lexer) isNull() bool {
	if l.position+4 > len(l.input) {
		return false
	}
	return l.input[l.position:l.position+4] == "null"
}
