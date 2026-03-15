# Demo App

Демонстрационное Go-приложение (HTTP-сервер) для изучения DevOps-пайплайна.
Используется как рабочая нагрузка в связке с [infra](https://github.com/animee744/infra) — стеком мониторинга и оркестрации.

## Что делает

Простой HTTP-сервер на Go stdlib с эндпоинтами:

| Эндпоинт | Описание |
|---|---|
| `GET /` | Информация о сервисе (имя, версия, статус) |
| `GET /health` | Health-check (статус + uptime) |
| `GET /metrics` | Prometheus-совместимые метрики |
| `GET /api/users` | Имитация списка пользователей |
| `GET /api/orders` | Имитация заказов (случайные данные) |
| `GET /api/slow` | Медленный ответ (100–2000ms) — для тестирования таймаутов |
| `GET /api/error` | Ответ с 30% шансом ошибки 500 — для тестирования мониторинга |

### Метрики (Prometheus)

Эндпоинт `/metrics` отдаёт в формате Prometheus:

- `demo_requests_total` — счётчик всех HTTP-запросов
- `demo_errors_total` — счётчик ошибок
- `demo_uptime_seconds` — время работы в секундах

## Запуск локально

```bash
go run main.go
# Сервер запустится на :8080
```

Можно задать порт через переменную окружения:

```bash
PORT=3000 go run main.go
```

## Docker

```bash
docker build -t demo-app .
docker run -p 8080:8080 demo-app
```

## CI/CD

При пуше в `master` GitHub Actions автоматически:

1. Прогоняет `go build` и `go vet`
2. Собирает Docker-образ
3. Пушит в GitHub Container Registry (`ghcr.io/animee744/demo-app:latest`)

Конфигурация: `.github/workflows/ci.yml`

## Как это работает в связке

```
push в master → GitHub Actions → Docker-образ в ghcr.io
                                         ↓
               infra: docker compose pull → обновление контейнера
                                         ↓
               Prometheus скрейпит /metrics → Grafana показывает дашборд
```

## Связанные репозитории

- [infra](https://github.com/animee744/infra) — Docker Compose, Prometheus, Grafana, Portainer, cAdvisor
- [deployer](https://github.com/animee744/deployer) — оригинальный Go-деплоер (учебный проект, заменён infra-подходом)
