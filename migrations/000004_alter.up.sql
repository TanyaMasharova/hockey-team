-- =====================================================
-- Добавление колонки win_type (если её нет)
-- =====================================================

DO $$ 
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'matches' AND column_name = 'win_type'
    ) THEN
        ALTER TABLE matches 
        ADD COLUMN win_type VARCHAR(20) DEFAULT 'regular' 
        CHECK (win_type IN ('regular', 'overtime', 'penalty'));
    END IF;
END $$;

-- Создаём индексы (IF NOT EXISTS уже есть)
CREATE INDEX IF NOT EXISTS idx_matches_win_type ON matches(win_type);
CREATE INDEX IF NOT EXISTS idx_matches_status_win_type ON matches(status, win_type);

-- =====================================================
-- Обновление win_type для существующих матчей
-- =====================================================

-- Обновляем win_type только для тех матчей, где он ещё NULL
UPDATE matches SET win_type = 'regular' 
WHERE win_type IS NULL 
AND id IN (
    SELECT m.id FROM matches m
    JOIN opponents o ON m.opponent_id = o.id
    WHERE o.name = 'Локомотив' AND m.match_date = '2025-03-10 19:00:00+03'
);

UPDATE matches SET win_type = 'regular' 
WHERE win_type IS NULL 
AND id IN (
    SELECT m.id FROM matches m
    JOIN opponents o ON m.opponent_id = o.id
    WHERE o.name = 'Спартак' AND m.match_date = '2025-03-05 19:30:00+03'
);

UPDATE matches SET win_type = 'overtime' 
WHERE win_type IS NULL 
AND id IN (
    SELECT m.id FROM matches m
    JOIN opponents o ON m.opponent_id = o.id
    WHERE o.name = 'Динамо' AND m.match_date = '2025-02-28 19:00:00+03'
);

UPDATE matches SET win_type = 'regular' 
WHERE win_type IS NULL 
AND id IN (
    SELECT m.id FROM matches m
    JOIN opponents o ON m.opponent_id = o.id
    WHERE o.name = 'Металлург' AND m.match_date = '2025-02-25 17:00:00+03'
);

-- Для всех остальных матчей со статусом 'finished' устанавливаем 'regular'
UPDATE matches SET win_type = 'regular' 
WHERE win_type IS NULL AND status = 'finished';

-- Для предстоящих матчей устанавливаем DEFAULT (regular)
UPDATE matches SET win_type = DEFAULT 
WHERE win_type IS NULL AND status = 'scheduled';