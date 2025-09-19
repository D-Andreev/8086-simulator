package simulator

import (
	"os"
	"testing"

	"github.com/8086-simulator/internal/decoder"
)

func TestSimulatorListing43(t *testing.T) {
	content, err := os.ReadFile("../../listings/listing_0043_immediate_movs")
	if err != nil {
		t.Fatalf("Error reading file: %v", err)
	}
	decoder := decoder.NewDecoder()
	sim := NewSimulator()
	sim.Init()
	expectedLogs := []string{
		"mov ax, 1 ; ax:0x0->0x1",
		"mov bx, 2 ; bx:0x0->0x2",
		"mov cx, 3 ; cx:0x0->0x3",
		"mov dx, 4 ; dx:0x0->0x4",
		"mov sp, 5 ; sp:0x0->0x5",
		"mov bp, 6 ; bp:0x0->0x6",
		"mov si, 7 ; si:0x0->0x7",
		"mov di, 8 ; di:0x0->0x8",
	}
	expectedRegisters := map[string][]byte{
		"ax": {0x1, 0x0},
		"bx": {0x2, 0x0},
		"cx": {0x3, 0x0},
		"dx": {0x4, 0x0},
		"sp": {0x5, 0x0},
		"bp": {0x6, 0x0},
		"si": {0x7, 0x0},
		"di": {0x8, 0x0},
	}

	instructions, err := decoder.Decode(content)
	if err != nil {
		t.Fatalf("Error decoding data: %v", err)
	}
	results, err := sim.Run(instructions)
	if err != nil {
		t.Fatalf("Error running instructions: %v", err)
	}

	for i, result := range results {
		if result.Text != expectedLogs[i] {
			t.Fatalf("Expected instruction %s but got %s", expectedLogs[i], result.Text)
		}
	}

	for register, value := range sim.Registers {
		if !registerValueEquals(t, value, expectedRegisters[register]) {
			t.Fatalf("Expected register %s to be 0x%x but got 0x%x", register, expectedRegisters[register], value)
		}
	}
}

func registerValueEquals(t *testing.T, value []byte, expected []byte) bool {
	t.Helper()
	if len(value) != len(expected) {
		return false
	}
	for i, v := range value {
		if v != expected[i] {
			return false
		}
	}
	return true
}

func TestSimulatorListing44(t *testing.T) {
	content, err := os.ReadFile("../../listings/listing_0044_register_movs")
	if err != nil {
		t.Fatalf("Error reading file: %v", err)
	}
	decoder := decoder.NewDecoder()
	sim := NewSimulator()
	sim.Init()
	expectedLogs := []string{
		"mov ax, 1 ; ax:0x0->0x1",
		"mov bx, 2 ; bx:0x0->0x2",
		"mov cx, 3 ; cx:0x0->0x3",
		"mov dx, 4 ; dx:0x0->0x4",
		"mov sp, ax ; sp:0x0->0x1",
		"mov bp, bx ; bp:0x0->0x2",
		"mov si, cx ; si:0x0->0x3",
		"mov di, dx ; di:0x0->0x4",
		"mov dx, sp ; dx:0x4->0x1",
		"mov cx, bp ; cx:0x3->0x2",
		"mov bx, si ; bx:0x2->0x3",
		"mov ax, di ; ax:0x1->0x4",
	}
	expectedRegisters := map[string][]byte{
		"ax": {0x4, 0x0},
		"bx": {0x3, 0x0},
		"cx": {0x2, 0x0},
		"dx": {0x1, 0x0},
		"sp": {0x1, 0x0},
		"bp": {0x2, 0x0},
		"si": {0x3, 0x0},
		"di": {0x4, 0x0},
	}

	instructions, err := decoder.Decode(content)
	if err != nil {
		t.Fatalf("Error decoding data: %v", err)
	}
	results, err := sim.Run(instructions)
	if err != nil {
		t.Fatalf("Error running instructions: %v", err)
	}

	for i, result := range results {
		if result.Text != expectedLogs[i] {
			t.Fatalf("Expected instruction %s but got %s", expectedLogs[i], result.Text)
		}
	}

	for register, value := range sim.Registers {
		if !registerValueEquals(t, value, expectedRegisters[register]) {
			t.Fatalf("Expected register %s to be %d but got %d", register, expectedRegisters[register], value)
		}
	}
}

func TestSimulatorListing46(t *testing.T) {
	content, err := os.ReadFile("../../listings/listing_0046_add_sub_cmp")
	if err != nil {
		t.Fatalf("Error reading file: %v", err)
	}
	decoder := decoder.NewDecoder()
	sim := NewSimulator()
	sim.Init()
	expectedLogs := []string{
		"mov bx, 61443 ; bx:0x0->0xf003",
		"mov cx, 3841 ; cx:0x0->0xf01",
		"sub bx, cx ; bx:0xf003->0xe102 flags:->S",
		"mov sp, 998 ; sp:0x0->0x3e6",
		"mov bp, 999 ; bp:0x0->0x3e7",
		"cmp bp, sp ; flags:S->",
		"add bp, 1027 ; bp:0x3e7->0x7ea",
		"sub bp, 2026 ; bp:0x7ea->0x0 flags:->Z",
	}
	expectedRegisters := map[string][]byte{
		"ax": {0x0, 0x0},
		"bx": {0x02, 0xe1}, // 57602 = 0xe102
		"cx": {0x01, 0x0f}, // 3841 = 0x0f01
		"dx": {0x0, 0x0},
		"sp": {0xe6, 0x03}, // 998 = 0x03e6
		"bp": {0x0, 0x0},
		"si": {0x0, 0x0},
		"di": {0x0, 0x0},
	}

	instructions, err := decoder.Decode(content)
	if err != nil {
		t.Fatalf("Error decoding data: %v", err)
	}
	results, err := sim.Run(instructions)
	if err != nil {
		t.Fatalf("Error running instructions: %v", err)
	}

	for i, result := range results {
		if result.Text != expectedLogs[i] {
			t.Fatalf("Expected instruction %s but got %s", expectedLogs[i], result.Text)
		}
	}

	for register, value := range sim.Registers {
		if !registerValueEquals(t, value, expectedRegisters[register]) {
			t.Fatalf("Expected register %s to be %d but got %d", register, expectedRegisters[register], value)
		}
	}
}
