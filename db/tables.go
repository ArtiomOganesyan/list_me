package db

import (
	"fmt"

	_ "github.com/lib/pq"
)

func executeTableCreations(funcs ...func() error) error {
	for _, fn := range funcs {
		if err := fn(); err != nil {
			return err
		}
	}
	return nil
}

func (s *Storage) createListTable() error {
	query := `CREATE TABLE IF NOT EXISTS List (
		id uuid PRIMARY KEY,
		secret VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := s.db.Exec(query)

	if err != nil {
		fmt.Println("Error executing query: ", err)
		return err
	}

	return nil
}

func (s *Storage) createListRowTable() error {
	query := `CREATE TABLE IF NOT EXISTS ListRow (
		id serial PRIMARY KEY,
		list_id uuid REFERENCES List(id),
		title VARCHAR(255) NOT NULL,
		description VARCHAR(255),
		done BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := s.db.Exec(query)

	if err != nil {
		fmt.Println("Error executing query: ", err)
		return err
	}

	return nil
}
