package stunning

// Get a single bit
func GetBit(n int, pos uint) bool {
	return (n & (1 << pos)) != 0
}

// Set a bit to 1
func SetBit(n int, pos uint, val bool) int {

	var intVal int = 0
	if val {
		intVal = 1
	}

	return n | (intVal << pos)
}

// Clear a bit to 0
func ClearBit(n int, pos uint) int {
	return n &^ (1 << pos)
}

// Toggle a bit
func ToggleBit(n int, pos uint) int {
	return n ^ (1 << pos)
}

// Get multiple bits
func GetBits(n int, start, length uint) int {
	mask := (1 << length) - 1
	return (n >> start) & mask
}

// Set multiple bits
func SetBits(n int, start, length uint, value int) int {
	mask := ((1 << length) - 1) << start
	n &^= mask // Clear the bits
	return n | ((value & ((1 << length) - 1)) << start)
}
