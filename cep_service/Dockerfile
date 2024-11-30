FROM golang:1.22.4

WORKDIR /app

COPY . .

RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o cep_service ./cmd

CMD ["./cep_service"]