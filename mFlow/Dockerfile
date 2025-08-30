FROM golang:1.23 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /automation-app ./cmd/server

FROM gcr.io/distroless/base-debian12
COPY --from=builder /automation-app /usr/local/bin/automation-app
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/automation-app"]
