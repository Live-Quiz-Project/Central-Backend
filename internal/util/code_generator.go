package util

import (
	"fmt"
	"math/rand"
	"time"
)

func CodeGenerator(existingCodes []string) string {
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))

	for {
		first2 := make([]byte, 2)
		for i := 0; i < 2; i++ {
			first2[i] = byte('A' + rand.Intn(26))
		}

		randCode := rand.Intn(10000)
		last4 := fmt.Sprintf("%04d", randCode)

		newCode := string(first2) + last4

		unique := true
		for _, code := range existingCodes {
			if code == newCode {
				unique = false
				break
			}
		}

		if unique {
			return newCode
		}
	}
}
