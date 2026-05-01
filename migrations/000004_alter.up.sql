-- Добавляем столбец win_type
ALTER TABLE matches 
ADD COLUMN win_type VARCHAR(20) DEFAULT 'regular' 
CHECK (win_type IN ('regular', 'overtime', 'penalty'));

CREATE INDEX IF NOT EXISTS idx_matches_win_type ON matches(win_type);
CREATE INDEX IF NOT EXISTS idx_matches_status_win_type ON matches(status, win_type);

-- Полностью перезаписываем win_type для всех матчей

-- Удаляем существующие значения (если есть)
UPDATE matches SET win_type = NULL;

-- Устанавливаем правильные win_type для каждого матча
UPDATE matches SET win_type = 'regular' WHERE id IN (
    SELECT m.id FROM matches m
    JOIN opponents o ON m.opponent_id = o.id
    WHERE o.name = 'Локомотив' AND m.match_date = '2025-03-10 19:00:00+03'
);

UPDATE matches SET win_type = 'regular' WHERE id IN (
    SELECT m.id FROM matches m
    JOIN opponents o ON m.opponent_id = o.id
    WHERE o.name = 'Спартак' AND m.match_date = '2025-03-05 19:30:00+03'
);

UPDATE matches SET win_type = 'overtime' WHERE id IN (
    SELECT m.id FROM matches m
    JOIN opponents o ON m.opponent_id = o.id
    WHERE o.name = 'Динамо' AND m.match_date = '2025-02-28 19:00:00+03'
);

UPDATE matches SET win_type = 'regular' WHERE id IN (
    SELECT m.id FROM matches m
    JOIN opponents o ON m.opponent_id = o.id
    WHERE o.name = 'Металлург' AND m.match_date = '2025-02-25 17:00:00+03'
);

-- Все предстоящие матчи будут 'regular' (DEFAULT)
UPDATE matches SET win_type = DEFAULT WHERE status = 'scheduled';
