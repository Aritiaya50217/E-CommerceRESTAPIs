package servers

import (
	"github.com/gofiber/fiber/v2"

	appinfoHandlers "github.com/Aritiaya50217/E-CommerceRESTAPIs/modules/appinfo/appinfoHandlers"
	appinfoRepositories "github.com/Aritiaya50217/E-CommerceRESTAPIs/modules/appinfo/appinfoRepositories"
	appinfoUsecases "github.com/Aritiaya50217/E-CommerceRESTAPIs/modules/appinfo/appinfoUsecases"
	middlewarehandlers "github.com/Aritiaya50217/E-CommerceRESTAPIs/modules/middlewares/middlewareHandlers"
	middlewareusecase "github.com/Aritiaya50217/E-CommerceRESTAPIs/modules/middlewares/middlewareUsecase"
	middlewaresrepositories "github.com/Aritiaya50217/E-CommerceRESTAPIs/modules/middlewares/middlewaresRepositories"
	monitorhandlers "github.com/Aritiaya50217/E-CommerceRESTAPIs/modules/monitor/monitorHandlers"
	usershandlers "github.com/Aritiaya50217/E-CommerceRESTAPIs/modules/users/usersHandlers"
	usersRepositories "github.com/Aritiaya50217/E-CommerceRESTAPIs/modules/users/usersRepositories"
	usersusecases "github.com/Aritiaya50217/E-CommerceRESTAPIs/modules/users/usersUsecases"
)

type IModuleFactory interface {
	MonitorModule()
	UsersModule()
	AppinfoModule()
}

type moduleFactory struct {
	r           fiber.Router
	s           *server
	middlewares middlewarehandlers.IMiddlewaresHandler
}

func InitModule(r fiber.Router, s *server, mid middlewarehandlers.IMiddlewaresHandler) IModuleFactory {
	return &moduleFactory{
		r:           r,
		s:           s,
		middlewares: mid,
	}
}

func InitMiddlewares(s *server) middlewarehandlers.IMiddlewaresHandler {
	repository := middlewaresrepositories.MiddlewaresRepository(s.db)
	usecase := middlewareusecase.MiddlewaresUsecase(repository)
	return middlewarehandlers.MiddlewaresHandler(s.cfg, usecase)
}

func (m *moduleFactory) MonitorModule() {
	handler := monitorhandlers.MonitorHandler(m.s.cfg)

	m.r.Get("/", handler.HealthCheck)
}

func (m *moduleFactory) UsersModule() {
	repository := usersRepositories.UserRepository(m.s.db)
	usecase := usersusecases.UsersUsecase(m.s.cfg, repository)
	handler := usershandlers.UsersHandler(m.s.cfg, usecase)

	router := m.r.Group("/users")

	router.Post("/signup", handler.SignUpCustomer)
	router.Post("/signin", handler.SignIn)
	router.Post("/refresh", handler.RefreshPassport)
	router.Post("/signout", handler.SignOut)
	router.Post("/signup-admin", handler.SignUpAdmin)

	router.Get("/:user_id", handler.GetUserProfile)
	router.Get("/admin/secret", handler.GenerateAdminToken)
	router.Post("/change", handler.ChangePassword)
}

func (m *moduleFactory) AppinfoModule() {
	repository := appinfoRepositories.AppinfoRepository(m.s.db)
	usecase := appinfoUsecases.AppinfoUsecase(repository)
	handler := appinfoHandlers.AppinfoHandler(m.s.cfg, usecase)

	router := m.r.Group("/appinfo")

	router.Get("/categories", handler.FindCategory)
	router.Post("/categories", handler.AddCategory)
	router.Post("/categories/update", handler.UpdateCategory)
	router.Delete("/categories/:category_id", handler.RemoveCategory)

}
