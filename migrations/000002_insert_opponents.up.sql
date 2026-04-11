-- Добавляем соперников (команды КХЛ сезона 2025/2026, кроме Динамо Минск)
INSERT INTO opponents (name, city, logo_url) VALUES
    -- Западная конференция
    ('ЦСКА', 'Москва', '/logos/cska.png'),
    ('СКА', 'Санкт-Петербург', '/logos/ska.png'),
    ('Динамо', 'Москва', '/logos/dinamo_m.png'),
    ('Локомотив', 'Ярославль', '/logos/lokomotiv.png'),
    ('Северсталь', 'Череповец', '/logos/severstal.png'),
    ('Торпедо', 'Нижний Новгород', '/logos/torpedo.png'),
    ('Шанхай Дрэгонс', 'Шанхай', '/logos/chanhai.png'),
    ('Лада', 'Тольятти', '/logos/lada.png'),
    ('Спартак', 'Москва', '/logos/spartak.png'),
    ('Сочи', 'Сочи', '/logos/sochi.png'),
    
    -- Восточная конференция
    ('Ак Барс', 'Казань', '/logos/ak_bars.png'),
    ('Салават Юлаев', 'Уфа', '/logos/salavat.png'),
    ('Авангард', 'Омск', '/logos/avangard.png'),
    ('Металлург', 'Магнитогорск', '/logos/metallurg.png'),
    ('Автомобилист', 'Екатеринбург', '/logos/avtomobilist.png'),
    ('Трактор', 'Челябинск', '/logos/traktor.png'),
    ('Барыс', 'Астана', '/logos/barys.png'),
    ('Сибирь', 'Новосибирск', '/logos/sibir.png'),
    ('Амур', 'Хабаровск', '/logos/amur.png'),
    ('Адмирал', 'Владивосток', '/logos/admiral.png'),
    ('Нефтехимик', 'Нижнекамск', '/logos/neftekhimik.png')
    
ON CONFLICT (name) DO NOTHING;
