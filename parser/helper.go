package parser

import (
	"fmt"
	"errors"

	"dcc/types"
)

func program(tokens []types.Token, index int) error {
	var err error

	// { 関数インポート文 }
	for ;; {
		if index >= len(tokens) || tokens[index].Kind != types.SIMPORT {
			break
		}

		index, err = importFunc(tokens, index)
		if err != nil {
			return err
		}
	}

	// { 関数 }
	for ;; {
		if index >= len(tokens) || tokens[index].Kind != types.SIDENTIFIER {
			break
		}

		index, err = function(tokens, index)
		if err != nil {
			return err
		}
	}

	if index >= len(tokens) || tokens[index].Kind != types.SMAIN {
		return nil
	}
	
	// メイン部
	index, err = mainFunction(tokens, index)
	if err != nil {
		return err
	}

	return nil
}

// 関数インポート文
func importFunc(tokens []types.Token, index int) (int, error) {
	var err error
	// "import"
	if index >= len(tokens) || tokens[index].Kind != types.SIMPORT {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find 'import'"))
	}

	index++

	// 関数名
	index, err = functionName(tokens, index)
	if err != nil {
		return index, err
	}

	// "from"
	if index >= len(tokens) || tokens[index].Kind != types.SFROM {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find 'from'"))
	}

	index++

	// ファイル名
	index, err = fileName(tokens, index)
	if err != nil {
		return index, err
	}

	return index, nil
}

// ファイル名
func fileName(tokens []types.Token, index int) (int, error) {
	if index >= len(tokens) || tokens[index].Kind != types.SSTRING {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find an identifier"))
	}

	index++

	return index, nil
}

// 関数
func function(tokens []types.Token, index int) (int, error) {
	var err error

	// 関数名
	index, err = functionName(tokens, index)
	if err != nil {
		return index, err
	}

	// 引数宣言部
	index, err = argumentDecralation(tokens, index)
	if err != nil {
		return index, err
	}

	// 関数記述部
	index, err = functionDescription(tokens, index)
	if err != nil {
		return index, err
	}

	return index, nil
}

// メイン部
func mainFunction(tokens []types.Token, index int) (int, error) {
	var err error

	// "main"
	if index >= len(tokens) || tokens[index].Kind != types.SMAIN {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find 'main' to declare"))
	}

	index++

	// 引数宣言部
	index, err = argumentDecralation(tokens, index)
	if err != nil {
		return index, err
	}

	// 関数記述部
	index, err = functionDescription(tokens, index)
	if err != nil {
		return index, err
	}

	return index, nil
}

// 引数宣言部
func argumentDecralation(tokens []types.Token, index int) (int, error) {
	var err error

	// "("
	if index >= len(tokens) || tokens[index].Kind != types.SLPAREN {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find '('"))
	}

	index++

	// 引数群
	if 	index < len(tokens) && tokens[index].Kind == types.SIDENTIFIER {
		index, err = arguments(tokens, index)
		if err != nil {
			return index, err
		}
	}

	// ")"
	if index >= len(tokens) || tokens[index].Kind != types.SRPAREN {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find ')'"))
	}

	index++

	return index, nil
}

// 引数群
func arguments(tokens []types.Token, index int) (int, error) {
	var err error

	// 変数
	index, err = variable(tokens, index)
	if err != nil {
		return index, err
	}

	for ;; {
		// ","
		if index >= len(tokens) || tokens[index].Kind != types.SCOMMA {
			break
		}

		index++
		
		// 変数
		index, err = variable(tokens, index)
		if err != nil {
			return index, err
		}
	}

	return index, nil
}

// 変数
func variable(tokens []types.Token, index int) (int, error) {
	var err error

	// 変数名
	index, err = variableName(tokens, index)
	if err != nil {
		return index, err
	}

	// "[]"
	if index < len(tokens) && tokens[index].Kind == types.SARRANGE {
		index++
	}

	return index, nil
}

// 関数記述部
func functionDescription(tokens []types.Token, index int) (int, error) {
	var err error

	// "{"
	if index >= len(tokens) || tokens[index].Kind != types.SLBRACE {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find '{'"))
	}

	index++

	// 記述部
	index, err = description(tokens, index)
	if err != nil {
		return index, err
	}

	// "}"
	if index >= len(tokens) || tokens[index].Kind != types.SRBRACE {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find '}'"))
	}

	index++

	return index, nil
}

