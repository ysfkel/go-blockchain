package db

import (
	"github.com/dgraph-io/badger"
	blc "github.com/ysfkel/go-blockchain/blockchain/block"
	"github.com/ysfkel/go-blockchain/shared"
)

const (
	dbPath = "./tmp/blocks"
)

type IRepository interface {
	Get(key []byte) ([]byte, error)
	Update(block *blc.Block) error
	Close() error
}

type Repository struct {
	Database *badger.DB
}

func (repo *Repository) Get(key []byte) ([]byte, error) {

	var latestHash []byte

	err := repo.Database.View(func(txn *badger.Txn) error {

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

	err := repo.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(block.Hash, block.Serialize())
		shared.HandleError(err)
		err = txn.Set([]byte(shared.LatestHashKey), block.Hash)
		return err
	})

	return err

}

func NewRepository() *Repository {
	//open db
	db, err := badger.Open(badger.DefaultOptions(dbPath))

	shared.HandleError(err)

	return &Repository{Database: db}
}

func (repo *Repository) Close() error {

	err := repo.Database.Close()

	return err

}
