FROM golang:alpine as builder

WORKDIR /go/src/app

ENV GO111MODULE=on

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o /go/src/app/simple-api ./cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

#Copy executable from builder
COPY --from=builder /go/src/app/simple-api .

EXPOSE 1323
CMD ["./simple-api"]
#CMD ["/usr/local/bin/simple-api"]