// 記述部
func description(tokens []types.Token, index int) (int, error) {
	var err error

	// 記述ブロック
	index, err = descriptionBlock(tokens, index)
	if err != nil {
		return index, err
	}

	for ;; {
		if index >= len(tokens) || (tokens[index].Kind != types.SDFCOMMAND && tokens[index].Kind != types.SIDENTIFIER) {
			break
		}
			
		index, err = descriptionBlock(tokens, index)
		if err != nil {
			return index, err
		}
	}

	return index, nil
}

// 記述ブロック
func descriptionBlock(tokens []types.Token, index int) (int, error) {
	var err error

	if index < len(tokens) && tokens[index].Kind == types.SDFCOMMAND {
		// Dfile文
		index, err = dockerFile(tokens, index)
		if err != nil {
			return index, err
		}
	} else if index < len(tokens) && tokens[index].Kind == types.SIDENTIFIER {
		// 関数呼び出し文
		index, err = functionCall(tokens, index)
		if err != nil {
			return index, err
		}
	} else {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find a description block"))
	}

	return index, nil
}

// Dfile文
func dockerFile(tokens []types.Token, index int) (int, error) {
	var err error
	// Df命令
	if index >= len(tokens) || tokens[index].Kind != types.SDFCOMMAND {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find a Dockerfile comamnd"))
	}

	index++

	// Df引数部
	index, err = dfArgs(tokens, index)
	if err != nil {
		return index, err
	}

	return index, nil
}

// Df引数部
func dfArgs(tokens []types.Token, index int) (int, error) {
	var err error
	index, err = dfArg(tokens, index)
	if err != nil {
		return index, err
	}

	for ;; {
		if index < len(tokens) && (tokens[index].Kind == types.SDFARG || tokens[index].Kind == types.SASSIGNVARIABLE) {
			index, err = dfArg(tokens, index)
			if err != nil {
				return index, err
			}
		} else {
			break
		}
	}

	return index, nil
}

// Df引数
func dfArg(tokens []types.Token, index int) (int, error) {
	if index >= len(tokens) || (tokens[index].Kind != types.SDFARG && tokens[index].Kind != types.SASSIGNVARIABLE) {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find Df argument"))
	}

	index++

	return index, nil
}

// 関数呼び出し文
func functionCall(tokens []types.Token, index int) (int, error) {
	var err error

	// 関数名
	index, err = functionName(tokens, index)
	if err != nil {
		return index, err
	}

	// "("
	if index >= len(tokens) || tokens[index].Kind != types.SLPAREN {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find '('"))
	}

	index++

	// 文字列の並び
	if index < len(tokens) && tokens[index].Kind == types.SSTRING {
		index, err = rowOfStrings(tokens, index)
		if err != nil {
			return index, err
		}
	}

	// ")"
	if index >= len(tokens) || tokens[index].Kind != types.SRPAREN {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find ')'"))
	}

	index++

	return index, nil
}

// 文字列の並び
func rowOfStrings(tokens []types.Token, index int) (int, error) {
	// 文字列
	if index >= len(tokens) || tokens[index].Kind != types.SSTRING {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find a string"))
	}

	index++

	for ;; {
		if index >= len(tokens) || tokens[index].Kind != types.SCOMMA {
			break
		}

		index++

		if index >= len(tokens) || tokens[index].Kind != types.SSTRING {
			return index, errors.New(fmt.Sprintf("syntax error: cannot find a string"))
		}

		index++
	}

	return index, nil
}

// 関数名
func functionName(tokens []types.Token, index int) (int, error) {
	if index >= len(tokens) || tokens[index].Kind != types.SIDENTIFIER {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find an identifier"))
	}

	index++

	return index, nil
}

// 変数名
func variableName(tokens []types.Token, index int) (int, error) {
	if index >= len(tokens) || tokens[index].Kind != types.SIDENTIFIER {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find an identifier"))
	}

	index++

	return index, nil
}