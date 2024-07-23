## Setup Cliente MongoDB

**Resumen**

1. **Definición de Configuración**: Crear una estructura `MongoDBClientConfig` para almacenar los detalles de conexión y una función para generar la cadena de conexión (URI).
2. **Configuración del Cliente**: Implementar un cliente MongoDB (`MongoDBClient`) que utiliza la configuración para conectarse a la base de datos.
3. **Inyección de Dependencias**: Configurar e inicializar el cliente MongoDB a través de la función `NewMongoDBSetup`.
4. **Repositorio MongoDB**: Crear un repositorio (`mongoRepository`) que utiliza el cliente MongoDB para realizar operaciones CRUD en la base de datos.

- `func NewMongoDBSetup() (*mongodbdriver.MongoDBClient, error)`:
    - Usa `type MongoDBClientConfig struct` para crear la configuración.
    - Usa `func NewMongoDBClient(config MongoDBClientConfig) (*MongoDBClient, error)` para crear la instancia de MongoDB.
- Esta instancia se inyecta en el repositorio con `func NewMongoRepository(db *mongo.Database) ItemRepositoryPort`.

---

### Configuración MongoDB

Este código define una estructura en Go para configurar un cliente MongoDB y una función asociada para generar una cadena de conexión (URI).

La estructura `MongoDBClientConfig` contiene los parámetros necesarios para conectarse a una base de datos MongoDB. Estos parámetros incluyen el usuario, la contraseña, el host, el puerto y el nombre de la base de datos.

- **User**: El nombre de usuario que se utilizará para conectarse a la base de datos.
- **Password**: La contraseña correspondiente al usuario.
- **Host**: La dirección del host donde se encuentra la base de datos.
- **Port**: El puerto en el que la base de datos está escuchando.
- **Database**: El nombre de la base de datos a la cual se desea conectar.

La función `uri()` genera una cadena de conexión (URI) que se utiliza para conectarse a la base de datos MongoDB. Esta cadena incluye todos los parámetros necesarios en el formato adecuado.

```go
package mongodbdriver

import (
    "fmt"
)

// MongoDBClientConfig contiene la configuración necesaria para conectarse a una base de datos MongoDB
type MongoDBClientConfig struct {
    User     string // Usuario de la base de datos
    Password string // Contraseña del usuario
    Host     string // Host donde se encuentra la base de datos
    Port     string // Puerto en el que escucha la base de datos
    Database string // Nombre de la base de datos
}

// uri genera el URI de conexión a MongoDB a partir de la configuración proporcionada
func (config MongoDBClientConfig) uri() string {
    return fmt.Sprintf("mongodb://%s:%s@%s:%s/%s",
        config.User, config.Password, config.Host, config.Port, config.Database)
}
```

Cuando se crea una instancia de `MongoDBClientConfig` con los detalles de conexión a la base de datos, se puede llamar a la función `uri()` para obtener la cadena de conexión que se utilizará para conectarse a MongoDB.

### Configuración MongoDB

El paquete `mongodbsetup` se utiliza para configurar e inicializar un cliente MongoDB utilizando los detalles de conexión definidos en una estructura de configuración. Este código tiene una relación directa con la estructura y la función `uri()` del código anterior.

```go
package mongodbsetup

import (
    "context"
    "time"
    mongodriver "api/pkg/mongodbdriver"
)

// NewMongoDBSetup configura y devuelve un nuevo cliente MongoDB
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

- **Importación del Paquete**: Se importa el paquete `mongodriver` que contiene la implementación del cliente MongoDB.
- **Función `NewMongoDBSetup`**:
  - Se define una configuración `MongoDBClientConfig` con los detalles de conexión (usuario, contraseña, host, puerto y base de datos).
  - Se llama a `NewMongoDBClient` con la configuración creada, que utiliza la función `uri()` del código anterior para generar la cadena de conexión y establecer la conexión con la base de datos MongoDB.
  - La función devuelve una instancia del cliente MongoDB configurado y listo para ser utilizado en otras partes del código.

### Cliente MongoDB

Este código define un cliente MongoDB en Go que interactúa con una base de datos MongoDB utilizando la configuración proporcionada. A continuación, se explican los componentes y su relación con el código anterior.

```go
package mongodbdriver

