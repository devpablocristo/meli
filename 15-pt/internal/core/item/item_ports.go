package item

// ItemRepositoryPort define a interface para o repositório de itens
type ItemRepositoryPort interface {
	SaveItem(*Item) error
	ListItems() (MapRepo, error)
}
