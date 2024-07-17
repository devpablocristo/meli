#### Configuración del Cliente MySQL

**Resumen**

1. **Definición de Configuración**: Crear una estructura `MySQLClientConfig` para almacenar los detalles de conexión y una función para generar la cadena DSN.
2. **Configuración del Cliente**: Implementar un cliente MySQL (`MySQLClient`) que utiliza la configuración para conectarse a la base de datos.
3. **Inyección de Dependencias**: Configurar e inicializar el cliente MySQL a través de la función `NewMySQLSetup`.
4. **Repositorio MySQL**: Crear un repositorio (`mysqlRepository`) que utiliza el cliente MySQL para realizar operaciones CRUD en la base de datos.

- `func NewMySQLSetup() (*gosqldriver.MySQLClient, error)`:
    - Usa `type MySQLClientConfig struct` para crear la configuración.
    - Usa `func NewMySQLClient(config MySQLClientConfig) (*MySQLClient, error)` para crear la instancia de MySQL.
- Esta instancia se inyecta en el repositorio con `func NewMySqlRepository(db *sql.DB) ItemRepositoryPort`.

---

### Configuración MySQL

Este código define una estructura en Go para configurar un cliente MySQL y una función asociada para generar una cadena de conexión (DSN - Data Source Name).

La estructura `MySQLClientConfig` contiene los parámetros necesarios para conectarse a una base de datos MySQL. Estos parámetros incluyen el usuario, la contraseña, el host, el puerto y el nombre de la base de datos.

- **User**: El nombre de usuario que se utilizará para conectarse a la base de datos.
- **Password**: La contraseña correspondiente al usuario.
- **Host**: La dirección del host donde se encuentra la base de datos.
- **Port**: El puerto en el que la base de datos está escuchando.
- **Database**: El nombre de la base de datos a la cual se desea conectar.

La función `dsn()` (Data Source Name) genera una cadena de conexión (DSN) que se utiliza para conectarse a la base de datos MySQL. Esta cadena incluye todos los parámetros necesarios en el formato adecuado.

```go
package gosqldriver

import (
    "fmt"
)

// MySQLClientConfig contiene la configuración necesaria para conectarse a una base de datos MySQL
type MySQLClientConfig struct {
    User     string // Usuario de la base de datos
    Password string // Contraseña del usuario
    Host     string // Host donde se encuentra la base de datos
    Port     string // Puerto en el que escucha la base de datos
    Database string // Nombre de la base de datos
}

// dsn genera el Data Source Name (DSN) a partir de la configuración proporcionada
func (config MySQLClientConfig) dsn() string {
    return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        config.User, config.Password, config.Host, config.Port, config.Database)
}
```

Cuando se crea una instancia de `MySQLClientConfig` con los detalles de conexión a la base de datos, se puede llamar a la función `dsn()` para obtener la cadena de conexión que se utilizará para conectarse a MySQL.

### Configuración MySQL

El paquete `mysqlsetup` se utiliza para configurar e inicializar un cliente MySQL utilizando los detalles de conexión definidos en una estructura de configuración. Este código tiene una relación directa con la estructura y la función `dsn()` del código anterior.

```go
package mysqlsetup

import (
    gosqldriver "api/pkg/mysql/go-sql-driver"
)

// NewMySQLSetup configura y devuelve un nuevo cliente MySQL
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

- **Importación del Paquete**: Se importa el paquete `gosqldriver` que contiene la implementación del cliente MySQL.
- **Función `NewMySQLSetup`**:
  - Se define una configuración `MySQLClientConfig` con los detalles de conexión (usuario, contraseña, host, puerto y base de datos).
  - Se llama a `NewMySQLClient` con la configuración creada, que utiliza la función `dsn()` del código anterior para generar la cadena de conexión y establecer la conexión con la base de datos MySQL.
  - La función devuelve una instancia del cliente MySQL configurado y listo para ser utilizado en otras partes del código.

### Cliente MySQL

Este código define un cliente MySQL en Go que interactúa con una base de datos MySQL utilizando la configuración proporcionada. A continuación, se explican los componentes y su relación con el código anterior.

```go
package gosqldriver

import (
    "database/sql"
    "fmt"

    _ "github.com/go-sql-driver/mysql"
)

// MySQLClient representa un cliente para interactuar con una base de datos MySQL
type MySQLClient struct {
    config MySQLClientConfig // Configuración del cliente MySQL
    db     *sql.DB           // Conexión a la base de datos
}

// NewMySQLClient crea una nueva instancia de MySQLClient y establece la conexión a la base de datos
func NewMySQLClient(config MySQLClientConfig) (*MySQLClient, error) {
    client := &MySQLClient{config: config}
    err := client.connect()
    if err != nil {
        return nil, fmt.Errorf("failed to initialize MySQLClient: %v", err)
    }
    return client, nil
}

// connect establece la conexión a la base de datos MySQL utilizando la configuración proporcionada
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

// Close cierra la conexión a la base de datos MySQL
func (client *MySQLClient) Close() {
    if client.db != nil {
        client.db.Close()
    }
}

