### Entity

An entity encapsulates essential data related to an abstraction, such as a book in this example, and can be utilized across various parts of the system to represent and manage information about books.

```go
package entity

// Book represents the entity of a book in the system.
type Book struct {
    ID      int64
    Title   string
    Author  string
    Summary string
    Year    int
}
```

In this definition, the `Book` entity represents a book in the system and includes the following fields:

- `ID`: The unique identifier of the book.
- `Title`: The title of the book.
- `Author`: The author of the book.
- `Summary`: The summary of the book.
- `Year`: The publication year of the book.

### Data Transfer Object (DTO)

A Data Transfer Object (DTO) is a design pattern used to transfer data between systems or layers of an application in a simplified manner. DTOs encapsulate data and are commonly used to transfer relevant information between client and server, or between different parts of a system, without exposing internal details of how those data are stored or managed. This promotes greater flexibility and decoupling between the layers of an application, being especially useful in the development of APIs to model both incoming requests and outgoing responses, allowing to control exactly what data is exposed through the API.

#### Example of DTO in Go for a book management application:

```go
package dto

// BookCreateDTO is used for creating books, defining what is necessary to add a new book.
type BookCreateDTO struct {
    Title   string `json:"title"`
    Author  string `json:"author"`
    Summary string `json:"summary,omitempty"` // Optional
    Year    int    `json:"year"`
}

// BookResponseDTO is used to send book data to the client, defining how a book is presented in responses.
type BookResponseDTO struct {
    ID      int64  `json:"ID"`
    Title   string `json:"title"`
    Author  string `json:"author"`
    Summary string `json:"summary,omitempty"` // Optional
    Year    int    `json:"year"`
}
```

### Data Access Object (DAO)

The DAO (Data Access Object) is a design pattern that provides an abstract interface to access data stored in a database, file, or any other persistence medium. Its purpose is to separate the data access logic from the business logic of the application, allowing the latter to be independent of the underlying storage mechanism.

#### Example of DAO in Go for the Book entity:

```go
package dao

import (
    "context"
    "errors"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
)

// MongoDbDaoPort define las operaciones de acceso a datos para cualquier tipo de datos.
type MongoDbDaoPort interface {
    Create(any) error
    GetById(string) (any, error)
    Update(any) error
    Delete(string) error
}

// MongoDbDao implementa MongoDbDaoPort utilizando MongoDB.
type MongoDbDao struct {
    collection *mongo.Collection
}

// NewMongoDbDao crea una nueva instancia de MongoDbDao.
func NewMongoDbDao(collection *mongo.Collection) MongoDbDaoPort {
    return &MongoDbDao{
        collection: collection,
    }
}

// Create inserta un nuevo elemento en la colección MongoDB.
func (dao *MongoDbDao) Create(data any) error {
    _, err := dao.collection.InsertOne(context.Background(), data)
    if err != nil {
        return err
    }
    return nil
}

// GetById recupera un elemento por su ID de la colección MongoDB.
func (dao *MongoDbDao) GetById(ID string) (any, error) {
    var result any
    filter := bson.M{"_id": ID}
    err := dao.collection.FindOne(context.Background(), filter).Decode(&result)
    if err != nil {
        return nil, err
    }
    return result, nil
}

// Update actualiza un elemento existente en la colección MongoDB.
func (dao *MongoDbDao) Update(data any) error {
    filter := bson.M{"_id": data.ID} // Se asume que todos los tipos de datos tienen un campo ID
    update := bson.M{"$set": data}
    _, err := dao.collection.UpdateOne(context.Background(), filter, update)
    if err != nil {
        return err
    }
    return nil
}

// Delete elimina un elemento de la colección MongoDB.
func (dao *MongoDbDao) Delete(ID string) error {
    filter := bson.M{"_id": ID}
    _, err := dao.collection.DeleteOne(context.Background(), filter)
    if err != nil {
        return err
    }
    return nil
}
```

