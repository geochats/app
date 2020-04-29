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
	GetConn() *bolt.DB
	GetGroup(tx *bolt.Tx, chatID int64) (*types.Group, error)
	SaveGroup(tx *bolt.Tx, g *types.Group) error
	ListGroups(tx *bolt.Tx) ([]types.Group, error)
	GetPoint(tx *bolt.Tx, chatID int64) (*types.Point, error)
	SavePoint(tx *bolt.Tx, point *types.Point) error
	ListPoint(tx *bolt.Tx) ([]types.Point, error)
}

type BoltStorage struct {
	db           *bolt.DB
	groupsBucket []byte
	pointsBucket []byte
}

func (b *BoltStorage) GetConn() *bolt.DB {
	return b.db
}

func (b *BoltStorage) GetGroup(tx *bolt.Tx, chatID int64) (*types.Group, error) {
	var group *types.Group
	groupGob := tx.Bucket(b.groupsBucket).Get(b.chatIDToBytes(chatID))
	if groupGob == nil {
		return nil, nil
	}
	if err := gob.NewDecoder(bytes.NewBuffer(groupGob)).Decode(&group); err != nil {
		return nil, fmt.Errorf("failed gob.Decode point: %v", err)
	}
	return group, nil
}

func (b *BoltStorage) SaveGroup(tx *bolt.Tx, group *types.Group) error {
	buf := bytes.Buffer{}
	if err := gob.NewEncoder(&buf).Encode(group); err != nil {
		return fmt.Errorf("failed gob.Encode group: %v", err)
	}
	if err := tx.Bucket(b.groupsBucket).Put(b.chatIDToBytes(group.ChatID), buf.Bytes()); err != nil {
		return fmt.Errorf("can't put group value to bucket: %v", err)
	}
	return nil
}

func (b *BoltStorage) ListGroups(tx *bolt.Tx) ([]types.Group, error) {
	groups := make([]types.Group, 0)
	err := tx.Bucket(b.groupsBucket).ForEach(func(k, v []byte) error {
		var g types.Group
		if err := gob.NewDecoder(bytes.NewBuffer(v)).Decode(&g); err != nil {
			log.Errorf("failed gob.Decode Group with chatID `%s`: %v", string(k), err)
		} else {
			groups = append(groups, g)
		}
		return nil
	})
	return groups, err
}

func (b *BoltStorage) GetPoint(tx *bolt.Tx, chatID int64) (*types.Point, error) {
	var point *types.Point
	pointGob := tx.Bucket(b.pointsBucket).Get(b.chatIDToBytes(chatID))
	if pointGob == nil {
		return nil, nil
	}
	if err := gob.NewDecoder(bytes.NewBuffer(pointGob)).Decode(&point); err != nil {
		return nil, fmt.Errorf("failed gob decode Points: %v", err)
	}
	return point, nil
}

func (b *BoltStorage) SavePoint(tx *bolt.Tx, point *types.Point) error {
	buf := bytes.Buffer{}
	if err := gob.NewEncoder(&buf).Encode(point); err != nil {
		return fmt.Errorf("can't gob.Encode point to save: %v", err)
	}
	if err := tx.Bucket(b.pointsBucket).Put(b.chatIDToBytes(point.ChatID), buf.Bytes()); err != nil {
		return fmt.Errorf("can't put Points value to bucket: %v", err)
	}
	return nil
}

func (b *BoltStorage) ListPoint(tx *bolt.Tx) ([]types.Point, error) {
	points := make([]types.Point, 0)
	err := tx.Bucket(b.pointsBucket).ForEach(func(k, v []byte) error {
		var p types.Point
		if err := gob.NewDecoder(bytes.NewBuffer(v)).Decode(&p); err != nil {
			log.Errorf("failed gob.Decode points with for chatID `%s`: %v", string(k), err)
		} else {
			points = append(points, p)
		}
		return nil
	})
	return points, err
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
		db:           db,
		groupsBucket: bucketName,
		pointsBucket: pointBucketName,
	}, nil
}
