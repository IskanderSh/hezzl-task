package handlers

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/IskanderSh/hezzl-task/internal/lib/error/response"
	"github.com/IskanderSh/hezzl-task/internal/models"
	"github.com/gin-gonic/gin"
)

type GoodHandler struct {
	log             *slog.Logger
	serviceProvider ServiceProvider
}

type ServiceProvider interface {
	CreateGood(ctx context.Context, req *models.CreateRequest) (*models.Good, error)
	UpdateGood(ctx context.Context, req *models.UpdateRequest) (*models.Good, error)
	DeleteGood(ctx context.Context, req *models.DeleteRequest) (*models.DeleteResponse, error)
	GetGoods(ctx context.Context, limit, offset int) (*models.ListGoodsResponse, error)
}

func NewGoodHandler(log *slog.Logger, provider ServiceProvider) *GoodHandler {
	return &GoodHandler{log: log, serviceProvider: provider}
}

func (h *GoodHandler) InitRoutes() *gin.Engine {
	r := gin.New()

	good := r.Group("/good")
	{
		good.POST("/create/:projectId", h.CreateGood)
		good.PATCH("/update/:id&:projectId", h.UpdateGood)
		good.DELETE("/delete/:id&:projectId", h.DeleteGood)
		good.PATCH("/reprioritize:id&:projectId", h.ReprioritizeGood)
	}

	r.GET("/goods/list/:limit&:offset", h.ListGoods)

	return r
}

const (
	defaultLimit  = 10
	defaultOffset = 1

	projectCtx = "projectId"
	idCtx      = "id"
	limitCtx   = "limit"
	offsetCtx  = "offset"
)

func (h *GoodHandler) CreateGood(c *gin.Context) {
	const op = "handlers.CreateGood"

	log := h.log.With(slog.String("op", op))

	projectID, err := getID(c, projectCtx)
	if err != nil {
		response.NewErrorResponse(c, log, http.StatusBadRequest, err.Error())
	}

	log.Debug(fmt.Sprintf("successfully get project id: %d", projectID))

	var input models.CreateRequest

	if err := c.BindJSON(&input); err != nil {
		response.NewErrorResponse(c, log, http.StatusBadRequest, "invalid input body")
	}

	log.Debug(fmt.Sprintf("successfully bind input with name: %s", input.Name))

	input.ProjectID = projectID

	output, err := h.serviceProvider.CreateGood(c, &input)
	if err != nil {
		response.NewErrorResponse(c, log, http.StatusInternalServerError, "internal error")
	}

	c.JSON(http.StatusOK, output)
}

func (h *GoodHandler) UpdateGood(c *gin.Context) {
	const op = "handlers.UpdateGood"

	log := h.log.With(slog.String("op", op))

	id, err := getID(c, idCtx)
	if err != nil {
		response.NewErrorResponse(c, log, http.StatusBadRequest, err.Error())
	}

	projectId, err := getID(c, projectCtx)
	if err != nil {
		response.NewErrorResponse(c, log, http.StatusBadRequest, err.Error())
	}

	var input models.UpdateRequest
	if err := c.BindJSON(&input); err != nil {
		response.NewErrorResponse(c, log, http.StatusBadRequest, "invalid input body")
	}

	input.ID = id
	input.ProjectID = projectId

	output, err := h.serviceProvider.UpdateGood(c, &input)
	if err != nil {
		response.NewErrorResponse(c, log, http.StatusInternalServerError, "internal error")
	}

	c.JSON(http.StatusOK, output)
}

func (h *GoodHandler) DeleteGood(c *gin.Context) {
	const op = "handlers.DeleteGood"

	log := h.log.With(slog.String("op", op))

	id, err := getID(c, idCtx)
	if err != nil {
		response.NewErrorResponse(c, log, http.StatusBadRequest, err.Error())
	}

	projectId, err := getID(c, projectCtx)
	if err != nil {
		response.NewErrorResponse(c, log, http.StatusBadRequest, err.Error())
	}

	var input models.DeleteRequest
	if err := c.BindJSON(&input); err != nil {
		response.NewErrorResponse(c, log, http.StatusBadRequest, "invalid input body")
	}

	input.ID = id
	input.ProjectID = projectId

	output, err := h.serviceProvider.DeleteGood(c, &input)
	if err != nil {
		response.NewErrorResponse(c, log, http.StatusInternalServerError, "internal error")
	}

	c.JSON(http.StatusOK, output)
}

func (h *GoodHandler) ListGoods(c *gin.Context) {
	const op = "handlers.ListGoods"

	log := h.log.With(slog.String("op", op))

	limit := defaultLimit
	offset := defaultOffset

	limit, err := getID(c, limitCtx)
	if err != nil {
		response.NewErrorResponse(c, log, http.StatusBadRequest, err.Error())
	}

	offset, err = getID(c, offsetCtx)
	if err != nil {
		response.NewErrorResponse(c, log, http.StatusBadRequest, err.Error())
	}

	output, err := h.serviceProvider.GetGoods(c, limit, offset)
	if err != nil {
		response.NewErrorResponse(c, log, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, output)
}

func (h *GoodHandler) ReprioritizeGood(c *gin.Context) {

}

func getID(c *gin.Context, param string) (int, error) {
	id, ok := c.Get(param)
	if !ok {
		return 0, errors.New(fmt.Sprintf("no %s in query", param))
	}

	idInt, ok := id.(int)
	if !ok {
		return 0, errors.New(fmt.Sprintf("%s is of invalid type", param))
	}

	return idInt, nil
}
