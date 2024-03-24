package service

import (
	"context"
	"fmt"
	"github.com/Waelson/go-o11y/internal/model"
	"go.opentelemetry.io/otel/trace"
	"strconv"
	"strings"
)

type ApplicationService interface {
	GetTemperature(ctx context.Context, cep string) (model.ApplicationResponse, error)
}

type applicationService struct {
	integrationService IntegrationService
	tracer             trace.Tracer
}

func (s *applicationService) GetTemperature(ctx context.Context, cep string) (model.ApplicationResponse, error) {
	fmt.Println(fmt.Sprintf("Vamos pesquisar o CEP %s", cep))

	ctx, span := s.tracer.Start(ctx, "github.com/Waelson/go-o11y/service/GetTemperature")
	defer span.End()

	cep = strings.TrimSpace(cep)
	if cep == "" || len(cep) != 8 || !isNumber(cep) {
		return model.ApplicationResponse{}, model.InvalidCepError
	}

	cepResponse, err := s.integrationService.GetCep(ctx, cep)
	if err != nil {
		return model.ApplicationResponse{}, err
	}

	temperaturaResponse, err := s.integrationService.GetTemperatura(ctx, strings.TrimSpace(cepResponse.Localidade))
	if err != nil {
		return model.ApplicationResponse{}, err
	}

	tempF := temperaturaResponse.Current.TempC*1.8 + 32
	tempK := temperaturaResponse.Current.TempC + 273

	return model.ApplicationResponse{
		City:  cepResponse.Localidade,
		TempC: temperaturaResponse.Current.TempC,
		TempF: tempF,
		TempK: tempK,
	}, nil
}

func isNumber(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func NewApplicationService(service IntegrationService, tracer trace.Tracer) ApplicationService {
	return &applicationService{integrationService: service, tracer: tracer}
}
