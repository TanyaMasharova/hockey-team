-- Включаем расширения
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 1. Пользователи
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    phone VARCHAR(20) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 2. Соперники
CREATE TABLE IF NOT EXISTS opponents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE,
    city VARCHAR(100),
    logo_url TEXT
);

-- 3. Матчи
CREATE TABLE IF NOT EXISTS matches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    opponent_id UUID NOT NULL REFERENCES opponents(id) ON DELETE RESTRICT,
    match_date TIMESTAMP WITH TIME ZONE NOT NULL,
    home_away VARCHAR(10) CHECK (home_away IN ('home', 'away')), --домашний/выездной матч
    our_score SMALLINT DEFAULT 0, --наши голы
    opponent_score SMALLINT DEFAULT 0, --голы соперника
    season VARCHAR(9), --сезон
    status VARCHAR(20) DEFAULT 'scheduled' CHECK (status IN ('scheduled', 'live', 'finished', 'cancelled')),
    is_derby BOOLEAN DEFAULT FALSE
);

-- 4. Секторы стадиона
CREATE TABLE IF NOT EXISTS stadium_sectors (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sector_number VARCHAR(10) NOT NULL UNIQUE,
    capacity INT NOT NULL CHECK (capacity > 0),
    sector_type VARCHAR(30) NOT NULL DEFAULT 'standard' 
        CHECK (sector_type IN ('vip', 'standard', 'away_fans')),
    price_coefficient DECIMAL(4,2) DEFAULT 1.0 CHECK (price_coefficient >= 0.5),
    color_code VARCHAR(7)
);

-- 5. Места
CREATE TABLE IF NOT EXISTS seats (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sector_id UUID NOT NULL REFERENCES stadium_sectors(id) ON DELETE CASCADE,
    seat_row VARCHAR(10) NOT NULL,
    seat_number VARCHAR(10) NOT NULL,
    is_handicap_accessible BOOLEAN DEFAULT FALSE, --инвалидное 
    UNIQUE(sector_id, seat_row, seat_number)
);

-- 6. Билеты
CREATE TABLE IF NOT EXISTS tickets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    match_id UUID NOT NULL REFERENCES matches(id) ON DELETE RESTRICT,
    seat_id UUID NOT NULL REFERENCES seats(id),
    final_price DECIMAL(10,2) NOT NULL CHECK (final_price >= 0),
    purchase_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, --дата покупки
    qr_code_hash VARCHAR(255) UNIQUE, --уникальный qr код
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'used', 'refunded', 'cancelled'))
);

-- 7. Игроки
CREATE TABLE IF NOT EXISTS players (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    patronymic VARCHAR(100),
    birth_date DATE NOT NULL,
    age SMALLINT,
    handedness VARCHAR(10) CHECK (handedness IN ('left', 'right')), --хват
    citizenship VARCHAR(3), --гражданство
    height_cm SMALLINT CHECK (height_cm BETWEEN 140 AND 230),
    weight_kg SMALLINT CHECK (weight_kg BETWEEN 50 AND 150),
    position VARCHAR(20) CHECK (position IN ('goalie', 'defense', 'forward'))
);

-- 8. Статистика игроков за матч
CREATE TABLE IF NOT EXISTS player_match_stats (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE, --Id игрока
    match_id UUID NOT NULL REFERENCES matches(id) ON DELETE CASCADE, --id матча
    goals INT DEFAULT 0 CHECK (goals >= 0), --голы 
    shots INT DEFAULT 0 CHECK (shots >= 0), --броски
    assists INT DEFAULT 0 CHECK (assists >= 0), --передачи
    blocked_shots INT DEFAULT 0 CHECK (blocked_shots >= 0), --заблокировано бросков
    saves INT DEFAULT 0 CHECK (saves >= 0), --сэйвы
    goals_against INT DEFAULT 0 CHECK (goals_against >= 0), --пропущено голов
    plus_minus SMALLINT, --коэф полезности
    penalties INT DEFAULT 0, --кол-во штрафных минут
    UNIQUE(player_id, match_id)
);

-- 9. Итоговая статистика игрока за сезон
CREATE TABLE IF NOT EXISTS player_season_stats (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    season VARCHAR(9) NOT NULL,

    total_goals INT DEFAULT 0,
    total_assists INT DEFAULT 0,

    total_shots INT DEFAULT 0,
    total_blocked_shots INT DEFAULT 0,
    total_saves INT DEFAULT 0,
    total_goals_against INT DEFAULT 0,
    total_plus_minus INT DEFAULT 0,
    total_penalties INT DEFAULT 0,
    matches_played INT DEFAULT 0,
    UNIQUE(player_id, season)
);

-- 10. Правила цен
CREATE TABLE IF NOT EXISTS pricing_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    match_id UUID REFERENCES matches(id) ON DELETE CASCADE, --если правило Только для конкретного матча (например, фмнал)
    sector_type VARCHAR(30) ,
    rule_name VARCHAR(100) NOT NULL, --текст правила
    priority INT DEFAULT 0, --Чем МЕНЬШЕ число, тем РАНЬШЕ применяется
    condition_type VARCHAR(30) CHECK (condition_type IN ('early_bird', 'last_minute', 'derby', 'playoff', 'holiday')), --для иконки во вронт
    discount_percent DECIMAL(5,2) DEFAULT 0 CHECK (discount_percent >= -100 AND discount_percent <= 100), --процент скидки
    fixed_price DECIMAL(10,2), --не обязательно. Если указано fixed_price, то discount_percent игнорируется. 
    is_active BOOLEAN DEFAULT TRUE, 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 11. Применённые правила
CREATE TABLE IF NOT EXISTS applied_pricing_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ticket_id UUID NOT NULL REFERENCES tickets(id) ON DELETE CASCADE,
    pricing_rule_id UUID NOT NULL REFERENCES pricing_rules(id) ON DELETE RESTRICT,
    applied_price_impact DECIMAL(10,2) NOT NULL,
    UNIQUE(ticket_id, pricing_rule_id)
);

-- Создаём индексы для производительности
CREATE INDEX IF NOT EXISTS idx_tickets_user_id ON tickets(user_id);
CREATE INDEX IF NOT EXISTS idx_tickets_match_id ON tickets(match_id);
CREATE INDEX IF NOT EXISTS idx_tickets_status ON tickets(status);
CREATE INDEX IF NOT EXISTS idx_matches_date ON matches(match_date);
CREATE INDEX IF NOT EXISTS idx_matches_status ON matches(status);
CREATE INDEX IF NOT EXISTS idx_player_match_stats_match ON player_match_stats(match_id);
CREATE INDEX IF NOT EXISTS idx_pricing_rules_active ON pricing_rules(is_active) WHERE is_active = true;

--функция дл яавтоматического обновления updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$ --возвращает триггер в кавычках " "
BEGIN 
  NEW.updated_at = CURRENT_DATE;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS update_users_updated_at ON users;
CREATE TRIGGER update_users_updated_at
  BEFORE UPDATE ON users -- при запросе пользователя UPDATE таблицы users срабатывает данный триггер
  FOR EACH ROW --для каждой строки в таблице выполнится
  EXECUTE FUNCTION update_updated_at_column(); --указание выполняемой функции