// DB devuelve la conexión a la base de datos MySQL
func (client *MySQLClient) DB() *sql.DB {
    return client.db
}
```

#### Descripción de los Componentes

1. **Importaciones y Paquete**:
   - `database/sql`: Paquete estándar de Go para interactuar con bases de datos SQL.
   - `fmt`: Paquete estándar de Go para formatear cadenas.
   - `_ "github.com/go-sql-driver/mysql"`: Importa el driver MySQL para `database/sql`, necesario para conectar Go con MySQL.

2. **Estructura `MySQLClient`**:
   - `MySQLClientConfig config`: Configuración del cliente MySQL, que fue definida en el código anterior.
   - `*sql.DB db`: La conexión con la base de datos.

3. **Función `NewMySQLClient`**:
   - Toma una configuración `MySQLClientConfig` y crea una nueva instancia de `MySQLClient`.
   - Llama a `connect()` para establecer la conexión con la base de datos.
   - Si la conexión falla, devuelve un error; si tiene éxito, devuelve la instancia del cliente.

4. **Función `connect`**:
   - Utiliza la función `dsn()` definida en `MySQLClientConfig` (del código anterior) para obtener la cadena de conexión.
   - Abre la conexión con la base de datos con `sql.Open`.
   - Verifica la conexión con `conn.Ping()`.
   - Si todo es exitoso, asigna la conexión a `client.db`.

5. **Función `Close`**:
   - Cierra la conexión con la base de datos si está abierta.

6. **Función `DB`**:
   - Devuelve la instancia de la conexión con la base de datos.

### Repositorio MySQL

Este código define un repositorio en Go que utiliza una base de datos MySQL para almacenar y recuperar elementos. A continuación, se explican los componentes y su relación con el código anterior.

```go
package item

import (
    "database/sql"
    "time"
)

// mysqlRepository es una implementación del repositorio de elementos utilizando MySQL
type mysqlRepository struct {
    db *sql.DB // Conexión a la base de datos MySQL
}

// NewMySqlRepository crea una nueva instancia de mysqlRepository
func NewMySqlRepository(db *sql.DB) ItemRepositoryPort {
    return &mysqlRepository{
        db: db,
    }
}

// SaveItem guarda un nuevo elemento en la base de datos MySQL
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

// ListItems lista todos los elementos de la base de datos MySQL
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

#### Descripción de los Componentes

1. **Importación del paquete `database/sql`**:
   - `database/sql`: Paquete estándar de Go para interactuar con bases de datos SQL.

2. **Estructura `mysqlRepository`**:
   - `*sql.DB db`: La conexión con la base de datos MySQL.

3. **Función `NewMySqlRepository`**:
   - Crea una nueva instancia de `mysqlRepository` con la conexión a la base de datos proporcionada.
   - Devuelve una implementación de `ItemRepositoryPort`.

4. **Función `SaveItem`**:
   - Guarda un nuevo elemento en la base de datos MySQL.
   - Inicializa `CreatedAt` y `UpdatedAt` con la hora actual si están vacíos.
   - Utiliza una consulta SQL `INSERT` para insertar los datos del elemento en la tabla `items`.
   - Devuelve un error si la operación falla.

5. **Función `ListItems`**:
   - Lista todos los elementos de la base de datos MySQL.
   - Utiliza una consulta SQL `SELECT` para recuperar los datos de la tabla `items`.
   - Almacena los resultados en un mapa (`MapRepo`) y los devuelve.
   - Gestiona el cierre de las filas (`rows`) después de iterar sobre ellas.

### Caso de Uso para Elementos

El siguiente código define un caso de uso para los elementos (`ItemUsecase`) que utiliza dos repositorios (uno basado en MySQL y otro en un mapa en memoria) para almacenar y recuperar elementos. 

```go
package core

import (
    "fmt"
    "time"

    "api/internal/core/item"
)

// ItemUsecase representa el caso de uso para los elementos
type ItemUsecase struct {
    mysqlRepo item.ItemRepositoryPort // Repositorio de MySQL
    mapRepo   item.ItemRepositoryPort // Repositorio de Map
}

// NewItemUsecase crea una nueva instancia de ItemUsecase
func NewItemUsecase(mysqlRepo, mapRepo item.ItemRepositoryPort) ItemUsecasePort {
    return &ItemUsecase{
        mysqlRepo: mysqlRepo,
        mapRepo:   mapRepo,
    }
}

// SaveItem guarda un nuevo elemento en ambos repositorios
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

// ListItems lista todos los elementos de ambos repositorios y los combina
func (u *ItemUsecase) ListItems() (item.MapRepo, error) {
    mysqlItems, err := u.mysqlRepo.ListItems()
    if err != nil {
        return nil, fmt.Errorf("error listing items from MySQL: %w", err)
    }

    mapItems, err := u.mapRepo.ListItems()
    if err != nil {
        return nil, fmt.Errorf("error listing items from MapRepo: %w", err)
    }

    // Combina los resultados de ambos repositorios
    for k, v := range mapItems {
        mysqlItems[k] = v
    }

    return mysqlItems, nil
}
```

### Handlers HTTP

El siguiente código define manejadores HTTP (`handler`) que utilizan el caso de uso (`ItemUsecase`) para procesar solicitudes relacionadas con los elementos. Se utilizan para guardar nuevos elementos y listar todos los elementos.

```go
package handler

import (
    "net/http"

    "github.com/gin-gonic/gin"

    "api/internal/core"
    "api/internal/core/item"
    "api/pkg/config"
)

// handler es el manejador para las solicitudes HTTP relacionadas con los elementos
type handler struct {
    core core.ItemUsecasePort // Caso de uso de elementos
}

// NewHandler crea una nueva instancia de handler
func NewHandler(u core.ItemUsecasePort) *handler {
    return &handler{
        core: u,
    }
}

// SaveItem maneja la solicitud para guardar un nuevo elemento
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

// ListItems maneja la solicitud para listar todos los elementos
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

### Ejemplo de Uso con JSON

Para probar el código y verificar que funcione correctamente, puedes utilizar el siguiente JSON para guardar un nuevo elemento:

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