package item

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoDBRepository es una implementación del repositorio de elementos utilizando MongoDB
type MongoDBRepository struct {
	db *mongo.Database // Conexión a la base de datos MongoDB
}

// NewMongoDBRepository crea una nueva instancia de MongoDBRepository
func NewMongoDBRepository(db *mongo.Database) ItemRepositoryPort {
	return &MongoDBRepository{
		db: db,
	}
}

// SaveItem guarda un nuevo elemento en la base de datos MongoDB
func (r *MongoDBRepository) SaveItem(it *Item) error {
	it.CreatedAt = time.Now()
	it.UpdatedAt = time.Now()
	_, err := r.db.Collection("items").InsertOne(context.TODO(), it)
	return err
}

// ListItems lista todos los elementos de la base de datos MongoDB
func (r *MongoDBRepository) ListItems() (MapRepo, error) {
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
