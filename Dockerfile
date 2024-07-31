# Указываем базовый образ
FROM golang:1.17 as builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем файлы go.mod и go.sum для установки зависимостей
COPY go.mod go.sum ./

# Устанавливаем зависимости
RUN go mod download

# Копируем весь проект в контейнер
COPY . .

# Собираем бинарный файл
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/telegram-bot-suggest-news

# Создаем минимальный образ для запуска нашего приложения
FROM alpine:latest

# Копируем бинарный файл из предыдущего этапа
COPY --from=builder /app/telegram-bot-suggest-news /telegram-bot-suggest-news

# Копируем файл .env (если необходим)
COPY .env.example /app/.env

# Устанавливаем рабочую директорию и указываем команду запуска
WORKDIR /app
CMD ["/telegram-bot-suggest-news"]