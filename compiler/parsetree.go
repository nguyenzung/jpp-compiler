package compiler

import (
	"fmt"
)

type IModule interface {
	// verify(token *Token) bool
	// check(token []*Token) bool
	process(token []*Token)
	IsFinished() bool
	getStartNode() *Node
	getEndNode() *Node
	getCurrentNode() *Node
}

type BaseModule struct {
	startNode   *Node
	endNode     *Node
	currentNode *Node
}

func (module *BaseModule) addTerminalChar() {
	module.addNode(MakeLabelNode([]string{";"}, []string{SYMBOL}))
}

func (module *BaseModule) addNode(node *Node) {
	if module.startNode == nil {
		module.startNode = node
		module.currentNode = node
		module.endNode = node
	} else {
		module.endNode.addNode(node)
		module.endNode = node
	}
}

func (module *BaseModule) addModule(nextModule IModule) {
	startModuleNode := nextModule.getStartNode()
	endModuleNode := nextModule.getEndNode()
	if module.startNode == nil {
		module.startNode = startModuleNode
		module.currentNode = startModuleNode
		module.endNode = startModuleNode
	} else {
		module.endNode.addNode(startModuleNode)
	}
	module.endNode = endModuleNode
}

func (module *BaseModule) IsFinished() bool {
	return module.currentNode == module.endNode
}

func (module *BaseModule) process(tokens []*Token) {
	fmt.Println("[Process]", tokens[0], module.currentNode)
	node, err := module.currentNode.Process(tokens)
	if err == nil {
		module.currentNode = node
	} else {
		panic(fmt.Sprintln("Cannot process token ", err))
	}
}

func (module *BaseModule) getStartNode() *Node {
	return module.startNode
}

func (module *BaseModule) getEndNode() *Node {
	return module.endNode
}

func (module *BaseModule) getCurrentNode() *Node {
	return module.currentNode
}

func MakeBaseModule() *BaseModule {
	return &BaseModule{startNode: nil, endNode: nil, currentNode: nil}
}

// type OrModule struct {
// 	modules []IModule
// }

// type AndModule struct {
// 	modules []IModule
// }

// type AtleastOneModule struct {
// 	modules IModule
// }

type StarModule struct {
	*BaseModule
	module IModule
}

// func MakeStarModule(modules IModule) *StarModule {
// 	module := &StarModule{BaseModule: MakeBaseModule(), module: }
// }

type ClassVarDec struct {
	*BaseModule
}

func MakeClassVarDec() *ClassVarDec {
	module := &ClassVarDec{BaseModule: MakeBaseModule()}
	module.addNode(MakeLabelNode([]string{"static", "field"}, []string{KEYWORD, KEYWORD}))
	module.addNode(MakeLabelNode([]string{"int", "char", "boolean"}, []string{KEYWORD, KEYWORD, KEYWORD}))
	module.addNode(MakeIdentifierNode())
	module.addModule(MakeAdditionalVarDec())
	module.addTerminalChar()
	return module
}

type AdditionalVarDec struct {
	*BaseModule
}

func MakeAdditionalVarDec() *AdditionalVarDec {
	module := &AdditionalVarDec{BaseModule: MakeBaseModule()}
	module.addNode(MakeLabelNode([]string{","}, []string{SYMBOL}))
	module.addNode(MakeIdentifierNode())
	return module
}

type ClassModule struct {
	*BaseModule
}

func MakeClassModule() *ClassModule {
	module := &ClassModule{BaseModule: MakeBaseModule()}
	module.addNode(MakeDefaultNode())
	module.addNode(MakeLabelNode([]string{"class"}, []string{KEYWORD}))
	module.addNode(MakeIdentifierNode())
	module.addNode(MakeLabelNode([]string{"{"}, []string{SYMBOL}))
	// module.addModule(MakeClassVarDec())
	module.addNode(MakeLabelNode([]string{"}"}, []string{SYMBOL}))
	fmt.Println("[Module]")
	return module
}

type AST struct {
	module *ClassModule
}

func (ast *AST) Parse(tokens []*Token) {
	for _, token := range tokens {
		fmt.Println(token.token, token.tag)
	}
	for i := range tokens {
		if i+1 < len(tokens) {
			ast.module.process(tokens[i : i+2])
		} else {
			ast.module.process(tokens[i : i+1])
		}
	}
}

func (ast *AST) IsFinish() bool {
	return ast.module.IsFinished()
}

func MakeAST() *AST {
	return &AST{module: MakeClassModule()}
}
