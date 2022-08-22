package compiler

type Compiler struct {
	vocabulary *Vocabulary
	tokenizer  *Tokenizer
}

func (compiler *Compiler) CompileFile(fileName string) {
	compiler.tokenizer = MakeTokenizer(fileName, compiler.vocabulary)
	compiler.tokenizer.parse()
	for _, token := range compiler.tokenizer.tokens {
		token.print()
	}
}

func MakeCompiler() *Compiler {
	vocabulary := MakeVocabulary()
	compiler := &Compiler{vocabulary: vocabulary}
	return compiler
}
