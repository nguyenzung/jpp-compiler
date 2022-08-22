package main

import (
	"fmt"

	cp "github.com/nguyenzung/jack-compiler/compiler"
)

func main() {
	fmt.Println("Jack Compiler")

	compiler := cp.MakeCompiler()
	compiler.CompileFile("jackfiles/ArrayTest/Main.jack")
}
