FROM golang:1.24.1-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY .env client.go main.go types.go util.go ./

RUN go build -o main .

CMD ["./main"]