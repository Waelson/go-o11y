package service_test

import (
	"github.com/Waelson/go-o11y/internal/model"
	"github.com/Waelson/go-o11y/internal/service"
	mock_requester "github.com/Waelson/go-o11y/internal/tests/requester"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestIntegrationService_GetCep_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cep := "64770001"
	cepJson := `
		{
		  "erro": true
		}
	`

	urls := model.NewModel()
	url := urls.GetCep(cep)

	m := mock_requester.NewMockHttpRequest(ctrl)

	m.EXPECT().Normalize(cep).Return(cep, nil).AnyTimes()
	m.EXPECT().MakeRequest(url).Return(cepJson, http.StatusOK, nil).AnyTimes()

	a := service.NewIntegrationService(m, model.NewModel())
	_, err := a.GetCep(cep)

	assert.NotNil(t, err)
	assert.Equal(t, err, model.CepNotFoundError)
}

func TestIntegrationService_GetTemperature_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cidade := "64770000"
	temperaturaJson := `
		{"location":{"name":"Brasília","region":"Distrito Federal","country":"Brasilien","lat":-15.78,"lon":-47.92,"tz_id":"America/Sao_Paulo","localtime_epoch":1710368861,"localtime":"2024-03-13 19:27"},"current":{"last_updated_epoch":1710368100,"last_updated":"2024-03-13 19:15","temp_c":22.0,"temp_f":71.6,"is_day":0,"condition":{"text":"Patchy light rain with thunder","icon":"//cdn.weatherapi.com/weather/64x64/night/386.png","code":1273},"wind_mph":6.9,"wind_kph":11.2,"wind_degree":70,"wind_dir":"ENE","pressure_mb":1019.0,"pressure_in":30.09,"precip_mm":5.1,"precip_in":0.2,"humidity":88,"cloud":75,"feelslike_c":24.6,"feelslike_f":76.3,"vis_km":10.0,"vis_miles":6.0,"uv":1.0,"gust_mph":8.4,"gust_kph":13.6}}
	`

	urls := model.NewModel()
	url := urls.GetTemperatura(cidade)

	m := mock_requester.NewMockHttpRequest(ctrl)

	m.EXPECT().Normalize(cidade).Return(cidade, nil).AnyTimes()
	m.EXPECT().MakeRequest(url).Return(temperaturaJson, http.StatusOK, nil).AnyTimes()

	a := service.NewIntegrationService(m, model.NewModel())
	temperaturaResponse, err := a.GetTemperatura(cidade)

	assert.Nil(t, err)
	assert.NotNil(t, temperaturaResponse)
}

func TestIntegrationService_GetTemperature_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cidade := "64770000"
	temperaturaJson := ""

	urls := model.NewModel()
	url := urls.GetTemperatura(cidade)

	m := mock_requester.NewMockHttpRequest(ctrl)

	m.EXPECT().Normalize(cidade).Return(cidade, nil).AnyTimes()
	m.EXPECT().MakeRequest(url).Return(temperaturaJson, http.StatusNotFound, nil).AnyTimes()

	a := service.NewIntegrationService(m, model.NewModel())
	_, err := a.GetTemperatura(cidade)

	assert.NotNil(t, err)
}
