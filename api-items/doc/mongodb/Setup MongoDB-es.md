## Setup Cliente MongoDB

### Resumen

1. **Inicialización de MongoDB con Docker Compose: Crear base de datos, usuario y contraseña.
2. **Definición de Configuración**: Crear una estructura `MongoDBClientConfig` para almacenar los detalles de conexión y una función para generar la cadena de conexión (URI).
3. **Configuración del Cliente**: Implementar un cliente MongoDB (`MongoDBClient`) que utiliza la configuración para conectarse a la base de datos.
4. **Inyección de Dependencias**: Configurar e inicializar el cliente MongoDB a través de la función `NewMongoDBSetup`.
5. **Repositorio MongoDB**: Crear un repositorio (`mongoRepository`) que utiliza el cliente MongoDB para realizar operaciones CRUD en la base de datos.

### Pasos

#### Paso 1: Inicialización de MongoDB con Docker Compose

#### Configuración del Contenedor MongoDB:

- **MONGO_INITDB_ROOT_USERNAME**: Define el nombre de usuario root que se creará al inicializar la base de datos.
- **MONGO_INITDB_ROOT_PASSWORD**: Define la contraseña para el usuario root que se creará al inicializar la base de datos.

#### Configuración del Contenedor Mongo Express:

- **ME_CONFIG_MONGODB_ADMINUSERNAME**: Define el nombre de usuario que Mongo Express usará para conectarse a MongoDB como administrador.
- **ME_CONFIG_MONGODB_ADMINPASSWORD**: Define la contraseña que Mongo Express usará para conectarse a MongoDB como administrador.

**Nota:** Las credenciales deben ser las mismas tanto para MongoDB como para Mongo Express. Esto asegura que Mongo Express pueda conectarse a MongoDB con permisos de administrador.

### Ejemplo de configuración en `docker-compose.yml`:

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

En esta configuración:
- **MONGO_INITDB_ROOT_USERNAME** y **MONGO_INITDB_ROOT_PASSWORD** configuran el usuario root de MongoDB.
- **ME_CONFIG_MONGODB_ADMINUSERNAME** y **ME_CONFIG_MONGODB_ADMINPASSWORD** configuran Mongo Express para usar el usuario root de MongoDB para conectarse.

### Acceso a Mongo Express

Puedes ver los usuarios en Mongo Express accediendo a: [http://localhost:8082/db/admin/system.users](http://localhost:8082/db/admin/system.users)

### Acceso a la línea de comandos de MongoDB

Puedes acceder a la línea de comandos de MongoDB usando el siguiente comando:

```sh
$ docker exec -it mongodb mongo -u root -p rootpassword --authenticationDatabase admin
```

### Creación de un Usuario en la Base de Datos `inventory`

Para crear un usuario con permisos de lectura y escritura en la base de datos `inventory`, utiliza el siguiente comando:

```sh
use inventory
db.createUser({
    user: "api_user",
    pwd: "api_password",
    roles: [{ role: "readWrite", db: "inventory" }]
})
```

- NOTA: probablmente sea necesario volver reiniciar el contenedor de la api (`app` en este caso)

#### Paso 2: Definición de Configuración

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

#### Paso 3: Configuración del Cliente

El siguiente código define un cliente MongoDB en Go que interactúa con una base de datos MongoDB utilizando la configuración proporcionada.

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

#### Paso 4: Inyección de Dependencias

La función `NewMongoDBSetup` configura e inicializa el cliente MongoDB utilizando los detalles de conexión definidos en `MongoDBClientConfig`.

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

#### Paso 5: Repositorio MongoDB

Este código define un repositorio en Go que utiliza una base de datos MongoDB para almacenar y recuperar elementos.

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
        it.UpdatedAt = time.Now()
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

### Ejemplo Completo

El siguiente ejemplo combina todos los pasos anteriores para configurar una API REST en Go que utiliza MongoDB como base de datos.

#### Configuración de MySQL (opcional)

Si también estás utilizando MySQL en tu proyecto, puedes configurar el cliente MySQL de la siguiente manera:

```go
package mysqlsetup

import (
    gosqldriver "api/pkg/mysql/go-sql-driver"
)

// NewMySQLSetup configura y retorna un nuevo cliente MySQL
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

#### Configuración del Servidor

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
        log.Fatalf("no se pudo configurar MySQL: %v", err)
    }
    defer mysqlClient.Close()

    // Configurar MongoDB
    mongoDBClient, err := mongodbsetup.NewMongoDBSetup()
    if err != nil {
        log.Fatalf("no se pudo configurar MongoDB: %v", err)
    }
    defer mongoDBClient.Close()

    // Inicializar repositorios
    mysqlRepo := item.NewMySqlRepository(mysqlClient.DB())
    mongoDBRepo := item.NewMongoRepository(mongoDBClient.DB())

    // Inicializar caso de uso con ambos repositorios
    usecase := core.NewItemUsecase(mysqlRepo, mongoDBRepo)

    // Inicializar handlers
    handler := handler.NewHandler(usecase)

    // Configurar enrutador
    router := gin.Default()
    router.POST("/items", handler.SaveItem)
    router.GET("/items", handler.ListItems)

    // Iniciar servidor
    log.Println("Servidor iniciado en http://localhost:8080")
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