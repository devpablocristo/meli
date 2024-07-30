## Setup Cliente MySQL

**Resumo**

1. **Definição de Configuração**: Criar uma estrutura `MySQLClientConfig` para armazenar os detalhes de conexão e uma função para gerar a string DSN.
2. **Configuração do Cliente**: Implementar um cliente MySQL (`MySQLClient`) que utiliza a configuração para se conectar ao banco de dados.
3. **Injeção de Dependências**: Configurar e inicializar o cliente MySQL através da função `NewMySQLSetup`.
4. **Repositório MySQL**: Criar um repositório (`mysqlRepository`) que utiliza o cliente MySQL para realizar operações CRUD no banco de dados.

- `func NewMySQLSetup() (*gosqldriver.MySQLClient, error)`:
    - Usa `type MySQLClientConfig struct` para criar a configuração.
    - Usa `func NewMySQLClient(config MySQLClientConfig) (*MySQLClient, error)` para criar a instância de MySQL.
- Esta instância é injetada no repositório com `func NewMySqlRepository(db *sql.DB) RepositoryPort`.

---

### Configuração MySQL

Este código define uma estrutura em Go para configurar um cliente MySQL e uma função associada para gerar uma cadeia de conexão (DSN - Data Source Name).

A estrutura `MySQLClientConfig` contém os parâmetros necessários para se conectar a um banco de dados MySQL. Esses parâmetros incluem o usuário, a senha, o host, a porta e o nome do banco de dados.

- **User**: O nome de usuário que será utilizado para se conectar ao banco de dados.
- **Password**: A senha correspondente ao usuário.
- **Host**: O endereço do host onde o banco de dados está localizado.
- **Port**: A porta na qual o banco de dados está ouvindo.
- **Database**: O nome do banco de dados ao qual se deseja conectar.

A função `dsn()` (Data Source Name) gera uma cadeia de conexão (DSN) que é utilizada para se conectar ao banco de dados MySQL. Esta cadeia inclui todos os parâmetros necessários no formato adequado.

```go
package gosqldriver

import (
    "fmt"
)

// MySQLClientConfig contém a configuração necessária para se conectar a um banco de dados MySQL
type MySQLClientConfig struct {
    User     string // Usuário do banco de dados
    Password string // Senha do usuário
    Host     string // Host onde o banco de dados está localizado
    Port     string // Porta na qual o banco de dados está ouvindo
    Database string // Nome do banco de dados
}

// dsn gera o Data Source Name (DSN) a partir da configuração fornecida
func (config MySQLClientConfig) dsn() string {
    return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        config.User, config.Password, config.Host, config.Port, config.Database)
}
```

Quando se cria uma instância de `MySQLClientConfig` com os detalhes da conexão ao banco de dados, pode-se chamar a função `dsn()` para obter a cadeia de conexão que será utilizada para se conectar ao MySQL.

### Configuração MySQL

O pacote `mysqlsetup` é utilizado para configurar e inicializar um cliente MySQL utilizando os detalhes de conexão definidos em uma estrutura de configuração. Este código tem uma relação direta com a estrutura e a função `dsn()` do código anterior.

```go
package mysqlsetup

import (
    gosqldriver "api/pkg/mysql/go-sql-driver"
)

// NewMySQLSetup configura e retorna um novo cliente MySQL
func

 NewMySQLSetup() (*gosqldriver.MySQLClient, error) {
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

- **Importação do Pacote**: Importa-se o pacote `gosqldriver` que contém a implementação do cliente MySQL.
- **Função `NewMySQLSetup`**:
  - Define-se uma configuração `MySQLClientConfig` com os detalhes da conexão (usuário, senha, host, porta e banco de dados).
  - Chama-se `NewMySQLClient` com a configuração criada, que utiliza a função `dsn()` do código anterior para gerar a cadeia de conexão e estabelecer a conexão com o banco de dados MySQL.
  - A função retorna uma instância do cliente MySQL configurado e pronto para ser utilizado em outras partes do código.

### Cliente MySQL

Este código define um cliente MySQL em Go que interage com um banco de dados MySQL utilizando a configuração fornecida. A seguir, explicam-se os componentes e sua relação com o código anterior.

```go
package gosqldriver

