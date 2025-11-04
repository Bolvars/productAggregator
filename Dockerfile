# ---------- Этап 1: сборка ----------
FROM golang:1.23-alpine AS builder

# Устанавливаем зависимости для сборки
RUN apk add --no-cache git

# Рабочая директория для сборки
WORKDIR /app

# Копируем go.mod и go.sum (для кеша зависимостей)
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем бинарь (без отладочной информации для уменьшения размера)
RUN go build -ldflags="-s -w" -o /app/bot ./cmd/main.go

# ---------- Этап 2: минимальный образ ----------
FROM alpine:3.20

# Рабочая директория для запуска
WORKDIR /app

# Копируем только бинарник из builder
COPY --from=builder /app/bot /app/bot

# Создаём non-root пользователя
RUN adduser -D botuser
USER botuser

# Переменная окружения для токена (можно передавать через docker run)
ENV TOKEN=""

# Запускаем бота
ENTRYPOINT ["/app/bot"]
CMD ["-token=${TOKEN}"]
