package storage

import (
	"context"
	"fmt"
	"geochats/pkg/types"
	"github.com/jackc/pgx/v4"
)

type Storage struct {
	conn *pgx.Conn
}

func New(dsn string) (*Storage, error) {
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("can't connect to db: %v", err)
	}
	return &Storage{conn: conn}, nil
}

func (s *Storage) IsReady() bool {
	return s.conn.Ping(context.Background()) == nil
}

func (s *Storage) Begin(writable bool) (pgx.Tx, error) {
	return s.conn.Begin(context.Background())
	//access := pgx.ReadOnly
	//if writable {
	//	access = pgx.ReadWrite
	//}
	//return s.conn.BeginTx(context.Background(), pgx.TxOptions{
	//	IsoLevel:       pgx.RepeatableRead,
	//	AccessMode:     access,
	//	DeferrableMode: pgx.NotDeferrable,
	//})
}

func (s *Storage) InTransaction(writable bool, fn func(tx pgx.Tx) error) error {
	tx, err := s.Begin(writable)
	if err != nil {
		return fmt.Errorf("can't start transaction: %v", err)
	}
	defer func() { _ = tx.Rollback(context.Background()) }()

	err = fn(tx)
	if err != nil {
		return fmt.Errorf("in-transaction callback return an errors: %v", err)
	}

	if err := tx.Commit(context.Background()); err != nil {
		return fmt.Errorf("can't commit in-transaction: %v", err)
	}

	return nil
}

func (s *Storage) GetPoint(tx pgx.Tx, chatID int64) (*types.Point, error) {
	point := &types.Point{ChatID: chatID}
	err := tx.
		QueryRow(
			context.Background(),
			"SELECT username, text, latitude, longitude, members_count, is_published, is_single FROM points WHERE chat_id=$1",
			chatID).
		Scan(
			&point.Username,
			&point.Text,
			&point.Latitude,
			&point.Longitude,
			&point.MembersCount,
			&point.Published,
			&point.IsSingle)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("can't select point: %v", err)
	}
	return point, nil
}

func (s *Storage) AddPoint(tx pgx.Tx, chatID int64, isSingle bool) (*types.Point, error) {
	_, err := tx.Exec(
			context.Background(),
			"INSERT INTO points(chat_id, is_single) VALUES ($1, $2)",
			chatID,
			isSingle)
	if err != nil {
		return nil, fmt.Errorf("can't insert point: %v", err)
	}
	return &types.Point{
		ChatID: chatID,
		IsSingle: isSingle,
	}, nil
}

func (s *Storage) UpdatePoint(tx pgx.Tx, group *types.Point) error {
	_, err := tx.Exec(
		context.Background(),
		"UPDATE points SET username=$2, text=$3, latitude=$4, longitude=$5, members_count=$6, is_published=$7, is_single=$8 WHERE chat_id=$1",
		group.ChatID,
		group.Username,
		group.Text,
		group.Latitude,
		group.Longitude,
		group.MembersCount,
		group.Published,
		group.IsSingle,
		)
	if err != nil {
		return fmt.Errorf("can't update point")
	}
	return nil
}

func (s *Storage) ListGroups(tx pgx.Tx) ([]types.Point, error) {
	points := make([]types.Point, 0)
	rows, _ := tx.Query(
		context.Background(),
		"SELECT chat_id, username, text, latitude, longitude, members_count, is_published, is_single FROM points LIMIT 10000")
	for rows.Next() {
		p := types.Point{}
		err := rows.Scan(
			&p.ChatID,
			&p.Username,
			&p.Text,
			&p.Latitude,
			&p.Longitude,
			&p.MembersCount,
			&p.Published,
			&p.IsSingle)
		if err != nil {
			return nil, err
		}
		points = append(points, p)
	}
	return points, rows.Err()
}

