package models

import (
	"fmt"
	"strconv"
)

type HexUint64 uint64

func (h *HexUint64) UnmarshalJSON(b []byte) error {
	if b[0] == '"' && b[len(b)-1] == '"' {
		b = b[1 : len(b)-1]
	}
	num, err := strconv.ParseUint(string(b[2:]), 16, 64)
	if err != nil {
		return err
	}
	*h = HexUint64(num)
	return nil
}

func (h HexUint64) MarshalJSON() ([]byte, error) {
	hexString := fmt.Sprintf("\"0x%x\"", h)
	return []byte(hexString), nil
}
