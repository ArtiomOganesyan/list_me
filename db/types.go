package db

type List struct {
	ID        string    `json:"id"`
	Secret    string    `json:"secret"`
	Rows      []ListRow `json:"rows"`
	CreatedAt string    `json:"created_at"`
}

type ListRow struct {
	ID          string `json:"id"`
	ListID      string `json:"list_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
	CreatedAt   string `json:"created_at"`
}
