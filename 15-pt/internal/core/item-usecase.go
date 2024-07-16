package core

import (
	"fmt"

	"api/internal/core/item"
	"api/pkg/config"
)

// ItemUsecase representa o caso de uso para os itens
type ItemUsecase struct {
	repo item.ItemRepositoryPort // Repositório de itens
}

// NewItemUsecase cria uma nova instância de ItemUsecase
func NewItemUsecase(repo item.ItemRepositoryPort) ItemUsecasePort {
	return &ItemUsecase{
		repo: repo,
	}
}

// SaveItem salva um novo item no repositório
func (u *ItemUsecase) SaveItem(it item.Item) error {
	if err := u.repo.SaveItem(&it); err != nil {
		return fmt.Errorf("error saving item: %w", err)
	}
	return nil
}

// ListItems lista todos os itens do repositório
func (u *ItemUsecase) ListItems() (item.MapRepo, error) {
	its, err := u.repo.ListItems()
	if err != nil {
		return nil, fmt.Errorf("error in repository: %w", err)
	}
	if len(its) == 0 {
		return nil, config.ErrNotFound
	}
	return its, nil
}
