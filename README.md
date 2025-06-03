# Система автоматического размещения товаров на складе

Система состоит из нескольких микросервисов, каждый из которых отвечает за свой алгоритм размещения товаров.

## Структура проекта

```
warehouse/
├── services/
│   ├── fixed-placement/           # Микросервис фиксированного размещения
│   │   ├── cmd/
│   │   │   └── api/              # Точка входа приложения
│   │   ├── internal/
│   │   │   ├── domain/           # Модели домена
│   │   │   ├── repository/       # Работа с БД
│   │   │   ├── service/          # Бизнес-логика
│   │   │   └── handler/          # HTTP-обработчики
│   │   └── pkg/
│   │       └── database/         # Утилиты для работы с БД
│   │
│   ├── abc-placement/            # Микросервис ABC-размещения (в разработке)
│   ├── xyz-placement/            # Микросервис XYZ-размещения (в разработке)
│   ├── dynamic-placement/        # Микросервис динамического размещения (в разработке)
│   └── optimal-placement/        # Микросервис оптимального размещения (в разработке)
│
└── shared/                       # Общие компоненты
    ├── database/                 # Общие утилиты для работы с БД
    └── models/                   # Общие модели данных
```

## Микросервисы

### 1. Fixed Placement Service
Сервис фиксированного размещения товаров. Размещает товары в заранее закрепленных за ними ячейках.

**Эндпоинты:**
- `POST /analyze` - анализ возможности размещения
- `POST /place` - размещение товара

### 2. ABC Placement Service (в разработке)
Сервис ABC-размещения. Размещает товары на основе ABC-анализа оборачиваемости.

### 3. XYZ Placement Service (в разработке)
Сервис XYZ-размещения. Размещает товары на основе XYZ-анализа предсказуемости спроса.

### 4. Dynamic Placement Service (в разработке)
Сервис динамического размещения. Размещает товары с учетом текущей загрузки склада и логистических маршрутов.

### 5. Optimal Placement Service (в разработке)
Сервис оптимального размещения. Использует комбинацию всех алгоритмов для выбора наилучшего места размещения.

## Запуск Fixed Placement Service

1. Создать базу данных PostgreSQL и применить схему из `updated_warehouse_schema.sql`
2. Настроить параметры подключения к БД в `services/fixed-placement/cmd/api/main.go`
3. Запустить сервис:
```bash
cd services/fixed-placement
go run cmd/api/main.go
```

## API Documentation

### Fixed Placement Service

#### Analyze Placement
```http
POST /analyze
Content-Type: application/json

{
    "item_id": "ITEM123",
    "batch_id": "BATCH001",
    "quantity": 10,
    "command": "analyze"
}
```

#### Place Item
```http
POST /place
Content-Type: application/json

{
    "item_id": "ITEM123",
    "batch_id": "BATCH001",
    "quantity": 10,
    "command": "place"
}
```

## Getting Started

### Prerequisites

- Go 1.16 or higher
- PostgreSQL

### Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/warehouse.git
cd warehouse
```

2. Install dependencies:
```bash
go mod download
```

3. Set up environment variables:
```bash
export SERVER_PORT=8080
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=postgres
export DB_NAME=warehouse
```

4. Run the application:
```bash
go run cmd/api/main.go
```

## API Endpoints

- `GET /health` - Health check endpoint
- `GET /api/v1/products` - Get all products
- `GET /api/v1/warehouses` - Get all warehouses
- `GET /api/v1/inventory` - Get inventory information

## Development

### Running Tests

```bash
go test ./...
```

### Building

```bash
go build -o bin/api cmd/api/main.go
```

## License

This project is licensed under the MIT License - see the LICENSE file for details. 