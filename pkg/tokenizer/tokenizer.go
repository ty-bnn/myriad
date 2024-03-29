package tokenizer

import "github.com/ty-bnn/myriad/pkg/model/token"

type Tokenizer struct {
	data        string
	p           int
	isInDfBlock bool
	isInCommand bool
	commandPtr  string
	Tokens      []token.Token
	filePath    string
}

func NewTokenizer(data string, filePath string) *Tokenizer {
	return &Tokenizer{
		data:     data,
		filePath: filePath,
	}
}
