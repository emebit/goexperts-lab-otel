package services

import (
	"net/http"

	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/url"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
)

// Estrutura que guarda info do clima para a localidade
type WeatherAPIResponse struct {
	Current Current `json:"current"`
}

// Estrutura que guarda as temperaturas da localidade
type Current struct {
	CelsiusTemperature    float64 `json:"temp_c"`
	FarhenheitTemperature float64 `json:"temp_f"`
}

// Url a api de clima
var weatherApiURL = "http://api.weatherapi.com/v1/current.json"

// Chave da api de clima
const apiKey = "8d6494e11c514bb2882131719240511"

/*
==========================================================
  - Função: GetWeatherByCity
  - Descrição : Função que busca dados do clima para a
  - localidade informada.
  - Parametros :
  - ctx - contexto - tipo: context.Context
  - v - Localidade informada - tipo: ViaCEPResponse
  - Retorno: Informações do clima ou erro

==========================================================
*/
func GetWeatherByCity(ctx context.Context, v ViaCEPResponse) (WeatherAPIResponse, error) {
	var span trace.Span
	ctx, span = tracer.Start(ctx, "GetWeatherByCity")
	defer span.End()
	params := map[string]string{
		"key": apiKey,
		"q":   v.Localidade,
		"aqi": "no",
	}
	u, err := url.Parse(weatherApiURL)
	if err != nil {
		slog.Error("error parsing URL", "url", weatherApiURL, "error", err)
		return WeatherAPIResponse{}, err
	}
	q := u.Query()
	for key, value := range params {
		q.Set(key, value)
	}
	u.RawQuery = q.Encode()
	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		slog.Error("unable to make new request with context", "ctx", ctx, "error", err)
		return WeatherAPIResponse{}, err
	}
	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("error sending request", "query", u.RawQuery, "error", err)
		return WeatherAPIResponse{}, err
	}
	defer resp.Body.Close()
	result, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("error reading response", "response:", resp.Body, "error", err)
		return WeatherAPIResponse{}, err
	}
	var weather WeatherAPIResponse
	err = json.Unmarshal(result, &weather)
	if err != nil {
		slog.Error("error unmarshal result", "result:", string(result), "error", err)
		return WeatherAPIResponse{}, err
	}
	return weather, nil
}