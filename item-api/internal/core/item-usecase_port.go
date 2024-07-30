package core

import "api/internal/core/item"

// ItemUsecasePort define a interface para o caso de uso de itens
type ItemUsecasePort interface {
	SaveItem(item.Item) error
	ListItems() (item.MapRepo, error)
}
