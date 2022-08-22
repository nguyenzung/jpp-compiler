package compiler

type ParseTree struct {
	tokens []*Token
}

func (tree *ParseTree) parseClass() {

}

func MakeParseTree(tokens []*Token) *ParseTree {
	return &ParseTree{tokens: tokens}
}
