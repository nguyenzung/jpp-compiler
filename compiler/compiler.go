package compiler

import "fmt"

type Compiler struct {
	vocabulary *Vocabulary
}

func (compiler *Compiler) CompileFile(fileName string) {
	tokenizer := MakeTokenizer(fileName, compiler.vocabulary)
	tokens := tokenizer.parse()
	ast := MakeAST()
	ast.Parse(tokens[:len(tokens)-1])
	fmt.Println(ast.IsFinish())
}

func MakeCompiler() *Compiler {
	vocabulary := MakeVocabulary()
	compiler := &Compiler{vocabulary: vocabulary}
	return compiler
}
