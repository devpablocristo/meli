## Setup Cliente MongoDB

**Resumo**

1. **Definição de Configuração**: Criar uma estrutura `MongoDBClientConfig` para armazenar os detalhes de conexão e uma função para gerar a cadeia de conexão (URI).
2. **Configuração do Cliente**: Implementar um cliente MongoDB (`MongoDBClient`) que utiliza a configuração para conectar-se ao banco de dados.
3. **Injeção de Dependências**: Configurar e inicializar o cliente MongoDB através da função `NewMongoDBSetup`.
4. **Repositório MongoDB**: Criar um repositório (`mongoRepository`) que utiliza o cliente MongoDB para realizar operações CRUD no banco de dados.

- `func NewMongoDBSetup() (*mongodbdriver.MongoDBClient, error)`:
    - Usa `type MongoDBClientConfig struct` para criar a configuração.
    - Usa `func NewMongoDBClient(config MongoDBClientConfig) (*MongoDBClient, error)` para criar a instância do MongoDB.
- Esta instância é injetada no repositório com `func NewMongoRepository(db *mongo.Database) ItemRepositoryPort`.

---

### Configuração MongoDB

Este código define uma estrutura em Go para configurar um cliente MongoDB e uma função associada para gerar uma cadeia de conexão (URI).

A estrutura `MongoDBClientConfig` contém os parâmetros necessários para conectar-se a um banco de dados MongoDB. Esses parâmetros incluem o usuário, a senha, o host, a porta e o nome do banco de dados.

- **User**: O nome de usuário que será utilizado para conectar-se ao banco de dados.
- **Password**: A senha correspondente ao usuário.
- **Host**: O endereço do host onde o banco de dados está localizado.
- **Port**: A porta na qual o banco de dados está escutando.
- **Database**: O nome do banco de dados ao qual se deseja conectar.

A função `uri()` gera uma cadeia de conexão (URI) que é utilizada para conectar-se ao banco de dados MongoDB. Esta cadeia inclui todos os parâmetros necessários no formato adequado.

```go
package mongodbdriver

import (
    "fmt"
)

// MongoDBClientConfig contém a configuração necessária para conectar-se a um banco de dados MongoDB
type MongoDBClientConfig struct {
    User     string // Usuário do banco de dados
    Password string // Senha do usuário
    Host     string // Host onde o banco de dados está localizado
    Port     string // Porta na qual o banco de dados está escutando
    Database string // Nome do banco de dados
}

// uri gera o URI de conexão ao MongoDB a partir da configuração fornecida
func (config MongoDBClientConfig) uri() string {
    return fmt.Sprintf("mongodb://%s:%s@%s:%s/%s",
        config.User, config.Password, config.Host, config.Port, config.Database)
}
```

Quando uma instância de `MongoDBClientConfig` é criada com os detalhes de conexão ao banco de dados, pode-se chamar a função `uri()` para obter a cadeia de conexão que será utilizada para conectar-se ao MongoDB.

### Configuração MongoDB

O pacote `mongodbsetup` é utilizado para configurar e inicializar um cliente MongoDB utilizando os detalhes de conexão definidos em uma estrutura de configuração. Este código tem uma relação direta com a estrutura e a função `uri()` do código anterior.

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
        User:     "root",
        Password: "rootpassword",
        Host:     "mongodb",
        Port:     "27017",
        Database: "inventory",
    }
    return mongodriver.NewMongoDBClient(config)
}
```

- **Importação do Pacote**: Importa-se o pacote `mongodriver` que contém a implementação do cliente MongoDB.
- **Função `NewMongoDBSetup`**:
  - Define uma configuração `MongoDBClientConfig` com os detalhes de conexão (usuário, senha, host, porta e banco de dados).
  - Chama `NewMongoDBClient` com a configuração criada, que utiliza a função `uri()` do código anterior para gerar a cadeia de conexão e estabelecer a conexão com o banco de dados MongoDB.
  - A função retorna uma instância do cliente MongoDB configurado e pronto para ser utilizado em outras partes do código.

### Cliente MongoDB

Este código define um cliente MongoDB em Go que interage com um banco de dados MongoDB utilizando a configuração fornecida. A seguir, explicam-se os componentes e sua relação com o código anterior.

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

    uri := client.config.uri()
    mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
    if err != nil {
        return fmt.Errorf("failed to connect to MongoDB: %w", err)
    }
    if err := mongoClient.Ping(ctx, nil); err != nil {
        return fmt.Errorf("failed to ping MongoDB: %w", err)
    }
    client.db = mongoClient.Database(client.config.Database)
    return nil
}

// Close fecha a conexão com o banco de dados MongoDB
func (client *MongoDBClient) Close(ctx context.Context) {
    if client.db != nil {
        client.db.Client().Disconnect(ctx)
    }
}

// DB retorna a conexão com o banco de dados MongoDB
func (client *MongoDBClient) DB() *mongo.Database {
    return client.db
}
```

#### Descrição dos Componentes

1. **Importações e Pacote**:
   - `context`: Pacote padrão de Go para lidar com contextos.
   - `fmt`: Pacote padrão de Go para formatar cadeias de caracteres.
   - `time`: Pacote padrão de Go para lidar com tempos.
   - `go.mongodb.org/mongo-driver/mongo`: Pacote para interagir com MongoDB em Go.
   - `go.mongodb.org/mongo-driver/mongo/options`: Pacote para opções de conexão com MongoDB.

2. **Estrutura `MongoDBClient`**:
   - `MongoDBClientConfig config`: Configuração do cliente MongoDB, que foi definida no código anterior.
   - `*mongo.Database db`: A conexão com o banco de dados.

