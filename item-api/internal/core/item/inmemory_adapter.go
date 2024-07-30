package item

import (
	"fmt"
)

// inMemoryRepository é uma implementação em memória do repositório de itens
type inMemoryRepository struct {
	items MapRepo // Mapa de itens
}

// NewRepository cria uma nova instância de inMemoryRepository
func NewInMemoryRepository() RepositoryPort {
	return &inMemoryRepository{
		items: make(MapRepo),
	}
}

// SaveItem salva um novo item no repositório
func (r *inMemoryRepository) SaveItem(it *Item) error {
	if it.ID == 0 {
		return fmt.Errorf("item ID cannot be 0")
	}
	if _, exists := r.items[it.ID]; exists {
		return fmt.Errorf("item with ID %d already exists", it.ID)
	}
	r.items[it.ID] = *it
	return nil
}

// ListItems lista todos os itens no repositório
func (r *inMemoryRepository) ListItems() (MapRepo, error) {
	return r.items, nil
}
