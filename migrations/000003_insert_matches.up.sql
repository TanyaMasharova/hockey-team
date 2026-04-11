
INSERT INTO matches (opponent_id, match_date, home_away, season, status, our_score, opponent_score) VALUES
    -- Прошедшие матчи (с результатами)
    ((SELECT id FROM opponents WHERE name = 'Локомотив'), '2025-03-10 19:00:00+03', 'away', '2024/25', 'finished', 3, 2),
    ((SELECT id FROM opponents WHERE name = 'Спартак'), '2025-03-05 19:30:00+03', 'away', '2024/25', 'finished', 2, 1),
    ((SELECT id FROM opponents WHERE name = 'Динамо'), '2025-02-28 19:00:00+03', 'home', '2024/25', 'finished', 4, 3),
    ((SELECT id FROM opponents WHERE name = 'Металлург'), '2025-02-25 17:00:00+03', 'away', '2024/25', 'finished', 1, 3),
    
    -- Предстоящие матчи
    ((SELECT id FROM opponents WHERE name = 'ЦСКА'), '2025-03-15 19:00:00+03', 'home', '2024/25', 'scheduled', 0, 0),
    ((SELECT id FROM opponents WHERE name = 'СКА'), '2025-03-20 19:30:00+03', 'home', '2024/25', 'scheduled', 0, 0),
    ((SELECT id FROM opponents WHERE name = 'Ак Барс'), '2025-03-25 19:00:00+03', 'home', '2024/25', 'scheduled', 0, 0),
    ((SELECT id FROM opponents WHERE name = 'Авангард'), '2025-03-30 16:00:00+03', 'away', '2024/25', 'scheduled', 0, 0),
    ((SELECT id FROM opponents WHERE name = 'Салават Юлаев'), '2025-04-05 17:00:00+03', 'home', '2024/25', 'scheduled', 0, 0),
    ((SELECT id FROM opponents WHERE name = 'Торпедо'), '2025-04-10 19:00:00+03', 'away', '2024/25', 'scheduled', 0, 0),

    -- Дерби
    ((SELECT id FROM opponents WHERE name = 'Спартак'), '2025-04-15 19:00:00+03', 'home', '2024/25', 'scheduled', 0, 0);