package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(connStr string) (*Storage, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return &Storage{db: db}, nil
}

func (s *Storage) Init() error {
	fmt.Println("Initializing database")
	err := executeTableCreations(
		s.createListTable,
		s.createListRowTable,
	)
	if err != nil {
		return err
	}
	fmt.Println("Database initialized")
	return nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}
