# Highload Metrics API

Высоконагруженный сервис для обработки потоковых данных с аналитикой на Go, развёрнутый в Kubernetes.

## Возможности

- Приём метрик через HTTP API (timestamp, CPU, RPS)
- Аналитика в реальном времени:
  - Rolling average (скользящее среднее, окно 50 событий)
  - Z-score детекция аномалий (порог 2σ)
- Кэширование в Redis
- Prometheus-метрики для мониторинга
- Автомасштабирование через Kubernetes HPA
- Обработка 1000+ RPS

## Стек технологий

| Компонент | Технология |
|-----------|------------|
| Язык | Go 1.25 |
| Роутер | gorilla/mux |
| Кэш | Redis 7 |
| Метрики | prometheus/client_golang |
| Контейнеризация | Docker |
| Оркестрация | Kubernetes / Minikube |
| Мониторинг | Prometheus + Grafana |
| Нагрузочное тестирование | Locust |

## Структура проекта

```
highload-final/
├── main.go           # Точка входа
├── api.go            # HTTP-обработчики
├── stats.go          # Аналитический модуль (rolling avg, z-score)
├── types.go          # Структуры данных
├── metrics.go        # Prometheus-метрики
├── cache.go          # Redis-клиент
├── routes.go         # HTTP-маршрутизация
├── Dockerfile        # Multi-stage сборка образа
├── go.mod            # Зависимости Go
├── go.sum
├── k8s/              # Kubernetes манифесты
│   ├── namespace.yaml
│   ├── app.yaml      # Deployment + Service
│   ├── redis.yaml    # Redis Deployment + Service
│   ├── hpa.yaml      # HorizontalPodAutoscaler
│   └── ingress.yaml  # Ingress
└── locust/           # Нагрузочное тестирование
    └── locustfile.py
```

## API

### POST /metrics

Приём метрики от клиента.

**Request:**
```json
{
  "timestamp": 1735123456,
  "cpu": 45.5,
  "rps": 850
}
```

**Response:** `202 Accepted`
```
ok
```

### GET /analyze

Получение результатов аналитики.

**Response:**
```json
{
  "avg": 523.7,
  "zscore": 1.24,
  "anomaly": false
}
```

| Поле | Описание |
|------|----------|
| `avg` | Rolling average (скользящее среднее) |
| `zscore` | Z-score текущего значения |
| `anomaly` | `true` если \|z-score\| > 2 |

### GET /metrics/prometheus

Экспорт метрик для Prometheus.

**Метрики:**
- `http_requests_total` — общее число запросов
- `anomalies_total` — число обнаруженных аномалий
- `http_latency_ms` — гистограмма задержек (мс)

## Быстрый старт

### Локальный запуск

```bash
# Запустить Redis
docker run -d --name redis -p 6379:6379 redis:7

# Запустить сервис
go run .

# Тест
curl -X POST http://localhost:8080/metrics \
  -H "Content-Type: application/json" \
  -d '{"timestamp": 123, "cpu": 50, "rps": 500}'

curl http://localhost:8080/analyze
```

### Docker

```bash
# Сборка образа
docker build -t highload-final-app:latest .

# Запуск
docker run -d --name redis -p 6379:6379 redis:7
docker run -d -p 8080:8080 -e REDIS_ADDR=host.docker.internal:6379 highload-final-app:latest
```

## Развертывание в Kubernetes

### Требования

- Minikube или другой Kubernetes-кластер
- kubectl
- Helm (для Prometheus/Grafana)

### Шаги

```bash
# 1. Запуск Minikube
minikube start --cpus=2 --memory=4g
minikube addons enable metrics-server
minikube addons enable ingress

# 2. Сборка и загрузка образа
docker build -t highload-final-app:latest .
minikube image load highload-final-app:latest

# 3. Применение манифестов
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/redis.yaml
kubectl apply -f k8s/app.yaml
kubectl apply -f k8s/hpa.yaml
kubectl apply -f k8s/ingress.yaml

# 4. Проверка
kubectl get all -n metrics-api
kubectl get hpa -n metrics-api
```

