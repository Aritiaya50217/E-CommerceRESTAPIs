package servers

import (
	"github.com/gofiber/fiber/v2"

	monitorhandlers "github.com/Aritiaya50217/E-CommerceRESTAPIs/modules/monitor/monitorHandlers"
)

type IModuleFactory interface {
	MonitorModule()
}

type moduleFactory struct {
	r fiber.Router
	s *server
}

func InitModule(r fiber.Router, s *server) IModuleFactory {
	return &moduleFactory{
		r: r,
		s: s,
	}
}

func (m *moduleFactory) MonitorModule() {
	handler := monitorhandlers.MonitorHandler(m.s.cfg)

	m.r.Get("/", handler.HealthCheck)
}
