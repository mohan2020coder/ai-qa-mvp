# syntax=docker/dockerfile:1
FROM golang:1.22-alpine AS build
WORKDIR /app
COPY . .
RUN go build -o /bin/orchestrator

FROM alpine:3.19
WORKDIR /app
COPY --from=build /bin/orchestrator /app/orchestrator
COPY web /app/web

# Copy .env into the container
COPY .env /app/.env

EXPOSE 8080
ENTRYPOINT ["/app/orchestrator"]
