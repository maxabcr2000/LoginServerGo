package repository

import (
	"github.com/boltdb/bolt"
	"fmt"
)

type boltRepository struct{
	db *bolt.DB
}

func CreateBoltRepository(dbName string) (*boltRepository,error) {
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		return nil, err
	}

	repo := &boltRepository{
		db:db,
	}

	return repo, nil
}

func (repo *boltRepository) CreateBucket(bucketName string) error{
	return repo.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		return err
	})
}

func (repo *boltRepository) SaveMessage(bucketName, key, value string) error{
	err:=repo.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		err := b.Put([]byte(key), []byte(value))
		return err
	})
	if err!=nil{
		return err
	}
	return nil
}

func (repo *boltRepository) ReadMessage(bucketName, key string) (string, error){
	var resultBytes []byte
	err:=repo.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		v := b.Get([]byte(key))
		resultBytes = make([]byte, len(v))
		copy(resultBytes, v)
		fmt.Printf("key: %s , value: %s\n", key, resultBytes)
		return nil
	})
	if err!=nil{
		return "", err
	}

	return string(resultBytes), nil
}

func (repo *boltRepository) Close(){
	repo.db.Close()
}
