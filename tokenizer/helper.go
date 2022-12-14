package tokenizer

import(
	"fmt"
	"errors"

	"dcc/types"
)

func readReservedWords(index int, lineStr string, line int) (int, types.Token, error) {
	if (lineStr[index : index + 6] == "import") {
		return index + 6, types.Token{"import", types.SIMPORT, line + 1}, nil 
	} else if (lineStr[index : index + 4] == "from") {
		return index + 4, types.Token{"from", types.SFROM, line + 1}, nil
	} else if (lineStr[index : index + 4] == "main") {
		return index + 4, types.Token{"main", types.SMAIN, line + 1}, nil
	} else if (lineStr[index : index + 2] == "if") {
		return index + 2, types.Token{"if", types.SIF, line + 1}, nil
	} else if (lineStr[index : index + 7] == "else if") {
		return index + 7, types.Token{"else if", types.SELIF, line + 1}, nil
	} else if (lineStr[index : index + 4] == "else") {
		return index + 4, types.Token{"else", types.SELSE, line + 1}, nil
	}

	return index, types.Token{}, errors.New(fmt.Sprintf("ReservedWords'index: %d find invalid token in \"%s\".", index, lineStr)) 
}

func readSymbols(index int, lineStr string, line int) (int, types.Token, error) {
	if (lineStr[index] == '(') {
		return index + 1, types.Token{"(", types.SLPAREN, line + 1}, nil
	} else if (lineStr[index] == ')') {
		return index + 1, types.Token{")", types.SRPAREN, line + 1}, nil
	} else if (lineStr[index] == ',') {
		return index + 1, types.Token{",", types.SCOMMA, line + 1}, nil
	} else if (lineStr[index : index + 1] == "[]") {
		return index + 2, types.Token{"[]", types.SARRANGE, line + 1}, nil
	} else if (lineStr[index] == '{') {
		return index + 1, types.Token{"{", types.SLBRACE, line + 1}, nil
	} else if (lineStr[index] == '}') {
		return index + 1, types.Token{"}", types.SRBRACE, line + 1}, nil
	} else if (lineStr[index : index + 1] == "==") {
		return index + 2, types.Token{"==", types.SEQUAL, line + 1}, nil
	} else if (lineStr[index : index + 1] == "!=") {
		return index + 2, types.Token{"!=", types.SNOTEQUAL, line + 1}, nil
	}

	return index, types.Token{}, errors.New(fmt.Sprintf("Symbols'index: %d find invalid token in \"%s\".", index, lineStr))
}

func readDfCommands(index int, lineStr string, line int) (int, types.Token, error) {
	/*
	ADD, ARG, CMD, COPY, ENTRYPOINT, ENV, EXPOSE, FROM,
	HEALTHCHECK, LABEL, MAINTAINER, ONBUILD, RUN, SHELL,
	STOPSIGNAL, USER, VOLUME, WORKDIR
	*/
	switch lineStr[index : index + 4] {
		case "ADD ", "ARG ", "CMD ", "ENV ", "RUN ":
			return index + 4, types.Token{lineStr[index : index + 3], types.SDFCOMMAND, line + 1}, nil
	}
	switch lineStr[index : index + 5] {
		case "COPY ", "FROM ", "USER ":
			return index + 5, types.Token{lineStr[index : index + 4], types.SDFCOMMAND, line + 1}, nil
	}
	switch lineStr[index : index + 6] {
		case "LABEL ", "SHELL ":
			return index + 6, types.Token{lineStr[index : index + 5], types.SDFCOMMAND, line + 1}, nil
	}
	switch lineStr[index : index + 7] {
		case "EXPOSE ", "VOLUME ":
			return index + 7, types.Token{lineStr[index : index + 6], types.SDFCOMMAND, line + 1}, nil
	}
	switch lineStr[index : index + 8] {
		case "ONBUILD ", "WORKDIR ":
			return index + 8, types.Token{lineStr[index : index + 7], types.SDFCOMMAND, line + 1}, nil
	}
	switch lineStr[index : index + 11] {
		case "ENTRYPOINT ", "MAINTAINER ", "STOPSIGNAL ":
			return index + 11, types.Token{lineStr[index : index + 10], types.SDFCOMMAND, line + 1}, nil
	}
	switch lineStr[index : index + 12] {
		case "HEALTHCHECK ":
			return index + 11, types.Token{lineStr[index : index + 10], types.SDFCOMMAND, line + 1}, nil
	}

	return index, types.Token{}, errors.New(fmt.Sprintf("DfCommand'index: %d find invalid token in \"%s\".", index, lineStr))
}

func readDfArgs(index int, lineStr string, line int) (int, types.Token, error) {
	start := index

	if lineStr[index : index + 2] == "${" {
		index += 2
		for index < len(lineStr) {
			if lineStr[index] == '}' {
				break
			}
			index++
		}
		if index != len(lineStr) {
			return index + 1, types.Token{lineStr[start + 2: index], types.SASSIGNVARIABLE, line + 1}, nil
		} else {
			return index, types.Token{}, errors.New(fmt.Sprintf("Variable in Dfarg: %d find invalid token in \"%s\".", index, lineStr))
		}
	} else {
		for index < len(lineStr) - 1 {
			if lineStr[index : index + 2] == "${" {
				break
			}
			index++
		}
		if index == len(lineStr) - 1 {
			return len(lineStr), types.Token{lineStr[start : len(lineStr)], types.SDFARG, line + 1}, nil
		} else {
			return index, types.Token{lineStr[start : index], types.SDFARG, line + 1}, nil
		}
	}
}

func readString(index int, lineStr string, line int) (int, types.Token, error) {
	start := index
	for index < len(lineStr) {
		index++
		if (lineStr[index] == '"') {
			return index + 1, types.Token{lineStr[start+1 : index], types.SSTRING, line + 1}, nil
		}
	}
	return index, types.Token{}, errors.New(fmt.Sprintf("String's index: %d find invalid token in \"%s\".", index, lineStr))
}

func readIdentifier(index int, lineStr string, line int) (int, types.Token, error) {
	start := index
	if ('a' <= lineStr[index] && lineStr[index] <= 'z') {
		for index < len(lineStr) {
			index++
			if (lineStr[index] < 'A' || ('Z' < lineStr[index] && lineStr[index] < 'a') || 'z' < lineStr[index]) {
				return index, types.Token{lineStr[start : index], types.SIDENTIFIER, line + 1}, nil
			}
		}
	}
	return index, types.Token{}, errors.New(fmt.Sprintf("Others'index: %d find invalid token in \"%s\".", index, lineStr))
}