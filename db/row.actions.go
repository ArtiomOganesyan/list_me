package db

func (s *Storage) CreateListRow(row ListRow) (int, error) {
	query := `INSERT INTO ListRow (list_id, title, description) VALUES ($1, $2, $3) RETURNING id;`
	var id int
	err := s.db.QueryRow(query, row.ListID, row.Title, row.Description).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *Storage) GetListRows(listID string) ([]ListRow, error) {
	query := `SELECT * FROM ListRow WHERE list_id = $1;`
	rows, err := s.db.Query(query, listID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var listRows []ListRow

	for rows.Next() {
		var row ListRow
		err = rows.Scan(&row.ID, &row.ListID, &row.Title, &row.Description, &row.Done, &row.CreatedAt)
		if err != nil {
			return nil, err
		}
		listRows = append(listRows, row)
	}

	return listRows, nil
}

func (s *Storage) ChangeRowState(rowID string, done bool) error {
	query := `UPDATE ListRow SET done = $1 WHERE id = $2;`
	_, err := s.db.Exec(query, done, rowID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) DeleteRow(rowID string) error {
	query := `DELETE FROM ListRow WHERE id = $1;`
	_, err := s.db.Exec(query, rowID)
	if err != nil {
		return err
	}
	return nil
}
