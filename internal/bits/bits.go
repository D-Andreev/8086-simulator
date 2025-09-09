package bits

// GetBit returns the value of the bit at the given index.
func GetBit(b byte, index int) bool {
	if index < 0 || index > 7 {
		panic("index out of range")
	}
	return (b>>index)&1 == 1
}

// GetBits returns the value of the bits in the given range.
func GetBits(b byte, start, count int) byte {
	if start < 0 || start+count > 8 {
		panic("start or count out of range")
	}

	// Create mask with 'count' number of 1s
	mask := (byte(1) << count) - 1 // e.g., for 3 bits: 0b111

	// Shift mask to correct position and apply
	return (b >> start) & mask
}

// For 8-bit signed numbers
func ToSigned8(bits byte) int {
	if bits&0x80 != 0 { // Check if MSB is set
		return int(int8(bits)) // Convert to signed 8-bit
	}
	return int(bits)
}

// For 16-bit signed numbers
func ToSigned16(low, high byte) int {
	value := uint16(low) | (uint16(high) << 8)
	if value&0x8000 != 0 { // Check if MSB is set
		return int(int16(value)) // Convert to signed 16-bit
	}
	return int(value)
}
