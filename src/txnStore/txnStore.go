package txnStore

const MAX_KEYS = 1000

type TxnStore interface {
	Name() string

	// Start a transaction, and pass the return value to following methods.
	Begin() (tx interface{}, err error)

	GET(tx interface{}, key int) (value int, err error)
	PUT(tx interface{}, key, value int) error
	Commit(tx interface{}) error
	Rollback(tx interface{}) error
}
