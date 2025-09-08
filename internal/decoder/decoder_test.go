package decoder

import "testing"

func TestDecoder(t *testing.T) {
	decoder := NewDecoder()
	_, err := decoder.Decode([]byte("Hello, World!"))
	if err != nil {
		t.Fatalf("Error decoding data: %v", err)
	}
}
