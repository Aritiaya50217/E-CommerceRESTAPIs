package usershandlers

import (
	"strings"

	"github.com/Aritiaya50217/E-CommerceRESTAPIs/config"
	"github.com/Aritiaya50217/E-CommerceRESTAPIs/modules/entities"
	"github.com/Aritiaya50217/E-CommerceRESTAPIs/modules/users"
	usersUsecases "github.com/Aritiaya50217/E-CommerceRESTAPIs/modules/users/usersUsecases"
	"github.com/Aritiaya50217/E-CommerceRESTAPIs/pkg/auth"
	"github.com/gofiber/fiber/v2"
)

const (
	signUpCustomerErr     string = "users-001"
	signInErr             string = "users-002"
	refreshPassportErr    string = "users-003"
	signOutErr            string = "users-004"
	signUpAdminErr        string = "users-005"
	generateAdminTokenErr string = "users-006"
	getUserProfileErr     string = "users-007"
	changePasswordErr     string = "users-008"
)

type IUsersHandler interface {
	SignUpCustomer(c *fiber.Ctx) error
	SignIn(c *fiber.Ctx) error
	RefreshPassport(c *fiber.Ctx) error
	SignOut(c *fiber.Ctx) error
	SignUpAdmin(c *fiber.Ctx) error
	GenerateAdminToken(c *fiber.Ctx) error
	GetUserProfile(c *fiber.Ctx) error
	ChangePassword(c *fiber.Ctx) error
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

	return entities.NewResponse(c).Success(fiber.StatusCreated, result).Res()
}

func (h *userHandler) SignIn(c *fiber.Ctx) error {
	req := new(users.UserCredential)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signInErr),
			err.Error(),
		).Res()
	}

	passport, err := h.usersUsecase.GetPassport(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signInErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, passport).Res()
}

func (h *userHandler) RefreshPassport(c *fiber.Ctx) error {
	req := new(users.UserRefreshCredential)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(refreshPassportErr),
			err.Error(),
		).Res()
	}
	passport, err := h.usersUsecase.RefreshPassport(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(refreshPassportErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, passport).Res()
}

func (h *userHandler) SignOut(c *fiber.Ctx) error {
	req := new(users.UserRemoveCredential)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signOutErr),
			err.Error(),
		).Res()
	}

	if err := h.usersUsecase.DeleteOauth(req.OauthId); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signOutErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, nil).Res()
}

func (h *userHandler) SignUpAdmin(c *fiber.Ctx) error {
	// request body parser
	req := new(users.UserRegisterReq)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signUpAdminErr),
			err.Error(),
		).Res()
	}

	// email validation
	if !req.IsEmail() {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signUpCustomerErr),
			"email pattern is invalid",
		).Res()
	}
	//  Insert
	result, err := h.usersUsecase.InsertAdmin(req)
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

	return entities.NewResponse(c).Success(fiber.StatusCreated, result).Res()
}

func (h *userHandler) GenerateAdminToken(c *fiber.Ctx) error {
	adminToken, err := auth.NewAuth(
		auth.AdminToken,
		h.cfg.Jwt(),
		nil,
	)

	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(generateAdminTokenErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(
		fiber.StatusOK,
		&struct {
			Token string `json:"token"`
		}{
			Token: adminToken.SignToken(),
		},
	).Res()
}

func (h *userHandler) GetUserProfile(c *fiber.Ctx) error {
	// set params
	userId := strings.Trim(c.Params("user_id"), " ")

	// get profile
	result, err := h.usersUsecase.GetUserProfile(userId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(getUserProfileErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, result).Res()
}

func (h *userHandler) ChangePassword(c *fiber.Ctx) error {
	req := new(users.UserRegisterReq)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			changePasswordErr,
			err.Error(),
		).Res()
	}

	if err := h.usersUsecase.ChangePassword(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			changePasswordErr,
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, req).Res()

}
