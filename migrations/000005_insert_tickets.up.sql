-- =====================================================
-- 2. Создание секторов стадиона (12 секторов)
-- =====================================================
INSERT INTO stadium_sectors (id, sector_number, capacity, sector_type, price_coefficient, color_code) VALUES
-- Нижний ярус (VIP и стандартные)
(gen_random_uuid(), 'A1', 250, 'vip', 2.0, '#FFD700'),
(gen_random_uuid(), 'A2', 250, 'vip', 2.0, '#FFD700'),
(gen_random_uuid(), 'B1', 250, 'standard', 1.2, '#4CAF50'),
(gen_random_uuid(), 'B2', 250, 'standard', 1.2, '#4CAF50'),
(gen_random_uuid(), 'B3', 250, 'standard', 1.2, '#4CAF50'),
(gen_random_uuid(), 'C1', 250, 'standard', 1.0, '#2196F3'),
(gen_random_uuid(), 'C2', 250, 'standard', 1.0, '#2196F3'),
(gen_random_uuid(), 'C3', 250, 'standard', 1.0, '#2196F3'),
(gen_random_uuid(), 'C4', 250, 'standard', 1.0, '#2196F3'),
(gen_random_uuid(), 'D1', 250, 'standard', 0.8, '#9C27B0'),
(gen_random_uuid(), 'D2', 250, 'standard', 0.8, '#9C27B0'),
(gen_random_uuid(), 'E1', 250, 'away_fans', 0.7, '#F44336');

-- =====================================================
-- 3. Создание мест в секторах (12 секторов × 250 мест = 3000 мест)
-- =====================================================
DO $$
DECLARE
    sector_record RECORD;
    row_num INTEGER;
    seat_num INTEGER;
    row_letter VARCHAR(10);
    sector_id UUID;
    is_handicap BOOLEAN;
BEGIN
    FOR sector_record IN SELECT id, sector_number FROM stadium_sectors LOOP
        sector_id := sector_record.id;
        
        -- Для каждого сектора создаём 250 мест (ряды A-J по 25 мест в каждом ряду)
        FOR row_num IN 1..10 LOOP
            -- Преобразуем номер ряда в букву (1=A, 2=B, 3=C, ...)
            row_letter := CHR(64 + row_num);
            
            FOR seat_num IN 1..25 LOOP
                -- Каждое 50-е место для инвалидов (5-10-15-20-25 ряды)
                is_handicap := (seat_num % 5 = 0);
                
                INSERT INTO seats (id, sector_id, seat_row, seat_number, is_handicap_accessible)
                VALUES (gen_random_uuid(), sector_id, row_letter, seat_num::VARCHAR, is_handicap);
            END LOOP;
        END LOOP;
    END LOOP;
END $$;

-- =====================================================
-- 4. Создание билетов (20 билетов на разные матчи)
-- =====================================================

-- Сначала получим ID существующих матчей и секторов для билетов
DO $$
DECLARE
    user_ids UUID[] := ARRAY[
        '2ea67113-3d8b-45f7-96d9-b3315ba7d4db'::UUID,
        '2617e402-da9d-455e-8570-6524f36b7dc2'::UUID,
        'e6278c67-8e5d-4b30-b5b2-c39c3bc6c4b3'::UUID,
        'f22a2cdb-fa3c-4f6d-bfde-50f92e647a79'::UUID
    ];
    match_ids UUID[] := ARRAY[
        '39628289-6b6b-4ab6-8287-06bed3321692'::UUID,
        'e9f1760c-9960-4cd5-9c7c-a1626d5e14c1'::UUID,
        'dc5c32fd-27c9-4450-a360-185740ef7bda'::UUID,
        'fe9e7a31-40a2-4668-87eb-ad436e6d7ebd'::UUID,
        '4929751b-09d9-4a5c-a7a5-2266b05a08b3'::UUID,
        '28d81847-d35b-4d97-81f4-38401aed28a4'::UUID,
        '9f1b66c8-14da-487d-a380-dcbb044ba8b6'::UUID,
        '12a188dd-a37a-47aa-8c85-a1d8c848de0d'::UUID,
        'b496dbe8-e015-4375-8dbe-97b35a116c29'::UUID,
        'cbf739c2-526b-4766-a891-67e3493a3c54'::UUID
    ];
    sector_rec RECORD;
    seat_rec RECORD;
    i INTEGER;
    price DECIMAL(10,2);
    qr_hash VARCHAR(255);