### Port-Forward

```bash
# Доступ к сервису
kubectl port-forward svc/metrics-api -n metrics-api 8080:80
```

## Мониторинг

### Установка Prometheus и Grafana

```bash
# Добавление репозиториев Helm
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo add grafana https://grafana.github.io/helm-charts
helm repo update

# Установка Prometheus
helm install prometheus prometheus-community/prometheus \
  --namespace monitoring \
  --create-namespace

# Установка Grafana
helm install grafana grafana/grafana \
  --namespace monitoring \
  --set adminPassword=admin
```

### Доступ к UI

```bash
# Prometheus
kubectl port-forward svc/prometheus-server -n monitoring 9090:80

# Grafana
kubectl port-forward svc/grafana -n monitoring 3000:80
```

- Prometheus: http://localhost:9090
- Grafana: http://localhost:3000 (admin/admin)

### Настройка Grafana

1. Data Sources → Add data source → Prometheus
2. URL: `http://prometheus-server:80`
3. Save & Test

**Примеры запросов для дашборда:**

```promql
# RPS
rate(http_requests_total[1m])

# Общее число запросов
http_requests_total

# Аномалии
anomalies_total

# Latency p95
histogram_quantile(0.95, rate(http_latency_ms_bucket[1m]))
```

## Нагрузочное тестирование

### Установка Locust

```bash
cd locust
python3 -m venv venv
source venv/bin/activate
pip install locust
```

### Запуск теста

```bash
# Убедитесь, что port-forward активен
kubectl port-forward svc/metrics-api -n metrics-api 8080:80

# Запуск Locust
locust -f locustfile.py --host=http://localhost:8080 --csv=report
```

Откройте http://localhost:8089 и настройте:
- Number of users: 200
- Spawn rate: 50
- Нажмите START

### Мониторинг HPA во время теста

```bash
kubectl get hpa -n metrics-api -w
```

При превышении CPU > 70% HPA автоматически увеличит число реплик (2 → 4 → 5).

## Конфигурация

### Переменные окружения

| Переменная | Описание | По умолчанию |
|------------|----------|--------------|
| `REDIS_ADDR` | Адрес Redis | `localhost:6379` |

### HPA параметры

| Параметр | Значение |
|----------|----------|
| Min replicas | 2 |
| Max replicas | 5 |
| Target CPU | 70% |

### Ресурсы пода

| Ресурс | Request | Limit |
|--------|---------|-------|
| CPU | 100m | 500m |
| Memory | 128Mi | 256Mi |

## Архитектура

```
┌──────────────┐
│    Client    │
│   (Locust)   │
└──────┬───────┘
       │ POST /metrics
       ▼
┌──────────────┐     channel      ┌──────────────┐
│   HTTP API   │ ───────────────► │   Analyzer   │
│   (api.go)   │                  │  (stats.go)  │
└──────┬───────┘                  └──────────────┘
       │                                 │
       │ cache                           │ metrics
       ▼                                 ▼
┌──────────────┐                  ┌──────────────┐
│    Redis     │                  │  Prometheus  │
└──────────────┘                  └──────────────┘
```

**Поток данных:**
1. Клиент отправляет метрику на `/metrics`
2. API отправляет метрику в канал аналитического модуля
3. Analyzer обновляет скользящее окно и вычисляет статистики
4. Последняя метрика кэшируется в Redis
5. Результат доступен через `/analyze`
6. Prometheus собирает метрики с `/metrics/prometheus`

## Аналитический модуль

### Rolling Average

Скользящее среднее по окну из 50 последних значений:

```go
func avg(data []float64) float64 {
    sum := 0.0
    for _, v := range data {
        sum += v
    }
    return sum / float64(len(data))
}
```

### Z-Score и детекция аномалий

```go
z := (currentValue - mean) / stddev

if math.Abs(z) > 2.0 {
    // Аномалия обнаружена
}
```

## Лицензия

MIT

