FROM golang:alpine AS builder

WORKDIR /app

# dependency cache logic
COPY go.mod ./
# COPY go.sum ./  <-- Bizde bu dosya henuz yok, o yuzden yorum satirinda
RUN go mod download

COPY . .

# static binary build
RUN CGO_ENABLED=0 go build -o lite-redis main.go

# runner
FROM scratch

COPY --from=builder /app/lite-redis /lite-redis

EXPOSE 6379

ENTRYPOINT ["/lite-redis"]
