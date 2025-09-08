package decoder

import "testing"

func TestDecoder(t *testing.T) {
	decoder := NewDecoder()
	decoder.Decode([]byte("Hello, World!"))
}
