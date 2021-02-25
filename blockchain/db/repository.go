package db

import (
	"os"

	"github.com/dgraph-io/badger"
	blc "github.com/ysfkel/go-blockchain/blockchain/block"
	"github.com/ysfkel/go-blockchain/shared"
)

type IRepository interface {
	Get(key []byte) ([]byte, error)
	Update(block *blc.Block) error
	DBExists() bool
}

type Repository struct {
}

func (repo *Repository) Get(key []byte) ([]byte, error) {

	db := openDb()
	defer db.Close()

	var latestHash []byte

	err := db.View(func(txn *badger.Txn) error {

		item, err := txn.Get(key)

		if err != nil {
			return err
		}

		latestHash, err = item.ValueCopy(nil)

		return err
	})

	return latestHash, err

}

func (repo *Repository) Update(block *blc.Block) error {

	db := openDb()

	defer db.Close()

	err := db.Update(func(txn *badger.Txn) error {
		err := txn.Set(block.Hash, block.Serialize())
		shared.HandleError(err)
		err = txn.Set([]byte(shared.LatestHashKey), block.Hash)
		return err
	})

	return err

}

func NewRepository() *Repository {

	return &Repository{}
}

func openDb() *badger.DB {

	db, err := badger.Open(badger.DefaultOptions(shared.DbPath))

	shared.HandleError(err)

	return db
}

func (repo *Repository) DBExists() bool {

	if _, err := os.Stat(shared.DbFile); os.IsNotExist(err) {
		return false
	}

	return true
}
