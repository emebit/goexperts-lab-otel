package service

import (
	"context"
	"io"
	"net/http"
	"os"
	"regexp"

	"log/slog"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// Variáveis do Tracer e URL
var (
	tracer = otel.Tracer("cep-service")
	URL    string
)

/*
==========================================================
  - Função: init
  - Descrição : Função que cria a url do cep
  - Parametros :
  - Retorno:

==========================================================
*/
func init() {
	Host := os.Getenv("URL_TEMP")
	if Host == "" {
		Host = "localhost"
	}
	URL = "http://" + Host + ":9090/cep/"
}

/*
==========================================================
  - Função: GetWeatherByZipCode
  - Descrição : Função que busca dados do clima para o
  - CEP informado.
  - Parametros :
  - ctx - contexto - tipo: context.Context
  - zipCode - CEP informado tipo: string
  - Retorno: Informações do clima no Response HTTP

==========================================================
*/
func GetWeatherByZipCode(ctx context.Context, zipCode string) ([]byte, error) {
	var span trace.Span
	ctx, span = tracer.Start(ctx, "GetWeatherByZipCode")
	defer span.End()
	req, err := http.NewRequestWithContext(ctx, "GET", URL+zipCode, nil)
	if err != nil {
		slog.Error("unable to make new request with context", "ctx", ctx, "error", err)
		return nil, err
	}
	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("unable to do request", "req:", req.URL.Path, "error", err)
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

/*
==========================================================
  - Função: ValidZipCode
  - Descrição : Função que verifica se o CEP é válido
  - Parametros :
  - zipCode - CEP - tipo: string
  - Retorno: Booleano

==========================================================
*/
func ValidZipcode(zipCode string) bool {
	re := regexp.MustCompile(`^\d{8}$`)
	return re.MatchString(zipCode)
}
