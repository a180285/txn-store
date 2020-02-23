package main

import (
	"context"
	"log"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
	"txn-store/src/txnStore"
)

const (
	NUMBER_IN_TXN = 10
)

var successTxnCount int32 = 0
var failedTxnCount int32 = 0

func main() {
	rand.Seed(time.Now().Unix())

	log.Printf("Start test now.")

	runningSeconds := 20

	// Chanel to generate different keys
	keyChanel := randKeyChannel()

	threadCountList := []int{
		10, 20, 30, 40, 50, 60, 70, 80, 90,
	}
	//threadCountList = []int{5000}
	for _, numThreads := range threadCountList {
		doOneTest(numThreads, runningSeconds, keyChanel)
	}
}

func doOneTest(numThreads int, runningSeconds int, keyChanel <-chan []int) {
	successTxnCount = 0
	failedTxnCount = 0

	myStore := txnStore.NewMyTxnStore()

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	for i := 0; i < numThreads; i++ {
		wg.Add(1)
		go doTransactions(ctx, &wg, myStore, keyChanel)
	}

	time.Sleep(time.Duration(runningSeconds) * time.Second)

	cancel()
	wg.Wait()

	report(myStore, runningSeconds, numThreads)
}

func report(mystore txnStore.TxnStore, runningSeconds int, numThreads int) {
	tx, err := mystore.Begin()
	if err != nil {
		println(err)
		return
	}
	sum := 0
	nonZeroCount := 0

	for i := 0; i < txnStore.MAX_KEYS; i++ {
		value, err := mystore.GET(tx, i)
		if err != nil {
			println(err)
			return
		}

		sum += value

		if value != 0 {
			nonZeroCount++
		}
	}

	log.Printf("theads: %4d, txn success QPS: %f, total QPS: %f, sum: %d, non zero count: %d\n",
		numThreads,
		float64(successTxnCount)/float64(runningSeconds),
		float64(successTxnCount+failedTxnCount)/float64(runningSeconds),
		sum,
		nonZeroCount)
}

func doTransactions(ctx context.Context, wg *sync.WaitGroup, myStore txnStore.TxnStore, keyChanel <-chan []int) {
	defer wg.Done()

	nextKeys := <-keyChanel
	for {
		select {
		case <-ctx.Done():
			return
		default:
			_doTransactions(myStore, nextKeys)
			nextKeys = <-keyChanel
		}
	}
}

func randKeyChannel() <-chan []int {
	keyChannel := make(chan []int, 1024)

	go func() {
		for {
			keyChannel <- randKeys(NUMBER_IN_TXN)
		}
	}()

	return keyChannel
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

func _doTransactions(myStore txnStore.TxnStore, keys []int) error {
	tx, err := myStore.Begin()
	if err != nil {
		return err
	}

	values := make([]int, NUMBER_IN_TXN)
	for i := 0; i < NUMBER_IN_TXN; i++ {
		values[i], err = myStore.GET(tx, keys[i])
		if err != nil {
			myStore.Rollback(tx)
			atomic.AddInt32(&failedTxnCount, 1)
			return err
		}
	}

	time.Sleep(time.Duration(100) * time.Millisecond)

	for i := 0; i < NUMBER_IN_TXN; i++ {
		newValue := values[i]
		if i < NUMBER_IN_TXN/2 {
			newValue++
		} else {
			newValue--
		}
		err = myStore.PUT(tx, keys[i], newValue)
		if err != nil {
			myStore.Rollback(tx)
			return err
		}
	}

	err = myStore.Commit(tx)
	if err == nil {
		atomic.AddInt32(&successTxnCount, 1)
	} else {
		atomic.AddInt32(&failedTxnCount, 1)
	}

	return err
}
