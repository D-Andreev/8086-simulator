package jsonparser

import (
	"fmt"
	"strconv"
)

type NodeType string

const (
	NODE_OBJECT    NodeType = "object"
	NODE_ARRAY     NodeType = "array"
	NODE_STRING    NodeType = "string"
	NODE_NUMBER    NodeType = "number"
	NODE_BOOL      NodeType = "bool"
	NODE_NULL      NodeType = "null"
	NODE_KEY_VALUE NodeType = "key_value"
)

type Node struct {
	Type     NodeType
	Value    any
	Children []*Node
	Key      *Node // Only used for NODE_KEY_VALUE
	Val      *Node // Only used for NODE_KEY_VALUE
}

func NewNode(nodeType NodeType, value any, children []*Node) *Node {
	switch nodeType {
	case NODE_BOOL:
		value = value.(bool)
	case NODE_NUMBER:
		value, _ = strconv.ParseFloat(value.(string), 64)
	case NODE_NULL:
		value = nil
	}

	return &Node{Type: nodeType, Value: value, Children: children}
}

func NewKeyValueNode(key, val *Node) *Node {
	return &Node{Type: NODE_KEY_VALUE, Key: key, Val: val}
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
	p.position++ // consume opening brace

	for p.position < len(p.tokens) && p.tokens[p.position].Type != CLOSE_BRACE {
		// Parse key-value pair
		if p.tokens[p.position].Type == STRING {
			key := p.Parse() // Parse the key (string)
			if p.position < len(p.tokens) && p.tokens[p.position].Type == COLON {
				p.position++     // consume colon
				val := p.Parse() // Parse the value
				kvPair := NewKeyValueNode(key, val)
				n.Children = append(n.Children, kvPair)
			}
		}

		// Skip comma if present
		if p.position < len(p.tokens) && p.tokens[p.position].Type == COMMA {
			p.position++
		}
	}

	if p.position < len(p.tokens) {
		p.position++ // consume closing brace
	}
	return n
}

func (p *Parser) parseArray() *Node {
	n := NewNode(NODE_ARRAY, nil, []*Node{})
	p.position++ // consume opening bracket

	for p.position < len(p.tokens) && p.tokens[p.position].Type != CLOSE_BRACKET {
		element := p.Parse() // Parse array element
		n.Children = append(n.Children, element)

		// Skip comma if present
		if p.position < len(p.tokens) && p.tokens[p.position].Type == COMMA {
			p.position++
		}
	}

	if p.position < len(p.tokens) {
		p.position++ // consume closing bracket
	}
	return n
}
