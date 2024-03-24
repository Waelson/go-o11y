package requester

import (
	"context"
	"fmt"
	"github.com/Waelson/go-o11y/internal/model"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
	"io/ioutil"
	"net/http"
	"net/url"
)

type HttpRequest interface {
	MakeRequest(ctx context.Context, url string) (string, int, error)
	Normalize(str string) (string, error)
}

type httpRequest struct {
	tracer trace.Tracer
}

func (h *httpRequest) Normalize(str string) (string, error) {
	parsedURL, err := url.Parse(str)
	if err != nil {
		return "", err
	}
	parsedURL.RawQuery = url.QueryEscape(parsedURL.Query().Encode())
	return parsedURL.String(), nil
}

func (h *httpRequest) MakeRequest(ctx context.Context, urlStr string) (string, int, error) {

	resp, err := otelhttp.Get(ctx, urlStr)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", resp.StatusCode, model.InternalError
	}

	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	if err != nil {
		return "", http.StatusInternalServerError, err
	}
	return string(body), http.StatusOK, nil
}

func NewHttpRequest(tracer trace.Tracer) HttpRequest {
	return &httpRequest{tracer: tracer}
}
