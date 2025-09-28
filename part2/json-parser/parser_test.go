package jsonparser

import (
	"testing"
)

func TestParserWithLexer(t *testing.T) {
	testCases := []struct {
		jsonString   string
		expectedNode *Node
	}{
		{
			jsonString:   "",
			expectedNode: nil,
		},
		{
			jsonString: `{"name": "john", "age": 30}`,
			expectedNode: &Node{
				Type: NODE_OBJECT, Value: nil, Children: []*Node{
					{Type: NODE_KEY_VALUE, Key: &Node{Type: NODE_STRING, Value: "name"}, Val: &Node{Type: NODE_STRING, Value: "john"}},
					{Type: NODE_KEY_VALUE, Key: &Node{Type: NODE_STRING, Value: "age"}, Val: &Node{Type: NODE_NUMBER, Value: "30"}},
				}},
		},
		{
			jsonString: `["jack", "john"]`,
			expectedNode: &Node{
				Type: NODE_ARRAY, Value: nil, Children: []*Node{
					{Type: NODE_STRING, Value: "jack"},
					{Type: NODE_STRING, Value: "john"},
				}},
		},
		{
			jsonString: `{"id": 1, "address": {"street": "main", "postcode": 123}}`,
			expectedNode: &Node{
				Type: NODE_OBJECT, Value: nil, Children: []*Node{
					{Type: NODE_KEY_VALUE, Key: &Node{Type: NODE_STRING, Value: "id"}, Val: &Node{Type: NODE_NUMBER, Value: "1"}},
					{Type: NODE_KEY_VALUE, Key: &Node{Type: NODE_STRING, Value: "address"}, Val: &Node{Type: NODE_OBJECT, Value: nil, Children: []*Node{
						{Type: NODE_KEY_VALUE, Key: &Node{Type: NODE_STRING, Value: "street"}, Val: &Node{Type: NODE_STRING, Value: "main"}},
						{Type: NODE_KEY_VALUE, Key: &Node{Type: NODE_STRING, Value: "postcode"}, Val: &Node{Type: NODE_NUMBER, Value: "123"}},
					}}},
				}},
		},
		{
			jsonString: `[{"id": 1}, {"id": 2}]`,
			expectedNode: &Node{
				Type: NODE_ARRAY, Value: nil, Children: []*Node{
					{Type: NODE_OBJECT, Value: nil, Children: []*Node{
						{Type: NODE_KEY_VALUE, Key: &Node{Type: NODE_STRING, Value: "id"}, Val: &Node{Type: NODE_NUMBER, Value: "1"}},
					}},
					{Type: NODE_OBJECT, Value: nil, Children: []*Node{
						{Type: NODE_KEY_VALUE, Key: &Node{Type: NODE_STRING, Value: "id"}, Val: &Node{Type: NODE_NUMBER, Value: "2"}},
					}},
				}},
		},
		{
			jsonString: `{"pairs": [{"x0": -21.907810617638056, "y0": -90, "x1": 22.221643690625143, "y1": -82.10060933854065}, {"x0": -158.7979515172962, "y0": 13.619545329636122, "x1": -113.67448769312786, "y1": 35.318722339808645}]}`,
			expectedNode: &Node{
				Type: NODE_OBJECT, Value: nil, Children: []*Node{
					{Type: NODE_KEY_VALUE, Key: &Node{Type: NODE_STRING, Value: "pairs"}, Val: &Node{Type: NODE_ARRAY, Value: nil, Children: []*Node{
						{Type: NODE_OBJECT, Value: nil, Children: []*Node{
							{Type: NODE_KEY_VALUE, Key: &Node{Type: NODE_STRING, Value: "x0"}, Val: &Node{Type: NODE_NUMBER, Value: "-21.907810617638056"}},
							{Type: NODE_KEY_VALUE, Key: &Node{Type: NODE_STRING, Value: "y0"}, Val: &Node{Type: NODE_NUMBER, Value: "-90"}},
							{Type: NODE_KEY_VALUE, Key: &Node{Type: NODE_STRING, Value: "x1"}, Val: &Node{Type: NODE_NUMBER, Value: "22.221643690625143"}},
							{Type: NODE_KEY_VALUE, Key: &Node{Type: NODE_STRING, Value: "y1"}, Val: &Node{Type: NODE_NUMBER, Value: "-82.10060933854065"}},
						}},
						{Type: NODE_OBJECT, Value: nil, Children: []*Node{
							{Type: NODE_KEY_VALUE, Key: &Node{Type: NODE_STRING, Value: "x0"}, Val: &Node{Type: NODE_NUMBER, Value: "-158.7979515172962"}},
							{Type: NODE_KEY_VALUE, Key: &Node{Type: NODE_STRING, Value: "y0"}, Val: &Node{Type: NODE_NUMBER, Value: "13.619545329636122"}},
							{Type: NODE_KEY_VALUE, Key: &Node{Type: NODE_STRING, Value: "x1"}, Val: &Node{Type: NODE_NUMBER, Value: "-113.67448769312786"}},
							{Type: NODE_KEY_VALUE, Key: &Node{Type: NODE_STRING, Value: "y1"}, Val: &Node{Type: NODE_NUMBER, Value: "35.318722339808645"}},
						}},
					}}},
				}},
		},
	}

	for _, testCase := range testCases {
		lexer := NewLexer(testCase.jsonString)
		tokens := lexer.Tokenize()

		parser := NewParser(tokens)
		node := parser.Parse()

		compareNodesWithLexer(t, testCase.jsonString, node, testCase.expectedNode)
	}
}

func compareNodesWithLexer(t *testing.T, jsonString string, node *Node, expectedNode *Node) {
	t.Helper()

	if expectedNode == nil || node == nil {
		if expectedNode != node {
			t.Errorf("Expected node to be %v, got: %v in test case with JSON: %s", expectedNode, node, jsonString)
		}
		return
	}

	if node.Type != expectedNode.Type {
		t.Errorf("Expected node type: %v, got: %v in test case with JSON: %s", expectedNode.Type, node.Type, jsonString)
	}
	if node.Value != expectedNode.Value {
		t.Errorf("Expected node value: %v, got: %v in test case with JSON: %s", expectedNode.Value, node.Value, jsonString)
	}
	if node.Children != nil {
		if len(node.Children) != len(expectedNode.Children) {
			t.Errorf("Expected %d children, got %d in test case with JSON: %s", len(expectedNode.Children), len(node.Children), jsonString)
			return
		}
		for i := range node.Children {
			compareNodesWithLexer(t, jsonString, node.Children[i], expectedNode.Children[i])
		}
	}
}
