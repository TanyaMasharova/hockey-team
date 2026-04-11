-- Откат миграции 000001_init_schema

-- Удаляем триггер
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

-- Удаляем функцию
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Удаляем индексы
DROP INDEX IF EXISTS idx_tickets_user_id;
DROP INDEX IF EXISTS idx_tickets_match_id;
DROP INDEX IF EXISTS idx_tickets_status;
DROP INDEX IF EXISTS idx_matches_date;
DROP INDEX IF EXISTS idx_matches_status;
DROP INDEX IF EXISTS idx_player_match_stats_match;
DROP INDEX IF EXISTS idx_pricing_rules_active;

-- Удаляем таблицы (в порядке обратном созданию из-за зависимостей)
DROP TABLE IF EXISTS applied_pricing_rules;
DROP TABLE IF EXISTS pricing_rules;
DROP TABLE IF EXISTS player_season_stats;
DROP TABLE IF EXISTS player_match_stats;
DROP TABLE IF EXISTS players;
DROP TABLE IF EXISTS tickets;
DROP TABLE IF EXISTS seats;
DROP TABLE IF EXISTS stadium_sectors;
DROP TABLE IF EXISTS matches;
DROP TABLE IF EXISTS opponents;
DROP TABLE IF EXISTS users;

-- Отключаем расширение (опционально, если хотите полностью очистить)
-- DROP EXTENSION IF EXISTS "uuid-ossp";