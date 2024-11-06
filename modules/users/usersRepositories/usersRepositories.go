package usersrepositories

import (
	"context"
	"fmt"
	"time"

	"github.com/Aritiaya50217/E-CommerceRESTAPIs/modules/users"
	usersPatterns "github.com/Aritiaya50217/E-CommerceRESTAPIs/modules/users/usersPatterns"
	"github.com/jmoiron/sqlx"
)

type IUsersRepository interface {
	InsertUser(req *users.UserRegisterReq, isAdmin bool) (*users.UserPassport, error)
	GetUserByEamil(email string) bool
	FindOneUserByEmail(email string) (*users.UserCredentialCheck, error)
	InsertOauth(req *users.UserPassport) error
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

func (r *userRepository) FindOneUserByEmail(email string) (*users.UserCredentialCheck, error) {
	word := "'%" + email + "%'"
	query := "select * from users where email like " + word
	user := new(users.UserCredentialCheck)
	if err := r.db.Get(user, query); err != nil {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func (r *userRepository) InsertOauth(req *users.UserPassport) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	id := "'" + req.User.Id + "'"
	access := "'" + req.Token.AccessToken + "'"
	refresh := "'" + req.Token.RefreshToken + "'"
	create := "'" + time.Now().Format("2006-01-02 15:04:05.000 ") + "'"
	update := "'" + time.Now().Format("2006-01-02 15:04:05.000 ") + "'"

	values := "values (" + id + "," + access + "," + refresh + "," + create + "," + update + ")"

	query := "insert into oauth ( " +
		"user_id,access_token,refresh_token,created_at,updated_at) " +
		values

	if err := r.db.QueryRowContext(ctx,
		query,
	); err != nil {
		return err.Err()
	}
	return nil
}
