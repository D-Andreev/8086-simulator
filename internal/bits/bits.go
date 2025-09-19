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

// ToSigned8 converts a byte to a signed 8-bit number.
func ToSigned8(bits byte) int16 {
	if bits&0x80 != 0 { // Check if MSB is set
		return int16(int8(bits)) // Convert to signed 8-bit
	}
	return int16(bits)
}

// ToSigned16 converts two bytes to a signed 16-bit number.
func ToSigned16(low, high byte) int16 {
	value := uint16(low) | (uint16(high) << 8)
	if value&0x8000 != 0 { // Check if MSB is set
		return int16(value) // Convert to signed 16-bit
	}
	return int16(value)
}

// ToUnsigned16 converts two bytes to an unsigned 16-bit number.
func ToUnsigned16(low, high byte) uint16 {
	return uint16(low) | (uint16(high) << 8)
}

// ToUnsigned8 converts a byte to an unsigned 8-bit number.
func ToUnsigned8(bits byte) uint16 {
	return uint16(bits)
}

// IsZero checks if all bits are zero.
func IsZero(bits []byte) bool {
	for _, bit := range bits {
		if bit != 0 {
			return false
		}
	}
	return true
}

// IsNegative checks if the most significant bit is set.
func IsNegative(bits []byte) bool {
	if len(bits) == 2 {
		return bits[1]&0x80 != 0 // Check high byte for 16-bit values
	}
	return bits[0]&0x80 != 0 // Check low byte for 8-bit values
}

// Add adds two 16-bit values represented as byte arrays.
func Add(a, b []byte) []byte {
	valA := ToUnsigned16(a[0], a[1])
	valB := ToUnsigned16(b[0], b[1])
	result := valA + valB
	return []byte{byte(result & 0xFF), byte((result >> 8) & 0xFF)}
}

// Sub subtracts two 16-bit values represented as byte arrays.
func Sub(a, b []byte) []byte {
	valA := ToUnsigned16(a[0], a[1])
	valB := ToUnsigned16(b[0], b[1])
	result := valA - valB
	return []byte{byte(result & 0xFF), byte((result >> 8) & 0xFF)}
}
