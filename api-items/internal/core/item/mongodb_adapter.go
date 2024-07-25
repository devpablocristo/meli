package item

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// mongoDbRepository es una implementación del repositorio de elementos utilizando MongoDB
type mongoDbRepository struct {
	db *mongo.Database // Conexión a la base de datos MongoDB
}

// NewMongoDBRepository crea una nueva instancia de mongoDbRepository
func NewMongoDbRepository(db *mongo.Database) RepositoryPort {
	return &mongoDbRepository{
		db: db,
	}
}

// SaveItem guarda un nuevo elemento en la base de datos MongoDB
func (r *mongoDbRepository) SaveItem(it *Item) error {
	it.CreatedAt = time.Now()
	it.UpdatedAt = time.Now()
	_, err := r.db.Collection("items").InsertOne(context.TODO(), it)
	return err
}

// ListItems lista todos los elementos de la base de datos MongoDB
func (r *mongoDbRepository) ListItems() (MapRepo, error) {
	cursor, err := r.db.Collection("items").Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	items := make(MapRepo)
	for cursor.Next(context.TODO()) {
		var it Item
		if err := cursor.Decode(&it); err != nil {
			return nil, err
		}
		items[it.ID] = it
	}

	return items, nil
}
