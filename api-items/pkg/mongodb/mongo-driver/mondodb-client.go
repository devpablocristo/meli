package mongodbdriver

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDBClient representa um cliente para interagir com um banco de dados MongoDB
type MongoDBClient struct {
	config MongoDBClientConfig // Configuração do cliente MongoDB
	db     *mongo.Database     // Conexão com o banco de dados
}

// NewMongoDBClient cria uma nova instância de MongoDBClient e estabelece a conexão com o banco de dados
func NewMongoDBClient(config MongoDBClientConfig) (*MongoDBClient, error) {
	client := &MongoDBClient{config: config}
	err := client.connect()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize MongoDBClient: %v", err)
	}
	return client, nil
}

// connect estabelece a conexão com o banco de dados MongoDB utilizando a configuração fornecida
func (client *MongoDBClient) connect() error {
	dsn := client.config.dsn()
	clientOptions := options.Client().ApplyURI(dsn)

	conn, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Verificar a conexão
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = conn.Ping(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	client.db = conn.Database(client.config.Database)
	return nil
}

// Close fecha a conexão com o banco de dados MongoDB
func (client *MongoDBClient) Close() {
	if client.db != nil {
		client.db.Client().Disconnect(context.TODO())
	}
}

// DB retorna a conexão com o banco de dados MongoDB
func (client *MongoDBClient) DB() *mongo.Database {
	return client.db
}
