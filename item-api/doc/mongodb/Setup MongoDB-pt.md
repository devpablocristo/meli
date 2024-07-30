## Configuração do Cliente MongoDB

### Resumo

1. **Inicialização do MongoDB com Docker Compose**: Crie banco de dados, usuário e senha.
2. **Definição de Configuração**: Crie uma estrutura `MongoDBClientConfig` para armazenar os detalhes de conexão e uma função para gerar a string de conexão (URI).
3. **Configuração do Cliente**: Implemente um cliente MongoDB (`MongoDBClient`) que utiliza a configuração para se conectar ao banco de dados.
4. **Injeção de Dependências**: Configure e inicialize o cliente MongoDB através da função `NewMongoDBSetup`.
5. **Repositório MongoDB**: Crie um repositório (`mongoRepository`) que utiliza o cliente MongoDB para realizar operações CRUD no banco de dados.

### Passos

#### Passo 1: Inicialização do MongoDB com Docker Compose

#### Configuração do Contêiner MongoDB:

- **MONGO_INITDB_ROOT_USERNAME**: Define o nome do usuário root que será criado ao inicializar o banco de dados.
- **MONGO_INITDB_ROOT_PASSWORD**: Define a senha para o usuário root que será criado ao inicializar o banco de dados.

#### Configuração do Contêiner Mongo Express:

- **ME_CONFIG_MONGODB_ADMINUSERNAME**: Define o nome do usuário que Mongo Express usará para se conectar ao MongoDB como administrador.
- **ME_CONFIG_MONGODB_ADMINPASSWORD**: Define a senha que Mongo Express usará para se conectar ao MongoDB como administrador.

**Nota:** As credenciais devem ser as mesmas tanto para MongoDB quanto para Mongo Express. Isso garante que Mongo Express possa se conectar ao MongoDB com permissões de administrador.

### Exemplo de configuração em `docker-compose.yml`:

```yaml
version: '3.8'

services:
  mongodb:
    image: mongo:latest
    container_name: mongodb
    environment:
      - MONGO_INITDB_ROOT_USERNAME=root
      - MONGO_INITDB_ROOT_PASSWORD=rootpassword
    ports:
      - "27017:27017"

  mongo-express:
    image: mongo-express:latest
    container_name: mongo-express
    environment:
      - ME_CONFIG_MONGODB_ADMINUSERNAME=root
      - ME_CONFIG_MONGODB_ADMINPASSWORD=rootpassword
      - ME_CONFIG_MONGODB_URL=mongodb://root:rootpassword@mongodb:27017/
    ports:
      - "8082:8081"
```

Nesta configuração:
- **MONGO_INITDB_ROOT_USERNAME** e **MONGO_INITDB_ROOT_PASSWORD** configuram o usuário root do MongoDB.
- **ME_CONFIG_MONGODB_ADMINUSERNAME** e **ME_CONFIG_MONGODB_ADMINPASSWORD** configuram Mongo Express para usar o usuário root do MongoDB para se conectar.

### Acesso ao Mongo Express

