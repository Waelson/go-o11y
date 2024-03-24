package controller

import (
	"errors"
	"github.com/Waelson/go-o11y/internal/model"
	"github.com/Waelson/go-o11y/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

type Input struct {
	Cep string `form:"cep"`
}

type ApplicationController interface {
	Handler(a *gin.Context)
}
type applicationController struct {
	service service.ApplicationService
}

func (a *applicationController) Handler(c *gin.Context) {
	var input Input
	if c.ShouldBind(&input) != nil {
		c.Error(errors.New("erro ao obter o CEP"))
	}

	_, err := strconv.ParseInt(input.Cep, 10, 64)
	if err != nil {
		c.String(http.StatusUnprocessableEntity, "invalid zipcode")
		return
	}

	response, err := a.service.GetTemperature(c.Request.Context(), strings.TrimSpace(input.Cep))

	if errors.Is(err, model.InvalidCepError) {
		c.String(http.StatusUnprocessableEntity, "invalid zipcode")
	} else if errors.Is(err, model.CepNotFoundError) {
		c.String(http.StatusNotFound, "can not find zipcode")
	} else if errors.Is(err, model.InternalError) {
		c.String(http.StatusInternalServerError, "internal error")
	} else {
		c.JSON(http.StatusOK, response)
	}
}

func NewApplicationController(service service.ApplicationService) ApplicationController {
	return &applicationController{service: service}
}
