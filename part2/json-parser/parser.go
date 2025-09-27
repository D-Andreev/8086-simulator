package jsonparser

import "fmt"

type NodeType string

const (
	NODE_OBJECT NodeType = "object"
	NODE_ARRAY  NodeType = "array"
	NODE_STRING NodeType = "string"
	NODE_NUMBER NodeType = "number"
	NODE_BOOL   NodeType = "bool"
	NODE_NULL   NodeType = "null"
)

type Node struct {
	Type     NodeType
	Value    any
	Children []*Node
}

func NewNode(nodeType NodeType, value any, children []*Node) *Node {
	return &Node{Type: nodeType, Value: value, Children: children}
}

type Parser struct {
	position int
	tokens   []Token
}

func NewParser(tokens []Token) *Parser {
	return &Parser{tokens: tokens}
}

func (p *Parser) Parse() *Node {
	if len(p.tokens) == 0 {
		return nil
	}

	token := p.tokens[p.position]
	switch token.Type {
	case OPEN_BRACE:
		return p.parseObject()
	case OPEN_BRACKET:
		return p.parseArray()
	case STRING:
		p.position++
		return NewNode(NODE_STRING, token.Literal, nil)
	case NUMBER:
		p.position++
		return NewNode(NODE_NUMBER, token.Literal, nil)
	case BOOL:
		p.position++
		return NewNode(NODE_BOOL, token.Literal, nil)
	case NULL:
		p.position++
		return NewNode(NODE_NULL, token.Literal, nil)
	default:
		fmt.Println("Unexpected token: ", token.Type)
		return nil
	}
}

func (p *Parser) parseObject() *Node {
	n := NewNode(NODE_OBJECT, nil, []*Node{})
	p.position++
	for p.tokens[p.position].Type != CLOSE_BRACE {
		n.Children = append(n.Children, p.Parse())
	}
	p.position++
	return n
}

func (p *Parser) parseArray() *Node {
	n := NewNode(NODE_ARRAY, nil, []*Node{})
	p.position++
	for p.tokens[p.position].Type != CLOSE_BRACKET {
		n.Children = append(n.Children, p.Parse())
	}
	p.position++
	return n
}
