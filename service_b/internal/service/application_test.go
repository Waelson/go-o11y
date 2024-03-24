package service_test

import (
	"github.com/Waelson/go-o11y/internal/model"
	"github.com/Waelson/go-o11y/internal/service"
	mock_service "github.com/Waelson/go-o11y/internal/tests/service"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestApplicationService_GetTemperature_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedTempC := 40.0
	expectedTempF := 104.0
	expectedTempK := 313.0

	mock := mock_service.NewMockIntegrationService(ctrl)

	cep := "64770000"
	cepApiResponse := model.CepResponse{
		Cep:        cep,
		Localidade: "São Raimundo Nonato",
	}
	temperatureApiResponse := model.TemperaturaResponse{
		Current: model.CurrentResponse{
			TempC: expectedTempC,
		},
	}

	mock.EXPECT().GetCep(cep).Return(cepApiResponse, nil).AnyTimes()
	mock.EXPECT().GetTemperatura(cepApiResponse.Localidade).Return(temperatureApiResponse, nil).AnyTimes()

	s := service.NewApplicationService(mock)
	response, err := s.GetTemperature(cep)

	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, response.TempC, expectedTempC)
	assert.Equal(t, response.TempF, expectedTempF)
	assert.Equal(t, response.TempK, expectedTempK)
}

func TestApplicationService_GetTemperature_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cep := "64770000"

	expectedResponse := model.CepResponse{}
	expectedError := model.CepNotFoundError

	mock := mock_service.NewMockIntegrationService(ctrl)

	mock.EXPECT().GetCep(cep).Return(expectedResponse, expectedError).AnyTimes()

	s := service.NewApplicationService(mock)
	response, err := s.GetTemperature(cep)

	assert.NotNil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, err, expectedError)

}
