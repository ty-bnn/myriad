package tokenizer

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ty-bnn/myriad/pkg/model/token"
)

func (t *Tokenizer) Tokenize() error {
	fmt.Printf("Tokenizing %s ...\n", t.filePath)

	for t.p < len(t.data) {
		if !t.isInDfBlock {
			tok, err := t.TokenizeMyriad()
			if err != nil {
				return err
			}
			if tok != (token.Token{}) {
				t.Tokens = append(t.Tokens, tok)
			}
		} else {
			tok, err := t.TokenizeDockerfile()
			if err != nil {
				return err
			}
			if tok != (token.Token{}) {
				t.Tokens = append(t.Tokens, tok)
			}
		}
	}

	fmt.Printf("Tokenize %s Done.\n", t.filePath)

	return nil
}

func (t *Tokenizer) TokenizeMyriad() (token.Token, error) {
	switch t.data[t.p] {
	case '(':
		t.p++
		return token.Token{Kind: token.LPAREN, Content: "("}, nil
	case ')':
		t.p++
		return token.Token{Kind: token.RPAREN, Content: ")"}, nil
	case ',':
		t.p++
		return token.Token{Kind: token.COMMA, Content: ","}, nil
	case '[':
		t.p++
		return token.Token{Kind: token.LBRACKET, Content: "["}, nil
	case ']':
		t.p++
		return token.Token{Kind: token.RBRACKET, Content: "]"}, nil
	case '{':
		if t.p+2 < len(t.data) && t.data[t.p:t.p+3] == "{{-" {
			t.p += 3
			for t.p < len(t.data) {
				if !isNewLine(t.data[t.p]) && !isWhiteSpace(t.data[t.p]) {
					break
				}
				t.p++
			}
			t.isInDfBlock = true
			return token.Token{Kind: token.DFBEGIN, Content: "{{-"}, nil
		} else {
			t.p++
			return token.Token{Kind: token.LBRACE, Content: "{"}, nil
		}
	case '}':
		if t.p+1 < len(t.data) && t.data[t.p:t.p+2] == "}}" {
			t.p += 2
			t.isInDfBlock = true
			return token.Token{Kind: token.RDOUBLEBRA, Content: "}}"}, nil
		} else {
			t.p++
			return token.Token{Kind: token.RBRACE, Content: "}"}, nil
		}
	case '.':
		t.p++
		return token.Token{Kind: token.DOT, Content: "."}, nil
	case ':':
		if t.p+1 < len(t.data) && t.data[t.p:t.p+2] == ":=" {
			t.p += 2
			return token.Token{Kind: token.DEFINE, Content: ":="}, nil
		} else {
			return token.Token{}, errors.New(fmt.Sprintf("tokenize error: invalid token ':'"))
		}
	case '=':
		if t.p+1 < len(t.data) && t.data[t.p:t.p+2] == "==" {
			t.p += 2
			return token.Token{Kind: token.EQUAL, Content: "=="}, nil
		} else {
			t.p++
			return token.Token{Kind: token.ASSIGN, Content: "="}, nil
		}
	case '!':
		if t.p+1 < len(t.data) && t.data[t.p:t.p+2] == "!=" {
			t.p += 2
			return token.Token{Kind: token.NOTEQUAL, Content: "!="}, nil
		} else {
			t.p++
			return token.Token{Kind: token.NOT, Content: "!"}, nil
		}
	case '&':
		if t.p+1 < len(t.data) && t.data[t.p:t.p+2] == "&&" {
			t.p += 2
			return token.Token{Kind: token.AND, Content: "&&"}, nil
		} else {
			return token.Token{}, errors.New(fmt.Sprintf("tokenize error: invalid token '&'"))
		}
	case '|':
		if t.p+1 < len(t.data) && t.data[t.p:t.p+2] == "||" {
			t.p += 2
			return token.Token{Kind: token.OR, Content: "||"}, nil
		} else {
			return token.Token{}, errors.New(fmt.Sprintf("tokenize error: invalid token '|'"))
		}
	case '"':
		// TODO: パース時に判断したい
		t.p++
		start := t.p
		for {
			if t.p >= len(t.data) || (t.data[t.p] == '"' && t.data[t.p-1] != '\\') {
				break
			}
			t.p++
			if t.p == len(t.data) {
				return token.Token{}, errors.New(fmt.Sprintf("tokenize error: cannot find '\"'"))
			}
		}
		content := t.data[start:t.p]
		t.p++
		return token.Token{Kind: token.STRING, Content: strings.Replace(content, "\\\"", "\"", -1)}, nil
	case '+':
		t.p++
		return token.Token{Kind: token.PLUS, Content: "+"}, nil
	case '<':
		if t.p+1 < len(t.data) && t.data[t.p:t.p+2] == "<<" {
			t.p += 2
			return token.Token{Kind: token.DOUBLELESS, Content: "<<"}, nil
		}
	default:
		if isWhiteSpace(t.data[t.p]) || isNewLine(t.data[t.p]) {
			t.p++
			return token.Token{}, nil
		} else if isLetter(t.data[t.p]) {
			start := t.p
			for t.p < len(t.data) && (isLetter(t.data[t.p]) || isDigit(t.data[t.p])) {
				t.p++
			}
			kind := getIdentKind(t.data[start:t.p])
			return token.Token{Kind: kind, Content: t.data[start:t.p]}, nil
		} else if isDigit(t.data[t.p]) {
			start := t.p
			for t.p < len(t.data) && isDigit(t.data[t.p]) {
				t.p++
			}
			return token.Token{Kind: token.NUMBER, Content: t.data[start:t.p]}, nil
		} else {
			return token.Token{}, errors.New(fmt.Sprintf("tokenize error: invalid token %c", t.data[t.p]))
		}
	}
	return token.Token{}, errors.New(fmt.Sprintf("tokenize error: invalid token"))
}

func (t *Tokenizer) TokenizeDockerfile() (token.Token, error) {
	if isNewLine(t.data[t.p]) {
		t.p++
		for t.p < len(t.data) {
			if !isNewLine(t.data[t.p]) && !isWhiteSpace(t.data[t.p]) {
				break
			}
			t.p++
		}
		return token.Token{Kind: token.DFARG, Content: "\n"}, nil
	}
	if t.nextTokenIs("-}}") {
		preContent := t.Tokens[len(t.Tokens)-1].Content
		if preContent[len(preContent)-1] != '\n' {
			return token.Token{Kind: token.DFARG, Content: "\n"}, nil
		}
		t.p += 3
		t.isInDfBlock = false
		return token.Token{Kind: token.DFEND, Content: "-}}"}, nil
	}
	if t.nextTokenIs("{{") {
		t.p += 2
		t.isInDfBlock = false
		return token.Token{Kind: token.LDOUBLEBRA, Content: "{{"}, nil
	}

	start := t.p
	for t.p < len(t.data) {
		if isNewLine(t.data[t.p]) || t.nextTokenIs("-}}") || t.nextTokenIs("{{") {
			break
		}
		t.p++
	}
	if t.nextTokenIs("-}}") {
		trimmed := strings.TrimRight(t.data[start:t.p], " \t")
		if trimmed != "" {
			return token.Token{Kind: token.DFARG, Content: trimmed}, nil
		}
	}
	return token.Token{Kind: token.DFARG, Content: t.data[start:t.p]}, nil
}
