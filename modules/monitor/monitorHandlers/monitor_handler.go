package monitorhandlers

import (
	"github.com/Aritiaya50217/E-CommerceRESTAPIs/config"
	"github.com/Aritiaya50217/E-CommerceRESTAPIs/modules/entities"
	"github.com/Aritiaya50217/E-CommerceRESTAPIs/modules/monitor"
	"github.com/gofiber/fiber/v2"
)

type IMonitorHandler interface {
	HealthCheck(c *fiber.Ctx) error
}

type monitorHandler struct {
	cfg config.IConfig
}

func MonitorHandler(cfg config.IConfig) IMonitorHandler {
	return &monitorHandler{
		cfg: cfg,
	}
}

func (h *monitorHandler) HealthCheck(c *fiber.Ctx) error {
	res := &monitor.Monitor{
		Name:    h.cfg.App().Name(),
		Version: h.cfg.App().Version(),
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, res).Res()
}
