FROM golang:1.21 AS builder

WORKDIR /code/infrastructure
COPY ./infrastructure/go.mod .
COPY ./infrastructure/go.sum .
RUN go mod download all
COPY    ./infrastructure/ .

WORKDIR /code/service
COPY ./contact-service/go.mod .
COPY ./contact-service/go.sum .
RUN go mod download all
COPY ./contact-service/ .

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
    go build -a -ldflags '-extldflags "-static"' \
    -o svc \
    ./cmd/userssrvweb


FROM alpine:3

RUN apk update && apk add ca-certificates curl && rm -rf /var/cache/apk/*

WORKDIR /service

COPY --from=builder /code/service/svc               svc

ENV ENV development
ENV SERVICE_NAME contact-service

ENV PORT ":3001"
EXPOSE $PORT

CMD ./svc \
    -port=$PORT \
