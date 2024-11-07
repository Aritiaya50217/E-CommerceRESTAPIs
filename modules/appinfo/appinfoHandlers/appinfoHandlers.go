package appinfohandlers

import (
	"strconv"
	"strings"

	"github.com/Aritiaya50217/E-CommerceRESTAPIs/config"
	"github.com/Aritiaya50217/E-CommerceRESTAPIs/modules/appinfo"
	appinfoUsecases "github.com/Aritiaya50217/E-CommerceRESTAPIs/modules/appinfo/appinfoUsecases"
	"github.com/Aritiaya50217/E-CommerceRESTAPIs/modules/entities"
	"github.com/Aritiaya50217/E-CommerceRESTAPIs/pkg/auth"
	"github.com/gofiber/fiber/v2"
)

type appinfoHandlersErrCode string

const (
	generateApiKeyErr appinfoHandlersErrCode = "appinfo-001"
	findCategoryErr   appinfoHandlersErrCode = "appinfo-002"
	addCategoryErr    appinfoHandlersErrCode = "appinfo-003"
	updateCategoryErr appinfoHandlersErrCode = "appinfo-004"
	removeCategoryErr appinfoHandlersErrCode = "appinfo-005"
)

type IAppinfoHandler interface {
	GenerateApiKey(c *fiber.Ctx) error
	FindCategory(c *fiber.Ctx) error
	AddCategory(c *fiber.Ctx) error
	RemoveCategory(c *fiber.Ctx) error
	UpdateCategory(c *fiber.Ctx) error
}

type appinfoHandler struct {
	cfg             config.IConfig
	appinfoUsecases appinfoUsecases.IAppinfoUsecase
}

func AppinfoHandler(cfg config.IConfig, appinfoUsecase appinfoUsecases.IAppinfoUsecase) IAppinfoHandler {
	return &appinfoHandler{
		cfg:             cfg,
		appinfoUsecases: appinfoUsecase,
	}
}

func (h *appinfoHandler) GenerateApiKey(c *fiber.Ctx) error {
	apiKey, err := auth.NewAuth(
		auth.ApiKeyToken,
		h.cfg.Jwt(),
		nil,
	)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(generateApiKeyErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(
		fiber.StatusOK,
		&struct {
			Key string `json:"key"`
		}{
			Key: apiKey.SignToken(),
		},
	).Res()
}

func (h *appinfoHandler) FindCategory(c *fiber.Ctx) error {
	req := new(appinfo.CategoryFilter)
	if err := c.QueryParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(findCategoryErr),
			err.Error(),
		).Res()
	}

	category, err := h.appinfoUsecases.FindCategory(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(findCategoryErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, category).Res()
}

func (h *appinfoHandler) AddCategory(c *fiber.Ctx) error {
	req := make([]*appinfo.Category, 0)
	if err := c.BodyParser(&req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(addCategoryErr),
			err.Error(),
		).Res()
	}

	// check empty
	if len(req) == 0 {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(addCategoryErr),
			"categories request are empty",
		).Res()
	}

	if err := h.appinfoUsecases.InsertCategory(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(addCategoryErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusCreated, nil).Res()
}

func (h *appinfoHandler) RemoveCategory(c *fiber.Ctx) error {
	categoryIdStr := strings.Trim(c.Params("category_id"), " ")
	id, err := strconv.Atoi(categoryIdStr)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(removeCategoryErr),
			"id type is invalid",
		).Res()
	}

	if err := h.appinfoUsecases.DeleteCategory(id); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(removeCategoryErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(
		fiber.StatusOK,
		&struct {
			CategoryId int `json:"category_id"`
		}{
			CategoryId: id,
		},
	).Res()
}

func (h *appinfoHandler) UpdateCategory(c *fiber.Ctx) error {
	var req *appinfo.Category
	if err := c.BodyParser(&req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(updateCategoryErr),
			err.Error(),
		).Res()
	}

	if err := h.appinfoUsecases.UpdateCategory(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(updateCategoryErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, nil).Res()
}
