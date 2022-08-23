package compiler

import (
	"fmt"
)

type ASTProcessor interface {
	process(token *Token)
	isFinished() bool
	Iterate()
}

type AST struct {
	label         string
	finished      bool
	astProcessors []ASTProcessor
}

func (ast *AST) isFinished() bool {
	return ast.finished
}

func (ast *AST) isChildProcessing() bool {
	if len(ast.astProcessors) == 0 {
		return false
	} else {
		astProcessor := ast.astProcessors[len(ast.astProcessors)-1]
		return astProcessor.isFinished()
	}
}

func (ast *AST) isClassWord() {
	
}

type ClassAST struct {
	AST
	className string
}

func (classAST *ClassAST) Iterate() {
	fmt.Println("<class>")
	fmt.Println("<keyword>class</keyword>")
	fmt.Println("<identifier>", classAST.className, "</identifier>")
	fmt.Println("<symbol> { </symbol>")
	for _, ast := range classAST.astProcessors {
		ast.Iterate()
	}
	fmt.Println("<symbol> } </symbol>")
	fmt.Println("</class>")
}

func (classAST *ClassAST) process(token *Token) {

}

func MakeClassAST() *ClassAST {
	return &ClassAST{}
}

type FileAST struct {
	AST
	tokens       []*Token
	currentIndex int
}

func (ast *FileAST) Iterate() {

}

func (parser *FileAST) buildAST() {
	for parser.hasMoreToken() {
		parser.advance()
	}
}

func (parser *FileAST) hasMoreToken() bool {
	return len(parser.tokens) > parser.currentIndex
}

func (parser *FileAST) advance() {
	parser.process(parser.tokens[parser.currentIndex])
	parser.currentIndex += 1
}

func (parser *FileAST) process(token *Token) {
	fmt.Println("[Process]", token)

}

func MakeFileAST(tokens []*Token) *FileAST {
	return &FileAST{tokens: tokens, currentIndex: 0}
}
