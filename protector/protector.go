package protector

import (
	"strconv"
	"strings"
	"unicode"
)

type Protector struct {
	Hash string
}

func New(hash string) *Protector {
	var p = new(Protector)
	p.Hash = hash
	return p
}

func (p Protector) GenerateNextSessionKey(sessionKey string) string {
	if p.Hash == "" {
		panic("Hash is empty")
	}

	for i := 0; i < len(p.Hash); i++ {
		var digit = rune(p.Hash[i])
		if !unicode.IsDigit(digit) {
			panic("Hash code contains non-digital letter")
		}
	}

	var result = 0
	for i := 0; i < len(p.Hash); i++ {
		var digit = int(p.Hash[i])
		var hash, err = strconv.Atoi(p.calculateHash(sessionKey, digit))
		if err != nil {
			panic(err)
		}
		result += hash
	}

	var sb strings.Builder
	sb.WriteString(strings.Repeat("0", 10))
	sb.WriteString(strconv.Itoa(result)[0:10])
	return sb.String()[len(sb.String()) - 10:]
}

func (p Protector) calculateHash(sessionKey string, value int) string {
	var sb strings.Builder
	switch value {
	case 1:
		sb.WriteString("00")
		var integer, err = strconv.Atoi(sessionKey[0:5])
		if err != nil {
			panic(err)
		}
		sb.WriteString(string(integer % 97))
		return sb.String()[len(sb.String()) - 2:]
	case 2:
		for i := 0; i < len(sessionKey); i++ {
			sb.WriteByte(sessionKey[i * (-1)])
		}
		sb.WriteByte(sessionKey[0])
		return sb.String()
	case 3:
		sb.WriteString(sessionKey[len(sessionKey) - 5:])
		sb.WriteString(sessionKey[0:5])
		return sb.String()
	case 4:
		var num = 0
		for i := 0; i < 9; i++ {
			num += (int(sessionKey[i]) - 48) + 41
		}
		return strconv.Itoa(num)
	case 5:
		var num = 0
		for i := 0; i < len(sessionKey); i++ {
			var char = rune(sessionKey[i] ^ 43)
			if !unicode.IsDigit(char) {
				char = rune(sessionKey[i])
			}
			num += int(char)
		}
		return strconv.Itoa(num)
	default:
		var num, err = strconv.Atoi(sessionKey)
		if err != nil {
			panic(err)
		}
		return strconv.Itoa(num + value)
	}
}
