package usersusecases

import (
	"github.com/Aritiaya50217/E-CommerceRESTAPIs/config"
	"github.com/Aritiaya50217/E-CommerceRESTAPIs/modules/users"
	usersRepositories "github.com/Aritiaya50217/E-CommerceRESTAPIs/modules/users/usersRepositories"
)

type IUsersUsecase interface {
	InsertCustomer(req *users.UserRegisterReq) (*users.UserPassport, error)
	GetUserByEamil(email string) bool
}

type usersUsecase struct {
	cfg            config.IConfig
	userRepository usersRepositories.IUsersRepository
}

func UsersUsecase(cfg config.IConfig, usersRepository usersRepositories.IUsersRepository) IUsersUsecase {
	return &usersUsecase{
		cfg:            cfg,
		userRepository: usersRepository,
	}
}

func (u *usersUsecase) InsertCustomer(req *users.UserRegisterReq) (*users.UserPassport, error) {
	// hashing a password
	if err := req.BcryptHashing(); err != nil {
		return nil, err
	}
	// insert user
	result, err := u.userRepository.InsertUser(req, false)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u *usersUsecase) GetUserByEamil(email string) bool {
	//  check  duplicate email
	isDuplicate := u.userRepository.GetUserByEamil(email)
	if isDuplicate {
		return true
	}
	return false
}
