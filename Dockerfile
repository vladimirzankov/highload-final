# build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# dependencies
COPY go.mod go.sum ./
RUN go mod download

# sources
COPY . .

# static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o server ./

# runtime stage
FROM gcr.io/distroless/base-debian12

WORKDIR /app

COPY --from=builder /app/server .

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/app/server"]
