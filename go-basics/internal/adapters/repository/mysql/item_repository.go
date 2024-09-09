package mysqlrepo

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/osalomon89/go-basics/internal/core/domain"
	"github.com/osalomon89/go-basics/internal/core/ports"
)

type itemRepositoryImpl struct {
	conn *sqlx.DB
}

func NewMySQLRepository(conn *sqlx.DB) ports.ItemRepository {
	return &itemRepositoryImpl{
		conn: conn,
	}
}

func (s *itemRepositoryImpl) GetAllItems(ctx context.Context, limit int, searchAfter []interface{}) ([]domain.Item, []interface{}, error) {
	var items []domain.Item

	err := s.conn.Select(&items, "SELECT * FROM items LIMIT 10")
	if err != nil {
		return items, nil, fmt.Errorf("error getting all items: %w", err)
	}

	return items, nil, nil
}

func (s *itemRepositoryImpl) AddItem(ctx context.Context, item *domain.Item) error {
	createdAt := time.Now()

	result, err := s.conn.ExecContext(ctx, `INSERT INTO items 
		(code, title, description, categories, price, stock, available, created_at, updated_at) 
		VALUES(?,?,?,?,?,?,?,?,?)`, item.Code, item.Title, item.Description, item.Categories, item.Price,
		item.Stock, item.Available, createdAt, createdAt)

	if err != nil {
		return fmt.Errorf("error inserting item: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("error saving item: %w", err)
	}

	item.ID = strconv.FormatInt(id, 10)
	item.CreatedAt = &createdAt
	item.UpdatedAt = &createdAt

	return nil
}

func (s *itemRepositoryImpl) ReadItem(ctx context.Context, id string) (*domain.Item, error) {
	var item domain.Item

	idSql, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	err = s.conn.GetContext(ctx, &item, "SELECT * FROM items WHERE id=?", idSql)
	if err != nil {
		return nil, fmt.Errorf("error getting an item: %w", err)
	}

	return &item, nil
}

func (s *itemRepositoryImpl) Update(ctx context.Context, itemNew domain.Item) (*domain.Item, error) {
	return nil, fmt.Errorf("item not found")
}

func (s *itemRepositoryImpl) Delete(ctx context.Context, id string) error {
	return fmt.Errorf("item not found")
}
