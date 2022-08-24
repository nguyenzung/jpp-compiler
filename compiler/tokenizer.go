package compiler

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"
)

const (
	KEYWORD    string = "keyword"
	SYMBOL     string = "symbol"
	NUMBER     string = "number"
	STRING     string = "string"
	IDENTIFIER string = "identifier"
	UNKNOWN    string = ""
)

type Vocabulary struct {
	keywords  []string
	symbols   []string
	keyTokens map[string]string
}

func (vocabulary *Vocabulary) init() {
	vocabulary.keywords = []string{"class", "constructor", "function", "method", "field", "static", "var", "int", "char", "boolean", "void", "do", "true", "false", "null", "this", "let", "if", "else", "while", "return"}
	vocabulary.symbols = []string{"{", "}", "(", ")", "[", "]", ".", ",", ";", "+", "-", "*", "/", "&", "|", "<", ">", "=", "~"}
	vocabulary.keyTokens = make(map[string]string)
	for _, keyword := range vocabulary.keywords {
		vocabulary.keyTokens[keyword] = KEYWORD
	}
	for _, symbol := range vocabulary.symbols {
		vocabulary.keyTokens[symbol] = SYMBOL
	}
	// fmt.Println(vocabulary.keyTokens)
}

func MakeVocabulary() *Vocabulary {
	vocabulary := &Vocabulary{}
	vocabulary.init()
	return vocabulary
}

type Token struct {
	token string
	tag   string
}

func (token *Token) print() {
	fmt.Println("[", token.tag, "]", token.token)
}

func MakeCodeItem(token string, tag string) *Token {
	return &Token{token: token, tag: tag}
}

type Tokenizer struct {
	vocabulary   *Vocabulary
	fileName     string
	currentToken *bytes.Buffer
	tokens       []*Token
}

func (tokenizer *Tokenizer) processNewLine(line string) {
	for i := range line {
		if string(line[i]) != " " {
			tokenizer.processCharacter(line[i])
		}
	}
}

func (tokenizer *Tokenizer) processCharacter(c byte) {
	if val, ok := tokenizer.vocabulary.keyTokens[string(c)]; ok && val == SYMBOL {
		tokenizer.processCurrentToken(tokenizer.currentToken.String(), getTag(tokenizer.currentToken.String()))
		tokenizer.processCurrentToken(string(c), SYMBOL)
		tokenizer.currentToken.Reset()
	} else {
		tokenizer.currentToken.WriteByte(c)
		val, ok := tokenizer.vocabulary.keyTokens[tokenizer.currentToken.String()]
		if ok {
			if val == KEYWORD {
				tokenizer.processCurrentToken(tokenizer.currentToken.String(), KEYWORD)
				tokenizer.currentToken.Reset()
			}
		}
	}
}

func (tokenizer *Tokenizer) processCurrentToken(token string, tag string) {
	if len(token) > 0 && tag != UNKNOWN {
		tokenizer.tagToken(token, tag)
	}
}

func getTag(token string) string {
	if len(token) > 0 {
		if isNumber(token) {
			return NUMBER
		} else if isString(token) {
			return STRING
		} else if isIdentifier(token) {
			return IDENTIFIER
		}
	}
	return UNKNOWN
}

var numberRegex, _ = regexp.Compile("^[0-9]+$")

func isNumber(token string) bool {
	return numberRegex.MatchString(token)
}

var stringRegex, _ = regexp.Compile("^\"[^\n]+\"$")

func isString(token string) bool {
	return stringRegex.MatchString(token)
}

var idenRegex, _ = regexp.Compile("^([a-zA-Z_]|[a-zA-Z0-9_])+$")

func isIdentifier(token string) bool {
	return idenRegex.MatchString(token)
}

func (tokenizer *Tokenizer) tagToken(token string, tag string) {
	fmt.Println("[PROCESS TOKENIZER]", token, tag)
	tokenizer.tokens = append(tokenizer.tokens, MakeCodeItem(token, tag))
}

func (tokenizer *Tokenizer) parse() []*Token {
	file, err := os.Open(tokenizer.fileName)
	if err != nil {
		panic((err))
	}
	defer file.Close()
	s := bufio.NewScanner(file)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if len(line) >= 3 {
			if line[0:3] == "/**" {
				continue
			}
		}
		if commentIndex := strings.Index(line, "//"); commentIndex != -1 {
			line = line[0:commentIndex]
		}
		line = strings.TrimSuffix(line, "//")
		if len(line) > 0 {
			tokenizer.processNewLine(line)
		}
	}
	return tokenizer.tokens
}

func MakeTokenizer(fileName string, vocabulary *Vocabulary) *Tokenizer {
	return &Tokenizer{fileName: fileName, vocabulary: vocabulary, currentToken: &bytes.Buffer{}, tokens: make([]*Token, 0)}
}
