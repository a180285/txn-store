package txnStore

import (
	"github.com/LK4D4/trylock"
	"github.com/pkg/errors"
	"sync"
	"sync/atomic"
)

type OperationType int

const (
	OP_GET OperationType = iota
	OP_PUT
)

type KvOperation struct {
	opType   OperationType
	key      int
	value    int
	keyMutex *trylock.Mutex
}

type KvValue struct {
	value int
	m     *trylock.Mutex
}

type MyTxnStore struct {
	nextTxnId *int64

	kvStore sync.Map
	kvMutex sync.Mutex

	txnOperations sync.Map
}

func (txnStore *MyTxnStore) Name() string {
	return "MyTxnStore"
}

func NewMyTxnStore() *MyTxnStore {
	txnStore := &MyTxnStore{
		nextTxnId: new(int64),
	}

	keysCount := 1000
	for i := 0; i < keysCount; i++ {
		txnStore.kvStore.Store(i, KvValue{
			value: 0,
			m:     &trylock.Mutex{},
		})
	}

	return txnStore
}

func (txnStore *MyTxnStore) GET(tx interface{}, key int) (value int, err error) {
	txnId := tx.(int64)

	rawKvValue, ok := txnStore.kvStore.Load(key)
	if !ok {
		return 0, errors.Errorf("Could not get key: %s", key)
	}
	kvValue := rawKvValue.(KvValue)

	if !kvValue.m.TryLock() {
		return 0, errors.Errorf("Cannot get lock")
	}

	rawKvValue, ok = txnStore.kvStore.Load(key)
	if !ok {
		return 0, errors.Errorf("Could not get key: %s", key)
	}
	kvValue = rawKvValue.(KvValue)

	ops := txnStore.getOperationByTxnId(txnId)
	ops = append(ops, KvOperation{
		opType:   OP_GET,
		key:      key,
		value:    kvValue.value,
		keyMutex: kvValue.m,
	})

	txnStore.txnOperations.Store(txnId, ops)

	return kvValue.value, nil
}

func (txnStore *MyTxnStore) PUT(tx interface{}, key, value int) error {
	txnId := tx.(int64)

	ops := txnStore.getOperationByTxnId(txnId)
	ops = append(ops, KvOperation{
		opType: OP_PUT,
		key:    key,
		value:  value,
	})

	txnStore.txnOperations.Store(txnId, ops)

	return nil
}

func (txnStore *MyTxnStore) Begin() (tx interface{}, err error) {
	newTxnId := atomic.AddInt64(txnStore.nextTxnId, 1)

	txnStore.txnOperations.Store(newTxnId, []KvOperation{})

	return newTxnId, nil
}

func (txnStore *MyTxnStore) Commit(tx interface{}) error {
	txnId := tx.(int64)

	txnOperations := txnStore.getOperationByTxnId(txnId)
	txnStore.txnOperations.Delete(txnId)

	txnStore.kvMutex.Lock()
	defer txnStore.kvMutex.Unlock()

	for _, operation := range txnOperations {
		if operation.opType == OP_PUT {
			rawKvValue, ok := txnStore.kvStore.Load(operation.key)
			if !ok {
				return errors.Errorf("Could not get key: %s", operation.key)
			}
			oldKvValue := rawKvValue.(KvValue)

			keyMutex := oldKvValue.m

			newValue := KvValue{
				value: operation.value,
				m:     keyMutex,
			}
			txnStore.kvStore.Store(operation.key, newValue)

			keyMutex.Unlock()
		}
	}

	return nil
}

func (txnStore *MyTxnStore) getOperationByTxnId(txnId int64) []KvOperation {
	value, ok := txnStore.txnOperations.Load(txnId)
	if !ok {
		return []KvOperation{}
	}
	return value.([]KvOperation)
}

func (txnStore *MyTxnStore) Rollback(tx interface{}) error {
	txnId := tx.(int64)

	txnOperations := txnStore.getOperationByTxnId(txnId)
	txnStore.txnOperations.Delete(txnId)

	for _, operation := range txnOperations {
		if operation.opType == OP_GET {
			operation.keyMutex.Unlock()
		}
	}

	return nil
}
