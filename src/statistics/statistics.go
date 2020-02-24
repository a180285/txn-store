package main

import (
	"fmt"
	"math/rand"
	"time"
	"txn-store/src/txnStore"
)

const (
	NUMBER_IN_TXN = 10
)

func printStatics() {
	total := 1.0
	for i := 0; i < 10; i++ {
		total *= float64(1000 - i)
	}

	expected := float64(0)

	for i := 1; i < 60; i++ {
		selected := 1.0
		for j := 0; j < 10; j++ {
			selected *= float64(1000 - (i-1)*10 - j)
		}

		expected += total / selected
		fmt.Printf("success txn: %3d, expected rand count: %f\n", i, expected)
	}

}

func main() {
	rand.Seed(time.Now().Unix())
	printStatics()
}

func randKeys(length int) []int {
	keys := make([]int, length)
	for idx := 0; idx < length; idx++ {
	retry:

		newKey := rand.Intn(txnStore.MAX_KEYS)
		for i := 0; i < idx; i++ {
			if keys[i] == newKey {
				goto retry
			}
		}

		keys[idx] = newKey
	}

	return keys
}
