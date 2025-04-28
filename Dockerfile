# Базовый образ
FROM golang:1.23.8

# Установка необходимых пакетов (если они нужны)
#RUN apt-get update && apt-get install -y \
#   golang \
#   ca-certificates \
#   git \
#   wget \
#   build-essential \
#   libsqlite3-dev \
#   sqlite3 

# Рабочая директория
#RUN mkdir /app
WORKDIR /app

# Копирование файлов в контейнер
COPY . .

# Установка go mod 
RUN go mod download

# Создание файла app
RUN go build -o app ./cmd/main.go

# Настройка переменных окружения (если нужно будет)
#ENV TODO_PORT=7540
#ENV TODO_DBFILE=/host/scheduler.db
#ENV TODO_PASSWORD=12345

# Определение порта сервера
#EXPOSE 7540

# Другой вариант определения порта сервера
EXPOSE ${TODO_PORT}

# Запуск программы
CMD ["./app"]