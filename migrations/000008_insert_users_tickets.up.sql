-- =====================================================
-- Добавление новых пользователей
-- =====================================================

DO $$
DECLARE
    new_user_ids UUID[] := ARRAY[]::UUID[];
    user_email VARCHAR;
    user_phone VARCHAR;
    user_name VARCHAR;
    user_id UUID;
    i INTEGER;
    user_record RECORD;
BEGIN
    -- Список новых пользователей
    FOR i IN 1..15 LOOP
        user_email := CONCAT('fan', i, '@example.com');
        user_phone := CONCAT('+7999', LPAD(i::TEXT, 7, '0'));
        user_name := CONCAT('Болельщик ', i);
        
        -- Проверяем, существует ли пользователь с таким email
        IF NOT EXISTS (SELECT 1 FROM users WHERE email = user_email) THEN
            INSERT INTO users (id, phone, email, full_name, password_hash, role, created_at, updated_at)
            VALUES (
                gen_random_uuid(),
                user_phone,
                user_email,
                user_name,
                'hashed_password_123',
                'user',
                NOW() - (random() * INTERVAL '90 days'),
                NOW()
            )
            RETURNING id INTO user_id;
            
            new_user_ids := array_append(new_user_ids, user_id);
            RAISE NOTICE 'Добавлен пользователь: % (%)', user_name, user_email;
        ELSE
            -- Если пользователь существует, добавляем его ID в список
            SELECT id INTO user_id FROM users WHERE email = user_email;
            new_user_ids := array_append(new_user_ids, user_id);
        END IF;
    END LOOP;
    
    -- Также добавим существующих пользователей, если их ещё нет в списке
    FOR user_record IN SELECT id FROM users WHERE email IN (
        'user@example.com',
        'admin@example.com',
        'tanya@mail.com',
        'lalala1@example.com',
        'ala@example.com'
    ) LOOP
        IF NOT (user_record.id = ANY(new_user_ids)) THEN
            new_user_ids := array_append(new_user_ids, user_record.id);
        END IF;
    END LOOP;
    
    RAISE NOTICE 'Всего пользователей для добавления билетов: %', array_length(new_user_ids, 1);
END $$;


-- =====================================================
-- Добавление билетов для всех пользователей на матчи
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
    user_count INTEGER := 0;
    user_id_val UUID;
BEGIN
    -- Для каждого пользователя
    FOR user_record IN 
        SELECT id FROM users 
        WHERE email LIKE '%@example.com' 
           OR email IN ('user@example.com', 'admin@example.com', 'tanya@mail.com', 'lalala1@example.com', 'ala@example.com')
    LOOP
        user_id_val := user_record.id;
        user_count := user_count + 1;
        
        -- Для каждого матча за последние 8 месяцев
        FOR match_record IN 
            SELECT id, match_date FROM matches 
            WHERE match_date BETWEEN '2025-10-01' AND '2026-05-31'
            ORDER BY match_date
        LOOP
            -- Разное количество билетов на матч (от 1 до 8)
            tickets_per_match := 1 + floor(random() * 8);
            
            FOR j IN 1..tickets_per_match LOOP
                -- Выбираем случайный сектор
                SELECT * INTO sector_record FROM stadium_sectors 
                ORDER BY random() LIMIT 1;
                
                -- Выбираем случайное место в секторе
                SELECT * INTO seat_record FROM seats 
                WHERE sector_id = sector_record.id 
                ORDER BY random() LIMIT 1;
                
                -- Рассчитываем цену
                price := 1000 * sector_record.price_coefficient * (0.7 + (random() * 0.6));
                price := ROUND(price / 50) * 50;
                
                -- Статус в зависимости от даты матча
                IF match_record.match_date < NOW() THEN
                    -- Прошедшие матчи
                    IF random() < 0.65 THEN
                        status_val := 'used';
                    ELSIF random() < 0.85 THEN
                        status_val := 'active';
                    ELSE
                        status_val := 'cancelled';
                    END IF;
                ELSE
                    -- Будущие матчи
                    IF random() < 0.75 THEN
                        status_val := 'active';
                    ELSE
                        status_val := 'cancelled';
                    END IF;
                END IF;
                
                -- Дата покупки: от 45 дней до матча до дня матча
                purchase_date_val := match_record.match_date - ((1 + floor(random() * 45)) * INTERVAL '1 day');
                
                -- Вставляем билет
                INSERT INTO tickets (id, user_id, match_id, seat_id, final_price, qr_code_hash, status, purchase_date)
                VALUES (
                    gen_random_uuid(),
                    user_id_val,
                    match_record.id,
                    seat_record.id,
                    price,
                    md5(random()::text || clock_timestamp()::text || user_id_val::text),
                    status_val,
                    purchase_date_val
                );
                
                tickets_created := tickets_created + 1;
            END LOOP;
        END LOOP;
        
        IF user_count % 5 = 0 THEN
            RAISE NOTICE 'Обработано % пользователей, создано % билетов', user_count, tickets_created;
        END IF;
    END LOOP;
    
    RAISE NOTICE '✅ ВСЕГО СОЗДАНО БИЛЕТОВ: % для % пользователей', tickets_created, user_count;
END $$;


-- =====================================================
-- Проверка топ-10 покупателей
-- =====================================================

-- Топ-10 покупателей по количеству билетов
SELECT 
    ROW_NUMBER() OVER (ORDER BY COUNT(t.id) DESC) as rank,
    u.full_name,
    u.email,
    COUNT(t.id) as tickets_count,
    COALESCE(SUM(t.final_price), 0) as total_spent,
    ROUND(AVG(t.final_price), 0) as avg_price,
    COUNT(DISTINCT t.match_id) as unique_matches
FROM users u
JOIN tickets t ON t.user_id = u.id
WHERE t.status != 'cancelled'
GROUP BY u.id, u.full_name, u.email
ORDER BY tickets_count DESC
LIMIT 10;