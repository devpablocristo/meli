package application

type UserService {
	repo <---- difenir que es ese repo????
}

func NewUserService(repo /*type?*/) *UserService {
	return &UserService{
		repo: repo,
	}
}

// logica primeira letra mayuscula SEMPRE
func (u *UserService) CreateUser() {
// logica	
}

func (u *UserService) DeleteUser() { 	
// logica	
}