import (
    "context"
    "fmt"
    "time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDBClient representa un cliente para interactuar con una base de datos MongoDB
type MongoDBClient struct {
    config MongoDBClientConfig // Configuración del cliente MongoDB
    db     *mongo.Database     // Conexión a la base de datos
}

// NewMongoDBClient crea una nueva instancia de MongoDBClient y establece la conexión a la base de datos
func NewMongoDBClient(config MongoDBClientConfig) (*MongoDBClient, error) {
    client := &MongoDBClient{config: config}
    err := client.connect()
    if err != nil {
        return nil, fmt.Errorf("failed to initialize MongoDBClient: %v", err)
    }
    return client, nil
}

// connect establece la conexión a la base de datos MongoDB utilizando la configuración proporcionada
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

// Close cierra la conexión a la base de datos MongoDB
func (client *MongoDBClient) Close(ctx context.Context) {
    if client.db != nil {
        client.db.Client().Disconnect(ctx)
    }
}

// DB devuelve la conexión a la base de datos MongoDB
func (client *MongoDBClient) DB() *mongo.Database {
    return client.db
}
```

#### Descripción de los Componentes

1. **Importaciones y Paquete**:
   - `context`: Paquete estándar de Go para manejar contextos.
   - `fmt`: Paquete estándar de Go para formatear cadenas.
   - `time`: Paquete estándar de Go para manejar tiempos.
   - `go.mongodb.org/mongo-driver/mongo`: Paquete para interactuar con MongoDB en Go.
   - `go.mongodb.org/mongo-driver/mongo/options`: Paquete para opciones de conexión con MongoDB.

2. **Estructura `MongoDBClient`**:
   - `MongoDBClientConfig config`: Configuración del cliente MongoDB, que fue definida en el código anterior.
   - `*mongo.Database db`: La conexión con la base de datos.

3. **Función `NewMongoDBClient`**:
   - Toma una configuración `MongoDBClientConfig` y crea una nueva instancia de `MongoDBClient`.
   - Llama a `connect()` para establecer la conexión con la base de datos.
   - Si la conexión falla, devuelve un error; si tiene éxito, devuelve la instancia del cliente.

4. **Función `connect`**:
   - Utiliza la función `uri()` definida en `MongoDBClientConfig` (del código anterior) para obtener la cadena de conexión.
   - Abre la conexión con la base de datos con `mongo.Connect`.
   - Verifica la conexión con `mongoClient.Ping`.
   - Si todo es exitoso, asigna la conexión a `client.db`.

5. **Función `Close`**:
   - Cierra la conexión con la base de datos si está abierta.

6. **Función `DB`**:
   - Devuelve la instancia de la conexión con la base de datos.

### Repositorio MongoDB

Este código define un repositorio en Go que utiliza una base de datos MongoDB para almacenar y recuperar elementos. A continuación, se explican los componentes y su relación con el código anterior.

```go
package item

import (
    "context"
    "time"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
)

// mongoRepository es una implementación del repositorio de elementos utilizando MongoDB
type mongoRepository struct {
    db *mongo.Database // Conexión a la base de datos MongoDB
}

// NewMongoRepository crea una nueva instancia de mongoRepository
func NewMongoRepository(db *mongo.Database) ItemRepositoryPort {
    return &mongoRepository{
        db: db,
    }
}

// SaveItem guarda un nuevo elemento en la base de datos MongoDB
func (r *mongoRepository) SaveItem(ctx context.Context, it *Item) error {
    if it.CreatedAt.IsZero() {
        it.CreatedAt = time.Now()
    }
    if it.UpdatedAt.IsZero() {
        it.Updated

At = time.Now()
    }
    _, err := r.db.Collection("items").InsertOne(ctx, it)
    return err
}

// ListItems lista todos los elementos de la base de datos MongoDB
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

#### Descripción de los Componentes

1. **Importaciones y Paquete**:
   - `context`: Paquete estándar de Go para manejar contextos.
   - `time`: Paquete estándar de Go para manejar tiempos.
   - `go.mongodb.org/mongo-driver/bson`: Paquete para manejar BSON (Binary JSON) en MongoDB.
   - `go.mongodb.org/mongo-driver/mongo`: Paquete para interactuar con MongoDB en Go.

2. **Estructura `mongoRepository`**:
   - `*mongo.Database db`: La conexión con la base de datos MongoDB.

3. **Función `NewMongoRepository`**:
   - Crea una nueva instancia de `mongoRepository` con la conexión a la base de datos proporcionada.
   - Devuelve una implementación de `ItemRepositoryPort`.

4. **Función `SaveItem`**:
   - Guarda un nuevo elemento en la base de datos MongoDB.
   - Inicializa `CreatedAt` y `UpdatedAt` con la hora actual si están vacíos.
   - Utiliza `InsertOne` para insertar los datos del elemento en la colección `items`.
   - Devuelve un error si la operación falla.

5. **Función `ListItems`**:
   - Lista todos los elementos de la base de datos MongoDB.
   - Utiliza `Find` para recuperar los datos de la colección `items`.
   - Itera sobre el cursor y decodifica cada documento en un mapa (`MapRepo`).
   - Gestiona el cierre del cursor y devuelve los elementos listados.

### Caso de Uso para Elementos

El siguiente código define un caso de uso para los elementos (`ItemUsecase`) que utiliza dos repositorios (uno basado en MongoDB y otro en un mapa en memoria) para almacenar y recuperar elementos.

```go
package core

import (
    "context"
    "fmt"
    "time"

    "api/internal/core/item"
)

// ItemUsecase representa el caso de uso para los elementos
type ItemUsecase struct {
    mongoRepo item.ItemRepositoryPort // Repositorio de MongoDB
    mapRepo   item.ItemRepositoryPort // Repositorio de Map
}

// NewItemUsecase crea una nueva instancia de ItemUsecase
func NewItemUsecase(mongoRepo, mapRepo item.ItemRepositoryPort) ItemUsecasePort {
    return &ItemUsecase{
        mongoRepo: mongoRepo,
        mapRepo:   mapRepo,
    }
}

// SaveItem guarda un nuevo elemento en ambos repositorios
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

// ListItems lista todos los elementos de ambos repositorios y los combina
func (u *ItemUsecase) ListItems(ctx context.Context) (item.MapRepo, error) {
    mongoItems, err := u.mongoRepo.ListItems(ctx)
    if (err != nil) {
        return nil, fmt.Errorf("error listing items from MongoDB: %w", err)
    }

    mapItems, err := u.mapRepo.ListItems(ctx)
    if err != nil {
        return nil, fmt.Errorf("error listing items from MapRepo: %w", err)
    }

    // Combina los resultados de ambos repositorios
    for k, v := range mapItems {
        mongoItems[k] = v
    }

    return mongoItems, nil
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

    ctx := c.Request.Context()
    if err := h.core.SaveItem(ctx, it); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, "item saved successfully")
}

// ListItems maneja la solicitud para listar todos los elementos
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