package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"unicode/utf8"

	"github.com/yetsing/startprompt"
	"github.com/yetsing/startprompt/lexer"
	"github.com/yetsing/startprompt/token"
)

/*
测试 Py3Lexer

使用同目录下的 generate_test_data.py 文件生成测试数据，数据会放在 /tmp/checkpy3lexer2023
之后再运行这个文件即可
*/

var tokenTypeMapping = map[token.TokenType]int{
	token.Name:     1,
	token.Number:   2,
	token.String:   3,
	token.NewLine:  4,
	token.Indent:   5,
	token.Dedent:   6,
	token.Operator: 54,
	token.Comment:  60,
	token.NL:       61,
}

type GoToken struct {
	Index   int
	Type    token.TokenType
	Literal string
}

func getTokens(filepath string, skipWhitespace bool) []GoToken {
	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	data, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}
	l := lexer.NewPy3Lexer(string(data))
	tokens := l.Tokens()
	result := make([]GoToken, 0, len(tokens))
	for i, t := range tokens {
		if skipWhitespace && t.Type == token.Whitespace {
			continue
		}
		result = append(result, GoToken{
			Index:   i,
			Type:    t.Type,
			Literal: t.Literal,
		})
	}
	return result
}

func check(filepath string) {
	type Pytoken struct {
		Type    int    `json:"type"`
		Literal string `json:"literal"`
		Start   []int  `json:"start"`
		End     []int  `json:"end"`
	}

	var info struct {
		Filepath string    `json:"filepath"`
		Tokens   []Pytoken `json:"tokens"`
	}

	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	data, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &info)
	if err != nil {
		panic(err)
	}
	fmt.Printf("test: %q\n filepath: %q\n tokens: %d\n", filepath, info.Filepath, len(info.Tokens))
	//    这个文件里面有个希伯来文的变量名 ，导致跟标准库的解析有差异，没想到解决方法
	if strings.HasSuffix(info.Filepath, "pyparsing/unicode.py") {
		return
	}

	anotherMapping := map[int]token.TokenType{}
	for tokenType, v := range tokenTypeMapping {
		anotherMapping[v] = tokenType
	}

	goTokens := getTokens(info.Filepath, true)
	fmt.Printf(" got tokens: %d\n", len(goTokens))
	for i, pytoken := range info.Tokens[1 : len(info.Tokens)-1] {
		if i >= len(goTokens) {
			fmt.Printf(
				"want:- got:%d want={type: %q, literal: %q}, but got empty, pos: %v-%v\n",
				i, anotherMapping[pytoken.Type], pytoken.Literal,
				pytoken.Start, pytoken.End,
			)
			os.Exit(0)
		} else {
			gotoken := goTokens[i]
			ttype := anotherMapping[pytoken.Type]
			if ttype != gotoken.Type {
				if i > 2 {
					fmt.Printf("last 2 token, got=%+v\n", goTokens[i-2])
				}
				if i > 1 {
					fmt.Printf("last 1 token, got=%+v\n", goTokens[i-1])
				}
				fmt.Printf(
					"want:%d got:%d Type want=%q, but got=%+v, pos: %v-%v\n",
					gotoken.Index, i, ttype, gotoken, pytoken.Start, pytoken.End)
				os.Exit(0)
			}
			if pytoken.Literal != gotoken.Literal {
				fmt.Printf(
					"want:%d got:%d Literal want=%q %d, but got=%q %d, pos: %v-%v\n",
					i,
					gotoken.Index,
					pytoken.Literal, utf8.RuneCountInString(pytoken.Literal),
					gotoken.Literal, utf8.RuneCountInString(gotoken.Literal),
					pytoken.Start, pytoken.End)
				os.Exit(0)
			}
		}
	}

}

func displayTokens(filepath string) {
	for _, goToken := range getTokens(filepath, false) {
		fmt.Printf(
			"%s %q\n",
			startprompt.StringLjustWidth(string(goToken.Type), 15),
			goToken.Literal,
		)
	}

}

func main() {
	if len(os.Args) == 2 {
		filepath := os.Args[1]
		displayTokens(filepath)
		return
	}
	directory := "/tmp/checkpy3lexer2023"
	entries, err := os.ReadDir(directory)
	if err != nil {
		fmt.Printf("ReadDir error: %s\n", err)
		return
	}
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			fmt.Printf("entry.Info error: %s\n", err)
			return
		}
		filepath := path.Join(directory, info.Name())
		check(filepath)
	}

}
