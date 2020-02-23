package txnStore

import (
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
	opType       OperationType
	key          int
	value        int
	valueVersion int
}

type KvValue struct {
	value   int
	version int
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
		txnStore.kvStore.Store(i, KvValue{})
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

	ops := txnStore.getOperationByTxnId(txnId)
	ops = append(ops, KvOperation{
		opType:       OP_GET,
		key:          key,
		value:        kvValue.value,
		valueVersion: kvValue.version,
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
		rawKvValue, ok := txnStore.kvStore.Load(operation.key)
		if !ok {
			return errors.Errorf("Could not get key: %s", operation.key)
		}
		oldKvValue := rawKvValue.(KvValue)

		if operation.opType == OP_GET && operation.valueVersion != oldKvValue.version {
			return errors.Errorf("Data has been modified, transaction [%d] cann't be commited.", txnId)
		}
	}

	for _, operation := range txnOperations {
		if operation.opType == OP_PUT {
			rawKvValue, ok := txnStore.kvStore.Load(operation.key)
			if !ok {
				return errors.Errorf("Could not get key: %s", operation.key)
			}
			oldKvValue := rawKvValue.(KvValue)

			newValue := KvValue{
				value:   operation.value,
				version: oldKvValue.version + 1,
			}
			txnStore.kvStore.Store(operation.key, newValue)
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
