# Система автоматического размещения товаров на складе

Система состоит из нескольких микросервисов, каждый из которых отвечает за свой алгоритм размещения товаров.

## Структура проекта

```
warehouse/
└── services/
    ├── fixed-placement/           # Микросервис фиксированного размещения
    │   ├── cmd/
    │   │   └── api/              # Точка входа приложения
    │   ├── internal/
    │   │   ├── domain/           # Модели домена
    │   │   ├── repository/       # Работа с БД
    │   │   ├── service/          # Бизнес-логика
    │   │   └── handler/          # HTTP-обработчики
    │   └── pkg/
    │       └── database/         # Утилиты для работы с БД
    │
    ├── abc-placement/            # Микросервис ABC анализа
    ├── xyz-placement/            # Микросервис XYZ анализа
    ├── freе-placement/           # Микросервис свободного размещения
    ├── genetic-placement/        # Микросервис генетического размещения
    └── greedy-placement/         # Микросервис жадного алгоритма
```

## API Endpoints

- `POST /analyze` - анализ возможности размещения
- `POST /place` - размещение товара

## Запуск микросервисов

```bash
# Микросервис фиксированного размещения
go run services/fixed-placement/api/main.go

# Микросервис свободного размещения
go run services/free-placement/api/main.go

# Микросервис ABC анализа
go run services/abc-placement/cmd/api/main.go

# Микросервис XYZ анализа
go run services/xyz-placement/cmd/api/main.go

# Микросервис жадного алгоритма
go run services/greedy-placement/cmd/api/main.go

# Микросервис генетического алгоритма
go run services/genetic-placement/cmd/api/main.go

# Оркестратор
go run services/orchestrator/cmd/api/main.go
```

## Пример запроса к оркестратору

**Endpoint:** `http://localhost:8086/place`

**Request Body:**
```json
{
    "item_id": "ITEM001",
    "batch_id": "BATCH001",
    "quantity": 50,
    "weight": 5.0,
    "volume": 0.03,
    "turnover_rate": 0.85,
    "demand_rate": 0.05,
    "seasonality": 0.1,
    "abc_class": "A",
    "xyz_class": "X",
    "is_heavy": false,
    "is_fragile": false,
    "is_hazardous": false,
    "storage_temp": 20.0,
    "storage_humidity": 0.5,
    "warehouse_load": 0.6,
    "has_fixed_slot": true,
    "fast_access_zone": true
}
```