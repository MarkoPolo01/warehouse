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

## Описание параметров запроса

### Основные идентификаторы
- `item_id` - уникальный идентификатор товара
- `batch_id` - идентификатор партии товара

### Физические характеристики
- `quantity` - количество единиц товара
- `weight` - вес одной единицы товара (в кг)
- `volume` - объем одной единицы товара (в м³)

### Параметры анализа
- `turnover_rate` - коэффициент оборачиваемости (0-1), показывает как часто товар перемещается
- `demand_rate` - коэффициент спроса (0-1), показывает стабильность спроса
- `seasonality` - коэффициент сезонности (0-1), показывает зависимость от сезона
- `abc_class` - класс товара по ABC-анализу (A, B, C)
  - A - наиболее ценные товары
  - B - товары средней ценности
  - C - наименее ценные товары
- `xyz_class` - класс товара по XYZ-анализу (X, Y, Z)
  - X - товары со стабильным спросом
  - Y - товары с умеренно стабильным спросом
  - Z - товары с нестабильным спросом

### Специальные характеристики
- `is_heavy` - является ли товар тяжелым
- `is_fragile` - является ли товар хрупким
- `is_hazardous` - является ли товар опасным

### Условия хранения
- `storage_temp` - требуемая температура хранения (в °C)
- `storage_humidity` - требуемая влажность хранения (0-1)

### Параметры склада
- `warehouse_load` - текущая загрузка склада (0-1)
- `has_fixed_slot` - есть ли у товара закрепленное место
- `fast_access_zone` - требуется ли размещение в зоне быстрого доступа