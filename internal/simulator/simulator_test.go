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
	expectedRegisters := map[string]int16{
		"ax": 1,
		"bx": 2,
		"cx": 3,
		"dx": 4,
		"sp": 5,
		"bp": 6,
		"si": 7,
		"di": 8,
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
		if value != expectedRegisters[register] {
			t.Fatalf("Expected register %s to be %d but got %d", register, expectedRegisters[register], value)
		}
	}
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
	expectedRegisters := map[string]int16{
		"ax": 4,
		"bx": 3,
		"cx": 2,
		"dx": 1,
		"sp": 1,
		"bp": 2,
		"si": 3,
		"di": 4,
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
		if value != expectedRegisters[register] {
			t.Fatalf("Expected register %s to be %d but got %d", register, expectedRegisters[register], value)
		}
	}
}
