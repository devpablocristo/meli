package mongodb

import (
	mongodbdriver "api/pkg/mongodb/mongo-driver"
)

// NewMongoDBSetup configura e inicializa uma conexão com MongoDB
func NewMongoDBSetup() (*mongodbdriver.MongoDBClient, error) {
	config := mongodbdriver.MongoDBClientConfig{
		User:     "root",              // Usuário do banco de dados
		Password: "rootpassword",      // Senha do usuário
		Host:     "mongodb",           // Host onde o banco de dados está localizado
		Port:     "27017",             // Porta na qual o banco de dados está ouvindo
		Database: "inventory_mongodb", // Nome do banco de dados
	}

	return mongodbdriver.NewMongoDBClient(config)
}
