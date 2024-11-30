/*
=====================================================================================================

  - main.go : Desenvolver um sistema em Go que receba um CEP, identifica a cidade e retorna o clima

  - atual (temperatura em graus celsius, fahrenheit e kelvin) juntamente com a cidade. Esse sistema

  - deverá implementar OTEL(Open Telemetry) e Zipkin.

  - Basedo no cenário conhecido "Sistema de temperatura por CEP" denominado Serviço B, será incluso

  - um novo projeto, denominado Serviço A.
    -

  - Requisitos - Serviço B (responsável pela orquestração):
    -

  - O sistema deve receber um CEP válido de 8 digitos

  - O sistema deve realizar a pesquisa do CEP e encontrar o nome da localização, a partir disso,

  - deverá retornar as temperaturas e formata-lás em: Celsius, Fahrenheit, Kelvin juntamente com o

  - nome da localização.

  - O sistema deve responder adequadamente nos seguintes cenários:

  - Em caso de sucesso:

  - Código HTTP: 200

  - Response Body: { "city: "São Paulo", "temp_C": 28.5, "temp_F": 28.5, "temp_K": 28.5 }

  - Em caso de falha, caso o CEP não seja válido (com formato correto):

  - Código HTTP: 422

  - Mensagem: invalid zipcode

​ ​​ - Em caso de falha, caso o CEP não seja encontrado:
  - Código HTTP: 404
  - Mensagem: can not find zipcode
    -
  - Dicas:
    -
  - Utilize a API WeatherAPI (ou similar) para consultar as temperaturas desejadas:
  - https://www.weatherapi.com/
  - Para realizar a conversão de Celsius para Fahrenheit, utilize a seguinte fórmula: F = C * 1,8 + 32
  - Para realizar a conversão de Celsius para Kelvin, utilize a seguinte fórmula: K = C + 273
  - Sendo F = Fahrenheit
  - Sendo C = Celsius
  - Sendo K = Kelvin

=====================================================================================================
*/
package main

import (
	"github.com/emebit/goexperts-lab-otel/weather_service/internal/handler"
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
const portNum string = ":9090"

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
			semconv.ServiceNameKey.String("service-B"),
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
	log.Println("Starting http server.")
	initTracer()
	http.HandleFunc("/cep/{zipcode}", handler.HandleTemperature)
	err := http.ListenAndServe(portNum, otelhttp.NewHandler(http.DefaultServeMux, "http-server"))
	if err != nil {
		log.Fatal("error listen and serve", "error:", err)
	}
	log.Println("Started on port", portNum)
}
