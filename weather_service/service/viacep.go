package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"io"
	"log/slog"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// Estrutura que guarda a localidade do CEP
type ViaCEPResponse struct {
	Localidade string `json:"localidade"`
}

// Url a api de cep
var viaCepApiURL = "http://viacep.com.br/ws/"

var (
	tracer = otel.Tracer("temp-service")
)

/*
==========================================================
  - Função: GetlocationByCEP
  - Descrição : Função que busca dados da localidade para o
  - CEP informado.
  - Parametros :
  - ctx - contexto - tipo: context.Context
  - zipCode - CEP informado tipo: string
  - Retorno: Informações da localidade ou erro

==========================================================
*/
func GetLocationByCEP(ctx context.Context, zipCode string) (ViaCEPResponse, error) {
	var span trace.Span
	ctx, span = tracer.Start(ctx, "GetLocationByCEP")
	defer span.End()
	req, err := http.NewRequestWithContext(ctx, "GET", viaCepApiURL+zipCode+"/json/", nil)
	if err != nil {
		slog.Error("unable to make new request with context", "ctx", ctx, "error", err)
		return ViaCEPResponse{}, err
	}
	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("unable to do request", "req:", req.URL.Path, "error", err)
		return ViaCEPResponse{}, err
	}
	defer resp.Body.Close()
	result, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("error reading response", "response:", resp.Body, "error", err)
		return ViaCEPResponse{}, err
	}
	var viaCepData ViaCEPResponse
	err = json.Unmarshal(result, &viaCepData)
	if err != nil {
		slog.Error("error unmarshal result", "result:", string(result), "error", err)
		return ViaCEPResponse{}, err
	}
	if viaCepData.Localidade == "" {
		err = fmt.Errorf("error validating location: %s", viaCepData.Localidade)
		slog.Error("location is empty", "error", err)
		return ViaCEPResponse{}, err
	}
	return viaCepData, nil
}
