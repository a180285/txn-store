package txnStore

import (
	"github.com/pkg/errors"
	"sync"
	"sync/atomic"
	"unsafe"
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

	kvStore []unsafe.Pointer
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
	txnStore.kvStore = make([]unsafe.Pointer, 1000)
	for i := 0; i < keysCount; i++ {
		atomic.StorePointer(&txnStore.kvStore[i], unsafe.Pointer(&KvValue{}))
	}

	return txnStore
}

func (txnStore *MyTxnStore) GET(tx interface{}, key int) (value int, err error) {
	txnId := tx.(int64)

	kvValue := (*KvValue)(atomic.LoadPointer(&txnStore.kvStore[key]))

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
		oldKvValue := (*KvValue)(atomic.LoadPointer(&txnStore.kvStore[operation.key]))

		if operation.opType == OP_GET && operation.valueVersion != oldKvValue.version {
			return errors.Errorf("Data has been modified, transaction [%d] cann't be commited.", txnId)
		}
	}

	for _, operation := range txnOperations {
		if operation.opType == OP_PUT {
			oldKvValue := (*KvValue)(atomic.LoadPointer(&txnStore.kvStore[operation.key]))

			newValue := &KvValue{
				value:   operation.value,
				version: oldKvValue.version + 1,
			}
			atomic.StorePointer(&txnStore.kvStore[operation.key], unsafe.Pointer(newValue))
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
