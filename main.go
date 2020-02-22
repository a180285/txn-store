package main

import (
	"context"
	"fmt"
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
	log.Printf("Start test now.")
	numThreads := 10000
	runningSeconds := 15

	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())

	myStore := txnStore.NewMyTxnStore()

	keyChanel := randKeyChannel()

	for i := 0; i < numThreads; i++ {
		wg.Add(1)
		go doTransactions(ctx, &wg, myStore, keyChanel)
	}

	//cpuprofile := "cpu.prof"
	//f, err := os.Create(cpuprofile)
	//if err != nil {
	//	log.Fatal("could not create CPU profile: ", err)
	//}
	//if err := pprof.StartCPUProfile(f); err != nil {
	//	log.Fatal("could not start CPU profile: ", err)
	//}

	time.Sleep(time.Duration(runningSeconds) * time.Second)

	//pprof.StopCPUProfile()

	cancel()
	wg.Wait()

	checkStore(myStore, runningSeconds)
}

func checkStore(mystore txnStore.TxnStore, runningSeconds int) {
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

	fmt.Printf("sucess txn count: %d, failed count: %d\n", successTxnCount, failedTxnCount)
	fmt.Printf("txn success QPS: %f, sum: %d, non zero count: %d\n", float64(successTxnCount)/float64(runningSeconds), sum, nonZeroCount)
}

func doTransactions(ctx context.Context, wg *sync.WaitGroup, myStore txnStore.TxnStore, keyChanel <-chan []int) {
	defer wg.Done()

	nextKeys := <-keyChanel
	for {
		select {
		case <-ctx.Done():
			return
		default:
			err := _doTransactions(myStore, nextKeys)
			if err == nil { // only try new keys when success
				nextKeys = <-keyChanel
			}
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
