package core

import (
	"fmt"
	"time"

	"api/internal/core/item"
)

// itemUseCase representa o caso de uso para os itens
type itemUseCase struct {
	mysqlRepo item.RepositoryPort // Repositório de MySQL
	mapRepo   item.RepositoryPort // Repositório de Map
}

// NewItemUsecase cria uma nova instância de itemUseCase
func NewItemUsecase(mysqlRepo, mapRepo item.RepositoryPort) ItemUsecasePort {
	return &itemUseCase{
		mysqlRepo: mysqlRepo,
		mapRepo:   mapRepo,
	}
}

// SaveItem salva um novo item em ambos os repositórios
func (u *itemUseCase) SaveItem(it item.Item) error {
	now := time.Now()
	it.CreatedAt = now
	it.UpdatedAt = now

	if err := u.mysqlRepo.SaveItem(&it); err != nil {
		return fmt.Errorf("error saving item in MySQL: %w", err)
	}
	if err := u.mapRepo.SaveItem(&it); err != nil {
		return fmt.Errorf("error saving item in MapRepo: %w", err)
	}
	return nil
}

// ListItems lista todos os itens de ambos os repositórios e os combina
func (u *itemUseCase) ListItems() (item.MapRepo, error) {
	mysqlItems, err := u.mysqlRepo.ListItems()
	if err != nil {
		return nil, fmt.Errorf("error listing items from MySQL: %w", err)
	}

	mapItems, err := u.mapRepo.ListItems()
	if err != nil {
		return nil, fmt.Errorf("error listing items from MapRepo: %w", err)
	}

	// Combina os resultados de ambos os repositórios
	for k, v := range mapItems {
		mysqlItems[k] = v
	}

	return mysqlItems, nil
}
