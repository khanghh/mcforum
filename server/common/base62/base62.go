package base62

import (
	"math"
	"strings"
)

const CODE62 = "xPZcjYlmQMXT2KG68ut4gNSd19HfkVDwabCzUvBLRoWJrqyhpisEAO53eI0nF7"
const CodeLength = 62

var offset int64 = 3645458574
var CodeMap = map[string]int64{
	"x": 0, "P": 1, "Z": 2, "c": 3, "j": 4, "Y": 5, "l": 6, "m": 7, "Q": 8, "M": 9,
	"X": 10, "T": 11, "2": 12, "K": 13, "G": 14, "6": 15, "8": 16, "u": 17, "t": 18, "4": 19,
	"g": 20, "N": 21, "S": 22, "d": 23, "1": 24, "9": 25, "H": 26, "f": 27, "k": 28, "V": 29,
	"D": 30, "w": 31, "a": 32, "b": 33, "C": 34, "z": 35, "U": 36, "v": 37, "B": 38, "L": 39,
	"R": 40, "o": 41, "W": 42, "J": 43, "r": 44, "q": 45, "y": 46, "h": 47, "p": 48, "i": 49,
	"s": 50, "E": 51, "A": 52, "O": 53, "5": 54, "3": 55, "e": 56, "I": 57, "0": 58, "n": 59,
	"F": 60, "7": 61,
}

func Encode(number int64) string {
	if number == 0 {
		return "0"
	}
	number = number + offset
	result := make([]byte, 0)
	for number > 0 {
		round := number / CodeLength
		remain := number % CodeLength
		result = append(result, CODE62[remain])
		number = round
	}
	return string(result)
}

func Decode(str string) int64 {
	if str == "" {
		return 0
	}
	str = strings.TrimSpace(str)
	var result int64 = 0
	for index, char := range []byte(str) {
		result += CodeMap[string(char)] * int64(math.Pow(CodeLength, float64(index)))
	}
	return result - offset
}
