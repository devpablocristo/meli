package core

import (
	"fmt"

	"api/internal/core/item"
)

// ItemUsecase representa el caso de uso para los elementos
type ItemUsecase struct {
	mysqlRepo item.ItemRepositoryPort // Repositorio de MySQL
	mapRepo   item.ItemRepositoryPort // Repositorio de Map
}

// NewItemUsecase crea una nueva instancia de ItemUsecase
func NewItemUsecase(mysqlRepo, mapRepo item.ItemRepositoryPort) ItemUsecasePort {
	return &ItemUsecase{
		mysqlRepo: mysqlRepo,
		mapRepo:   mapRepo,
	}
}

// SaveItem guarda un nuevo elemento en ambos repositorios
func (u *ItemUsecase) SaveItem(it item.Item) error {
	if err := u.mysqlRepo.SaveItem(&it); err != nil {
		return fmt.Errorf("error saving item in MySQL: %w", err)
	}
	if err := u.mapRepo.SaveItem(&it); err != nil {
		return fmt.Errorf("error saving item in MapRepo: %w", err)
	}
	return nil
}

// ListItems lista todos los elementos de ambos repositorios y los combina
func (u *ItemUsecase) ListItems() (item.MapRepo, error) {
	mysqlItems, err := u.mysqlRepo.ListItems()
	if err != nil {
		return nil, fmt.Errorf("error listing items from MySQL: %w", err)
	}

	mapItems, err := u.mapRepo.ListItems()
	if err != nil {
		return nil, fmt.Errorf("error listing items from MapRepo: %w", err)
	}

	// Combina los resultados de ambos repositorios
	for k, v := range mapItems {
		mysqlItems[k] = v
	}

	return mysqlItems, nil
}
