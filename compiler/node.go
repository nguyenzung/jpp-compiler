package compiler

import (
	"fmt"
)

type INode interface {
	verify(token *Token) error // Verify valid transite from previous node
	setup(token *Token)
	kToken() int
}

type Node struct {
	Epsilon bool
	item    INode
	nexts   []*Node
	prevs   []*Node
}

func (node *Node) addNode(next *Node) {
	// node.next = next
	node.nexts = append(node.nexts, next)
	next.prevs = append(next.prevs, node)
}

func (node *Node) setup(token *Token) {
	node.item.setup(token)
}

func (node *Node) verify(tokens []*Token) (*Node, error) {
	if node.Epsilon {
		return node.check(tokens)
	} else {
		if node.item.verify(tokens[0]) == nil {
			return node, nil
		} else {
			return nil, fmt.Errorf("Cannot verify")
		}
	}
}

// Check if can process tokens from this node
func (node *Node) check(tokens []*Token) (*Node, error) {
	tokenLength := len(tokens)
	for _, next := range node.nexts {
		if next != nil && tokenLength > next.item.kToken() {
			tokens = tokens[:node.item.kToken()]
		}
		if len(tokens) == 0 {
			return node, nil
		} else {
			if err := next.item.verify(tokens[0]); err == nil {
				return next.check(tokens[1:])
			} else {
				return nil, err
			}
		}
	}
	return nil, fmt.Errorf("Cannot handle token")
}

// func (node *Node) transite(token *Token, next *Node) *Node {
// 	if next != nil {
// 		next.setup(token)
// 	}
// 	return node.next
// }

func (node *Node) Process(tokens []*Token) (*Node, error) {
	return node.check(tokens)

}

type DefaultNode struct{}

func (node *DefaultNode) verify(token *Token) error {
	return nil
}

func (node *DefaultNode) setup(token *Token) {
}

func (node *DefaultNode) kToken() int {
	return 1
}

func MakeDefaultNode() *Node {
	return &Node{item: &DefaultNode{}, nexts: make([]*Node, 0), prevs: make([]*Node, 0)}
}

type StartNode DefaultNode

type EndNode DefaultNode

/** Handle KEYWORD and SYMBOL */
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

func (node *LabelNode) kToken() int {
	return 1
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

func (node *IdentifierNode) kToken() int {
	return 1
}

func MakeIdentifierNode() *Node {
	dataNode := &IdentifierNode{}
	return &Node{item: dataNode}
}
