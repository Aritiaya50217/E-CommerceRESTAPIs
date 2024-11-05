package usersrepositories

import (
	"github.com/Aritiaya50217/E-CommerceRESTAPIs/modules/users"
	usersPatterns "github.com/Aritiaya50217/E-CommerceRESTAPIs/modules/users/usersPatterns"
	"github.com/jmoiron/sqlx"
)

type IUsersRepository interface {
	InsertUser(req *users.UserRegisterReq, isAdmin bool) (*users.UserPassport, error)
	GetUserByEamil(email string) bool
}

type userRepository struct {
	db *sqlx.DB
}

func UserRepository(db *sqlx.DB) IUsersRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) InsertUser(req *users.UserRegisterReq, isAdmin bool) (*users.UserPassport, error) {
	result := usersPatterns.InsertUser(r.db, req, isAdmin)
	var err error
	if isAdmin {
		result, err = result.Admin()
		if err != nil {
			return nil, err
		}
	} else {
		result, err = result.Customer()
		if err != nil {
			return nil, err
		}
	}

	// Get result from inserting
	user, err := result.Result()
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) GetUserByEamil(email string) bool {
	//  check  duplicate email
	_, isUnique := usersPatterns.GetUserByEamil(r.db, email)
	if isUnique {
		return true
	}
	return false
}
