package tokenizer

import (
	"fmt"
)

type TokenKind int
const (
	SIMPORT TokenKind = iota
	SFROM
	SMAIN
	SIF
	SELIF
	SELSE
	SLPAREN
	SRPAREN
	SCOMMA
	SARRANGE
	SLBRACE
	SRBRACE
	SEQUAL
	SNOTEQUAL
	SSTRING
	SDFCOMMAND
	SDFARG
	SIDENTIFIER
	SASSIGNVARIABLE
)

type Token struct {
	Content string
	Kind TokenKind
	Line int
}

// for Debug.
func printTokens(tokens []Token) {
	typeToToken := map[TokenKind]string{
		SIMPORT: "import",
		SFROM: "from",
		SMAIN: "main",
		SIF: "if",
		SELSE: "else",
		SLPAREN: "(",
		SRPAREN: ")",
		SCOMMA: ",",
		SARRANGE: "[]",
		SLBRACE: "{",
		SRBRACE: "}",
		SEQUAL: "=",
		SNOTEQUAL: "!=",
		SSTRING: "string",
		SDFCOMMAND: "DfCommand",
		SDFARG: "DfArg",
		SIDENTIFIER: "identifier",
		SASSIGNVARIABLE: "assignvariable",
	}

	for _, token := range tokens {
		fmt.Println("{")
		fmt.Printf("  Content: %s\n", token.Content)
		fmt.Printf("     Kind: %s\n", typeToToken[token.Kind])
		fmt.Println("},")
	}
}