-- Создание таблиц
CREATE TABLE IF NOT EXISTS items (
    item_id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    item_type VARCHAR(50) NOT NULL,
    weight FLOAT NOT NULL,
    length FLOAT NOT NULL,
    width FLOAT NOT NULL,
    height FLOAT NOT NULL,
    storage_conditions VARCHAR(100),
    label_type VARCHAR(50),
    turnover FLOAT NOT NULL, -- для ABC анализа
    mr FLOAT NOT NULL, -- для XYZ анализа
    is_heavy BOOLEAN DEFAULT false,
    is_fragile BOOLEAN DEFAULT false,
    is_hazardous BOOLEAN DEFAULT false,
    storage_temp FLOAT,
    storage_humidity FLOAT
);

CREATE TABLE IF NOT EXISTS batches (
    batch_id VARCHAR(50) PRIMARY KEY,
    item_id VARCHAR(50) REFERENCES items(item_id),
    quantity INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS slots (
    slot_id VARCHAR(50) PRIMARY KEY,
    location_description VARCHAR(100),
    max_weight FLOAT NOT NULL,
    max_length FLOAT NOT NULL,
    max_width FLOAT NOT NULL,
    max_height FLOAT NOT NULL,
    storage_conditions VARCHAR(100),
    is_occupied BOOLEAN DEFAULT false,
    zone_type VARCHAR(50) NOT NULL, -- 'fast-access', 'regular', 'deep'
    level INTEGER NOT NULL,
    distance_from_exit INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS item_slot_map (
    item_id VARCHAR(50) REFERENCES items(item_id),
    slot_id VARCHAR(50) REFERENCES slots(slot_id),
    PRIMARY KEY (item_id, slot_id)
);

CREATE TABLE IF NOT EXISTS placement_requests (
    request_id SERIAL PRIMARY KEY,
    item_id VARCHAR(50) REFERENCES items(item_id),
    batch_id VARCHAR(50) REFERENCES batches(batch_id),
    quantity INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS placement_responses (
    response_id SERIAL PRIMARY KEY,
    request_id INTEGER REFERENCES placement_requests(request_id),
    success BOOLEAN NOT NULL,
    slot_id VARCHAR(50) REFERENCES slots(slot_id),
    algorithm_used VARCHAR(50) NOT NULL,
    score FLOAT NOT NULL,
    comment TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS placement_logs (
    log_id SERIAL PRIMARY KEY,
    slot_id VARCHAR(50) REFERENCES slots(slot_id),
    item_id VARCHAR(50) REFERENCES items(item_id),
    batch_id VARCHAR(50) REFERENCES batches(batch_id),
    algorithm VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Вставка тестовых данных

-- Товары с разными характеристиками
INSERT INTO items (item_id, name, item_type, weight, length, width, height, storage_conditions, label_type, turnover, mr, is_heavy, is_fragile, is_hazardous, storage_temp, storage_humidity) VALUES
-- ABC класс A (высокая оборачиваемость)
('ITEM001', 'Популярный товар A', 'regular', 5.0, 0.5, 0.3, 0.2, 'normal', 'standard', 0.85, 0.05, false, false, false, 20.0, 0.5),
('ITEM002', 'Популярный товар B', 'regular', 8.0, 0.6, 0.4, 0.3, 'normal', 'standard', 0.82, 0.08, false, false, false, 20.0, 0.5),

-- ABC класс B (средняя оборачиваемость)
('ITEM003', 'Средний товар A', 'regular', 3.0, 0.4, 0.3, 0.2, 'normal', 'standard', 0.45, 0.15, false, false, false, 20.0, 0.5),
('ITEM004', 'Средний товар B', 'regular', 6.0, 0.5, 0.4, 0.3, 'normal', 'standard', 0.40, 0.18, false, false, false, 20.0, 0.5),

-- ABC класс C (низкая оборачиваемость)
('ITEM005', 'Редкий товар A', 'regular', 2.0, 0.3, 0.2, 0.1, 'normal', 'standard', 0.10, 0.30, false, false, false, 20.0, 0.5),
('ITEM006', 'Редкий товар B', 'regular', 4.0, 0.4, 0.3, 0.2, 'normal', 'standard', 0.08, 0.35, false, false, false, 20.0, 0.5),

-- XYZ класс X (стабильный спрос)
('ITEM007', 'Стабильный товар A', 'regular', 5.0, 0.5, 0.3, 0.2, 'normal', 'standard', 0.50, 0.05, false, false, false, 20.0, 0.5),
('ITEM008', 'Стабильный товар B', 'regular', 7.0, 0.6, 0.4, 0.3, 'normal', 'standard', 0.45, 0.08, false, false, false, 20.0, 0.5),

-- XYZ класс Y (умеренно изменчивый спрос)
('ITEM009', 'Изменчивый товар A', 'regular', 4.0, 0.4, 0.3, 0.2, 'normal', 'standard', 0.40, 0.20, false, false, false, 20.0, 0.5),
('ITEM010', 'Изменчивый товар B', 'regular', 6.0, 0.5, 0.4, 0.3, 'normal', 'standard', 0.35, 0.22, false, false, false, 20.0, 0.5),

-- XYZ класс Z (нестабильный спрос)
('ITEM011', 'Нестабильный товар A', 'regular', 3.0, 0.3, 0.2, 0.1, 'normal', 'standard', 0.30, 0.40, false, false, false, 20.0, 0.5),
('ITEM012', 'Нестабильный товар B', 'regular', 5.0, 0.4, 0.3, 0.2, 'normal', 'standard', 0.25, 0.45, false, false, false, 20.0, 0.5),

-- Специальные товары
('ITEM013', 'Тяжелый товар', 'heavy', 50.0, 1.0, 1.0, 1.0, 'normal', 'heavy', 0.60, 0.15, true, false, false, 20.0, 0.5),
('ITEM014', 'Хрупкий товар', 'fragile', 2.0, 0.3, 0.2, 0.1, 'fragile', 'fragile', 0.40, 0.20, false, true, false, 20.0, 0.5),
('ITEM015', 'Опасный товар', 'hazardous', 5.0, 0.5, 0.3, 0.2, 'hazardous', 'hazardous', 0.30, 0.25, false, false, true, 20.0, 0.5),
('ITEM016', 'Температурный товар', 'temperature', 3.0, 0.4, 0.3, 0.2, 'temperature', 'temperature', 0.50, 0.10, false, false, false, 5.0, 0.3);

-- Создание партий товаров
INSERT INTO batches (batch_id, item_id, quantity) VALUES
('BATCH001', 'ITEM001', 100),
('BATCH002', 'ITEM002', 50),
('BATCH003', 'ITEM003', 75),
('BATCH004', 'ITEM004', 25),
('BATCH005', 'ITEM005', 10),
('BATCH006', 'ITEM006', 5),
('BATCH007', 'ITEM007', 80),
('BATCH008', 'ITEM008', 40),
('BATCH009', 'ITEM009', 60),
('BATCH010', 'ITEM010', 30),
('BATCH011', 'ITEM011', 20),
('BATCH012', 'ITEM012', 15),
('BATCH013', 'ITEM013', 10),
('BATCH014', 'ITEM014', 25),
('BATCH015', 'ITEM015', 15),
('BATCH016', 'ITEM016', 30);

-- Создание ячеек склада
INSERT INTO slots (slot_id, location_description, max_weight, max_length, max_width, max_height, storage_conditions, is_occupied, zone_type, level, distance_from_exit) VALUES
-- Fast-access зона (близко к выходу)
('SLOT001', 'Fast-Access Zone 1', 10.0, 1.0, 1.0, 1.0, 'normal', false, 'fast-access', 1, 2),
('SLOT002', 'Fast-Access Zone 2', 10.0, 1.0, 1.0, 1.0, 'normal', false, 'fast-access', 1, 3),
('SLOT003', 'Fast-Access Zone 3', 10.0, 1.0, 1.0, 1.0, 'normal', false, 'fast-access', 1, 4),

-- Regular зона (средняя удаленность)
('SLOT004', 'Regular Zone 1', 10.0, 1.0, 1.0, 1.0, 'normal', false, 'regular', 1, 8),
('SLOT005', 'Regular Zone 2', 10.0, 1.0, 1.0, 1.0, 'normal', false, 'regular', 1, 9),
('SLOT006', 'Regular Zone 3', 10.0, 1.0, 1.0, 1.0, 'normal', false, 'regular', 1, 10),

-- Deep зона (далеко от выхода)
('SLOT007', 'Deep Zone 1', 10.0, 1.0, 1.0, 1.0, 'normal', false, 'deep', 1, 15),
('SLOT008', 'Deep Zone 2', 10.0, 1.0, 1.0, 1.0, 'normal', false, 'deep', 1, 16),
('SLOT009', 'Deep Zone 3', 10.0, 1.0, 1.0, 1.0, 'normal', false, 'deep', 1, 17),

-- Специальные зоны
('SLOT010', 'Heavy Zone', 100.0, 2.0, 2.0, 2.0, 'normal', false, 'regular', 1, 5),
('SLOT011', 'Fragile Zone', 5.0, 0.5, 0.5, 0.5, 'fragile', false, 'regular', 1, 6),
('SLOT012', 'Hazardous Zone', 10.0, 1.0, 1.0, 1.0, 'hazardous', false, 'regular', 1, 7),
('SLOT013', 'Temperature Zone', 10.0, 1.0, 1.0, 1.0, 'temperature', false, 'regular', 1, 8);

-- Фиксированные ячейки для некоторых товаров
INSERT INTO item_slot_map (item_id, slot_id) VALUES
('ITEM001', 'SLOT001'), -- Популярный товар A в fast-access зоне
('ITEM007', 'SLOT002'), -- Стабильный товар A в fast-access зоне
('ITEM013', 'SLOT010'), -- Тяжелый товар в специальной зоне
('ITEM014', 'SLOT011'), -- Хрупкий товар в специальной зоне
('ITEM015', 'SLOT012'), -- Опасный товар в специальной зоне
('ITEM016', 'SLOT013'); -- Температурный товар в специальной зоне 