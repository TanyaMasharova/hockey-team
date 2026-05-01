-- Удаляем индексы
DROP INDEX IF EXISTS idx_matches_win_type;
DROP INDEX IF EXISTS idx_matches_status_win_type;

-- Удаляем столбец
ALTER TABLE matches DROP COLUMN IF EXISTS win_type;