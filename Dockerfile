# syntax=docker/dockerfile:1

##
## Build
##
FROM golang:alpine as builder

WORKDIR /app

ENV GO111MODULE=on

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /server-api ./cmd/main.go

##
## Deploy
##
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /

#Copy executable from builder
COPY --from=builder /server-api /server-api

EXPOSE 1323

CMD [ "/server-api" ]