BEGIN
    -- Билет 1: Пользователь 1, Матч 1, Сектор VIP
    SELECT id INTO sector_rec FROM stadium_sectors WHERE sector_number = 'A1' LIMIT 1;
    SELECT id INTO seat_rec FROM seats WHERE sector_id = sector_rec.id AND seat_row = 'A' AND seat_number = '1' LIMIT 1;
    INSERT INTO tickets (id, user_id, match_id, seat_id, final_price, qr_code_hash, status, purchase_date)
    VALUES (gen_random_uuid(), user_ids[1], match_ids[1], seat_rec.id, 5000.00, md5(random()::text), 'active', NOW() - INTERVAL '10 days');
    
    -- Билет 2: Пользователь 1, Матч 2, Сектор B1
    SELECT id INTO sector_rec FROM stadium_sectors WHERE sector_number = 'B1' LIMIT 1;
    SELECT id INTO seat_rec FROM seats WHERE sector_id = sector_rec.id AND seat_row = 'B' AND seat_number = '5' LIMIT 1;
    INSERT INTO tickets (id, user_id, match_id, seat_id, final_price, qr_code_hash, status, purchase_date)
    VALUES (gen_random_uuid(), user_ids[1], match_ids[2], seat_rec.id, 3000.00, md5(random()::text), 'used', NOW() - INTERVAL '5 days');
    
    -- Билет 3: Пользователь 2, Матч 3, Сектор C2
    SELECT id INTO sector_rec FROM stadium_sectors WHERE sector_number = 'C2' LIMIT 1;
    SELECT id INTO seat_rec FROM seats WHERE sector_id = sector_rec.id AND seat_row = 'C' AND seat_number = '10' LIMIT 1;
    INSERT INTO tickets (id, user_id, match_id, seat_id, final_price, qr_code_hash, status, purchase_date)
    VALUES (gen_random_uuid(), user_ids[2], match_ids[3], seat_rec.id, 2500.00, md5(random()::text), 'active', NOW() - INTERVAL '3 days');
    
    -- Билет 4: Пользователь 2, Матч 4, Сектор D1
    SELECT id INTO sector_rec FROM stadium_sectors WHERE sector_number = 'D1' LIMIT 1;
    SELECT id INTO seat_rec FROM seats WHERE sector_id = sector_rec.id AND seat_row = 'D' AND seat_number = '15' LIMIT 1;
    INSERT INTO tickets (id, user_id, match_id, seat_id, final_price, qr_code_hash, status, purchase_date)
    VALUES (gen_random_uuid(), user_ids[2], match_ids[4], seat_rec.id, 2000.00, md5(random()::text), 'cancelled', NOW() - INTERVAL '7 days');
    
    -- Билет 5: Пользователь 3, Матч 5 (будущий матч), Сектор A2
    SELECT id INTO sector_rec FROM stadium_sectors WHERE sector_number = 'A2' LIMIT 1;
    SELECT id INTO seat_rec FROM seats WHERE sector_id = sector_rec.id AND seat_row = 'A' AND seat_number = '20' LIMIT 1;
    INSERT INTO tickets (id, user_id, match_id, seat_id, final_price, qr_code_hash, status, purchase_date)
    VALUES (gen_random_uuid(), user_ids[3], match_ids[5], seat_rec.id, 4800.00, md5(random()::text), 'active', NOW() - INTERVAL '1 day');
    
    -- Билет 6: Пользователь 3, Матч 6, Сектор B2
    SELECT id INTO sector_rec FROM stadium_sectors WHERE sector_number = 'B2' LIMIT 1;
    SELECT id INTO seat_rec FROM seats WHERE sector_id = sector_rec.id AND seat_row = 'E' AND seat_number = '8' LIMIT 1;
    INSERT INTO tickets (id, user_id, match_id, seat_id, final_price, qr_code_hash, status, purchase_date)
    VALUES (gen_random_uuid(), user_ids[3], match_ids[6], seat_rec.id, 3200.00, md5(random()::text), 'active', NOW() - INTERVAL '2 days');
    
    -- Билет 7: Пользователь 4, Матч 7, Сектор E1 (гостевой)
    SELECT id INTO sector_rec FROM stadium_sectors WHERE sector_number = 'E1' LIMIT 1;
    SELECT id INTO seat_rec FROM seats WHERE sector_id = sector_rec.id AND seat_row = 'A' AND seat_number = '3' LIMIT 1;
    INSERT INTO tickets (id, user_id, match_id, seat_id, final_price, qr_code_hash, status, purchase_date)
    VALUES (gen_random_uuid(), user_ids[4], match_ids[7], seat_rec.id, 1500.00, md5(random()::text), 'active', NOW() - INTERVAL '4 days');
    
    -- Билет 8: Пользователь 4, Матч 8, Сектор C3
    SELECT id INTO sector_rec FROM stadium_sectors WHERE sector_number = 'C3' LIMIT 1;
    SELECT id INTO seat_rec FROM seats WHERE sector_id = sector_rec.id AND seat_row = 'F' AND seat_number = '12' LIMIT 1;
    INSERT INTO tickets (id, user_id, match_id, seat_id, final_price, qr_code_hash, status, purchase_date)
    VALUES (gen_random_uuid(), user_ids[4], match_ids[8], seat_rec.id, 2600.00, md5(random()::text), 'refunded', NOW() - INTERVAL '6 days');
    
    -- Билет 9: Пользователь 1, Матч 9, Сектор C1
    SELECT id INTO sector_rec FROM stadium_sectors WHERE sector_number = 'C1' LIMIT 1;
    SELECT id INTO seat_rec FROM seats WHERE sector_id = sector_rec.id AND seat_row = 'G' AND seat_number = '22' LIMIT 1;
    INSERT INTO tickets (id, user_id, match_id, seat_id, final_price, qr_code_hash, status, purchase_date)
    VALUES (gen_random_uuid(), user_ids[1], match_ids[9], seat_rec.id, 2400.00, md5(random()::text), 'used', NOW() - INTERVAL '8 days');
    
    -- Билет 10: Пользователь 2, Матч 10, Сектор D2
    SELECT id INTO sector_rec FROM stadium_sectors WHERE sector_number = 'D2' LIMIT 1;
    SELECT id INTO seat_rec FROM seats WHERE sector_id = sector_rec.id AND seat_row = 'H' AND seat_number = '18' LIMIT 1;
    INSERT INTO tickets (id, user_id, match_id, seat_id, final_price, qr_code_hash, status, purchase_date)
    VALUES (gen_random_uuid(), user_ids[2], match_ids[10], seat_rec.id, 1900.00, md5(random()::text), 'active', NOW() - INTERVAL '1 day');
    
    -- Билет 11: Пользователь 3, Матч 1, Сектор B3
    SELECT id INTO sector_rec FROM stadium_sectors WHERE sector_number = 'B3' LIMIT 1;
    SELECT id INTO seat_rec FROM seats WHERE sector_id = sector_rec.id AND seat_row = 'I' AND seat_number = '7' LIMIT 1;
    INSERT INTO tickets (id, user_id, match_id, seat_id, final_price, qr_code_hash, status, purchase_date)
    VALUES (gen_random_uuid(), user_ids[3], match_ids[1], seat_rec.id, 3100.00, md5(random()::text), 'active', NOW() - INTERVAL '12 days');
    
    -- Билет 12: Пользователь 4, Матч 2, Сектор C4
    SELECT id INTO sector_rec FROM stadium_sectors WHERE sector_number = 'C4' LIMIT 1;
    SELECT id INTO seat_rec FROM seats WHERE sector_id = sector_rec.id AND seat_row = 'J' AND seat_number = '25' LIMIT 1;
    INSERT INTO tickets (id, user_id, match_id, seat_id, final_price, qr_code_hash, status, purchase_date)
    VALUES (gen_random_uuid(), user_ids[4], match_ids[2], seat_rec.id, 2300.00, md5(random()::text), 'active', NOW() - INTERVAL '9 days');
    
    -- Билет 13: Пользователь 1, Матч 3, Сектор A1
    SELECT id INTO sector_rec FROM stadium_sectors WHERE sector_number = 'A1' LIMIT 1;
    SELECT id INTO seat_rec FROM seats WHERE sector_id = sector_rec.id AND seat_row = 'B' AND seat_number = '4' LIMIT 1;
    INSERT INTO tickets (id, user_id, match_id, seat_id, final_price, qr_code_hash, status, purchase_date)
    VALUES (gen_random_uuid(), user_ids[1], match_ids[3], seat_rec.id, 5200.00, md5(random()::text), 'cancelled', NOW() - INTERVAL '11 days');
    
    -- Билет 14: Пользователь 2, Матч 4, Сектор E1
    SELECT id INTO sector_rec FROM stadium_sectors WHERE sector_number = 'E1' LIMIT 1;
    SELECT id INTO seat_rec FROM seats WHERE sector_id = sector_rec.id AND seat_row = 'C' AND seat_number = '11' LIMIT 1;
    INSERT INTO tickets (id, user_id, match_id, seat_id, final_price, qr_code_hash, status, purchase_date)
    VALUES (gen_random_uuid(), user_ids[2], match_ids[4], seat_rec.id, 1600.00, md5(random()::text), 'active', NOW() - INTERVAL '2 days');
    
    -- Билет 15: Пользователь 3, Матч 5, Сектор B1
    SELECT id INTO sector_rec FROM stadium_sectors WHERE sector_number = 'B1' LIMIT 1;
    SELECT id INTO seat_rec FROM seats WHERE sector_id = sector_rec.id AND seat_row = 'D' AND seat_number = '19' LIMIT 1;
    INSERT INTO tickets (id, user_id, match_id, seat_id, final_price, qr_code_hash, status, purchase_date)
    VALUES (gen_random_uuid(), user_ids[3], match_ids[5], seat_rec.id, 3400.00, md5(random()::text), 'active', NOW() - INTERVAL '5 days');
    
    -- Билет 16: Пользователь 4, Матч 6, Сектор C2
    SELECT id INTO sector_rec FROM stadium_sectors WHERE sector_number = 'C2' LIMIT 1;
    SELECT id INTO seat_rec FROM seats WHERE sector_id = sector_rec.id AND seat_row = 'E' AND seat_number = '21' LIMIT 1;
    INSERT INTO tickets (id, user_id, match_id, seat_id, final_price, qr_code_hash, status, purchase_date)
    VALUES (gen_random_uuid(), user_ids[4], match_ids[6], seat_rec.id, 2550.00, md5(random()::text), 'refunded', NOW() - INTERVAL '3 days');
    
    -- Билет 17: Пользователь 1, Матч 7, Сектор D1
    SELECT id INTO sector_rec FROM stadium_sectors WHERE sector_number = 'D1' LIMIT 1;
    SELECT id INTO seat_rec FROM seats WHERE sector_id = sector_rec.id AND seat_row = 'F' AND seat_number = '14' LIMIT 1;
    INSERT INTO tickets (id, user_id, match_id, seat_id, final_price, qr_code_hash, status, purchase_date)
    VALUES (gen_random_uuid(), user_ids[1], match_ids[7], seat_rec.id, 2100.00, md5(random()::text), 'active', NOW() - INTERVAL '7 days');
    
    -- Билет 18: Пользователь 2, Матч 8, Сектор B2
    SELECT id INTO sector_rec FROM stadium_sectors WHERE sector_number = 'B2' LIMIT 1;
    SELECT id INTO seat_rec FROM seats WHERE sector_id = sector_rec.id AND seat_row = 'G' AND seat_number = '9' LIMIT 1;
    INSERT INTO tickets (id, user_id, match_id, seat_id, final_price, qr_code_hash, status, purchase_date)
    VALUES (gen_random_uuid(), user_ids[2], match_ids[8], seat_rec.id, 3300.00, md5(random()::text), 'used', NOW() - INTERVAL '15 days');
    
    -- Билет 19: Пользователь 3, Матч 9, Сектор C3
    SELECT id INTO sector_rec FROM stadium_sectors WHERE sector_number = 'C3' LIMIT 1;
    SELECT id INTO seat_rec FROM seats WHERE sector_id = sector_rec.id AND seat_row = 'H' AND seat_number = '6' LIMIT 1;
    INSERT INTO tickets (id, user_id, match_id, seat_id, final_price, qr_code_hash, status, purchase_date)
    VALUES (gen_random_uuid(), user_ids[3], match_ids[9], seat_rec.id, 2700.00, md5(random()::text), 'active', NOW() - INTERVAL '4 days');
    
    -- Билет 20: Пользователь 4, Матч 10, Сектор A2
    SELECT id INTO sector_rec FROM stadium_sectors WHERE sector_number = 'A2' LIMIT 1;
    SELECT id INTO seat_rec FROM seats WHERE sector_id = sector_rec.id AND seat_row = 'I' AND seat_number = '13' LIMIT 1;
    INSERT INTO tickets (id, user_id, match_id, seat_id, final_price, qr_code_hash, status, purchase_date)
    VALUES (gen_random_uuid(), user_ids[4], match_ids[10], seat_rec.id, 4900.00, md5(random()::text), 'active', NOW() - INTERVAL '0 days');
    
    RAISE NOTICE 'Создано 20 билетов';
END $$;

-- =====================================================
-- 5. Проверка созданных данных
-- =====================================================

-- Статистика по билетам
SELECT 
    COUNT(*) as total_tickets,
    COUNT(DISTINCT user_id) as unique_users,
    COUNT(DISTINCT match_id) as unique_matches,
    SUM(final_price) as total_revenue,
    AVG(final_price) as avg_ticket_price
FROM tickets;

-- Билеты по статусам
SELECT status, COUNT(*) as count 
FROM tickets 
GROUP BY status;

-- Билеты по секторам
SELECT 
    ss.sector_number,
    COUNT(t.id) as tickets_sold,
    SUM(t.final_price) as revenue
FROM tickets t
JOIN seats s ON t.seat_id = s.id
JOIN stadium_sectors ss ON s.sector_id = ss.id
GROUP BY ss.sector_number
ORDER BY ss.sector_number;