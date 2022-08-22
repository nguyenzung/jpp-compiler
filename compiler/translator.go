package compiler

type Compiler struct {
	vocabulary *Vocabulary
	tokenizer  *Tokenizer
}

func (compiler *Compiler) CompileFile(fileName string) {
	compiler.tokenizer = MakeTokenizer(fileName, compiler.vocabulary)
	compiler.tokenizer.parse()
}

func MakeCompiler() *Compiler {
	vocabulary := MakeVocabulary()
	compiler := &Compiler{vocabulary: vocabulary}
	return compiler
}
