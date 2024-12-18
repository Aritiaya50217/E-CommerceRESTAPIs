package usersusecases

import (
	"fmt"

	"github.com/Aritiaya50217/E-CommerceRESTAPIs/config"
	"github.com/Aritiaya50217/E-CommerceRESTAPIs/modules/users"
	usersRepositories "github.com/Aritiaya50217/E-CommerceRESTAPIs/modules/users/usersRepositories"
	"github.com/Aritiaya50217/E-CommerceRESTAPIs/pkg/auth"
	"golang.org/x/crypto/bcrypt"
)

type IUsersUsecase interface {
	InsertCustomer(req *users.UserRegisterReq) (*users.UserPassport, error)
	GetUserByEamil(email string) bool
	GetPassport(req *users.UserCredential) (*users.UserPassport, error)
	RefreshPassport(req *users.UserRefreshCredential) (*users.UserPassport, error)
	DeleteOauth(oauthId string) error
	InsertAdmin(req *users.UserRegisterReq) (*users.UserPassport, error)
	GetUserProfile(userId string) (*users.User, error)
	ChangePassword(req *users.UserRegisterReq) error
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

func (u *usersUsecase) InsertAdmin(req *users.UserRegisterReq) (*users.UserPassport, error) {
	// hashing a password
	if err := req.BcryptHashing(); err != nil {
		return nil, err
	}

	// insert user
	result, err := u.userRepository.InsertUser(req, true)
	fmt.Println(result)
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

func (u *usersUsecase) GetPassport(req *users.UserCredential) (*users.UserPassport, error) {
	// Find user
	user, err := u.userRepository.FindOneUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	// compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("password is invalid")
	}

	// sign token
	accessToken, _ := auth.NewAuth(auth.Access, u.cfg.Jwt(), &users.UserClaims{
		Id:     user.Id,
		RoleId: user.RoleId,
	})

	refreshToken, _ := auth.NewAuth(auth.Refresh, u.cfg.Jwt(), &users.UserClaims{
		Id:     user.Id,
		RoleId: user.RoleId,
	})

	// set passport
	passport := &users.UserPassport{
		User: &users.User{
			Id:       user.Id,
			Email:    user.Email,
			Username: user.Username,
			RoleId:   user.RoleId,
		},
		Token: &users.UserToken{
			AccessToken:  accessToken.SignToken(),
			RefreshToken: refreshToken.SignToken(),
		},
	}

	if err := u.userRepository.InsertOauth(passport); err != nil {
		return nil, err
	}
	return passport, nil
}

func (u *usersUsecase) RefreshPassport(req *users.UserRefreshCredential) (*users.UserPassport, error) {
	// Parse token
	claims, err := auth.ParseToken(u.cfg.Jwt(), req.RefreshToken)
	if err != nil {
		return nil, err
	}

	// check oauth
	oauth, err := u.userRepository.FindOneOauth(req.RefreshToken)
	if err != nil {
		return nil, err
	}

	// find profile
	profile, err := u.userRepository.GetProfile(oauth.UserId)
	if err != nil {
		return nil, err
	}

	newClaims := &users.UserClaims{
		Id:     profile.Id,
		RoleId: profile.RoleId,
	}
	fmt.Println("newClaims : ", newClaims)

	accessToken, err := auth.NewAuth(
		auth.Access,
		u.cfg.Jwt(),
		newClaims,
	)

	if err != nil {
		return nil, err
	}

	refreshToken := auth.RepeatToken(
		u.cfg.Jwt(),
		newClaims,
		claims.ExpiresAt.Unix(),
	)

	passport := &users.UserPassport{
		User: profile,
		Token: &users.UserToken{
			Id:           oauth.Id,
			AccessToken:  accessToken.SignToken(),
			RefreshToken: refreshToken,
		},
	}

	if err := u.userRepository.UpdateOauth(passport.Token); err != nil {
		return nil, err
	}
	return passport, nil
}

func (u *usersUsecase) DeleteOauth(oauthId string) error {
	return u.userRepository.DeleteOauth(oauthId)
}

func (u *usersUsecase) GetUserProfile(userId string) (*users.User, error) {
	return u.userRepository.GetProfile(userId)
}

func (u *usersUsecase) ChangePassword(req *users.UserRegisterReq) error {
	return u.userRepository.ChangePassword(req)
}
