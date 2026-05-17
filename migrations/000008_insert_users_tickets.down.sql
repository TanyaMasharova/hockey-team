-- =====================================================
-- Откат: удаление добавленных пользователей и их билетов
-- =====================================================

DO $$
DECLARE
    deleted_tickets INTEGER;
    deleted_users INTEGER;
BEGIN
    -- Удаляем билеты добавленных пользователей
    WITH deleted AS (
        DELETE FROM tickets 
        WHERE user_id IN (SELECT id FROM users WHERE email LIKE '%@example.com')
        RETURNING id
    )
    SELECT COUNT(*) INTO deleted_tickets FROM deleted;
    
    RAISE NOTICE 'Удалено билетов: %', deleted_tickets;
    
    -- Удаляем добавленных пользователей (кроме существующих admin и user)
    WITH deleted AS (
        DELETE FROM users 
        WHERE email LIKE '%@example.com' 
        AND email NOT IN ('admin@example.com', 'user@example.com')
        RETURNING id
    )
    SELECT COUNT(*) INTO deleted_users FROM deleted;
    
    RAISE NOTICE 'Удалено пользователей: %', deleted_users;
    RAISE NOTICE 'Откат выполнен успешно!';
END $$;

-- Проверка после отката
SELECT 'Осталось пользователей:' as info, COUNT(*) as count FROM users
UNION ALL
SELECT 'Осталось билетов:', COUNT(*) FROM tickets;