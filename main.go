package main

import (
	"./protector"
	"math/rand"
	"strings"
)

func getSessionKey() string {
	var sb strings.Builder
	for i := 0; i < 10; i++ {
		var byte = byte('1' + rand.Intn(9))
		sb.WriteByte(byte)
	}

	return sb.String()
}

func getHashString() string {
	var sb strings.Builder
	for i := 0; i < 5; i++ {
		var byte = byte('1' + rand.Intn(6))
		sb.WriteByte(byte)
	}

	return sb.String()
}

func main() {
	// initial - start with random string
	var sKeyInitial = getSessionKey()
	var hashString = getHashString()

	var protector1 = protector.New(hashString)
	var protector2 = protector.New(hashString)

	// stage 0
	var sKey1 = protector1.GenerateNextSessionKey(sKeyInitial)
	var sKey2 = protector2.GenerateNextSessionKey(sKeyInitial)
	if sKey1 != sKey2 {
		panic("Keys are not equal")
	}

	// stage 1 - each protector use it's own skey and compare with another's side
	var sKey3 = protector1.GenerateNextSessionKey(sKey1)
	var sKey4 = protector2.GenerateNextSessionKey(sKey2)
	if sKey3 != sKey4 {
		panic("Keys are not equal")
	}

	// stage 2
	var sKey5 = protector1.GenerateNextSessionKey(sKey3)
	var sKey6 = protector2.GenerateNextSessionKey(sKey4)
	if sKey5 != sKey6 {
		panic("Keys are not equal")
	}

	print (sKey1, sKey2, sKey3, sKey4, sKey5, sKey6)
}