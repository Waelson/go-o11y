package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	resource2 "go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
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
			attribute.String("service.name", "input"),
			attribute.String("application", "service-a"),
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

type Endereco struct {
	CEP string `json:"cep"`
}

func main() {

	cleanup := initTracerAuto()
	defer cleanup(context.Background())

	otel.SetTextMapPropagator(propagation.TraceContext{})
	tracer := otel.Tracer("service-a")

	r := gin.Default()
	r.Use(otelgin.Middleware("service-a"))

	r.GET("/metrics", func(c *gin.Context) {
		h := promhttp.Handler()
		h.ServeHTTP(c.Writer, c.Request)
	})
	r.POST("/clima", func(c *gin.Context) {
		var endereco Endereco

		// Tenta fazer o bind do corpo da requisição para a struct EnderecoRequest
		if err := c.ShouldBindJSON(&endereco); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if len(endereco.CEP) != 8 {
			c.String(http.StatusUnprocessableEntity, "invalid zipcode")
			return
		}

		response, status, _ := obterClima(c.Request.Context(), endereco.CEP, tracer)
		if status == http.StatusOK {
			c.Header("Content-Type", "application/json")
			c.String(http.StatusOK, response)
		} else if status == http.StatusNotFound {
			c.String(http.StatusNotFound, "can not find zipcode")
		} else if status == http.StatusUnprocessableEntity {
			c.String(http.StatusUnprocessableEntity, "invalid zipcode")
		} else {
			c.String(http.StatusInternalServerError, "internal error")
		}

	})
	log.Println("Servidor iniciado na porta 8080")
	r.Run(":8080")
}

func obterClima(ctx context.Context, cep string, tracer trace.Tracer) (string, int, error) {
	ctx, span := tracer.Start(ctx, "github.com/Waelson/go-o11y/main.obterClima")
	defer span.End()

	baseURL := "http://service-b:8181/temperatura"
	//baseURL := "http://localhost:8181/temperatura"
	u, err := url.Parse(baseURL)
	if err != nil {
		log.Fatal(err)
	}
	q := u.Query()
	q.Set("cep", cep)
	u.RawQuery = q.Encode()

	fmt.Println("URL: ", u.String())

	resp, err := otelhttp.Get(ctx, u.String())

	if err != nil {
		return "", http.StatusInternalServerError, err
	}
	defer resp.Body.Close()

	// Ler o corpo da resposta
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", http.StatusInternalServerError, errors.New("erro interno")
	}

	return string(body), resp.StatusCode, nil
}
