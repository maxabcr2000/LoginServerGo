package repository

import (
	"github.com/boltdb/bolt"
	"fmt"
	"errors"
	"encoding/json"
	domain "github.com/maxabcr2000/LoginServerGo/pkg/domain"
)

type boltRepository struct{
	db *bolt.DB
}

const (
	BUCKET_NAME_USERS = "Users"
)

func CreateBoltRepository(dbName string) (*boltRepository,error) {
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(BUCKET_NAME_USERS))
		return err
	})
	if err!=nil{
		return nil, err
	}

	repo := &boltRepository{
		db:db,
	}

	return repo, nil
}

func (repo *boltRepository) SaveUser(user *domain.User, key string) error{
	return repo.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BUCKET_NAME_USERS))
		v := b.Get([]byte(key))
		if v!=nil{
			return errors.New("User account is already used")
		}

		body,_:= json.Marshal(user)
		err := b.Put([]byte(key), body)
		return err
	})
}

func (repo *boltRepository) ReadUser(key string) (*domain.User, error){
	user:= &domain.User{}
	
	err:=repo.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BUCKET_NAME_USERS))
		v := b.Get([]byte(key))
		err := json.Unmarshal(v, user)
		if err!=nil{
			return err
		}

		fmt.Printf("key: %s , value: %s\n", key, user)
		return nil
	})
	if err!=nil{
		return nil, err
	}

	return user, nil
}

func (repo *boltRepository) Close(){
	repo.db.Close()
}
