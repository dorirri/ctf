FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o server .

FROM alpine:3.19

RUN addgroup -S ctf && adduser -S ctf -G ctf

WORKDIR /app
COPY --from=builder /app/server .
COPY --from=builder /app/migrations ./migrations

USER ctf

EXPOSE 8080
CMD ["./server"]
