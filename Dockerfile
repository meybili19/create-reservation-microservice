FROM golang:1.19-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN go build -o create-reservation-microservice main.go

EXPOSE 8080

CMD ["./create-reservation-microservice"]
