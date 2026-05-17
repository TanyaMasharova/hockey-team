-- =====================================================
-- Добавление матчей за период Октябрь 2025 - Май 2026
-- =====================================================

DO $$
DECLARE
    opponent_record RECORD;
    match_date TIMESTAMP;
    is_home BOOLEAN;
    our_score_val SMALLINT;
    opponent_score_val SMALLINT;
    status_val VARCHAR(20);
    is_derby_val BOOLEAN;
    base_date TIMESTAMP;
BEGIN
    -- Октябрь 2025 (домашние матчи)
    FOR opponent_record IN SELECT id, name FROM opponents WHERE name IN ('ЦСКА', 'Спартак', 'Динамо') LOOP
        base_date := '2025-10-04 19:00:00+03'::TIMESTAMP;
        match_date := base_date + (random() * 7 * INTERVAL '1 day');
        our_score_val := floor(random() * 5)::SMALLINT;
        opponent_score_val := floor(random() * 4)::SMALLINT;
        status_val := 'finished';
        is_derby_val := (opponent_record.name IN ('Спартак', 'ЦСКА', 'Динамо'));
        
        INSERT INTO matches (id, opponent_id, match_date, home_away, our_score, opponent_score, season, status, is_derby)
        VALUES (gen_random_uuid(), opponent_record.id, match_date, 'home', our_score_val, opponent_score_val, '2025/26', status_val, is_derby_val);
    END LOOP;

    -- Ноябрь 2025 (выездные матчи)
    FOR opponent_record IN SELECT id, name FROM opponents WHERE name IN ('СКА', 'Локомотив', 'Ак Барс') LOOP
        base_date := '2025-11-02 17:00:00+03'::TIMESTAMP;
        match_date := base_date + (random() * 14 * INTERVAL '1 day');
        our_score_val := floor(random() * 4)::SMALLINT;
        opponent_score_val := floor(random() * 5)::SMALLINT;
        status_val := 'finished';
        is_derby_val := FALSE;
        
        INSERT INTO matches (id, opponent_id, match_date, home_away, our_score, opponent_score, season, status, is_derby)
        VALUES (gen_random_uuid(), opponent_record.id, match_date, 'away', our_score_val, opponent_score_val, '2025/26', status_val, is_derby_val);
    END LOOP;

    -- Декабрь 2025
    FOR opponent_record IN SELECT id, name FROM opponents WHERE name IN ('Авангард', 'Металлург', 'Салават Юлаев') LOOP
        base_date := '2025-12-05 19:00:00+03'::TIMESTAMP;
        match_date := base_date + (random() * 10 * INTERVAL '1 day');
        is_home := (random() > 0.5);
        our_score_val := floor(random() * 5)::SMALLINT;
        opponent_score_val := floor(random() * 5)::SMALLINT;
        status_val := 'finished';
        is_derby_val := FALSE;
        
        INSERT INTO matches (id, opponent_id, match_date, home_away, our_score, opponent_score, season, status, is_derby)
        VALUES (gen_random_uuid(), opponent_record.id, match_date, CASE WHEN is_home THEN 'home' ELSE 'away' END, 
                our_score_val, opponent_score_val, '2025/26', status_val, is_derby_val);
    END LOOP;

    -- Январь 2026 (перерыв, мало матчей)
    FOR opponent_record IN SELECT id, name FROM opponents WHERE name IN ('Торпедо', 'Северсталь') LOOP
        base_date := '2026-01-15 19:00:00+03'::TIMESTAMP;
        match_date := base_date + (random() * 7 * INTERVAL '1 day');
        our_score_val := floor(random() * 5)::SMALLINT;
        opponent_score_val := floor(random() * 4)::SMALLINT;
        status_val := 'finished';
        is_derby_val := FALSE;
        
        INSERT INTO matches (id, opponent_id, match_date, home_away, our_score, opponent_score, season, status, is_derby)
        VALUES (gen_random_uuid(), opponent_record.id, match_date, 'home', our_score_val, opponent_score_val, '2025/26', status_val, is_derby_val);
    END LOOP;

    -- Февраль 2026
    FOR opponent_record IN SELECT id, name FROM opponents WHERE name IN ('Трактор', 'Автомобилист', 'Барыс') LOOP
        base_date := '2026-02-07 17:00:00+03'::TIMESTAMP;
        match_date := base_date + (random() * 14 * INTERVAL '1 day');
        is_home := (random() > 0.5);
        our_score_val := floor(random() * 5)::SMALLINT;
        opponent_score_val := floor(random() * 5)::SMALLINT;
        status_val := 'finished';
        is_derby_val := FALSE;
        
        INSERT INTO matches (id, opponent_id, match_date, home_away, our_score, opponent_score, season, status, is_derby)
        VALUES (gen_random_uuid(), opponent_record.id, match_date, CASE WHEN is_home THEN 'home' ELSE 'away' END,
                our_score_val, opponent_score_val, '2025/26', status_val, is_derby_val);
    END LOOP;

    -- Март 2026 (включая дерби)
    FOR opponent_record IN SELECT id, name FROM opponents WHERE name IN ('Спартак', 'ЦСКА', 'Динамо', 'Локомотив') LOOP
        base_date := '2026-03-05 19:00:00+03'::TIMESTAMP;
        match_date := base_date + (random() * 20 * INTERVAL '1 day');
        is_home := (random() > 0.5);
        our_score_val := floor(random() * 5)::SMALLINT;
        opponent_score_val := floor(random() * 5)::SMALLINT;
        status_val := CASE WHEN match_date < NOW() THEN 'finished' ELSE 'scheduled' END;
        is_derby_val := (opponent_record.name IN ('Спартак', 'ЦСКА', 'Динамо'));
        
        INSERT INTO matches (id, opponent_id, match_date, home_away, our_score, opponent_score, season, status, is_derby)
        VALUES (gen_random_uuid(), opponent_record.id, match_date, CASE WHEN is_home THEN 'home' ELSE 'away' END,
                our_score_val, opponent_score_val, '2025/26', status_val, is_derby_val);
    END LOOP;

    -- Апрель 2026 (плей-офф начало)
    FOR opponent_record IN SELECT id, name FROM opponents WHERE name IN ('СКА', 'Авангард', 'Ак Барс', 'Металлург') LOOP
        base_date := '2026-04-02 19:00:00+03'::TIMESTAMP;
        match_date := base_date + (random() * 21 * INTERVAL '1 day');
        is_home := (random() > 0.5);
        our_score_val := floor(random() * 4)::SMALLINT;
        opponent_score_val := floor(random() * 4)::SMALLINT;
        status_val := CASE WHEN match_date < NOW() THEN 'finished' ELSE 'scheduled' END;
        is_derby_val := FALSE;
        
        INSERT INTO matches (id, opponent_id, match_date, home_away, our_score, opponent_score, season, status, is_derby)
        VALUES (gen_random_uuid(), opponent_record.id, match_date, CASE WHEN is_home THEN 'home' ELSE 'away' END,
                our_score_val, opponent_score_val, '2025/26', status_val, is_derby_val);
    END LOOP;

    -- Май 2026 (финалы, будущие матчи)
    FOR opponent_record IN SELECT id, name FROM opponents WHERE name IN ('СКА', 'Ак Барс', 'ЦСКА', 'Локомотив') LOOP
        base_date := '2026-05-10 17:00:00+03'::TIMESTAMP;
        match_date := base_date + (random() * 14 * INTERVAL '1 day');
        is_home := (random() > 0.5);
        our_score_val := 0;
        opponent_score_val := 0;
        status_val := 'scheduled';
        is_derby_val := (opponent_record.name IN ('Спартак', 'ЦСКА', 'Динамо'));
        
        INSERT INTO matches (id, opponent_id, match_date, home_away, our_score, opponent_score, season, status, is_derby)
        VALUES (gen_random_uuid(), opponent_record.id, match_date, CASE WHEN is_home THEN 'home' ELSE 'away' END,
                our_score_val, opponent_score_val, '2025/26', status_val, is_derby_val);
    END LOOP;
