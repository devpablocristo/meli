package item

// RepositoryPort define a interface para o repositório de itens
type RepositoryPort interface {
	SaveItem(*Item) error
	ListItems() (MapRepo, error)
}
