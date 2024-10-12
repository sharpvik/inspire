FROM golang:1.23-bullseye AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o /app/bin ./cmd/teach

FROM scratch
WORKDIR /app
COPY --from=builder /app/bin /app/bin
CMD ["/app/bin"]
