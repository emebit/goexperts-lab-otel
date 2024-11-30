package handler

import (
	"Labs/goexperts-lab-otel/cep_service/internal/service"
	"encoding/json"
	"log/slog"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// Estrutura DTO de entrada que guarda Zipcode da temperatura
type TemperatureInputDTO struct {
	Zipcode string `json:"cep"`
}

// Estrutura DTO de saída que guarda dados da temperatura
type TemperatureOutputDTO struct {
	City       string `json:"city"`
	Celsius    string `json:"temp_C"`
	Fahrenheit string `json:"temp_F"`
	Kelvin     string `json:"temp_K"`
}

var tracer = otel.Tracer("cep-service")

/*
==========================================================
  - Função: HandleZipcode
  - Descrição : Função que valida e busca dados do clima
  - CEP informado.
  - Parametros :
  - w - Resposta do HTTP - tipo: http.ResponseWriter
  - r - Ponteiro para a requisição do HTTP tipo: http.Request
  - Retorno: Informações do clima no Response HTTP

==========================================================
*/
func HandleZipcode(w http.ResponseWriter, r *http.Request) {
	var span trace.Span
	ctx, span := tracer.Start(r.Context(), "handleZipcode")
	defer span.End()
	var inputDto TemperatureInputDTO
	err := json.NewDecoder(r.Body).Decode(&inputDto)
	if err != nil {
		slog.Error("unable to decode", "zipcode", inputDto.Zipcode, "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !service.ValidZipcode(inputDto.Zipcode) {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}
	response, err := service.GetWeatherByZipCode(ctx, inputDto.Zipcode)
	if err != nil {
		http.Error(w, `"unable to fetch temperature by zipcode"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
