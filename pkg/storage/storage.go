package storage

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/boltdb/bolt"
	"geochats/pkg/types"
	"time"
)

type Storage interface {
	AddGroup(g *types.Group) error
	ListGroups() ([]types.Group, error)
}

type BoltStorage struct {
	db *bolt.DB
	bucketName []byte
}

func (b *BoltStorage) AddGroup(g *types.Group) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		buc := tx.Bucket(b.bucketName)
		b := bytes.Buffer{}
		if err := gob.NewEncoder(&b).Encode(g); err != nil {
			return fmt.Errorf("failed gob Encode: %v", err)
		}
		if err := buc.Put([]byte(fmt.Sprintf("%d", g.ChatID)), b.Bytes()); err != nil {
			return fmt.Errorf("can't put value to bucket: %v", err)
		}
		return nil
	})
}

func (b *BoltStorage) ListGroups() ([]types.Group, error) {
	grs := make([]types.Group, 0)
	err := b.db.View(func(tx *bolt.Tx) error {
		buc := tx.Bucket(b.bucketName)
		return buc.ForEach(func(k, v []byte) error {
			var g types.Group
			b := bytes.NewBuffer(v)
			if err := gob.NewDecoder(b).Decode(&g); err != nil {
				return fmt.Errorf("failed gob Decode: %v", err)
			}
			grs = append(grs, g)
			return nil
		})
	})
	return grs, err
}

func New(name string) (Storage, error) {
	db, err := bolt.Open(name, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, fmt.Errorf("can't open bolt db: %v", err)
	}

	bucketName := []byte("groups")
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucketName)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("can't create bucket: %v", err)
	}

	return &BoltStorage{
		db: db,
		bucketName: bucketName,
	}, nil
}