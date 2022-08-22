package compiler

type AST struct {
	tokens       []*Token
	currentIndex int
}

func (tree *AST) hasMoreToken() bool {
	return len(tree.tokens) > tree.currentIndex
}

func (tree *AST) advance() *Token {
	tree.currentIndex += 1
	return tree.tokens[tree.currentIndex-1]
}

func (tree *AST) buildAST() {

}

func (tree *AST) parseClass() bool {

	return false
}

func (tree *AST) parseClassVarDec() bool {
	return false
}

func (tree *AST) parseSubroutineDec() bool {
	return false
}

func MakeParseTree(tokens []*Token) *AST {
	return &AST{tokens: tokens, currentIndex: 0}
}
