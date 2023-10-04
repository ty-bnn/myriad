package token

type TokenKind int

const (
	IMPORT TokenKind = iota
	FROM
	MAIN
	IF
	ELIF
	ELSE
	FOR
	IN
	LPAREN
	RPAREN
	COMMA
	LBRACE
	RBRACE
	LBRACKET
	RBRACKET
	DEFINE
	ASSIGN
	EQUAL
	NOTEQUAL
	LDOUBLEBRA
	RDOUBLEBRA
	STRING
	DFCOMMAND
	DFARG
	IDENTIFIER
	NUMBER
)

type Token struct {
	Content string
	Kind    TokenKind
	Line    int
}