### The Repository Pattern

The Repository Pattern is a software design pattern used to separate business logic from data access logic. Its goal is to provide an abstraction layer between the application logic and the details of data persistence implementation.

Key features and benefits of the Repository pattern:

1. **Data Persistence Abstraction**: It provides a uniform interface for accessing data, regardless of how it is stored or retrieved.

2. **Separation of Concerns**: It allows separating business logic from data access logic, facilitating maintenance and independent evolution of both parts.

3. **Code Reusability**: It facilitates reusing data access logic across different parts of the application, promoting cohesion and code reuse.

4. **Facilitates Unit Testing**: By separating data access logic, it makes it easier to write unit tests for the application logic.

### Persistence Model

The Persistence Model represents the data structure used to store and manipulate information persistently in a database or other storage medium. This model defines how the information is organized and stored, as well as the relationships between different types of data.

#### Example of Persistence Model in Go for the Book entity:

```go
package repository

import "time"

// BookModel represents the data structure for persistent storage of books.
type BookDao struct {
    ID        int64     `db:"ID"`
    Title     string    `db:"title"`
    Author    string    `db:"author"`
    Summary   string    `db:"summary,omitempty"` // Optional
    Year      int       `db:"year"`
    CreatedAt time.Time `db:"created_at"`
    UpdatedAt time.Time `db:"updated_at"`
}
```
#### Example of the Repository Pattern with BookDao as persistence model:

```go
package repository

import (
    "context"
    "errors"
    "yourProject/dao"
    "yourProject/models"
)

type BookRepositoryPort interface {
    CreateBook(*BookDao) error 
    GetBookByID(string) (*BookDao, error) 
    UpdateBook(*BookDao) error 
    DeleteBookByID(string) error 
}

// BookRepository represents a repository for working with books.
type BookRepository struct {
    dao dao.MongoDbDaoPort
}

// NewBookRepository creates a new instance of BookRepository.
func NewBookRepository(dao dao.MongoDbDaoPort) BookRepositoryPort {
    return &BookRepository{
        dao: dao,
    }
}

// CreateBook creates a new book in the database.
func (r *BookRepository) CreateBook(book *BookDao) error {
    return r.dao.Create(book)
}

// GetBookByID retrieves a book by its ID.
func (r *BookRepository) GetBookByID(ID string) (*BookDao, error) {
    result, err := r.dao.GetById(ID)
    if err != nil {
        return nil, err
    }
   
    book, ok := result.(*BookDao)
    if !ok {
        return nil, errors.New("failed to convert the result to BookDao")
    }

    return book, nil
}

// UpdateBook updates a book in the database.
func (r *BookRepository) UpdateBook(book *BookDao) error {
    return r.dao.Update(book)
}

// DeleteBookByID deletes a book from the database by its ID.
func (r *BookRepository) DeleteBookByID(ID string) error {
    return r.dao.Delete(ID)
}
```
### Presenter

The Presenter is part of the Model-View-Presenter (MVP) pattern, acting as an intermediary between the view (UI) and the model (business data). Its main goal is to prepare the model data for presentation in the view, containing specific presentation logic that decides how the data should be shown to the user.

#### Example of Presenter in Go:

```go
package presenter

import (
    "fmt"
    "strings"
    "yourProject/dto"
)

// BookPresenter is responsible for preparing book data for the view.
type BookPresenter struct {
    Book dto.BookResponseDTO
}

// FormatTitle converts the book title to uppercase.
func (bp *BookPresenter) FormatTitle() string {
    return strings.ToUpper(bp.Book.Title)
}

// BookDetails composes a string with the formatted details of the book for presentation.
func (bp *BookPresenter) BookDetails() string {
    return fmt.Sprintf("Title: %s, Author: %s, Publication Year: %d",
        bp.FormatTitle(), bp.Book.Author, bp.Book.Year)
}
```