FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w -X main.version=$(git describe --tags --always 2>/dev/null || echo dev)" -o /demo-app .

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
COPY --from=builder /demo-app /usr/local/bin/demo-app
EXPOSE 8080
CMD ["demo-app"]