import (
    "database/sql"
    "fmt"

    _ "github.com/go-sql-driver/mysql"
)

// MySQLClient representa um cliente para interagir com um banco de dados MySQL
type MySQLClient struct {
    config MySQLClientConfig // Configuração do cliente MySQL
    db     *sql.DB           // Conexão com o banco de dados
}

// NewMySQLClient cria uma nova instância de MySQLClient e estabelece a conexão com o banco de dados
func NewMySQLClient(config MySQLClientConfig) (*MySQLClient, error) {
    client := &MySQLClient{config: config}
    err := client.connect()
    if err != nil {
        return nil, fmt.Errorf("failed to initialize MySQLClient: %v", err)
    }
    return client, nil
}

// connect estabelece a conexão com o banco de dados MySQL utilizando a configuração fornecida
func (client *MySQLClient) connect() error {
    dsn := client.config.dsn()
    conn, err := sql.Open("mysql", dsn)
    if err != nil {
        return fmt.Errorf("failed to connect to MySQL: %w", err)
    }
    if err := conn.Ping(); err != nil {
        return fmt.Errorf("failed to ping MySQL: %w", err)
    }
    client.db = conn
    return nil
}

// Close fecha a conexão com o banco de dados MySQL
func (client *MySQLClient) Close() {
    if client.db != nil {
        client.db.Close()
    }
}

// DB retorna a conexão com o banco de dados MySQL
func (client *MySQLClient) DB() *sql.DB {
    return client.db
}
```

#### Descrição dos Componentes

1. **Importações e Pacote**:
   - `database/sql`: Pacote padrão de Go para interagir com bancos de dados SQL.
   - `fmt`: Pacote padrão de Go para formatar strings.
   - `_ "github.com/go-sql-driver/mysql"`: Importa o driver MySQL para `database/sql`, necessário para conectar Go com MySQL.

2. **Estrutura `MySQLClient`**:
   - `MySQLClientConfig config`: Configuração do cliente MySQL, que foi definida no código anterior.
   - `*sql.DB db`: A conexão com o banco de dados.

3. **Função `NewMySQLClient`**:
   - Toma uma configuração `MySQLClientConfig` e cria uma nova instância de `MySQLClient`.
   - Chama `connect()` para estabelecer a conexão com o banco de dados.
   - Se a conexão falhar, retorna um erro; se tiver sucesso, retorna a instância do cliente.

4. **Função `connect`**:
   - Utiliza a função `dsn()` definida em `MySQLClientConfig` (do código anterior) para obter a string de conexão.
   - Abre a conexão com o banco de dados com `sql.Open`.
   - Verifica a conexão com `conn.Ping()`.
   - Se tudo for bem-sucedido, atribui a conexão a `client.db`.

5. **Função `Close`**:
   - Fecha a conexão com o banco de dados se estiver aberta.

6. **Função `DB`**:
   - Retorna a instância da conexão com o banco de dados.

### Repositório MySQL

Este código define um repositório em Go que utiliza um banco de dados MySQL para armazenar e recuperar itens. A seguir, explicam-se os componentes e sua relação com o código anterior.

```go
package item

import (
    "database/sql"
    "time"
)

// mysqlRepository é uma implementação do repositório de itens utilizando MySQL
type mysqlRepository struct {
    db *sql.DB // Conexão com o banco de dados MySQL
}

// NewMySqlRepository cria uma nova instância de mysqlRepository
func NewMySqlRepository(db *sql.DB) RepositoryPort {
    return &mysqlRepository{
        db: db,
    }
}

