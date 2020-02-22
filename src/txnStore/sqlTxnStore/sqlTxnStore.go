package sqlTxnStore

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"os"
	"path"
	"txn-store/src/txnStore"
)

type SqlTxnStore struct {
	name string
	db   *gorm.DB
}

func (txnStore *SqlTxnStore) Name() string {
	return txnStore.name
}

type KvValue struct {
	// use *int to avoid gorm framework ignore zero values
	// http://gorm.io/docs/create.html#Default-Values
	Key   *int `gorm:"primary_key;auto_increment:false"`
	Value *int
}

func newInt(v int) *int {
	return &v
}

func newSqlTxnStore(db *gorm.DB, name string) (*SqlTxnStore, error) {
	db.DropTableIfExists(&KvValue{})
	db.AutoMigrate(&KvValue{})

	tx := db.Begin()
	for i := 0; i < txnStore.MAX_KEYS; i++ {
		err := tx.Save(&KvValue{Key: newInt(i), Value: newInt(0)}).Error
		//println("i : ", i)
		if err != nil {
			return nil, err
		}
	}

	tx.Commit()

	log.Printf("%s init done.", name)

	return &SqlTxnStore{db: db, name: name}, nil
}

func NewSQLiteTxnStore() (txnStore.TxnStore, error) {
	dbFilePath := path.Join(os.TempDir(), "txn-test.db")
	os.Remove(dbFilePath) // try remove old file if exists

	db, err := gorm.Open("sqlite3", dbFilePath)
	if err != nil {
		return nil, err
	}
	return newSqlTxnStore(db, "sqlite3")
}

func NewMysqlTxnStore() (txnStore.TxnStore, error) {
	db, err := gorm.Open("mysql", "root:123456@(localhost:33306)/txn")
	if err != nil {
		return nil, err
	}
	return newSqlTxnStore(db, "mysql")
}

func (txnStore *SqlTxnStore) Begin() (tx interface{}, err error) {
	return txnStore.db.Begin(), nil
}

func (txnStore *SqlTxnStore) Commit(tx interface{}) error {
	rawTx := tx.(*gorm.DB)
	return rawTx.Commit().Error
}

func (txnStore *SqlTxnStore) GET(tx interface{}, key int) (value int, err error) {
	rawTx := tx.(*gorm.DB)
	kvValue := &KvValue{}
	if err := rawTx.Set("gorm:query_option", "FOR UPDATE").First(kvValue, key).Error; err != nil {
		return 0, err
	}

	return *kvValue.Value, nil
}

func (txnStore *SqlTxnStore) PUT(tx interface{}, key, value int) error {
	rawTx := tx.(*gorm.DB)

	err := rawTx.Model(&KvValue{}).Updates(&KvValue{
		Key:   &key,
		Value: &value,
	}).Error

	return err
}

// docker run --name txn-mysql --rm -e MYSQL_ROOT_PASSWORD=123456 -e MYSQL_DATABASE=txn -p 33306:3306 -it mysql:5.7.29
// docker stop txn-mysql
