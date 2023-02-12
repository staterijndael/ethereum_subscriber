package models

import (
	"errors"
	"fmt"
	"math/big"
)

// HexBigInt represent 256-bit type Quantity in Ethereum network.
// Usually this type using for a fields that could reach such a large size
// and should be used as little as possible for fields that does not have this requirement
// cause handling this type can take significant runtime resources
type HexBigInt big.Int

// UnmarshalJSON converts a JSON encoded string representation of a hexadecimal number
// into a `big.Int` and stores it in the `HexBigInt` type.
func (h *HexBigInt) UnmarshalJSON(b []byte) error {
	// Convert the incoming byte slice to a string
	str := string(b)

	// Remove the quotes from the string, if present
	if str[0] == '"' && str[len(b)-1] == '"' {
		str = str[1 : len(str)-1]
	}

	// Try to parse the hexadecimal string as a big.Int
	v, success := new(big.Int).SetString(str[2:], 16)
	if !success {
		// Return an error if the string could not be converted to a big.Int
		return errors.New("cannot convert hex value to bigint: " + string(str))
	}

	// Assign the converted big.Int to the HexBigInt type
	*h = HexBigInt(*v)

	// Return nil if there were no errors
	return nil
}

// MarshalJSON converts a `HexBigInt` type into its JSON encoded string representation as a hexadecimal number.
func (h HexBigInt) MarshalJSON() ([]byte, error) {
	// Convert the HexBigInt type to a big.Int
	bigInt := big.Int(h)

	// Format the big.Int as a hexadecimal string
	hexString := fmt.Sprintf("\"0x%x\"", &bigInt)

	// Return the hexadecimal string as a byte slice
	return []byte(hexString), nil
}
