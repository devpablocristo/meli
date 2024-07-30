# Go Software Design and Architecture Guide: Design Patterns and Hexagonal Architecture

## Introduction
In this guide, we'll discuss design patterns and architecture in software development, focusing on implementation in Go using the ORM GORM. We'll explore key patterns like Data Transfer Object (DTO), Data Access Object (DAO), and Entity, as well as the Repository pattern and how these integrate into Hexagonal Architecture.

## GORM and gorm.Model
GORM is a popular ORM (Object-Relational Mapper) in Go that provides a high-level API for database operations. In GORM, `gorm.Model` is a basic structure that includes common fields: ID, CreatedAt, UpdatedAt, DeletedAt. It can be embedded in your structures to automatically add these fields.

```go
type User struct {
    gorm.Model
    Name  string
    Email string
}
```

## DTO, Entity, DAO
These three patterns are commonly used in software design.

- **DTO (Data Transfer Object)**: Used to transfer data between processes or application components. Usually a simple structure that only groups data.

```go
type UserDTO struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}
```

- **Entity**: A representation of a domain object in the application. Contains both data and behavior.

```go
type User struct {
    ID    uint
    Name  string
    Email string
}
```

- **DAO (Data Access Object)**: Provides an interface for data access at the level of a specific entity. Mainly concerned with CRUD (Create, Read, Update, Delete) operations.

```go
type UserDAO struct {
    db *gorm.DB
}

func (dao *UserDAO) FindByID(ID uint) (*User, error) { /*...*/ }
func (dao *UserDAO) Save(user *User) error { /*...*/ }
```

In some cases, the structures of the domain entity and persistence entity may vary based on application needs, but in other cases, they might be the same structure.

## Using Custom Entities in Presenters and Persistence

In certain cases, it's beneficial to define custom entities for presenters and persistence.

**Presenter with Custom Entities**: A Presenter can have a custom entity. It's part of the Model-View-Presenter (MVP) design pattern, where data from the model is prepared for the view.

```go
type User struct {
    ID    uint
    Name  string
    Email string
    Password string
    Orders []Order
}

type UserPresenter struct {}

type UserDTO struct {
    ID    string `json:"ID"`
    Name  string `json:"name"`
    Email string `json:"email"`
    OrderCount int `json:"order_count"`
}

func (p *UserPresenter) Present(user *User) *UserDTO {
    return &UserDTO{
        ID:    strconv.FormatUint(uint64(user.ID), 10),
        Name:  user.Name,
        Email: user.Email,
        OrderCount: len(user.Orders),
    }
}
```

**Persistence with Customized Entities**: For handling data persistence, it's common to use a custom entity representing the table structure in the database.

```go
type UserModel struct {
    gorm.Model
    Name  string
    Email string
}
```

## Repository Pattern vs DAO Pattern
Both patterns provide abstraction over data access operations, but they differ in their approaches and typical uses.

**DAO (Data Access Object)**: Provides an abstraction of any type of persistence operation and is associated with specific table-level operations in an SQL database.

```go
type UserDAO struct {
    db *gorm.DB
}

func (dao *UserDAO) Insert(user *User) error { /*...*/ }
func (dao *UserDAO) Update(user *User) error { /*...*/ }
```

**Repository**: Adds a layer of abstraction over the storage and retrieval operations of domain objects and adheres more to an object-oriented style of entity manipulation.

```go
type UserRepository interface {
    Save(user *User) error
}

type GormUserRepository struct {
    db *gorm.DB
}

func (repo *GormUserRepository) Save(user *User) error { /*...*/ }
```

In the case of hexagonal architecture, the Repository pattern is often preferred because it offers greater decoupling between domain logic and persistence infrastructure.

## Conclusion
In this guide, we covered several common software design and architecture patterns and how they can be implemented in Go using GORM. Design patterns like DTO, DAO, and Repository, as well as the Presenter pattern, are useful tools for keeping your code clean, organized, and easy to understand. Finally, we discussed how these patterns integrate into hexagonal architecture, providing a robust application architecture that's easy to maintain and evolve.