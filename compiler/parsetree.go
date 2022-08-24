package compiler

import (
	"fmt"
)

type INode interface {
	verify(token *Token) error // Verify valid transite from previous node
	setup(token *Token)
}

type Node struct {
	item INode
	next *Node
	prev *Node
}

func (node *Node) addNode(next *Node) {
	node.next = next
	next.prev = node
}

func (node *Node) setup(token *Token) {
	node.item.setup(token)
}

func (node *Node) check(tokens []*Token) error {
	if len(tokens) == 0 {
		return nil
	} else {
		if err := node.next.item.verify(tokens[0]); err == nil {
			return node.next.check(tokens[1:])
		} else {
			return err
		}
	}
}

func (node *Node) transite(token *Token) *Node {
	next := node.next
	if next != nil {
		next.setup(token)
	}
	return node.next
}

func (node *Node) Process(tokens []*Token) (*Node, error) {
	err := node.check(tokens)
	if err == nil {
		return node.transite(tokens[0]), nil
	}
	return nil, err
}

type DefaultNode struct{}

func (node *DefaultNode) verify(token *Token) error {
	return nil
}

func (node *DefaultNode) setup(token *Token) {
}

func MakeDefaultNode() *Node {
	return &Node{item: &DefaultNode{}}
}

type StartNode DefaultNode

type EndNode DefaultNode

type LabelNode struct {
	labels []string
	tags   []string
	result *Token
}

func (node *LabelNode) verify(token *Token) error {
	// fmt.Println("[Label]", token.token, token.tag)
	for i := range node.labels {
		if node.labels[i] == token.token && node.tags[i] == token.tag {
			return nil
		}
	}
	return fmt.Errorf("[ERROR] %s", token.token)
}

func (node *LabelNode) setup(token *Token) {
	node.result = token
}

func MakeLabelNode(labels []string, tags []string) *Node {
	dataNode := &LabelNode{labels: labels, tags: tags}
	return &Node{item: dataNode}
}

type IdentifierNode struct {
	result *Token
}

func (node *IdentifierNode) verify(token *Token) error {
	if token.tag == IDENTIFIER {
		return nil
	} else {
		return fmt.Errorf("[ERROR] Expected IDENTIFIER instead of %s. Token: %s", token.tag, token.token)
	}
}

func (node *IdentifierNode) setup(token *Token) {
	node.result = token
}

func MakeIdentifierNode() *Node {
	dataNode := &IdentifierNode{}
	return &Node{item: dataNode}
}

type IModule interface {
	verify(token *Token) bool
	check(token []*Token) bool
	process(token *Token) bool
	finish()
	getStartNode() *Node
	getEndNode() *Node
	getCurrentNode() *Node
}

type BaseModule struct {
	startNode   *Node
	endNode     *Node
	currentNode *Node
}

func (module *BaseModule) addLinearNode(node *Node) {
	currentBuiltNode := module.endNode.prev
	currentBuiltNode.addNode(node)
	node.addNode(module.endNode)
}

func (module *BaseModule) IsFinished() bool {
	return module.currentNode == module.endNode
}

// func (module *BaseModule) getStartNode() *Node {
// 	return module.startNode
// }

// func (module *BaseModule) getEndNode() *Node {
// 	return module.endNode
// }

// func (module *BaseModule) getCurrentNode() *Node {
// 	return module.currentNode
// }

func (module *BaseModule) process(tokens []*Token) {
	fmt.Println("[Process]", tokens[0])
	node, err := module.currentNode.Process(tokens)
	if err == nil {
		module.currentNode = node
		if node.next == module.endNode {
			module.currentNode = node.next
		}
		fmt.Println("END", module.IsFinished())
	} else {
		panic(fmt.Sprintln("Cannot process token ", err))
	}
}

func MakeBaseModule() *BaseModule {
	startNode := MakeDefaultNode()
	endNode := MakeDefaultNode()
	startNode.addNode(endNode)
	return &BaseModule{startNode: startNode, endNode: endNode, currentNode: startNode}
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

// type AtleastZero struct {
// 	modules IModule
// }

type AdditionalVarDec struct {
	*BaseModule
}

func MakeAdditionalVarDec() *AdditionalVarDec {
	module := &AdditionalVarDec{BaseModule: MakeBaseModule()}
	module.addLinearNode(MakeLabelNode([]string{","}, []string{SYMBOL}))
	module.addLinearNode(MakeIdentifierNode())
	return module
}

type ClassModule struct {
	*BaseModule
}

func MakeClassModule() *ClassModule {
	module := &ClassModule{BaseModule: MakeBaseModule()}
	module.addLinearNode(MakeLabelNode([]string{"class"}, []string{KEYWORD}))
	module.addLinearNode(MakeIdentifierNode())
	module.addLinearNode(MakeLabelNode([]string{"{"}, []string{SYMBOL}))
	module.addLinearNode(MakeLabelNode([]string{"}"}, []string{SYMBOL}))
	fmt.Println("[Module]")
	return module
}

type AST struct {
	module *ClassModule
	vardec *AdditionalVarDec
}

func (ast *AST) Parse(tokens []*Token) {
	for _, token := range tokens {
		fmt.Println(token.token, token.tag)
	}
	for i := range tokens {
		if i+1 < len(tokens) {
			ast.vardec.process(tokens[i : i+2])
		} else {
			ast.vardec.process(tokens[i : i+1])
		}
	}
}

func (ast *AST) IsFinish() bool {
	return ast.module.IsFinished()
}

func MakeAST() *AST {
	return &AST{module: MakeClassModule(), vardec: MakeAdditionalVarDec()}
}
