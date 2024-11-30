# Goexperts-Lab-Otel

Sistema de consulta de temperatura por CEP.

##Como usar

### Pré-requisitos

- Go (1.22 ou maior) instalado
- Docker (opcional)

###Executando Localmente

1. Clone o repositório:

```bash
git clone git@github.com:emebit/goexperts-lab-otel.git
cd goexperts-lab-otel
```

2. Compile e execute os programas em terminais diferentes:

  ```bash
  go run cep_service/cmd/main.go 
```
```bash
    go run weather_service/cmd/main.go 
```

3. Fazer a requisição que está no arquivo api/api.http. O resultado será parecido com o seguinte:

    HTTP/1.1 200 OK

    Content-Type: application/json

    Date: Sat, 30 Nov 2024 22:09:43 GMT

    Content-Length: 68

    Connection: close
    
    {

      "city": "Curitiba",

      "temp_C": "22.2",

      "temp_F": "72.0",

      "temp_K": "295.2"

    }


### Executando com Docker

1. Dentro da pasta raiz, execute o comando para iniciar os serviços com Docker Compose:

```bash
docker compose up
```

2. Fazer a requisição que está no arquivo api/api.http. O resultado será parecido com o seguinte:

    HTTP/1.1 200 OK

    Content-Type: application/json

    Date: Sat, 30 Nov 2024 22:09:43 GMT

    Content-Length: 68

    Connection: close
    
    {

      "city": "Curitiba",

      "temp_C": "22.2",

      "temp_F": "72.0",

      "temp_K": "295.2"

    }

##Rastreamento

Este projeto utiliza OpenTelemetry e Zipkin para rastreamento distribuído, poed consultar as finrmações em: 
```bash
http://localhost:9411
```
