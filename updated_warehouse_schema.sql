-- SQL скрипт для создания базы данных автоматического размещения товаров (PostgreSQL)

-- 1. Таблица: Товары
CREATE TABLE items (
    item_id VARCHAR PRIMARY KEY,
    name VARCHAR NOT NULL,
    item_type VARCHAR,
    weight FLOAT NOT NULL,
    length FLOAT NOT NULL,
    width FLOAT NOT NULL,
    height FLOAT NOT NULL,
    storage_conditions TEXT,
    label_type VARCHAR,
    turnover FLOAT, -- для ABC-анализа
    мr FLOAT -- для XYZ-анализа (коэффициент вариации)
);

-- 2. Таблица: Партии товаров
CREATE TABLE batches (
    batch_id VARCHAR PRIMARY KEY,
    item_id VARCHAR REFERENCES items(item_id),
    quantity INT NOT NULL,
    production_date DATE,
    expiration_date DATE
);

-- 3. Таблица: Ячейки склада
CREATE TABLE slots (
    slot_id VARCHAR PRIMARY KEY,
    location_description TEXT,
    max_weight FLOAT,
    max_length FLOAT,
    max_width FLOAT,
    max_height FLOAT,
    storage_conditions TEXT,
    is_occupied BOOLEAN DEFAULT FALSE,
    zone_type VARCHAR, -- fast-access, regular, deep
    level INT, -- для определения верхних/нижних ярусов
    distance_from_exit INT -- для оценки логистики
);

-- 4. Таблица: Закреплённые ячейки (фиксированное размещение)
CREATE TABLE item_slot_map (
    id SERIAL PRIMARY KEY,
    item_id VARCHAR REFERENCES items(item_id),
    slot_id VARCHAR REFERENCES slots(slot_id)
);

-- 5. Таблица: Запросы на размещение
CREATE TABLE placement_requests (
    request_id SERIAL PRIMARY KEY,
    item_id VARCHAR REFERENCES items(item_id),
    batch_id VARCHAR REFERENCES batches(batch_id),
    quantity INT,
    special_conditions TEXT[],
    request_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 6. Таблица: Ответы системы на размещение
CREATE TABLE placement_responses (
    response_id SERIAL PRIMARY KEY,
    request_id INT REFERENCES placement_requests(request_id),
    success BOOLEAN,
    slot_id VARCHAR REFERENCES slots(slot_id),
    algorithm_used VARCHAR,
    score FLOAT,
    response_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    comment TEXT
);

-- 7. Таблица: История размещений
CREATE TABLE placement_logs (
    log_id SERIAL PRIMARY KEY,
    slot_id VARCHAR REFERENCES slots(slot_id),
    item_id VARCHAR REFERENCES items(item_id),
    batch_id VARCHAR REFERENCES batches(batch_id),
    algorithm VARCHAR,
    placed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    operator VARCHAR
);

-- 8. Таблица: Алгоритмы размещения
CREATE TABLE algorithms (
    algorithm_id SERIAL PRIMARY KEY,
    name VARCHAR NOT NULL,
    description TEXT,
    active BOOLEAN DEFAULT TRUE
);

-- Заполнение справочника товаров
INSERT INTO items (item_id, name, item_type, weight, length, width, height, storage_conditions, label_type, turnover, мr) VALUES
    ('ITEM001', 'Ноутбук Dell XPS', 'Электроника', 1.5, 30.0, 20.0, 2.0, 'Сухое помещение, комнатная температура', 'A', 0.8, 0.15),
    ('ITEM002', 'Смартфон iPhone', 'Электроника', 0.2, 15.0, 7.5, 1.0, 'Сухое помещение, комнатная температура', 'A', 0.9, 0.1),
    ('ITEM003', 'Молочные продукты', 'Продукты', 1.0, 20.0, 15.0, 10.0, 'Холодильная камера +4°C', 'B', 0.6, 0.3),
    ('ITEM004', 'Одежда летняя', 'Текстиль', 0.5, 25.0, 20.0, 5.0, 'Сухое помещение', 'C', 0.3, 0.4),
    ('ITEM005', 'Мебель офисная', 'Мебель', 25.0, 120.0, 60.0, 75.0, 'Сухое помещение', 'C', 0.2, 0.5);

-- Заполнение справочника ячеек склада
INSERT INTO slots (slot_id, location_description, max_weight, max_length, max_width, max_height, storage_conditions, is_occupied, zone_type, level, distance_from_exit) VALUES
    ('SLOT001', 'Зона быстрого доступа, уровень 1', 50.0, 100.0, 100.0, 200.0, 'Сухое помещение, комнатная температура', false, 'fast-access', 1, 10),
    ('SLOT002', 'Зона быстрого доступа, уровень 2', 50.0, 100.0, 100.0, 200.0, 'Сухое помещение, комнатная температура', false, 'fast-access', 2, 15),
    ('SLOT003', 'Холодильная камера, уровень 1', 30.0, 80.0, 80.0, 150.0, 'Холодильная камера +4°C', false, 'regular', 1, 30),
    ('SLOT004', 'Основное хранение, уровень 1', 100.0, 150.0, 150.0, 250.0, 'Сухое помещение', false, 'regular', 1, 50),
    ('SLOT005', 'Глубокое хранение, уровень 1', 200.0, 200.0, 200.0, 300.0, 'Сухое помещение', false, 'deep', 1, 100);

-- Заполнение справочника алгоритмов
INSERT INTO algorithms (name, description, active) VALUES
    ('Ближайший слот', 'Размещение в ближайшем доступном слоте', true),
    ('Оптимальный по размеру', 'Размещение с учетом оптимального использования пространства', true),
    ('ABC-анализ', 'Размещение на основе ABC-анализа товаров', true),
    ('XYZ-анализ', 'Размещение на основе XYZ-анализа товаров', true),
    ('Гибридный', 'Комбинированный алгоритм с учетом всех параметров', true);

INSERT INTO items (
    item_id, 
    name, 
    weight, 
    length, 
    width, 
    height, 
    turnover
) VALUES (
    'ITEM123',
    'Название товара',
    1.5,  -- вес
    30.0, -- длина
    20.0, -- ширина
    2.0,  -- высота
    0.85  -- turnover для ABC-анализа (>= 0.8 для категории A)
);

INSERT INTO batches (
    batch_id,
    item_id,
    quantity
) VALUES (
    'BATCH001',
    'ITEM123',
    10
);

INSERT INTO slots (
    slot_id,
    location_description,
    max_weight,
    max_length,
    max_width,
    max_height,
    storage_conditions,
    is_occupied,
    zone_type,
    level,
    distance_from_exit
) VALUES (
    'SLOT-A-001',
    'Зона быстрого доступа, уровень 1',
    50.0,  -- максимальный вес
    100.0, -- максимальная длина
    100.0, -- максимальная ширина
    200.0, -- максимальная высота
    'Сухое помещение, комнатная температура',
    false, -- не занят
    'fast-access',
    1,     -- уровень
    10     -- расстояние от выхода
);