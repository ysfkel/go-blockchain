package blockchain

import (
	"fmt"

	"github.com/dgraph-io/badger"
)

type IRepository interface {
	Get(key []byte) ([]byte, error)
	Update(block *Block) error
	Close() error
}

type Repository struct {
	Database *badger.DB
}

func (repo *Repository) Get(key []byte) ([]byte, error) {

	var latestHash []byte

	err := repo.Database.View(func(txn *badger.Txn) error {

		item, err := txn.Get(key)
		HandleError(err)
		latestHash, err = item.ValueCopy(nil)

		return err
	})

	return latestHash, err

}

func (repo *Repository) Update(block *Block) error {

	err := repo.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(block.Hash, block.Serialize())
		HandleError(err)
		err = txn.Set([]byte(latestHashKey), block.Hash)
		return err
	})

	return err

}

func initRepository() (*Repository, []byte) {
	var lastHash []byte
	//open db
	db, err := badger.Open(badger.DefaultOptions(dbPath))

	HandleError(err)

	err = db.Update(func(txn *badger.Txn) error {
		// if database is empty , add genesis block
		_, err := txn.Get([]byte(latestHashKey))

		if err == badger.ErrKeyNotFound {
			fmt.Println("No existing blockchain found")
			//generate genesis block
			genesis := Genesis()
			fmt.Println("Genesis proved")
			//add genesis to database
			err = txn.Set(genesis.Hash, genesis.Serialize())
			HandleError(err)
			//set genesis as last hash
			err = txn.Set([]byte(latestHashKey), genesis.Hash)

			lastHash = genesis.Hash

			return err

		} else {

			item, err := txn.Get([]byte(latestHashKey))
			HandleError(err)
			lastHash, err = item.ValueCopy(nil)

			return err
		}
	})

	HandleError(err)

	return &Repository{Database: db}, lastHash
}

func (repo *Repository) Close() error {

	err := repo.Database.Close()

	return err

}
