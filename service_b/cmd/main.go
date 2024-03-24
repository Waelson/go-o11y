package main

import (
	"context"
	"fmt"
	"github.com/Waelson/go-o11y/internal/controller"
	"github.com/Waelson/go-o11y/internal/model"
	"github.com/Waelson/go-o11y/internal/requester"
	"github.com/Waelson/go-o11y/internal/service"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	resource2 "go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"log"
)

func initTracerAuto() func(context.Context) error {

	client := otlptracegrpc.NewClient(
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint("otel-collector:4317"))

	exporter, err := otlptrace.New(context.Background(), client)

	if err != nil {
		log.Fatal("Could not set exporter: ", err)
	}
	resources, err := resource2.New(
		context.Background(),
		resource2.WithAttributes(
			attribute.String("service.name", "orquestrador"),
			attribute.String("application", "service-b"),
		),
	)
	if err != nil {
		log.Fatal("Could not set resources: ", err)
	}

	otel.SetTracerProvider(
		sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithSpanProcessor(sdktrace.NewBatchSpanProcessor(exporter)),
			sdktrace.WithSyncer(exporter),
			sdktrace.WithResource(resources),
		),
	)

	return exporter.Shutdown
}

func main() {
	fmt.Println("Iniciando aplicacao")

	cleanup := initTracerAuto()
	defer cleanup(context.Background())

	otel.SetTextMapPropagator(propagation.TraceContext{})
	tracer := otel.Tracer("service-b")

	httpRequest := requester.NewHttpRequest(tracer)
	urls := model.NewModel()
	integrationService := service.NewIntegrationService(httpRequest, urls, tracer)
	applicationService := service.NewApplicationService(integrationService, tracer)
	applicationController := controller.NewApplicationController(applicationService)

	r := gin.Default()
	r.Use(otelgin.Middleware("service-b"))
	r.GET("/temperatura", applicationController.Handler)

	log.Println("Iniciando o servidor na porta 8181...")
	_ = r.Run(":8181")
}
