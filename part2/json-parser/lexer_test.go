package jsonparser

import (
	"testing"
)

type LexerTestCase struct {
	input          string
	expectedTokens []Token
}

func TestLexer(t *testing.T) {
	testCases := []LexerTestCase{
		{input: "{}", expectedTokens: []Token{
			{Type: OPEN_BRACE, Literal: "{"},
			{Type: CLOSE_BRACE, Literal: "}"},
		}},
		{input: "[]", expectedTokens: []Token{
			{Type: OPEN_BRACKET, Literal: "["},
			{Type: CLOSE_BRACKET, Literal: "]"},
		}},
		{
			input: "{\"name\":\"John\",\"lastname\":\"Doe\",\"age\":30,\"is_student\":true,\"is_teacher\":false,\"is_admin\":null}",
			expectedTokens: []Token{
				{Type: OPEN_BRACE, Literal: "{"},
				{Type: STRING, Literal: "name"},
				{Type: COLON, Literal: ":"},
				{Type: STRING, Literal: "John"},
				{Type: COMMA, Literal: ","},
				{Type: STRING, Literal: "lastname"},
				{Type: COLON, Literal: ":"},
				{Type: STRING, Literal: "Doe"},
				{Type: COMMA, Literal: ","},
				{Type: STRING, Literal: "age"},
				{Type: COLON, Literal: ":"},
				{Type: NUMBER, Literal: "30"},
				{Type: COMMA, Literal: ","},
				{Type: STRING, Literal: "is_student"},
				{Type: COLON, Literal: ":"},
				{Type: BOOL, Literal: "true"},
				{Type: COMMA, Literal: ","},
				{Type: STRING, Literal: "is_teacher"},
				{Type: COLON, Literal: ":"},
				{Type: BOOL, Literal: "false"},
				{Type: COMMA, Literal: ","},
				{Type: STRING, Literal: "is_admin"},
				{Type: COLON, Literal: ":"},
				{Type: NULL, Literal: "null"},
				{Type: CLOSE_BRACE, Literal: "}"},
			}},
		// Array with flat values
		{
			input: "{\"stuff\":[\"Jane\",\"Jim\",\"Jill\",30,true,false,null]}", expectedTokens: []Token{
				{Type: OPEN_BRACE, Literal: "{"},
				{Type: STRING, Literal: "stuff"},
				{Type: COLON, Literal: ":"},
				{Type: OPEN_BRACKET, Literal: "["},
				{Type: STRING, Literal: "Jane"},
				{Type: COMMA, Literal: ","},
				{Type: STRING, Literal: "Jim"},
				{Type: COMMA, Literal: ","},
				{Type: STRING, Literal: "Jill"},
				{Type: COMMA, Literal: ","},
				{Type: NUMBER, Literal: "30"},
				{Type: COMMA, Literal: ","},
				{Type: BOOL, Literal: "true"},
				{Type: COMMA, Literal: ","},
				{Type: BOOL, Literal: "false"},
				{Type: COMMA, Literal: ","},
				{Type: NULL, Literal: "null"},
				{Type: CLOSE_BRACKET, Literal: "]"},
				{Type: CLOSE_BRACE, Literal: "}"},
			}},
		// Objects
		{
			input: "{\"user\":{\"name\":\"john\",\"age\":30}}",
			expectedTokens: []Token{
				{Type: OPEN_BRACE, Literal: "{"},
				{Type: STRING, Literal: "user"},
				{Type: COLON, Literal: ":"},
				{Type: OPEN_BRACE, Literal: "{"},
				{Type: STRING, Literal: "name"},
				{Type: COLON, Literal: ":"},
				{Type: STRING, Literal: "john"},
				{Type: COMMA, Literal: ","},
				{Type: STRING, Literal: "age"},
				{Type: COLON, Literal: ":"},
				{Type: NUMBER, Literal: "30"},
				{Type: CLOSE_BRACE, Literal: "}"},
				{Type: CLOSE_BRACE, Literal: "}"},
			},
		},
		// Array with nested objects
		{
			input: "{\"users\":[{\"name\":\"john\",\"age\":30},{\"name\":\"jane\",\"age\":25}]}",
			expectedTokens: []Token{
				{Type: OPEN_BRACE, Literal: "{"},
				{Type: STRING, Literal: "users"},
				{Type: COLON, Literal: ":"},
				{Type: OPEN_BRACKET, Literal: "["},
				{Type: OPEN_BRACE, Literal: "{"},
				{Type: STRING, Literal: "name"},
				{Type: COLON, Literal: ":"},
				{Type: STRING, Literal: "john"},
				{Type: COMMA, Literal: ","},
				{Type: STRING, Literal: "age"},
				{Type: COLON, Literal: ":"},
				{Type: NUMBER, Literal: "30"},
				{Type: CLOSE_BRACE, Literal: "}"},
				{Type: COMMA, Literal: ","},
				{Type: OPEN_BRACE, Literal: "{"},
				{Type: STRING, Literal: "name"},
				{Type: COLON, Literal: ":"},
				{Type: STRING, Literal: "jane"},
				{Type: COMMA, Literal: ","},
				{Type: STRING, Literal: "age"},
				{Type: COLON, Literal: ":"},
				{Type: NUMBER, Literal: "25"},
				{Type: CLOSE_BRACE, Literal: "}"},
				{Type: CLOSE_BRACKET, Literal: "]"},
				{Type: CLOSE_BRACE, Literal: "}"},
			},
		},
		// Array with objects
		{
			input: "{\"users\":[{\"name\":\"john\",\"age\":30},{\"name\":\"jane\",\"age\":25}]}",
			expectedTokens: []Token{
				{Type: OPEN_BRACE, Literal: "{"},
				{Type: STRING, Literal: "users"},
				{Type: COLON, Literal: ":"},
				{Type: OPEN_BRACKET, Literal: "["},
				{Type: OPEN_BRACE, Literal: "{"},
				{Type: STRING, Literal: "name"},
				{Type: COLON, Literal: ":"},
				{Type: STRING, Literal: "john"},
				{Type: COMMA, Literal: ","},
				{Type: STRING, Literal: "age"},
				{Type: COLON, Literal: ":"},
				{Type: NUMBER, Literal: "30"},
				{Type: CLOSE_BRACE, Literal: "}"},
				{Type: COMMA, Literal: ","},
				{Type: OPEN_BRACE, Literal: "{"},
				{Type: STRING, Literal: "name"},
				{Type: COLON, Literal: ":"},
				{Type: STRING, Literal: "jane"},
				{Type: COMMA, Literal: ","},
				{Type: STRING, Literal: "age"},
				{Type: COLON, Literal: ":"},
				{Type: NUMBER, Literal: "25"},
				{Type: CLOSE_BRACE, Literal: "}"},
				{Type: CLOSE_BRACKET, Literal: "]"},
				{Type: CLOSE_BRACE, Literal: "}"},
			},
		},
		// 3 levels of nested objects
		{
			input: "{\"user\":{\"name\":\"john\",\"age\":30,\"address\":{\"street\":\"123 Main St\",\"city\":\"Anytown\",\"state\":\"CA\"}}}",
			expectedTokens: []Token{
				{Type: OPEN_BRACE, Literal: "{"},
				{Type: STRING, Literal: "user"},
				{Type: COLON, Literal: ":"},
				{Type: OPEN_BRACE, Literal: "{"},
				{Type: STRING, Literal: "name"},
				{Type: COLON, Literal: ":"},
				{Type: STRING, Literal: "john"},
				{Type: COMMA, Literal: ","},
				{Type: STRING, Literal: "age"},
				{Type: COLON, Literal: ":"},
				{Type: NUMBER, Literal: "30"},
				{Type: COMMA, Literal: ","},
				{Type: STRING, Literal: "address"},
				{Type: COLON, Literal: ":"},
				{Type: OPEN_BRACE, Literal: "{"},
				{Type: STRING, Literal: "street"},
				{Type: COLON, Literal: ":"},
				{Type: STRING, Literal: "123 Main St"},
				{Type: COMMA, Literal: ","},
				{Type: STRING, Literal: "city"},
				{Type: COLON, Literal: ":"},
				{Type: STRING, Literal: "Anytown"},
				{Type: COMMA, Literal: ","},
				{Type: STRING, Literal: "state"},
				{Type: COLON, Literal: ":"},
				{Type: STRING, Literal: "CA"},
				{Type: CLOSE_BRACE, Literal: "}"},
				{Type: CLOSE_BRACE, Literal: "}"},
				{Type: CLOSE_BRACE, Literal: "}"},
			},
		},
	}

	for _, testCase := range testCases {
		lexer := NewLexer(testCase.input)
		tokens := lexer.Tokenize()
		compareTokens(t, testCase, tokens, testCase.expectedTokens)
	}
}

func compareTokens(t *testing.T, testCase LexerTestCase, tokens []Token, expectedTokens []Token) {
	t.Helper()
	for i, expectedToken := range expectedTokens {
		if tokens[i] != expectedToken {
			t.Errorf("Expected token: %v, got: %v in test case: %s", expectedToken, tokens[i], testCase.input)
		}
	}

	if len(tokens) != len(expectedTokens) {
		t.Errorf("Expected %d tokens, got %d in test case: %s", len(expectedTokens), len(tokens), testCase.input)
	}
}
