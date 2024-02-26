package db

func (s *Storage) CreateList(list *List) error {
	query := `INSERT INTO List (id, secret) VALUES ($1, $2);`
	_, err := s.db.Exec(query, list.ID, list.Secret)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetList(id, secret string) (List, error) {
	query := `SELECT * FROM List WHERE id=$1 AND secret=$2;`
	row := s.db.QueryRow(query, id, secret)
	list := List{}
	err := row.Scan(&list.ID, &list.Secret, &list.CreatedAt)
	if err != nil {
		return list, err
	}
	return list, nil
}
