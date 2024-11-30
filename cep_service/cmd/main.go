/*
=====================================================================================================

  - main.go : Desenvolver um sistema em Go que receba um CEP, identifica a cidade e retorna o clima
  - atual (temperatura em graus celsius, fahrenheit e kelvin) juntamente com a cidade. Esse sistema
  - deverá implementar OTEL(Open Telemetry) e Zipkin.

  - Basedo no cenário conhecido "Sistema de temperatura por CEP" denominado Serviço B, será incluso
  - um novo projeto, denominado Serviço A.
  -
  - Requisitos - Serviço A (responsável pelo input):
  -
  - O sistema deve receber um input de 8 dígitos via POST, através do schema:  { "cep": "29902555" }
  - O sistema deve validar se o input é valido (contem 8 dígitos) e é uma STRING
  - Caso seja válido, será encaminhado para o Serviço B via HTTP
  - Caso não seja válido, deve retornar:
  - Código HTTP: 422
  - Mensagem: invalid zipcode
  -
  - Dica:
  -
  - Utilize a API viaCEP (ou similar) para encontrar a localização que deseja consultar a
  - temperatura: https://viacep.com.br/
=====================================================================================================
*/

package main

import (
	"github.com/emebit/goexperts-lab-otel/cep_service/internal/handler"
	"context"
	"log"
	"net/http"
	"os"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

// Variaveis para a url do Zipkin
var Host string
var URL_ZIPKIN string

/*
==========================================================
  - Função: init
  - Descrição : Função que cria a url do Zipkin
  - Parametros :
  - Retorno:

==========================================================
*/
func init() {
	Host = os.Getenv("URL_ZIPKIN")
	if Host == "" {
		Host = "localhost"
	}
	URL_ZIPKIN = "http://" + Host + ":9411/api/v2/spans"
}

/*
==========================================================
  - Função: initTracer
  - Descrição : Função que cria exporter do Zipkin e trace
  - provider do Otel
  - Parametros :
  - Retorno:

==========================================================
*/
func initTracer() func() {
	exporter, err := zipkin.New(
		URL_ZIPKIN,
	)
	if err != nil {
		log.Fatalf("failed to create Zipkin exporter: %v", err)
	}
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("service-A"),
		)),
	)
	otel.SetTracerProvider(tp)
	return func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Fatalf("failed to shutdown TracerProvider: %v", err)
		}
	}
}

func main() {
	log.Println("Starting server.")
	initTracer()
	http.HandleFunc("/cep", handler.HandleZipcode)
	err := http.ListenAndServe(":8080", otelhttp.NewHandler(http.DefaultServeMux, "http-server"))
	if err != nil {
		log.Fatal("error listen and serve", "error:", err)
	}
	log.Println("Started on port", "8080")
}
