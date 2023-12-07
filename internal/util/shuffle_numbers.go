package util

import (
	"math/rand"
)

func ShuffleNumbers(max int) []int {
	numbers := make([]int, max)
	for i := 0; i < max; i++ {
		numbers[i] = i + 1
	}

	// Fisher-Yates shuffle algorithm
	for i := len(numbers) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		numbers[i], numbers[j] = numbers[j], numbers[i]
	}

	return numbers
}
