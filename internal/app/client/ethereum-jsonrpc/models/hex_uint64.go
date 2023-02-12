package models

import (
	"fmt"
	"strconv"
)

// HexUint64 represent 64-bit quantity that could be used for general fields
// in Ethereum Network whose size does larger than 64-bit, for larger fields
// using 256-bit HexBigInt will be more convenient for handling
type HexUint64 uint64

// UnmarshalJSON converts a JSON encoded string representation of a hexadecimal number
// into a `uint64` and stores it in the `HexUint64` type.
func (h *HexUint64) UnmarshalJSON(b []byte) error {
	// Remove the quotes from the incoming byte slice, if present
	if b[0] == '"' && b[len(b)-1] == '"' {
		b = b[1 : len(b)-1]
	}

	// Try to parse the hexadecimal string as a uint64
	num, err := strconv.ParseUint(string(b[2:]), 16, 64)
	if err != nil {
		// Return the error if the string could not be converted to a uint64
		return err
	}

	// Assign the converted uint64 to the HexUint64 type
	*h = HexUint64(num)

	// Return nil if there were no errors
	return nil
}

// MarshalJSON converts a `HexUint64` type into its JSON encoded string representation as a hexadecimal number.
func (h HexUint64) MarshalJSON() ([]byte, error) {
	// Format the HexUint64 type as a hexadecimal string
	hexString := fmt.Sprintf("\"0x%x\"", h)

	// Return the hexadecimal string as a byte slice
	return []byte(hexString), nil
}