END $$;


-- =====================================================
-- Создание билетов для всех матчей
-- =====================================================

DO $$
DECLARE
    user_record RECORD;
    match_record RECORD;
    sector_record RECORD;
    seat_record RECORD;
    tickets_per_match INTEGER;
    price DECIMAL(10,2);
    status_val VARCHAR(20);
    purchase_date_val TIMESTAMP;
    tickets_created INTEGER := 0;
    j INTEGER;
BEGIN
    -- Для каждого пользователя
    FOR user_record IN SELECT id FROM users WHERE id IN (
        '2ea67113-3d8b-45f7-96d9-b3315ba7d4db'::UUID,
        '2617e402-da9d-455e-8570-6524f36b7dc2'::UUID,
        'e6278c67-8e5d-4b30-b5b2-c39c3bc6c4b3'::UUID,
        'f22a2cdb-fa3c-4f6d-bfde-50f92e647a79'::UUID
    ) LOOP
        -- Для каждого матча
        FOR match_record IN SELECT id, match_date FROM matches WHERE match_date BETWEEN '2025-10-01' AND '2026-05-31' LOOP
            -- Разное количество билетов на матч (от 1 до 15)
            tickets_per_match := 1 + floor(random() * 15);
            
            FOR j IN 1..tickets_per_match LOOP
                -- Выбираем случайный сектор
                SELECT * INTO sector_record FROM stadium_sectors 
                ORDER BY random() LIMIT 1;
                
                -- Выбираем случайное место в секторе
                SELECT * INTO seat_record FROM seats 
                WHERE sector_id = sector_record.id 
                ORDER BY random() LIMIT 1;
                
                -- Рассчитываем цену
                price := 1000 * sector_record.price_coefficient * (0.8 + (random() * 0.4));
                price := ROUND(price / 50) * 50;
                
                -- Статус в зависимости от даты матча
                IF match_record.match_date < NOW() THEN
                    -- Прошедшие матчи: 70% used, 20% active, 10% cancelled
                    IF random() < 0.7 THEN
                        status_val := 'used';
                    ELSIF random() < 0.9 THEN
                        status_val := 'active';
                    ELSE
                        status_val := 'cancelled';
                    END IF;
                ELSE
                    -- Будущие матчи: 80% active, 20% cancelled
                    IF random() < 0.8 THEN
                        status_val := 'active';
                    ELSE
                        status_val := 'cancelled';
                    END IF;
                END IF;
                
                -- Дата покупки: от 30 дней до матча до дня матча
                purchase_date_val := match_record.match_date - ((1 + floor(random() * 30)) * INTERVAL '1 day');
                
                -- Вставляем билет
                INSERT INTO tickets (id, user_id, match_id, seat_id, final_price, qr_code_hash, status, purchase_date)
                VALUES (
                    gen_random_uuid(),
                    user_record.id,
                    match_record.id,
                    seat_record.id,
                    price,
                    md5(random()::text || clock_timestamp()::text),
                    status_val,
                    purchase_date_val
                );
                
                tickets_created := tickets_created + 1;
            END LOOP;
        END LOOP;
    END LOOP;
    
    RAISE NOTICE 'Создано % билетов', tickets_created;
