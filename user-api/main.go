package main

import (
	"fmt"

	app "user-api/application"
	infra "user-api/infrastructure"
)

func main() {
	// Configurar o repositório e o serviço
	repo := infra.NewMemoryUserRepository() // nao existe ainda

	userService := app.NewUserService(repo) ///<----- usecases o serviço

	// Criar um novo usuário
	err := userService.CreateUser("1", "John Doe")
	if err != nil {
		fmt.Println("Erro ao criar usuário:", err)
		return
	}

	// Obter o usuário criado
	user, err := userService.GetUser("1")
	if err != nil {
		fmt.Println("Erro ao obter usuário:", err)
		return
	}
	fmt.Println("Usuário obtido:", user.Name)
}