// SaveItem salva um novo item no banco de dados MySQL
func (r *mysqlRepository) SaveItem(it *Item) error {
    if it.CreatedAt.IsZero() {
        it.CreatedAt = time.Now()
    }
    if it.UpdatedAt.IsZero() {
        it.UpdatedAt = time.Now()
    }
    query := `INSERT INTO items (code, title, description, price, stock, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
    _, err := r.db.Exec(query, it.Code, it.Title, it.Description, it.Price, it.Stock, it.Status, it.CreatedAt, it.UpdatedAt)
    return err
}

// ListItems lista todos os itens do banco de dados MySQL
func (r *mysqlRepository) ListItems() (MapRepo, error) {
    query := `SELECT id, code, title, description, price, stock, status, created_at, updated_at FROM items`
    rows, err := r.db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    items := make(MapRepo)
    for rows.Next() {
        var it Item
        if err := rows.Scan(&it.ID, &it.Code, &it.Title, &it.Description, &it.Price, &it.Stock, &it.Status, &it.CreatedAt, &it.UpdatedAt); err != nil {
            return nil, err
        }
        items[it.ID] = it
    }

    return items, nil
}
```

#### Descrição dos Componentes

1. **Importação do pacote `database/sql`**:
   - `database/sql`: Pacote padrão de Go para interagir com bancos de dados SQL.

2. **Estrutura `mysqlRepository`**:
   - `*sql.DB db`: A conexão com o banco de dados MySQL.

3. **Função `NewMySqlRepository`**:
   - Cria uma nova instância de `mysqlRepository` com a conexão ao banco de dados fornecida.
   - Retorna uma implementação de `RepositoryPort`.

4. **Função `SaveItem`**:
   - Salva um novo item no banco de dados MySQL.
   - Inicializa `CreatedAt` e `UpdatedAt` com o horário atual se estiverem zerados.
   - Utiliza uma consulta SQL `INSERT` para inserir os dados do item na tabela `items`.
   - Retorna um erro se a operação falhar.

5. **Função `ListItems`**:
   - Lista todos os itens do banco de dados MySQL.
   - Utiliza uma consulta SQL `SELECT` para recuperar os dados da tabela `items`.
   - Armazena os resultados em um mapa (`MapRepo`) e os retorna.
   - Gerencia o fechamento das linhas (`rows`) após iterar sobre elas.

### Caso de Uso para Itens

O seguinte código define um caso de uso para os itens (`ItemUsecase`) que utiliza dois repositórios (um baseado em MySQL e outro em um mapa em memória) para armazenar e recuperar itens.

```go
package core

import (
    "fmt"
    "time"

    "api/internal/core/item"
)

// ItemUsecase representa o caso de uso para os itens
type ItemUsecase struct {
    mysqlRepo item.RepositoryPort // Repositório de MySQL
    mapRepo   item.RepositoryPort // Repositório de Map
}

// NewItemUsecase cria uma nova instância de ItemUsecase
func NewItemUsecase(mysqlRepo, mapRepo item.RepositoryPort) ItemUsecasePort {
    return &ItemUsecase{
        mysqlRepo: mysqlRepo,
        mapRepo:   mapRepo,
    }
}

// SaveItem salva um novo item em ambos os repositórios
func (u *ItemUsecase) SaveItem(it item.Item) error {
    now := time.Now()
    it.CreatedAt = now
    it.UpdatedAt = now

    if err := u.mysqlRepo.SaveItem(&it); err != nil {
        return fmt.Errorf("error saving item in MySQL: %w", err)
    }
    if err := u.mapRepo.SaveItem(&it); err != nil {
        return fmt.Errorf("error saving item in MapRepo: %w", err)
    }
    return nil
}

// ListItems lista todos os itens de ambos

 os repositórios e os combina
func (u *ItemUsecase) ListItems() (item.MapRepo, error) {
    mysqlItems, err := u.mysqlRepo.ListItems()
    if err != nil {
        return nil, fmt.Errorf("error listing items from MySQL: %w", err)
    }

    mapItems, err := u.mapRepo.ListItems()
    if err != nil {
        return nil, fmt.Errorf("error listing items from MapRepo: %w", err)
    }

    // Combina os resultados de ambos os repositórios
    for k, v := range mapItems {
        mysqlItems[k] = v
    }

    return mysqlItems, nil
}
```

### Handlers HTTP

O seguinte código define manipuladores HTTP (`handler`) que utilizam o caso de uso (`ItemUsecase`) para processar solicitações relacionadas aos itens. Eles são utilizados para salvar novos itens e listar todos os itens.

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

    if err := h.core.SaveItem(it); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, "item saved successfully")
}

// ListItems manipula a solicitação para listar todos os itens
func (h *handler) ListItems(c *gin.Context) {
    its, err := h.core.ListItems()
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

Para testar o código e verificar se ele funciona corretamente, você pode utilizar o seguinte JSON para salvar um novo item:

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