FROM golang:1.22-alpine as builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN go build -o eventom-backend main.go

RUN chmod +x /app/eventom-backend

# build tiny docker image
FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/eventom-backend /app/private-key.pem /app/

CMD [ "/app/eventom-backend" ]