package dijkstra

import (
	"bytes"
	"errors"
	"strconv"

	"github.com/mr-tron/base58"
)

const (
	AddrLen   = 20
	EdgeIdLen = 40
)

type EdgeId [EdgeIdLen]byte

type Address [AddrLen]byte

func AddressContains(addrs []Address, address Address) bool {
	if addrs == nil {
		return false
	}
	for _, addr := range addrs {
		if AddressEqual(addr, address) {
			return true
		}
	}
	return false
}

func AddressEqual(address1 Address, address2 Address) bool {
	result := true

	for i := 0; i < AddrLen; i++ {
		if address1[i] != address2[i] {
			result = false
			break
		}
	}

	return result
}

func (self Address) MarshalText() (text []byte, err error) {
	var scratch [64]byte
	var e bytes.Buffer

	e.WriteByte('[')
	for i := 0; i < AddrLen; i++ {
		b := strconv.AppendUint(scratch[:0], uint64(self[i]), 10)
		e.Write(b)
		if i < AddrLen-1 {
			e.WriteByte(' ')
		}

	}
	e.WriteByte(']')

	return e.Bytes(), nil
}

func (self *Address) UnmarshalText(text []byte) error {
	newText := text[1:]

	startIdx := 0
	for i := 0; i < AddrLen; i++ {
		for newText[startIdx] == ' ' || newText[startIdx] == '[' {
			startIdx++
		}

		toIdx := startIdx
		for newText[toIdx] >= '0' && newText[toIdx] <= '9' {
			toIdx++
		}

		res, err := strconv.ParseUint(string(newText[startIdx:toIdx]), 10, 8)

		if err != nil {
			return errors.New("PaymentNetworkID TextUnmarshaler error!")
		} else {
			self[i] = byte(res)
		}

		startIdx = toIdx
	}

	return nil
}

func (self EdgeId) GetAddr1() Address {
	var tmp Address
	copy(tmp[:], self[0:AddrLen])
	return tmp
}

func (self EdgeId) GetAddr2() Address {
	var tmp Address
	copy(tmp[:], self[AddrLen:])
	return tmp
}

func (self EdgeId) MarshalText() (text []byte, err error) {
	var scratch [64]byte
	var e bytes.Buffer

	e.WriteByte('[')
	for i := 0; i < EdgeIdLen; i++ {
		b := strconv.AppendUint(scratch[:0], uint64(self[i]), 10)
		e.Write(b)
		if i < EdgeIdLen-1 {
			e.WriteByte(' ')
		}

	}
	e.WriteByte(']')

	return e.Bytes(), nil
}

func (self *EdgeId) UnmarshalText(text []byte) error {
	newText := text[1:]

	startIdx := 0
	for i := 0; i < EdgeIdLen; i++ {
		for newText[startIdx] == ' ' || newText[startIdx] == '[' {
			startIdx++
		}

		toIdx := startIdx
		for newText[toIdx] >= '0' && newText[toIdx] <= '9' {
			toIdx++
		}

		res, err := strconv.ParseUint(string(newText[startIdx:toIdx]), 10, 8)

		if err != nil {
			return errors.New("EdgeId TextUnmarshaler error!")
		} else {
			self[i] = byte(res)
		}

		startIdx = toIdx
	}

	return nil
}

// ToBase58 Encode Address to base58 string
func ToBase58(address Address) string {
	return base58.Encode(address[:])
}

// FromBase58 Decode base58 string into Address
func FromBase58(s string) (Address, error) {
	decoded, err := base58.Decode(s)
	if err != nil {
		return Address{}, err
	}
	var addr Address
	copy(addr[:], decoded[:])
	return addr, nil
}
