package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Waelson/go-o11y/internal/model"
	"github.com/Waelson/go-o11y/internal/requester"
	"go.opentelemetry.io/otel/trace"
	"net/http"
)

type IntegrationService interface {
	GetCep(ctx context.Context, cep string) (model.CepResponse, error)
	GetTemperatura(ctx context.Context, cidade string) (model.TemperaturaResponse, error)
}

type integrationService struct {
	httpRequest requester.HttpRequest
	urls        model.Url
	tracer      trace.Tracer
}

func (s *integrationService) GetCep(ctx context.Context, cep string) (model.CepResponse, error) {
	fmt.Println(fmt.Sprintf("Pesquisando CEP %s", cep))

	ctx, span := s.tracer.Start(ctx, "github.com/Waelson/go-o11y/service/GetCep")
	defer span.End()

	cep, err := s.httpRequest.Normalize(cep)
	if err != nil {
		return model.CepResponse{}, model.InternalError
	}

	res, status, err := s.httpRequest.MakeRequest(ctx, s.urls.GetCep(cep))

	if status == http.StatusBadRequest {
		return model.CepResponse{}, model.InvalidCepError
	}

	if status == http.StatusNotFound {
		return model.CepResponse{}, model.CepNotFoundError
	}

	if status != http.StatusOK {
		return model.CepResponse{}, model.InternalError
	}

	var cepResponse model.CepResponse
	if err = json.Unmarshal([]byte(res), &cepResponse); err != nil {
		return model.CepResponse{}, err
	}

	if cepResponse.Error {
		return model.CepResponse{}, model.CepNotFoundError
	}

	return cepResponse, nil
}

func (s *integrationService) GetTemperatura(ctx context.Context, cidade string) (model.TemperaturaResponse, error) {
	fmt.Println(fmt.Sprintf("Pesquisando Temperatura %s", cidade))

	ctx, span := s.tracer.Start(ctx, "github.com/Waelson/go-o11y/service/GetTemperatura")
	defer span.End()

	cidade, err := s.httpRequest.Normalize(cidade)
	if err != nil {
		return model.TemperaturaResponse{}, model.InternalError
	}

	res, status, err := s.httpRequest.MakeRequest(ctx, s.urls.GetTemperatura(cidade))

	if status == http.StatusNotFound {
		return model.TemperaturaResponse{}, model.CepNotFoundError
	}

	if status != http.StatusOK {
		return model.TemperaturaResponse{}, model.InternalError
	}

	var temperaturaResponse model.TemperaturaResponse
	if err = json.Unmarshal([]byte(res), &temperaturaResponse); err != nil {
		return model.TemperaturaResponse{}, err
	}
	return temperaturaResponse, nil
}

func NewIntegrationService(httpRequest requester.HttpRequest, urls model.Url, tracer trace.Tracer) IntegrationService {
	return &integrationService{
		httpRequest: httpRequest,
		urls:        urls,
		tracer:      tracer,
	}
}
