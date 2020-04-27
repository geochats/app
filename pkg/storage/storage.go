package storage

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"geochats/pkg/types"
	"github.com/boltdb/bolt"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

type Storage interface {
	AddGroup(g *types.Group) error
	ListGroups() ([]types.Group, error)
	UpdatePoint(diff *types.Point) (*types.Point, error)
	GetPoint(chatID int64) (*types.Point, error)
	ListPoint() ([]types.Point, error)
}

type BoltStorage struct {
	db               *bolt.DB
	groupsBucketName []byte
	pointsBucketName []byte
}

func (b *BoltStorage) AddGroup(g *types.Group) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		buc := tx.Bucket(b.groupsBucketName)
		buf := bytes.Buffer{}
		if err := gob.NewEncoder(&buf).Encode(g); err != nil {
			return fmt.Errorf("failed gob Encode: %v", err)
		}
		if err := buc.Put(b.chatIDToBytes(g.ChatID), buf.Bytes()); err != nil {
			return fmt.Errorf("can't put value to bucket: %v", err)
		}
		return nil
	})
}

func (b *BoltStorage) ListGroups() ([]types.Group, error) {
	grs := make([]types.Group, 0)
	err := b.db.View(func(tx *bolt.Tx) error {
		buc := tx.Bucket(b.groupsBucketName)
		return buc.ForEach(func(k, v []byte) error {
			var g types.Group
			b := bytes.NewBuffer(v)
			if err := gob.NewDecoder(b).Decode(&g); err != nil {
				log.Errorf("failed gob.Decode Group with key `%s`: %v", k, err)
			} else {
				grs = append(grs, g)
			}
			return nil
		})
	})
	return grs, err
}

func (b *BoltStorage) UpdatePoint(diff *types.Point) (*types.Point, error) {
	var merged *types.Point
	return merged, b.db.Update(func(tx *bolt.Tx) error {
		buc := tx.Bucket(b.pointsBucketName)
		v := buc.Get(b.chatIDToBytes(diff.ChatID))
		if v == nil {
			merged = diff
		} else {
			b := bytes.NewBuffer(v)
			if err := gob.NewDecoder(b).Decode(&merged); err != nil {
				return fmt.Errorf("failed gob decode old Points: %v", err)
			}
			if diff.Photo.Width != 0 {
				merged.Photo = diff.Photo
			}
			if diff.Latitude != 0 {
				merged.Latitude = diff.Latitude
			}
			if diff.Longitude != 0 {
				merged.Longitude = diff.Longitude
			}
		}

		buf := bytes.Buffer{}
		if err := gob.NewEncoder(&buf).Encode(merged); err != nil {
			return fmt.Errorf("failed gob encode Points after merge: %v", err)
		}
		if err := buc.Put(b.chatIDToBytes(diff.ChatID), buf.Bytes()); err != nil {
			return fmt.Errorf("can't put Points value to bucket: %v", err)
		}
		return nil
	})
}

func (b *BoltStorage) GetPoint(chatID int64) (*types.Point, error) {
	var g types.Point
	err := b.db.View(func(tx *bolt.Tx) error {
		buc := tx.Bucket(b.pointsBucketName)
		v := buc.Get(b.chatIDToBytes(chatID))
		b := bytes.NewBuffer(v)
		if err := gob.NewDecoder(b).Decode(&g); err != nil {
			return fmt.Errorf("failed gob decode Points: %v", err)
		}
		return nil
	})
	return &g, err
}

func (b *BoltStorage) ListPoint() ([]types.Point, error) {
	grs := make([]types.Point, 0)
	err := b.db.View(func(tx *bolt.Tx) error {
		buc := tx.Bucket(b.pointsBucketName)
		return buc.ForEach(func(k, v []byte) error {
			var g types.Point
			b := bytes.NewBuffer(v)
			if err := gob.NewDecoder(b).Decode(&g); err != nil {
				log.Errorf("failed gob.Decode Points with key `%s`: %v", k, err)
			} else {
				grs = append(grs, g)
			}
			return nil
		})
	})
	return grs, err
}

func (b *BoltStorage) ListMottos() ([]types.Motto, error) {
	return []types.Motto{
		{1, "Требую введения ЧС!"},
	}, nil
}

func (b *BoltStorage) chatIDToBytes(chatID int64) []byte {
	return []byte(strconv.FormatInt(chatID, 10))
}

func New(name string) (Storage, error) {
	db, err := bolt.Open(name, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, fmt.Errorf("can't open bolt db: %v", err)
	}

	bucketName := []byte("groups")
	pointBucketName := []byte("points")
	err = db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists(bucketName); err != nil {
			return err
		}
		if _, err := tx.CreateBucketIfNotExists(pointBucketName); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("can't create buckets: %v", err)
	}

	return &BoltStorage{
		db:               db,
		groupsBucketName: bucketName,
		pointsBucketName: pointBucketName,
	}, nil
}