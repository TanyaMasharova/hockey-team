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
-- 4. Билеты создаются позже, в `000008_insert_users_tickets.up.sql`
-- =====================================================
