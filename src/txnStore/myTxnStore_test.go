package txnStore_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"txn-store/src/txnStore"
)

func TestNewMyTxnStore(t *testing.T) {
	_ = txnStore.NewMyTxnStore()
}

func TestMyTxnStore_EmptyTransaction(t *testing.T) {
	myStore := txnStore.NewMyTxnStore()

	tx, err := myStore.Begin()
	assert.NoError(t, err)

	err = myStore.Commit(tx)
	assert.NoError(t, err)
}

func TestMyTxnStore_OneTransaction(t *testing.T) {
	myStore := txnStore.NewMyTxnStore()

	tx, err := myStore.Begin()
	assert.NoError(t, err)

	err = myStore.PUT(tx, 1, 2)
	assert.NoError(t, err)

	err = myStore.Commit(tx)
	assert.NoError(t, err)

	value := getKey(myStore, 1)
	assert.Equal(t, 2, value)
}

func TestMyTxnStore_NotReadUncommitedData(t *testing.T) {
	myStore := txnStore.NewMyTxnStore()

	tx1, err := myStore.Begin()
	assert.NoError(t, err)
	tx2, err := myStore.Begin()
	assert.NoError(t, err)

	err = myStore.PUT(tx1, 1, 2)
	assert.NoError(t, err)

	value2, err := myStore.GET(tx2, 1)
	assert.NoError(t, err)
	assert.Equal(t, 0, value2)
}

func TestMyTxnStore_TransactionConflict1(t *testing.T) {
	myStore := txnStore.NewMyTxnStore()

	tx1, err := myStore.Begin()
	assert.NoError(t, err)
	tx2, err := myStore.Begin()
	assert.NoError(t, err)

	myStore.GET(tx1, 1)
	err = myStore.PUT(tx1, 1, 2)
	assert.NoError(t, err)

	value2, err := myStore.GET(tx2, 1)
	assert.NoError(t, err)
	assert.Equal(t, 0, value2)

	err = myStore.PUT(tx2, 1, 3)
	assert.NoError(t, err)

	err = myStore.Commit(tx1)
	assert.NoError(t, err)

	err = myStore.Commit(tx2)
	assert.Error(t, err)

	assert.Equal(t, 2, getKey(myStore, 1))
}

func TestMyTxnStore_TransactionConflict2(t *testing.T) {
	myStore := txnStore.NewMyTxnStore()

	tx1, err := myStore.Begin()
	assert.NoError(t, err)
	tx2, err := myStore.Begin()
	assert.NoError(t, err)

	myStore.GET(tx1, 1)
	err = myStore.PUT(tx1, 1, 2)
	assert.NoError(t, err)

	value2, err := myStore.GET(tx2, 1)
	assert.NoError(t, err)
	assert.Equal(t, 0, value2)

	err = myStore.PUT(tx2, 1, 3)
	assert.NoError(t, err)

	err = myStore.Commit(tx2) // commit tx2 first
	assert.NoError(t, err)

	err = myStore.Commit(tx1)
	assert.Error(t, err)

	assert.Equal(t, 3, getKey(myStore, 1))
}

func TestMyTxnStore_FailedTransactionWillRollback(t *testing.T) {
	myStore := txnStore.NewMyTxnStore()

	tx1, err := myStore.Begin()
	assert.NoError(t, err)
	tx2, err := myStore.Begin()
	assert.NoError(t, err)

	myStore.GET(tx1, 1)
	err = myStore.PUT(tx1, 1, 2)
	assert.NoError(t, err)

	value2, err := myStore.GET(tx2, 1)
	assert.NoError(t, err)
	assert.Equal(t, 0, value2)

	err = myStore.PUT(tx2, 3, 3)
	err = myStore.PUT(tx2, 2, 3)
	err = myStore.PUT(tx2, 1, 3)
	assert.NoError(t, err)

	err = myStore.Commit(tx1)
	assert.NoError(t, err)

	err = myStore.Commit(tx2)
	assert.Error(t, err)

	assert.Equal(t, 2, getKey(myStore, 1))
	assert.Equal(t, 0, getKey(myStore, 2))
	assert.Equal(t, 0, getKey(myStore, 3))
}

func getKey(myStore *txnStore.MyTxnStore, key int) int {
	tx, _ := myStore.Begin()

	v, _ := myStore.GET(tx, key)

	myStore.Commit(tx)

	return v
}
