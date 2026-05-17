-- =====================================================
-- Откат данных (удаление добавленных матчей и билетов)
-- =====================================================

-- Удаляем билеты, созданные за период октябрь 2025 - май 2026
DELETE FROM tickets 
WHERE match_id IN (
    SELECT id FROM matches 
    WHERE match_date BETWEEN '2025-10-01' AND '2026-05-31'
);

-- Удаляем матчи, созданные за период октябрь 2025 - май 2026
DELETE FROM matches 
WHERE match_date BETWEEN '2025-10-01' AND '2026-05-31';

-- =====================================================
-- Если нужно также удалить добавленных пользователей (осторожно!)
-- =====================================================
-- DELETE FROM users WHERE email IN ('test_user_1@example.com', 'test_user_2@example.com');
-- DELETE FROM users WHERE created_at > '2025-10-01';

-- =====================================================
-- Проверка после отката
-- =====================================================
-- SELECT COUNT(*) as remaining_matches FROM matches;
-- SELECT COUNT(*) as remaining_tickets FROM tickets;