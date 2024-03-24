package controller_test

import (
	"fmt"
	"github.com/Waelson/go-o11y/internal/controller"
	"github.com/Waelson/go-o11y/internal/model"
	mock_service "github.com/Waelson/go-o11y/internal/tests/service"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestApplicationController_GetTemperature_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cep := "64770000"
	url := fmt.Sprintf("/?cep=%s", cep)
	mockResponse := model.ApplicationResponse{
		TempC: 40,
		TempF: 104,
		TempK: 313,
	}
	expectedResponse := "{\"temp_C\":40,\"temp_F\":104,\"temp_K\":313}\n"

	req := httptest.NewRequest(http.MethodGet, url, nil)
	w := httptest.NewRecorder()

	mock := mock_service.NewMockApplicationService(ctrl)
	mock.EXPECT().GetTemperature(cep).Return(mockResponse, nil).AnyTimes()

	c := controller.NewApplicationController(mock)
	c.Handler(w, req)
	res := w.Result()

	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)

	assert.Nil(t, err)
	assert.Equal(t, string(data), expectedResponse)
}

func TestApplicationController_GetTemperature_Error422(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	invalidCep := "6477000"
	url := fmt.Sprintf("/?cep=%s", invalidCep)

	req := httptest.NewRequest(http.MethodGet, url, nil)
	w := httptest.NewRecorder()

	mock := mock_service.NewMockApplicationService(ctrl)
	mock.EXPECT().GetTemperature(invalidCep).Return(model.ApplicationResponse{}, model.InvalidCepError).AnyTimes()

	c := controller.NewApplicationController(mock)
	c.Handler(w, req)
	res := w.Result()

	defer res.Body.Close()

	assert.Equal(t, res.StatusCode, http.StatusUnprocessableEntity)
}

func TestApplicationController_GetTemperature_Error404(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	invalidCep := "00000000"
	url := fmt.Sprintf("/?cep=%s", invalidCep)

	req := httptest.NewRequest(http.MethodGet, url, nil)
	w := httptest.NewRecorder()

	mock := mock_service.NewMockApplicationService(ctrl)
	mock.EXPECT().GetTemperature(invalidCep).Return(model.ApplicationResponse{}, model.CepNotFoundError).AnyTimes()

	c := controller.NewApplicationController(mock)
	c.Handler(w, req)
	res := w.Result()

	defer res.Body.Close()

	assert.Equal(t, res.StatusCode, http.StatusNotFound)
}

func TestApplicationController_GetTemperature_Error500(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cep := "64770000"
	url := fmt.Sprintf("/?cep=%s", cep)

	req := httptest.NewRequest(http.MethodGet, url, nil)
	w := httptest.NewRecorder()

	mock := mock_service.NewMockApplicationService(ctrl)
	mock.EXPECT().GetTemperature(cep).Return(model.ApplicationResponse{}, model.InternalError).AnyTimes()

	c := controller.NewApplicationController(mock)
	c.Handler(w, req)
	res := w.Result()

	defer res.Body.Close()

	assert.Equal(t, res.StatusCode, http.StatusInternalServerError)
}
