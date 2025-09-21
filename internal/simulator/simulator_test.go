package simulator

import (
	"os"
	"testing"

	"github.com/8086-simulator/internal/bits"
	"github.com/8086-simulator/internal/decoder"
)

func TestSimulatorListing43(t *testing.T) {
	content, err := os.ReadFile("../../listings/listing_0043_immediate_movs")
	if err != nil {
		t.Fatalf("Error reading file: %v", err)
	}
	decoder := decoder.NewDecoder()
	sim := NewSimulator(false)
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
		"ax": bits.Uint16ToBytes(1),
		"bx": bits.Uint16ToBytes(2),
		"cx": bits.Uint16ToBytes(3),
		"dx": bits.Uint16ToBytes(4),
		"sp": bits.Uint16ToBytes(5),
		"bp": bits.Uint16ToBytes(6),
		"si": bits.Uint16ToBytes(7),
		"di": bits.Uint16ToBytes(8),
		"ip": bits.Uint16ToBytes(0),
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
	sim := NewSimulator(false)
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
		"ax": bits.Uint16ToBytes(4),
		"bx": bits.Uint16ToBytes(3),
		"cx": bits.Uint16ToBytes(2),
		"dx": bits.Uint16ToBytes(1),
		"sp": bits.Uint16ToBytes(1),
		"bp": bits.Uint16ToBytes(2),
		"si": bits.Uint16ToBytes(3),
		"di": bits.Uint16ToBytes(4),
		"ip": bits.Uint16ToBytes(0),
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
	sim := NewSimulator(false)
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
		"ax": bits.Uint16ToBytes(0),
		"bx": bits.Uint16ToBytes(57602),
		"cx": bits.Uint16ToBytes(3841),
		"dx": bits.Uint16ToBytes(0),
		"sp": bits.Uint16ToBytes(998),
		"bp": bits.Uint16ToBytes(0),
		"si": bits.Uint16ToBytes(0),
		"di": bits.Uint16ToBytes(0),
		"ip": bits.Uint16ToBytes(0),
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

func TestSimulatorListing48(t *testing.T) {
	content, err := os.ReadFile("../../listings/listing_0048_ip_register")
	if err != nil {
		t.Fatalf("Error reading file: %v", err)
	}
	decoder := decoder.NewDecoder()
	sim := NewSimulator(true)
	sim.Init()
	expectedLogs := []string{
		"mov cx, 200 ; cx:0x0->0xc8 ip:0x0->0x3",
		"mov bx, cx ; bx:0x0->0xc8 ip:0x3->0x5",
		"add cx, 1000 ; cx:0xc8->0x4b0 ip:0x5->0x9",
		"mov bx, 2000 ; bx:0xc8->0x7d0 ip:0x9->0xc",
		"sub cx, bx ; cx:0x4b0->0xfce0 ip:0xc->0xe flags:->S",
	}
	expectedRegisters := map[string][]byte{
		"ax": bits.Uint16ToBytes(0),
		"bx": bits.Uint16ToBytes(2000),
		"cx": bits.Uint16ToBytes(64736),
		"dx": bits.Uint16ToBytes(0),
		"sp": bits.Uint16ToBytes(0),
		"bp": bits.Uint16ToBytes(0),
		"si": bits.Uint16ToBytes(0),
		"di": bits.Uint16ToBytes(0),
		"ip": bits.Uint16ToBytes(14),
	}
	expectedFlags := map[string]bool{
		"Z": false,
		"S": true,
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

	for flag, value := range sim.flags {
		if value != expectedFlags[flag] {
			t.Fatalf("Expected flag %s to be %t but got %t", flag, expectedFlags[flag], value)
		}
	}
}
