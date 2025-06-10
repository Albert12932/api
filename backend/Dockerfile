# Используем официальный образ Go
FROM golang:1.21

# Рабочая директория
WORKDIR /app

# Копируем модули
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходники
COPY . .

# Собираем приложение
RUN go build -o main .

# Указываем команду запуска
CMD ["./main"]
