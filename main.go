package main

import (
	cp "github.com/nguyenzung/jack-compiler/compiler"
)

func main() {
	compiler := cp.MakeCompiler()
	compiler.CompileFile("jackfiles/ArrayTest/VarDec.jpp")
}
