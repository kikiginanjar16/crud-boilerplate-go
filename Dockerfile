# Build stage
FROM golang:1.22 as builder

WORKDIR /app
COPY go.mod go.sum* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/server

# Runtime stage
FROM gcr.io/distroless/base-debian12

WORKDIR /
ENV APP_PORT=8080
COPY --from=builder /app/server /server

EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/server"]
