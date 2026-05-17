# Building
FROM golang:tip-alpine3.23 AS builder

WORKDIR /app
COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o modemwatchdog .

# Runtime
FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/resolv.conf /etc/resolv.conf
COPY --from=builder /app/modemwatchdog /modemwatchdog

ENTRYPOINT ["/modemwatchdog"]
