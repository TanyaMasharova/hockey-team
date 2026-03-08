package database

//файл для подключения к базе данных postgresql

import (
	"fmt"  //Стандартная библиотека Go для форматированного вывода и создания строк (используется в fmt.Sprintf и fmt.Errorf).
	"time" //Стандартная библиотека для работы с временем (используется для установки таймаутов соединения).

	"github.com/jmoiron/sqlx"    //Это основная библиотека для работы с БД. Она расширяет стандартный database/sql, делая работу с данными удобнее (например, позволяет сразу "сканировать" строки из БД в ваши Go-структуры).
	_ "github.com/lib/pq"        // Это драйвер PostgreSQL. Обратите внимание на символ _ (пустой идентификатор). Он означает, что мы импортируем пакет ради его побочного эффекта: в данном случае, чтобы он зарегистрировал себя в стандартной библиотеке database/sql как драйвер для работы с PostgreSQL. Нам не нужно вызывать его функции напрямую.
	"github.com/sirupsen/logrus" // Библиотека для логирования, чтобы мы видели, что происходит при подключении.
)

type Config struct { //структура, которую мы потом будем передавать как cfg. Значения полей берутся из config.yaml или их переменных окружения
    Host     string
    Port     int
    User     string
    Password string
    DBName   string
    SSLMode  string
}
func NewPostgresDB(cfg Config, logger *logrus.Logger) (*sqlx.DB, error) { //конструктор, который возвращает готовое подклчение в БД или ошибку
    logger.WithFields(logrus.Fields{
        "host": cfg.Host,
        "port": cfg.Port,
        "user": cfg.User,
        "dbname": cfg.DBName,
    }).Info("Connecting to database") //логирование начала подключения

    dsn := fmt.Sprintf(
        "host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
        cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
    ) //формирование строки для подключения

    db, err := sqlx.Connect("postgres", dsn) //открывает соединение с БД с помошью переданного драйвера и проверяет, живо ли оно
    if err != nil {
        return nil, fmt.Errorf("failed to connect to database: %w", err)
    } //обёртка для ошибки

    // Настройка пула соединений
    db.SetMaxOpenConns(25) //максимальное количество открытых соединений
    db.SetMaxIdleConns(25) //максимальное количество соединений, котрые хранятся в пуле для быстрого использования
    db.SetConnMaxLifetime(5 * time.Minute) //максимальное время жизни соединения. Через 5 минут соединение будет закрыто и открыто заново. Это помогает бороться с "утечками" памяти на стороне БД и сетевыми проблемами с "долгоживущими" соединениями.

    // Проверяем соединение
    if err = db.Ping(); err != nil {
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }

    logger.Info("Successfully connected to database")
    return db, nil
}