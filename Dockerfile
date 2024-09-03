# Dockerfile for building a Docker image for the application

FROM golang:1.16.3-alpine3.13 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /app/sus . 


FROM alpine:3.13

COPY --from=builder app/sus /go/bin/sus 
RUN chmod +x go/bin/sus
CMD ["/go/bin/sus"]