END $$;


-- =====================================================
-- Проверка созданных данных
-- =====================================================

-- Матчи по месяцам
SELECT 
    TO_CHAR(match_date, 'YYYY-MM') as month,
    COUNT(*) as matches_count,
    COUNT(CASE WHEN status = 'finished' THEN 1 END) as finished,
    COUNT(CASE WHEN status = 'scheduled' THEN 1 END) as scheduled
FROM matches
WHERE match_date >= '2025-10-01'
GROUP BY TO_CHAR(match_date, 'YYYY-MM')
ORDER BY month;

-- Билеты по месяцам
SELECT 
    TO_CHAR(m.match_date, 'YYYY-MM') as month,
    COUNT(t.id) as tickets_sold,
    SUM(t.final_price) as revenue,
    COUNT(DISTINCT t.user_id) as unique_fans
FROM tickets t
JOIN matches m ON t.match_id = m.id
WHERE t.status != 'cancelled'
GROUP BY TO_CHAR(m.match_date, 'YYYY-MM')
ORDER BY month;

-- Общая статистика
SELECT 
    COUNT(DISTINCT t.id) as total_tickets,
    COUNT(DISTINCT t.user_id) as active_fans,
    SUM(t.final_price) as total_revenue,
    ROUND(AVG(t.final_price), 0) as avg_ticket_price,
    COUNT(DISTINCT m.id) as matches_played
FROM tickets t
JOIN matches m ON t.match_id = m.id
WHERE t.status != 'cancelled';