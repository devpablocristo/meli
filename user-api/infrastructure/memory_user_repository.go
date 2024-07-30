package infrastructure

import "user-api/domain"

type MemoryUserRepository struct {
	users map[string]domain.User
}

func NewMemoryUserRepository() {

}
