package item

import (
	"database/sql"
)

// mysqlRepository é uma implementação do repositório de itens utilizando MySQL
type mysqlRepository struct {
	db *sql.DB // Conexão com o banco de dados MySQL
}

// NewMySqlRepository cria uma nova instância de mysqlRepository
func NewMySqlRepository(db *sql.DB) ItemRepositoryPort {
	return &mysqlRepository{
		db: db,
	}
}

// SaveItem salva um novo item no banco de dados MySQL
func (r *mysqlRepository) SaveItem(it *Item) error {
	query := `INSERT INTO items (code, title, description, price, stock, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.Exec(query, it.Code, it.Title, it.Description, it.Price, it.Stock, it.Status, it.CreatedAt, it.UpdatedAt)
	return err
}

// ListItems lista todos os itens do banco de dados MySQL
func (r *mysqlRepository) ListItems() (MapRepo, error) {
	query := `SELECT id, code, title, description, price, stock, status, created_at, updated_at FROM items`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make(MapRepo)
	for rows.Next() {
		var it Item
		if err := rows.Scan(&it.ID, &it.Code, &it.Title, &it.Description, &it.Price, &it.Stock, &it.Status, &it.CreatedAt, &it.UpdatedAt); err != nil {
			return nil, err
		}
		items[it.ID] = it
	}

	return items, nil
}
