package txnStore

const MAX_KEYS = 1000

type TxnStore interface {
	Name() string

	Begin() (tx interface{}, err error)

	Commit(tx interface{}) error

	GET(tx interface{}, key int) (value int, err error)

	PUT(tx interface{}, key, value int) error
}
