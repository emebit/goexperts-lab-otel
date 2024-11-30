package handler

import (
	services "github.com/emebit/goexperts-lab-otel/weather_service/service"

	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// Estrutura DTO de saída que guarda dados da temperatura
type TemperatureOutputDTO struct {
	City       string `json:"city"`
	Celsius    string `json:"temp_C"`
	Fahrenheit string `json:"temp_F"`
	Kelvin     string `json:"temp_K"`
}

var (
	tracer = otel.Tracer("temp-service")
)

/*
==========================================================
  - Função: HandleTemperature
  - Descrição : Função que busca dados do clima para o
  - CEP informado e formata e converte as temperaturas.
  - Parametros :
  - w - Resposta do HTTP - tipo: http.ResponseWriter
  - r - Ponteiro para a requisição do HTTP tipo: http.Request
  - Retorno: Informações do clima no Response HTTP

==========================================================
*/
func HandleTemperature(w http.ResponseWriter, r *http.Request) {
	var span trace.Span
	ctx, span := tracer.Start(r.Context(), "HandleTemperature")
	defer span.End()
	zipCode := strings.Trim(r.URL.Path, "/cep/")
	location, err := services.GetLocationByCEP(ctx, zipCode)
	if err != nil {
		slog.Error("failed to fetch location by zipCode", "input:", zipCode, "error", err)
		http.Error(w, "cannot find zipCode", http.StatusNotFound)
		return
	}
	temperature, err := services.GetWeatherByCity(ctx, location)
	if err != nil {
		slog.Error("failed to fetch temperature by location", "input:", location.Localidade, "error", err)
		http.Error(w, "could not get weather", http.StatusInternalServerError)
		return
	}
	formatCelcius := fmt.Sprintf("%.1f", temperature.Current.CelsiusTemperature)
	var dto TemperatureOutputDTO = TemperatureOutputDTO{
		City:       location.Localidade,
		Celsius:    formatCelcius,
		Fahrenheit: ConvertCelsiusToFahrenheit(temperature.Current.CelsiusTemperature),
		Kelvin:     ConvertCelsiusToKelvin(temperature.Current.CelsiusTemperature),
	}
	byteJson, err := json.Marshal(dto)
	if err != nil {
		slog.Error("failed to marshal dto", dto, "error", err)
		http.Error(w, "could not get temperature", http.StatusInternalServerError)
		return
	}
	w.Write(byteJson)
}

/*
==========================================================
  - Função: ConvertCelsiusToFahrenheit
  - Descrição : Função que calcula a temperatura em
  - Fahrenheit a partir da temperatura em Celsius.
  - Parametros :
  - celsius - temperatura em Celsius - tipo: float64
  - Retorno: tempertura em Fahrenheit - tipo: string

==========================================================
*/
func ConvertCelsiusToFahrenheit(celsius float64) string {
	var f = celsius*1.8 + 32
	return fmt.Sprintf("%.1f", f)
}

/*
==========================================================
  - Função: ConvertCelsiusToKelvin
  - Descrição : Função que calcula a temperatura em
  - Kelvin a partir da temperatura em Celsius.
  - Parametros :
  - celsius - temperatura em Celsius - tipo: float64
  - Retorno: tempertura em Kelvin - tipo: string

==========================================================
*/
func ConvertCelsiusToKelvin(celsius float64) string {
	var k = celsius + 273
	return fmt.Sprintf("%.1f", k)
}
