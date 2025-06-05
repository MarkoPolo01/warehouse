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
│   ├── abc-placement/            
│   ├── xyz-placement/            
│   ├── dynamic-placement/     
│   └── optimal-placement/      
│
└── shared/                       # Общие компоненты
    ├── database/                 # Общие утилиты для работы с БД
    └── models/                   # Общие модели данных
```

**Эндпоинты:**
- `POST /analyze` - анализ возможности размещения
- `POST /place` - размещение товара

старт микросервиса фиксированного размещения 
go run services/fixed-placement/api/main.go
старт микросервиса свободного размещения 
go run services/fixed-placement/api/main.go
старт микросервиса ABC анализа
go run services/abc-placement/cmd/api/main.go
старт микросервиса XYZ анализа
go run services/xyz-placement/cmd/api/main.go
старт микросервиса жадного алгоритма
go run services/greedy-placement/cmd/api/main.go
 старт микросервиса генетического алгоритма
go run services/genetic-placement/cmd/api/main.go