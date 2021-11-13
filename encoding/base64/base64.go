package base64

import (
	"fmt"
	"strings"
)

var codes = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_-"

func EncodeID(id uint64) string {
	str := make([]byte, 0, 12)
	if id == 0 {
		return "0"
	}
	for id > 0 {
		ch := codes[id%64]
		str = append(str, byte(ch))
		id /= 64
	}
	return string(str)
}

func DecodeID(token string) (uint64, error) {
	res := uint64(0)

	for i := len(token); i > 0; i-- {
		ch := token[i-1]
		res *= 64
		mod := strings.IndexRune(codes, rune(ch))
		if mod == -1 {
			return 0, fmt.Errorf("Invalid ltoken character: '%c'", ch)
		}
		res += uint64(mod)
	}
	return res, nil
}
