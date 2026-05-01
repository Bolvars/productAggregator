# ---------- Этап 1: сборка ----------
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git
WORKDIR /app

# Кешируем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем код
COPY . .

# Собираем бинарник. Обрати внимание на путь к main.go
RUN go build -ldflags="-s -w" -o /app/bot ./cmd/main.go

# ---------- Этап 2: минимальный образ ----------
FROM alpine:3.20

# Устанавливаем сертификаты (обязательно для https-запросов к API)
RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /app/bot /app/bot

RUN adduser -D botuser
USER botuser

# Значения по умолчанию для флагов (можно переопределить при запуске)
# Теперь запускаем бинарник напрямую, чтобы сигналы (SIGINT/SIGTERM) доходили до Go-процесса
ENTRYPOINT ["/app/bot"]
CMD ["-isTg=false", "-isWebhook=true"]