Você pode ver os usuários no Mongo Express acessando: [http://localhost:8082/db/admin/system.users](http://localhost:8082/db/admin/system.users)

### Acesso à linha de comando do MongoDB

Você pode acessar a linha de comando do MongoDB usando o seguinte comando:

```sh
$ docker exec -it mongodb mongo -u root -p rootpassword --authenticationDatabase admin
```

### Criação de um Usuário no Banco de Dados `inventory`

Para criar um usuário com permissões de leitura e escrita no banco de dados `inventory`, utilize o seguinte comando:

```sh
use inventory
db.createUser({
    user: "api_user",
    pwd: "api_password",
    roles: [{ role: "readWrite", db: "inventory" }]
})
```

- NOTA: provavelmente seja necessário reiniciar o contêiner da API (`app` neste caso)

#### Passo 2: Definição de Configuração

A estrutura `MongoDBClientConfig` contém os parâmetros necessários para se conectar a um banco de dados MongoDB. Esses parâmetros incluem o usuário, a senha, o host, a porta e o nome do banco de dados.

- **User**: O nome do usuário que será utilizado para se conectar ao banco de dados.
- **Password**: A senha correspondente ao usuário.
- **Host**: O endereço do host onde está o banco de dados.
- **Port**: A porta em que o banco de dados está ouvindo.
- **Database**: O nome do banco de dados ao qual se deseja conectar.

A função `dns()` gera uma string de conexão (URI) que é utilizada para se conectar ao banco de dados MongoDB. Esta string inclui todos os parâmetros necessários no formato adequado.

```go
package mongodbdriver

import (
    "fmt"
)

// MongoDBClientConfig contém a configuração necessária para se conectar a um banco de dados MongoDB
type MongoDBClientConfig struct {
    User     string // Usuário do banco de dados
    Password string // Senha do usuário
    Host     string // Host onde está o banco de dados
    Port     string // Porta em que o banco de dados está ouvindo
    Database string // Nome do banco de dados
}

// dns gera a URI de conexão ao MongoDB a partir da configuração fornecida
func (config MongoDBClientConfig) dns() string {
    return fmt.Sprintf("mongodb://%s:%s@%s:%s/%s",
        config.User, config.Password, config.Host, config.Port, config.Database)
}
```

#### Passo 3: Configuração do Cliente

O código a seguir define um cliente MongoDB em Go que interage com um banco de dados MongoDB utilizando a configuração fornecida.

```go
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
    db     *mongo.Database     // Conexão ao banco de dados
}

// NewMongoDBClient cria uma nova instância de MongoDBClient e estabelece a conexão ao banco de dados
func NewMongoDBClient(config MongoDBClientConfig) (*MongoDBClient, error) {
    client := &MongoDBClient{config: config}
    err := client.connect()
    if err != nil {
        return nil, fmt.Errorf("failed to initialize MongoDBClient: %v", err)
    }
    return client, nil
}

// connect estabelece a conexão ao banco de dados MongoDB utilizando a configuração fornecida
func (client *MongoDBClient) connect() error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    dns := client.config.dns()
    mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(dns))
    if err != nil {
        return fmt.Errorf("failed to connect to MongoDB: %w", err)
    }
    if err := mongoClient.Ping(ctx, nil); err != nil {
        return fmt.Errorf("failed to ping MongoDB: %w", err)
    }
    client.db = mongoClient.Database(client.config.Database)
    return nil
}

// Close fecha a conexão ao banco de dados MongoDB
func (client *MongoDBClient) Close(ctx context.Context) {
    if client.db != nil {
        client.db.Client().Disconnect(ctx)
    }
}

// DB retorna a conexão ao banco de dados MongoDB
func (client *MongoDBClient) DB() *mongo.Database {
    return client.db
}
```

#### Passo 4: Injeção de Dependências

A função `NewMongoDBSetup` configura e inicializa o cliente MongoDB utilizando os detalhes de conexão definidos em `MongoDBClientConfig`.

```go
package mongodbsetup

import (
    "context"
    "time"
    mongodriver "api/pkg/mongodbdriver"
)

// NewMongoDBSetup configura e retorna um novo cliente MongoDB
func NewMongoDBSetup() (*mongodriver.MongoDBClient, error) {
    config := mongodriver.MongoDBClientConfig{
        User:     "api_user",
        Password: "api_password",
        Host:     "mongodb",
        Port:     "27017",
        Database: "inventory",
    }
    return mongodriver.NewMongoDBClient(config)
}
```

#### Passo 5: Repositório MongoDB

Este código define um repositório em Go que utiliza um banco de dados MongoDB para armazenar e recuperar itens.

```go
package item

import (
    "context"
    "time"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
)

// mongoRepository é uma implementação do repositório de itens utilizando MongoDB
type mongoRepository struct {
    db *mongo.Database // Conexão ao banco de dados MongoDB
}

// NewMongoRepository cria uma nova instância de mongoRepository
func NewMongoRepository(db *mongo.Database) RepositoryPort {
    return &mongoRepository{
        db: db,
    }
}

// SaveItem salva um novo item no banco de dados MongoDB
func (r *mongoRepository) SaveItem(ctx context.Context, it *Item) error {
    if it.CreatedAt.IsZero() {
        it.CreatedAt = time.Now()
    }
    if it.UpdatedAt.IsZero() {
        it.UpdatedAt = time.Now()
    }
    _, err := r.db.Collection("items").InsertOne(ctx, it)
    return err
}

// ListItems lista todos os itens do banco de dados MongoDB
func (r *mongoRepository) ListItems(ctx context.Context) (MapRepo, error) {
    cursor, err := r.db.Collection("items").Find(ctx, bson.D{})
    if err != nil {
        return nil, err
    }
    defer cursor

.Close(ctx)

    items := make(MapRepo)
    for cursor.Next(ctx) {
        var it Item
        if err := cursor.Decode(&it); err != nil {
            return nil, err
        }
        items[it.ID] = it
    }
    if err := cursor.Err(); err != nil {
        return nil, err
    }

    return items, nil
}
```

### Exemplo Completo

O exemplo a seguir combina todos os passos anteriores para configurar uma API REST em Go que utiliza MongoDB como banco de dados.

#### Configuração de MySQL (opcional)

Se também estiver utilizando MySQL em seu projeto, você pode configurar o cliente MySQL da seguinte forma:

```go
package mysqlsetup

import (
    gosqldriver "api/pkg/mysql/go-sql-driver"
)

// NewMySQLSetup configura e retorna um novo cliente MySQL
func NewMySQLSetup() (*gosqldriver.MySQLClient, error) {
    config := gosqldriver.MySQLClientConfig{
        User:     "api_user",
        Password: "api_password",
        Host:     "mysql",
        Port:     "3306",
        Database: "inventory",
    }
    return gosqldriver.NewMySQLClient(config)
}
```

#### Configuração do Servidor

```go
package main

import (
    "log"

    "github.com/gin-gonic/gin"

    handler "api/cmd/rest/handlers"
    core "api/internal/core"
    item "api/internal/core/item"
    mongodbsetup "api/internal/platform/mongodb"
    mysqlsetup "api/internal/platform/mysql"
)

func main() {
    // Configurar MySQL
    mysqlClient, err := mysqlsetup.NewMySQLSetup()
    if err != nil {
        log.Fatalf("não foi possível configurar MySQL: %v", err)
    }
    defer mysqlClient.Close()

    // Configurar MongoDB
    mongoDBClient, err := mongodbsetup.NewMongoDBSetup()
    if err != nil {
        log.Fatalf("não foi possível configurar MongoDB: %v", err)
    }
    defer mongoDBClient.Close()

    // Inicializar repositórios
    mysqlRepo := item.NewMySqlRepository(mysqlClient.DB())
    mongoDBRepo := item.NewMongoRepository(mongoDBClient.DB())

    // Inicializar caso de uso com ambos repositórios
    usecase := core.NewItemUsecase(mysqlRepo, mongoDBRepo)

    // Inicializar handlers
    handler := handler.NewHandler(usecase)

    // Configurar roteador
    router := gin.Default()
    router.POST("/items", handler.SaveItem)
    router.GET("/items", handler.ListItems)

    // Iniciar servidor
    log.Println("Servidor iniciado em http://localhost:8080")
    if err := router.Run(":8080"); err != nil {
        log.Fatal(err)
    }
}
```

### Handlers HTTP

```go
package handler

import (
    "net/http"

    "github.com/gin-gonic/gin"

    "api/internal/core"
    "api/internal/core/item"
    "api/pkg/config"
)

// handler é o manipulador para as solicitações HTTP relacionadas aos itens
type handler struct {
    core core.ItemUsecasePort // Caso de uso dos itens
}

// NewHandler cria uma nova instância de handler
func NewHandler(u core.ItemUsecasePort) *handler {
    return &handler{
        core: u,
    }
}

// SaveItem manipula a solicitação para salvar um novo item
func (h *handler) SaveItem(c *gin.Context) {
    var it item.Item

    err := c.BindJSON(&it)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    ctx := c.Request.Context()
    if err := h.core.SaveItem(ctx, it); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, "item saved successfully")
}

// ListItems manipula a solicitação para listar todos os itens
func (h *handler) ListItems(c *gin.Context) {
    ctx := c.Request.Context()
    its, err := h.core.ListItems(ctx)
    if err != nil {
        if err == config.ErrNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        }
        return
    }

    c.JSON(http.StatusOK, its)
}
```

### Exemplo de Uso com JSON

Para testar o código e verificar se funciona corretamente, você pode utilizar o seguinte JSON para salvar um novo item:

```json
{
  "id": 100,
  "code": "ABC123",
  "title": "Sample Item",
  "description": "This is a sample item.",
  "price": 19.99,
  "stock": 100,
  "status": "Available",
  "created_at": "2024-07-17T10:53:22.123456789Z",
  "updated_at": "2024-07-17T10:53:22.123456789Z"
}
```