package usecase

import "math/rand"

func Index(s string, sub rune, offset int) int {
	for i := offset; i < len(s); i++ {
		if rune(s[i]) == sub {
			return i
		}
	}
	return -1
}

func PowInt(a, b int) int {
	res := 1

	for ; b > 0; b-- {
		res *= a
	}

	return res
}

func RandRange(min, max int) int {
	return rand.Intn(max-min) + min
}
