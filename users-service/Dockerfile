FROM golang:1.21 AS builder

WORKDIR /code/service
COPY ./users-service/go.mod .
#COPY ./users-service/go.sum .
RUN go mod download all
COPY ./users-service/ .

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
    go build -a -ldflags '-extldflags "-static"' \
    -o svc \
    ./cmd/userssrvweb


FROM alpine:3

RUN apk update && apk add ca-certificates curl && rm -rf /var/cache/apk/*

WORKDIR /service

COPY --from=builder /code/service/svc               svc
#COPY --from=builder /code/service/config            config

ENV ENV development
#ENV CONFIG_FOLDER_PATH config/development
ENV SERVICE_NAME users-service

ENV PORT ":3001"
EXPOSE $PORT

CMD ./svc \
    -port=$PORT \
#    -configFolder=$CONFIG_FOLDER_PATH \