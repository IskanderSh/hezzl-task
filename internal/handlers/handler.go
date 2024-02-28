package handlers

import (
	"log/slog"

	"github.com/gin-gonic/gin"
)

type GoodHandler struct {
	log             *slog.Logger
	serviceProvider ServiceProvider
}

type ServiceProvider interface {
}

func NewGoodHandler(log *slog.Logger, provider ServiceProvider) *GoodHandler {
	return &GoodHandler{log: log, serviceProvider: provider}
}

func (h *GoodHandler) InitRoutes() *gin.Engine {
	r := gin.New()

	good := r.Group("/good")
	{
		good.POST("/create", h.CreateGood)
		good.PATCH("/update", h.UpdateGood)
		good.DELETE("/delete", h.DeleteGood)
		good.PATCH("/reprioritize", h.ReprioritizeGood)
	}

	goods := r.Group("/goods")
	{
		goods.GET("/list", h.ListGoods)
	}

	return r
}
