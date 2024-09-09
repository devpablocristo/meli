package ports

import (
	"context"

	"github.com/osalomon89/go-basics/internal/core/domain"
)

type ItemService interface {
	GetAllItems(ctx context.Context, limit int, cursor []interface{}) ([]domain.Item, []interface{}, error)
	AddItem(ctx context.Context, item domain.Item) (*domain.Item, error)
	ReadItem(ctx context.Context, id string) *domain.Item
	UpdateItem(ctx context.Context, itemNew domain.Item) *domain.Item
}
