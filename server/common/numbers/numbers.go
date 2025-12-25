package numbers

import (
	"strconv"

	"golang.org/x/exp/constraints"
)

// ToInt64 converts a string to int64. Returns 0 if conversion fails.
// str input string
func ToInt64(str string) int64 {
	return ToInt64ByDefault(str, 0)
}

// ToInt64ByDefault converts a string to int64.
// str input string
// def default value if conversion fails
func ToInt64ByDefault(str string, def int64) int64 {
	val, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		val = def
	}
	return val
}

// ToInt converts a string to int. Returns 0 if conversion fails.
// str input string
func ToInt(str string) int {
	return ToIntByDefault(str, 0)
}

// ToIntByDefault converts a string to int.
// str input string
// def default value if conversion fails
func ToIntByDefault(str string, def int) int {
	val, err := strconv.Atoi(str)
	if err != nil {
		val = def
	}
	return val
}

func Min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func Max[T constraints.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}
