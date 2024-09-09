package ports

import (
	"context"

	"github.com/osalomon89/go-basics/internal/core/domain"
)

type ItemRepository interface {
	GetAllItems(ctx context.Context, limit int, searchAfter []interface{}) ([]domain.Item, []interface{}, error)
	AddItem(ctx context.Context, item *domain.Item) error
	ReadItem(ctx context.Context, id string) (*domain.Item, error)
	Update(ctx context.Context, itemNew domain.Item) (*domain.Item, error)
	Delete(ctx context.Context, id string) error
}

type ProviderClient interface {
	GetProvider(id int) (domain.Provider, error)
}
