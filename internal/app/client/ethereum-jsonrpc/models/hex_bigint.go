package models

import (
	"errors"
	"fmt"
	"math/big"
)

type HexBigInt big.Int

func (h *HexBigInt) UnmarshalJSON(b []byte) error {
	str := string(b)
	if str[0] == '"' && str[len(b)-1] == '"' {
		str = str[1 : len(str)-1]
	}
	v, success := new(big.Int).SetString(str[2:], 16)
	if !success {
		return errors.New("cannot convert hex value to bigint: " + string(str))
	}

	*h = HexBigInt(*v)

	return nil
}

func (h HexBigInt) MarshalJSON() ([]byte, error) {
	bigInt := big.Int(h)
	hexString := fmt.Sprintf("\"0x%x\"", &bigInt)
	return []byte(hexString), nil
}
