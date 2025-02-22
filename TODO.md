# Пилим Авиасейлс

<!-- 1. Роздать реп кинуть инвайт mtvy
2.
    Запилить ручки

Post /flight body{filght_id, destination_from, destination_to …} rest 203, id
Get /flight/:id response 200, body{…}
Put /flight body{…} reps 200, id
Delete /flight/:id  -->


<!-- 3. Проверить в postman -->

4. Спайс vs массив, структура массива, что происходит при append https://www.youtube.com/watch?v=1vAIvqDo5LE&t=6s
5. Маппа, бакеты, миграции, коллизии
6. gmp, горутины, шедулер

7. Перейти на БД в комментах есть
8. Добавить docker compose
   ./docker-compose.yaml

```yaml
services:
    go-api:
    build: ./
    ports:
      - "8080:8080"  
      - "8081:8081"
    depends_on:
      postgres:
        condition: service_healthy

    postgres:
    image: postgres:14.10-alpine
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_HOST_AUTH_METHOD: trust
    volumes:
      - pgdata:/var/lib/postgresql/data  
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 5s
      timeout: 10s
      retries: 5
```

Для сборки сервиса прописать ./Dockerfile

```Dockerfile
FROM golang:1.23-alpine AS builder

WORKDIR /

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /app/api-gateway ./cmd/
RUN ls -l /app

FROM alpine:latest 

WORKDIR /app
COPY --from=builder /app/api-gateway .

ENTRYPOINT ["./api-gateway"]
```

1. docker compsoe build
2. docker compose up (docker compose up -d)

docker compsoe up --build