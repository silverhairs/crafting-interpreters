package lexer

import (
	"fmt"
	"glox/exception"
	"glox/token"
	"strconv"
)

type Lexer struct {
	Source  string
	tokens  []token.Token
	start   int
	current int
	line    int
}

func New(Source string) *Lexer {
	return &Lexer{
		Source:  Source,
		tokens:  make([]token.Token, 0),
		start:   0,
		current: 0,
		line:    1,
	}
}

func (lxr *Lexer) Tokenize() ([]token.Token, error) {
	var err error
	for !lxr.isAtEnd() {
		lxr.start = lxr.current
		err = lxr.lex()
	}

	lxr.tokens = append(lxr.tokens, token.Token{Type: token.EOF, Literal: nil, Lexeme: "", Line: lxr.line})
	return lxr.tokens, err
}

func (s *Lexer) lex() error {
	var err error
	char := s.advance()

	switch char {
	case '(':
		s.addTokenType(token.L_PAREN)
	case ')':
		s.addTokenType(token.R_PAREN)
	case '{':
		s.addTokenType(token.L_BRACE)
	case '}':
		s.addTokenType(token.R_BRACE)
	case ',':
		s.addTokenType(token.COMMA)
	case '.':
		s.addTokenType(token.DOT)
	case '-':
		s.addTokenType(token.MINUS)
	case '+':
		s.addTokenType(token.PLUS)
	case ';':
		s.addTokenType(token.SEMICOLON)
	case '*':
		s.addTokenType(token.ASTERISK)
	case '?':
		s.addTokenType(token.QUESTION_MARK)
	case ':':
		s.addTokenType(token.COLON)
	case '!':
		s.operator(struct {
			char     byte
			unique   token.TokenType
			twoChars token.TokenType
		}{'=', token.BANG, token.BANG_EQ},
		)
	case '=':
		s.operator(struct {
			char     byte
			unique   token.TokenType
			twoChars token.TokenType
		}{'=', token.EQUAL, token.EQ_EQ},
		)
	case '<':
		s.operator(
			struct {
				char     byte
				unique   token.TokenType
				twoChars token.TokenType
			}{'=', token.LESS, token.LESS_EQ},
		)
	case '>':
		s.operator(struct {
			char     byte
			unique   token.TokenType
			twoChars token.TokenType
		}{'=', token.GREATER, token.GREATER_EQ},
		)
	case '/':
		err = s.slash()
	case ' ':
	case '\r':
	case '\t':
		break
	case '\n':
		s.line++
	case '"':
		err = s.string()
	default:
		if isDigit(char) {
			err = s.number()
		} else if isAlpha(char) {
			s.identifier()
		} else {
			s.addToken(token.ILLEGAL, s.peek())
		}
	}

	return err
}

func (s *Lexer) isAtEnd() bool {
	return s.current >= len(s.Source)
}

func (s *Lexer) advance() byte {
	prev := s.current
	s.current++
	return s.Source[prev]
}

func (s *Lexer) addTokenType(tokenType token.TokenType) {
	s.addToken(tokenType, nil)
}

func (s *Lexer) addToken(tokenType token.TokenType, literal any) {
	lexeme := s.Source[s.start:s.current]
	tok := token.Token{Type: tokenType, Literal: literal, Lexeme: lexeme, Line: s.line}
	s.tokens = append(s.tokens, tok)
}

func (s *Lexer) match(expect byte) bool {
	if s.isAtEnd() || s.Source[s.current] != expect {
		return false
	}

	s.current++
	return true
}

func (s *Lexer) operator(props struct {
	char     byte
	unique   token.TokenType // If the lexeme has only one character, which token type should be recorded.
	twoChars token.TokenType // If the lexeme has two characters, which token type should be recorded.
}) {
	var tok token.TokenType
	if s.match(props.char) {
		tok = props.twoChars
	} else {
		tok = props.unique
	}

	s.addTokenType(tok)
}

func (s *Lexer) peek() byte {
	if s.isAtEnd() {
		return '\000'
	}
	return s.Source[s.current]
}

// tokenizes a string literal.
func (s *Lexer) string() error {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			// Multi-line string literals are allowed
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		return exception.Short(s.line, "Please add a double-quote at the end of the string.")

	}

	s.advance()

	value := s.Source[s.start+1 : s.current-1]
	s.addToken(token.STRING, value)
	return nil
}

// Scans number literals, this handles all floating-point numbers with or without decimals.
func (s *Lexer) number() error {
	for isDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == '.' && isDigit(s.peekNext()) {
		s.advance()
		for isDigit(s.peek()) {
			s.advance()
		}
	}
	literal := s.Source[s.start:s.current]
	value, err := strconv.ParseFloat(literal, 64)
	if err != nil {
		return exception.Short(s.line, fmt.Sprintf("%q is an invalid %q", literal, token.NUMBER))

	}
	s.addToken(token.NUMBER, value)
	return nil
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (s *Lexer) peekNext() byte {
	if s.current+1 >= len(s.Source) {
		return '\000'
	}

	return s.Source[s.current+1]
}

func isAlpha(char byte) bool {
	return char >= 'a' && char <= 'z' || char >= 'A' && char <= 'Z' || char == '_'
}

func isAlphaNumeric(char byte) bool {
	return isAlpha(char) || isDigit(char)
}

func (s *Lexer) identifier() {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}

	literal := s.Source[s.start:s.current]
	tok := token.LookupIdentifier(literal)
	s.addTokenType(tok)
}

func (s *Lexer) slash() error {
	if s.match('/') {
		for s.peek() != '\n' && !s.isAtEnd() {
			s.advance()
		}
		literal := s.Source[s.start+2 : s.current]
		s.addToken(token.SLASH_SLASH, literal)

	} else if s.match('*') {
		for s.peek() != '*' && !s.isAtEnd() {
			s.advance()
		}

		if s.match('/') {
			literal := s.Source[s.start+2 : s.current-2]
			s.addToken(token.SLASK_ASTERISK, literal)
		} else {
			return exception.Short(s.line, "opened multi-line comment has not been closed.")
		}

	} else {
		s.addTokenType(token.SLASH)
	}
	return nil
}
