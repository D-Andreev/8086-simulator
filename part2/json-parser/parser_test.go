package jsonparser

import (
	"testing"
)

type ParserTestCase struct {
	input        []Token
	expectedNode *Node
}

func TestParser(t *testing.T) {
	testCases := []ParserTestCase{
		{
			input:        []Token{},
			expectedNode: nil,
		},
		{
			input: []Token{
				{Type: OPEN_BRACE, Literal: "{"},
				{Type: STRING, Literal: "name"},
				{Type: STRING, Literal: "john"},
				{Type: STRING, Literal: "age"},
				{Type: NUMBER, Literal: "30"},
				{Type: CLOSE_BRACE, Literal: "}"},
			},
			expectedNode: &Node{
				Type: NODE_OBJECT, Value: nil, Children: []*Node{
					{Type: NODE_STRING, Value: "name"},
					{Type: NODE_STRING, Value: "john"},
					{Type: NODE_STRING, Value: "age"},
					{Type: NODE_NUMBER, Value: "30"},
				}}},
		{
			input: []Token{
				{Type: OPEN_BRACKET, Literal: "["},
				{Type: STRING, Literal: "jack"},
				{Type: STRING, Literal: "john"},
				{Type: CLOSE_BRACKET, Literal: "]"},
			},
			expectedNode: &Node{
				Type: NODE_ARRAY, Value: nil, Children: []*Node{
					{Type: NODE_STRING, Value: "jack"},
					{Type: NODE_STRING, Value: "john"},
				}},
		},
		// nested objects
		{
			input: []Token{
				{Type: OPEN_BRACE, Literal: "{"},
				{Type: STRING, Literal: "id"},
				{Type: NUMBER, Literal: "1"},
				{Type: STRING, Literal: "address"},
				{Type: OPEN_BRACE, Literal: "{"},
				{Type: STRING, Literal: "street"},
				{Type: STRING, Literal: "main"},
				{Type: STRING, Literal: "postcode"},
				{Type: NUMBER, Literal: "123"},
				{Type: CLOSE_BRACE, Literal: "}"},
				{Type: CLOSE_BRACE, Literal: "}"},
			},
			expectedNode: &Node{
				Type: NODE_OBJECT, Value: nil, Children: []*Node{
					{Type: NODE_STRING, Value: "id"},
					{Type: NODE_NUMBER, Value: "1"},
					{Type: NODE_STRING, Value: "address"},
					{Type: NODE_OBJECT, Value: nil, Children: []*Node{
						{Type: NODE_STRING, Value: "street"},
						{Type: NODE_STRING, Value: "main"},
						{Type: NODE_STRING, Value: "postcode"},
						{Type: NODE_NUMBER, Value: "123"},
					}},
				}},
		},
		// array of objects
		{
			input: []Token{
				{Type: OPEN_BRACKET, Literal: "["},
				{Type: OPEN_BRACE, Literal: "{"},
				{Type: STRING, Literal: "id"},
				{Type: NUMBER, Literal: "1"},
				{Type: CLOSE_BRACE, Literal: "}"},
				{Type: OPEN_BRACE, Literal: "{"},
				{Type: STRING, Literal: "id"},
				{Type: NUMBER, Literal: "2"},
				{Type: CLOSE_BRACE, Literal: "}"},
				{Type: CLOSE_BRACKET, Literal: "]"},
			},
			expectedNode: &Node{
				Type: NODE_ARRAY, Value: nil, Children: []*Node{
					{Type: NODE_OBJECT, Value: nil, Children: []*Node{
						{Type: NODE_STRING, Value: "id"},
						{Type: NODE_NUMBER, Value: "1"},
					}},
					{Type: NODE_OBJECT, Value: nil, Children: []*Node{
						{Type: NODE_STRING, Value: "id"},
						{Type: NODE_NUMBER, Value: "2"},
					}},
				}},
		},
		// 3 levels nested objects
		{
			input: []Token{
				{Type: OPEN_BRACE, Literal: "{"},
				{Type: STRING, Literal: "id"},
				{Type: NUMBER, Literal: "1"},
				{Type: STRING, Literal: "address"},
				{Type: OPEN_BRACE, Literal: "{"},
				{Type: STRING, Literal: "street"},
				{Type: STRING, Literal: "main"},
				{Type: STRING, Literal: "block"},
				{Type: OPEN_BRACE, Literal: "{"},
				{Type: STRING, Literal: "n"},
				{Type: NUMBER, Literal: "1"},
				{Type: STRING, Literal: "ap"},
				{Type: NUMBER, Literal: "27"},
				{Type: CLOSE_BRACE, Literal: "}"},
				{Type: CLOSE_BRACE, Literal: "}"},
				{Type: CLOSE_BRACE, Literal: "}"},
			},
			expectedNode: &Node{
				Type: NODE_OBJECT, Value: nil, Children: []*Node{
					{Type: NODE_STRING, Value: "id"},
					{Type: NODE_NUMBER, Value: "1"},
					{Type: NODE_STRING, Value: "address"},
					{Type: NODE_OBJECT, Value: nil, Children: []*Node{
						{Type: NODE_STRING, Value: "street"},
						{Type: NODE_STRING, Value: "main"},
						{Type: NODE_STRING, Value: "block"},
						{Type: NODE_OBJECT, Value: nil, Children: []*Node{
							{Type: NODE_STRING, Value: "n"},
							{Type: NODE_NUMBER, Value: "1"},
							{Type: NODE_STRING, Value: "ap"},
							{Type: NODE_NUMBER, Value: "27"},
						}},
					}},
				}},
		},
	}

	for _, testCase := range testCases {
		parser := NewParser(testCase.input)
		node := parser.Parse()

		compareNodes(t, testCase, node, testCase.expectedNode)
	}
}

func compareNodes(t *testing.T, testCase ParserTestCase, node *Node, expectedNode *Node) {
	t.Helper()

	if expectedNode == nil || node == nil {
		if expectedNode != node {
			t.Errorf("Expected node to be %v, got: %v in test case: %s", expectedNode, node, testCase.input)
		}
		return
	}

	if node.Type != expectedNode.Type {
		t.Errorf("Expected node type: %v, got: %v in test case: %s", expectedNode.Type, node.Type, testCase.input)
	}
	if node.Value != expectedNode.Value {
		t.Errorf("Expected node value: %v, got: %v in test case: %s", expectedNode.Value, node.Value, testCase.input)
	}
	if node.Children != nil {
		if len(node.Children) != len(expectedNode.Children) {
			t.Errorf("Expected %d children, got %d in test case: %s", len(expectedNode.Children), len(node.Children), testCase.input)
			return
		}
		for i := range node.Children {
			compareNodes(t, testCase, node.Children[i], expectedNode.Children[i])
		}
	}
}
