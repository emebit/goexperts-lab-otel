FROM golang:1.22.4

WORKDIR /app

COPY . .

RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o weather_service ./cmd

CMD ["./weather_service"]