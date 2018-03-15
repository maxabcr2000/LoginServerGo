package repository

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/boltdb/bolt"
	domain "github.com/maxabcr2000/LoginServerGo/pkg/domain"
)

type boltRepository struct {
	db *bolt.DB
}

const (
	BUCKET_NAME_USERS = "Users"
)

var (
	usersBucket = []byte(BUCKET_NAME_USERS)
)

func CreateBoltRepository(dbName string) (*boltRepository, error) {
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(usersBucket)
		return err
	})
	if err != nil {
		return nil, err
	}

	repo := &boltRepository{
		db: db,
	}

	return repo, nil
}

func (r *boltRepository) SaveUser(user *domain.User, key string) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(usersBucket)
		v := b.Get([]byte(key))
		if v != nil {
			return errors.New("User account is already used")
		}

		body, _ := json.Marshal(user)
		err := b.Put([]byte(key), body)
		return err
	})
}

func (r *boltRepository) ReadUser(key string) (*domain.User, error) {
	user := &domain.User{}

	err := r.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(usersBucket)
		v := b.Get([]byte(key))
		err := json.Unmarshal(v, user)
		if err != nil {
			return err
		}

		fmt.Printf("key: %s , value: %s\n", key, user)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *boltRepository) Close() {
	r.db.Close()
}
