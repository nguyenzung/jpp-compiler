package compiler

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
)

const (
	KEYWORD string = "keyword"
	SYMBOL  string = "symbol"
)

type Vocabulary struct {
	keywords  []string
	symbols   []string
	keyTokens map[string]string
}

func (vocabulary *Vocabulary) init() {
	vocabulary.keywords = []string{"class", "constructor", "function", "method", "field", "static", "var", "int", "char", "boolean", "void", "do", "true", "false", "null", "this", "let", "if", "else", "white", "return"}
	// fmt.Println("[WORD]", len(token.keywords))
	vocabulary.symbols = []string{"{", "}", "(", ")", "[", "]", ".", ",", ";", "+", "-", "*", "/", "&", "|", "<", ">", "=", "~"}
	// fmt.Println("[SYM]", len(token.symbols))
	vocabulary.keyTokens = make(map[string]string)
	for _, keyword := range vocabulary.keywords {
		vocabulary.keyTokens[keyword] = KEYWORD
	}
	for _, symbol := range vocabulary.symbols {
		vocabulary.keyTokens[symbol] = SYMBOL
	}
	fmt.Println(vocabulary.keyTokens)
}

func MakeVocabulary() *Vocabulary {
	vocabulary := &Vocabulary{}
	vocabulary.init()
	return vocabulary
}

type Tokenizer struct {
	vocabulary   *Vocabulary
	fileName     string
	currentToken *bytes.Buffer
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
		tokenizer.processCurrentToken(tokenizer.currentToken.String())
		tokenizer.processCurrentToken(string(c))
		tokenizer.currentToken.Reset()
	} else {
		tokenizer.currentToken.WriteByte(c)
		val, ok := tokenizer.vocabulary.keyTokens[tokenizer.currentToken.String()]
		if ok {
			if val == KEYWORD {
				tokenizer.processCurrentToken(tokenizer.currentToken.String())
				tokenizer.currentToken.Reset()
			}
		}
	}

}

func (tokenizer *Tokenizer) processCurrentToken(candidate string) {
	if len(candidate) > 0 {
		fmt.Println("[CURR Token]", candidate)
	}

}

func (tokenizer *Tokenizer) tagToken(token string, tag string) {

}

func (tokenizer *Tokenizer) parse() {
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
}

func MakeTokenizer(fileName string, vocabulary *Vocabulary) *Tokenizer {
	return &Tokenizer{fileName: fileName, vocabulary: vocabulary, currentToken: &bytes.Buffer{}}
}
