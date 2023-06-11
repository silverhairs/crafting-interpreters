package token

type TokenType string

const (
	// Single-character tokens
	L_PAREN   = "LEFT_PARENT"
	R_PAREN   = "RIGHT_PARENT"
	L_BRACE   = "LEFT_BRACE"
	R_BRACE   = "RIGHT_BRACE"
	COMMA     = "COMMA"
	DOT       = "DOT"
	MINUS     = "MINUS"
	PLUS      = "PLUS"
	SEMICOLON = "SEMICOLON"
	SLASH     = "SLASH"
	ASTERISK  = "ASTERISK"

	// One or two characters tokens
	BANG       = "BANG"
	BANG_EQ    = "BANG_EQUAL"
	EQUAL      = "EQUAL"
	EQ_EQ      = "EQUAL_EQUAL"
	GREATER    = "GREATER"
	GREATER_EQ = "GREATER_EQUAL"
	LESS       = "LESS"
	LESS_EQ    = "LESS_EQUAL"

	// Literals
	IDENTIFIER = "IDENT"
	STRING     = "STR"
	NUMBER     = "NUM"

	// Keywords
	AND       = "AND"
	CLASS     = "CLASS"
	ELSE      = "ELSE"
	IF        = "IF"
	FALSE     = "FALSE"
	FUNCTION  = "FUNCTION"
	FOR       = "FOR"
	OR        = "OR"
	NIL       = "NIL"
	PRINT     = "PRINT"
	RETURN    = "RETURN"
	SUPER     = "SUPER"
	THIS      = "THIS"
	TRUE      = "TRUE"
	LET       = "LET"
	WHILE     = "WHILE"
	COMMENT_L = "COMMENT_LINE"
	COMMENT_B = "COMMENT_BLOCK"

	EOF = "EOF"
)

type Token struct {
	Type    TokenType
	Literal any
	Lexeme  string
	Line    int
}

func New(Type TokenType, Lexeme string, Literal any, Line int) *Token {
	return &Token{
		Literal: Literal,
		Type:    Type,
		Lexeme:  Lexeme,
		Line:    Line,
	}
}

var keywords = map[string]TokenType{
	"and":    AND,
	"class":  CLASS,
	"else":   ELSE,
	"false":  FALSE,
	"true":   TRUE,
	"for":    FOR,
	"fn":     FUNCTION,
	"if":     IF,
	"nil":    NIL,
	"or":     OR,
	"print":  PRINT,
	"return": RETURN,
	"super":  SUPER,
	"this":   THIS,
	"let":    LET,
	"while":  WHILE,
}

func LookupIdentifier(keyword string) TokenType {
	tok, ok := keywords[keyword]
	if ok {
		return tok
	}

	return IDENTIFIER
}