3. **Função `NewMongoDBClient`**:
   - Toma uma configuração `MongoDBClientConfig` e cria uma nova instância de `MongoDBClient`.
   - Chama `connect()` para estabelecer a conexão com o banco de dados.
   - Se a conexão falhar, retorna um erro; se tiver sucesso, retorna a instância do cliente.

4. **Função `connect`**:
   - Utiliza a função `uri()` definida em `MongoDBClientConfig` (do código anterior) para obter a cadeia de conexão.
   - Abre a conexão com o banco de dados com `mongo.Connect`.
   - Verifica a conexão com `mongoClient.Ping`.
   - Se tudo for bem-sucedido, atribui a conexão a `client.db`.

5. **Função `Close`**:
   - Fecha a conexão com o banco de dados se estiver aberta.

6. **Função `DB`**:
   - Retorna a instância da conexão com o banco de dados.

### Repositório MongoDB

Este código define um repositório em Go que utiliza um banco de dados MongoDB para armazenar e recuperar itens. A seguir, explicam-se os componentes e sua relação com o código anterior.

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
func NewMongoRepository(db *mongo.Database) ItemRepositoryPort {
    return &mongoRepository{
        db: db,
    }
}

// SaveItem salva um novo item no banco de dados MongoDB
func (r *mongoRepository) SaveItem(ctx context.Context, it *Item) error {
    if it.CreatedAt.IsZero() {
        it.CreatedAt = time.Now()
    }
    if it.UpdatedAt.IsZero()

 {
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
    defer cursor.Close(ctx)

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

#### Descrição dos Componentes

1. **Importações e Pacote**:
   - `context`: Pacote padrão de Go para lidar com contextos.
   - `time`: Pacote padrão de Go para lidar com tempos.
   - `go.mongodb.org/mongo-driver/bson`: Pacote para lidar com BSON (Binary JSON) no MongoDB.
   - `go.mongodb.org/mongo-driver/mongo`: Pacote para interagir com MongoDB em Go.

2. **Estrutura `mongoRepository`**:
   - `*mongo.Database db`: A conexão com o banco de dados MongoDB.

3. **Função `NewMongoRepository`**:
   - Cria uma nova instância de `mongoRepository` com a conexão ao banco de dados fornecida.
   - Retorna uma implementação de `ItemRepositoryPort`.

4. **Função `SaveItem`**:
   - Salva um novo item no banco de dados MongoDB.
   - Inicializa `CreatedAt` e `UpdatedAt` com a hora atual se estiverem vazios.
   - Utiliza `InsertOne` para inserir os dados do item na coleção `items`.
   - Retorna um erro se a operação falhar.

5. **Função `ListItems`**:
   - Lista todos os itens do banco de dados MongoDB.
   - Utiliza `Find` para recuperar os dados da coleção `items`.
   - Itera sobre o cursor e decodifica cada documento em um mapa (`MapRepo`).
   - Gerencia o fechamento do cursor e retorna os itens listados.

### Caso de Uso para Itens

O código a seguir define um caso de uso para os itens (`ItemUsecase`) que utiliza dois repositórios (um baseado em MongoDB e outro em um mapa na memória) para armazenar e recuperar itens.

```go
package core

import (
    "context"
    "fmt"
    "time"

    "api/internal/core/item"
)

// ItemUsecase representa o caso de uso para os itens
type ItemUsecase struct {
    mongoRepo item.ItemRepositoryPort // Repositório de MongoDB
    mapRepo   item.ItemRepositoryPort // Repositório de Map
}

// NewItemUsecase cria uma nova instância de ItemUsecase
func NewItemUsecase(mongoRepo, mapRepo item.ItemRepositoryPort) ItemUsecasePort {
    return &ItemUsecase{
        mongoRepo: mongoRepo,
        mapRepo:   mapRepo,
    }
}

// SaveItem salva um novo item em ambos os repositórios
func (u *ItemUsecase) SaveItem(ctx context.Context, it item.Item) error {
    now := time.Now()
    it.CreatedAt = now
    it.UpdatedAt = now

    if err := u.mongoRepo.SaveItem(ctx, &it); err != nil {
        return fmt.Errorf("error saving item in MongoDB: %w", err)
    }
    if err := u.mapRepo.SaveItem(ctx, &it); err != nil {
        return fmt.Errorf("error saving item in MapRepo: %w", err)
    }
    return nil
}

// ListItems lista todos os itens de ambos os repositórios e os combina
func (u *ItemUsecase) ListItems(ctx context.Context) (item.MapRepo, error) {
    mongoItems, err := u.mongoRepo.ListItems(ctx)
    if err != nil {
        return nil, fmt.Errorf("error listing items from MongoDB: %w", err)
    }

    mapItems, err := u.mapRepo.ListItems(ctx)
    if err != nil {
        return nil, fmt.Errorf("error listing items from MapRepo: %w", err)
    }

    // Combina os resultados de ambos os repositórios
    for k, v := range mapItems {
        mongoItems[k] = v
    }

    return mongoItems, nil
}
```

### Handlers HTTP

O código a seguir define manipuladores HTTP (`handler`) que utilizam o caso de uso (`ItemUsecase`) para processar solicitações relacionadas aos itens. Eles são utilizados para salvar novos itens e listar todos os itens.

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
    core core.ItemUsecasePort // Caso de uso de itens
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

---

Essa documentação fornece uma visão completa de como configurar e utilizar um cliente MongoDB em Go, incluindo a configuração, criação do cliente, injeção de dependências, repositórios, casos de uso, manipuladores HTTP e exemplos de uso com JSON.