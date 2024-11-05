package usershandlers

import (
	"fmt"

	"github.com/Aritiaya50217/E-CommerceRESTAPIs/config"
	"github.com/Aritiaya50217/E-CommerceRESTAPIs/modules/entities"
	"github.com/Aritiaya50217/E-CommerceRESTAPIs/modules/users"
	usersUsecases "github.com/Aritiaya50217/E-CommerceRESTAPIs/modules/users/usersUsecases"
	"github.com/gofiber/fiber/v2"
)

type userHandlersErrCode string

const (
	signUpCustomerErr     userHandlersErrCode = "users-001"
	signInErr             userHandlersErrCode = "users-002"
	refreshPassportErr    userHandlersErrCode = "users-003"
	signOutErr            userHandlersErrCode = "users-004"
	signUpAdminErr        userHandlersErrCode = "users-005"
	generateAdminTokenErr userHandlersErrCode = "users-006"
	getUserProfileErr     userHandlersErrCode = "users-007"
)

type IUsersHandler interface {
	SignUpCustomer(c *fiber.Ctx) error
}

type userHandler struct {
	cfg          config.IConfig
	usersUsecase usersUsecases.IUsersUsecase
}

func UsersHandler(cfg config.IConfig, usersUsecase usersUsecases.IUsersUsecase) IUsersHandler {
	return &userHandler{
		cfg:          cfg,
		usersUsecase: usersUsecase,
	}
}

func (h *userHandler) SignUpCustomer(c *fiber.Ctx) error {
	// Request body parser
	req := new(users.UserRegisterReq)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signUpCustomerErr),
			err.Error(),
		).Res()
	}

	// Email validation
	if !req.IsEmail() {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signUpCustomerErr),
			"email pattern is invalid",
		).Res()
	}

	// duplicated
	isDuplicate := h.usersUsecase.GetUserByEamil(req.Email)
	var result *users.UserPassport
	var err error
	fmt.Println(isDuplicate)
	if !isDuplicate {
		//  Insert
		result, err = h.usersUsecase.InsertCustomer(req)
		if err != nil {
			switch err.Error() {
			case "username has been used":
				return entities.NewResponse(c).Error(
					fiber.ErrBadRequest.Code,
					string(signUpCustomerErr),
					err.Error(),
				).Res()
			case "email has been used":
				return entities.NewResponse(c).Error(
					fiber.ErrBadRequest.Code,
					string(signUpCustomerErr),
					err.Error(),
				).Res()
			}
		}
	}
	//  else {
	// 	return entities.NewResponse(c).Error(
	// 		fiber.ErrBadRequest.Code,
	// 		string(signUpCustomerErr),
	// 		"email pattern is duplicated",
	// 	).Res()
	// }

	return entities.NewResponse(c).Success(fiber.StatusCreated, result).Res()
}
