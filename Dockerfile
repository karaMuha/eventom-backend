FROM golang:1.22-alpine

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o eventom-backend main.go

EXPOSE 8080

CMD ["/app/eventom-backend"]