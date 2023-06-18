package main

import (
	"unicode/utf8"
)

func revers_string(stringToRevers string) string {

	var reversedString string = ""
	var reversedByte []byte
	var runeToReverse []rune

	for len(stringToRevers) > 0 {
		runeChar, size := utf8.DecodeRuneInString(stringToRevers)
		runeToReverse = append(runeToReverse, runeChar)
		stringToRevers = stringToRevers[size:]
	}

	runeLen := len(runeToReverse)
	for runePos := runeLen - 1; runePos >= 0; runePos-- {
		reversedByte = utf8.AppendRune(reversedByte, runeToReverse[runePos])
	}

	reversedString = string(reversedByte)

	return reversedString
}
