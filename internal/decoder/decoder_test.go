package decoder

import (
	"os"
	"testing"
)

func TestDecoderListing37(t *testing.T) {
	content, err := os.ReadFile("../../listings/listing_0037_single_register_mov")
	if err != nil {
		t.Fatalf("Error reading file: %v", err)
	}
	expectedInstructions := []string{
		"mov cx, bx",
	}

	decoder := NewDecoder()
	instructions, err := decoder.Decode(content)
	if err != nil {
		t.Fatalf("Error decoding data: %v", err)
	}

	for i, instruction := range instructions {
		if instruction.Text != expectedInstructions[i] {
			t.Fatalf("Expected instruction %s but got %s", expectedInstructions[i], instruction.Text)
		}
	}
}

func TestDecoderListing38(t *testing.T) {
	content, err := os.ReadFile("../../listings/listing_0038_many_register_mov")
	if err != nil {
		t.Fatalf("Error reading file: %v", err)
	}
	expectedInstructions := []string{
		"mov cx, bx",
		"mov ch, ah",
		"mov dx, bx",
		"mov si, bx",
		"mov bx, di",
		"mov al, cl",
		"mov ch, ch",
		"mov bx, ax",
		"mov bx, si",
		"mov sp, di",
		"mov bp, ax",
	}

	decoder := NewDecoder()
	instructions, err := decoder.Decode(content)
	if err != nil {
		t.Fatalf("Error decoding data: %v", err)
	}

	for i, instruction := range instructions {
		if instruction.Text != expectedInstructions[i] {
			t.Fatalf("Expected instruction %s but got %s", expectedInstructions[i], instruction.Text)
		}
	}
}

func TestDecoderListing39(t *testing.T) {
	content, err := os.ReadFile("../../listings/listing_0039_more_movs")
	if err != nil {
		t.Fatalf("Error reading file: %v", err)
	}
	expectedInstructions := []string{
		"mov si, bx",
		"mov dh, al",
		"mov cl, 12",
		"mov ch, -12",
		"mov cx, 12",
		"mov cx, -12",
		"mov dx, 3948",
		"mov dx, -3948",
		"mov al, [bx + si]",
		"mov bx, [bp + di]",
		"mov dx, [bp + 0]",
		"mov ah, [bx + si + 4]",
		"mov al, [bx + si + 4999]",
		"mov [bx + di], cx",
		"mov [bp + si], cl",
		"mov [bp + 0], ch",
	}

	decoder := NewDecoder()
	instructions, err := decoder.Decode(content)
	if err != nil {
		t.Fatalf("Error decoding data: %v", err)
	}

	for i, instruction := range instructions {
		if instruction.Text != expectedInstructions[i] {
			t.Fatalf("Expected instruction %s but got %s", expectedInstructions[i], instruction.Text)
		}
	}
}
