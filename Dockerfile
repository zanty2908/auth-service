# Build stage
FROM golang:1.19-alpine3.16 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

# Run stage
FROM alpine:3.16
WORKDIR /app
COPY --from=builder /app/main .
COPY .env .env
COPY db/migration ./db/migration
COPY language/localizations_src ./language/localizations_src
COPY ./config/config.docker.yml /app/config/config.yml

EXPOSE 8080 7080
CMD [ "/app/main" ]
