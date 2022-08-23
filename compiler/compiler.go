package compiler

type Compiler struct {
	vocabulary *Vocabulary
}

func (compiler *Compiler) CompileFile(fileName string) {
	tokenizer := MakeTokenizer(fileName, compiler.vocabulary)
	tokens := tokenizer.parse()
	for _, token := range tokens {
		token.print()
	}
	ast := MakeFileAST(tokens)
	ast.buildAST()
}

func MakeCompiler() *Compiler {
	vocabulary := MakeVocabulary()
	compiler := &Compiler{vocabulary: vocabulary}
	return compiler
}
