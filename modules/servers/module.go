package servers

import (
	"github.com/gofiber/fiber/v2"

	middlewarehandlers "github.com/Aritiaya50217/E-CommerceRESTAPIs/modules/middlewares/middlewareHandlers"
	middlewareusecase "github.com/Aritiaya50217/E-CommerceRESTAPIs/modules/middlewares/middlewareUsecase"
	middlewaresrepositories "github.com/Aritiaya50217/E-CommerceRESTAPIs/modules/middlewares/middlewaresRepositories"
	monitorhandlers "github.com/Aritiaya50217/E-CommerceRESTAPIs/modules/monitor/monitorHandlers"
	// usersRepositories "github.com/Aritiaya50217/E-CommerceRESTAPIs/modules/users/usersRepositories"
)

type IModuleFactory interface {
	MonitorModule()
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
	// repository := usersRepositories.UsersRepositories()
	// return repository
}
