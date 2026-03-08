package main

//чисто для тестирования, всё скопировано.
import (
	"fmt"
	"os"
	"strconv"

	"github.com/TanyaMasharova/hockey-team/pkg/database"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
    // 1. Настраиваем логгер (чтобы было красиво)
    logger := logrus.New()
    logger.SetFormatter(&logrus.TextFormatter{
        FullTimestamp: true,
        ForceColors:   true,
    })
    
    logger.Info("🚀 Запуск приложения")

    // 2. Загружаем .env
    logger.Info("📁 Загрузка .env файла")
    if err := godotenv.Load(); err != nil {
        logger.Warn("⚠️ .env файл не найден, используем переменные окружения")
    }

    // 3. Читаем переменные окружения и создаём конфиг
    cfg := database.Config{
        Host:     getEnv("DB_HOST", "localhost"),
        Port:     getEnvAsInt("DB_PORT", 5432),
        User:     getEnv("DB_USER", "postgres"),
        Password: getEnv("DB_PASSWORD", ""),
        DBName:   getEnv("DB_NAME", "hockey_db"),
        SSLMode:  getEnv("DB_SSLMODE", "disable"),
    }

    // 4. Подключаемся к БД через твою функцию
    logger.WithFields(logrus.Fields{
        "host": cfg.Host,
        "port": cfg.Port,
        "db":   cfg.DBName,
        "user": cfg.User,
    }).Info("📊 Подключение к базе данных")

    db, err := database.NewPostgresDB(cfg, logger)
    if err != nil {
        logger.WithError(err).Fatal("❌ Не удалось подключиться к БД")
    }
    defer db.Close()

    logger.Info("✅ База данных успешно подключена!")
    
    // Здесь будет остальной код: репозитории, сервисы, HTTP сервер
    logger.Info("🎉 Приложение готово к работе")
    
    // Чтобы программа не завершалась сразу (для теста)
    fmt.Println("\nНажмите Ctrl+C для выхода")
    select {} // блокируем главную горутину
}

// Вспомогательные функции для чтения переменных окружения
func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
    valueStr := getEnv(key, "")
    if valueStr == "" {
        return defaultValue
    }
    
    value, err := strconv.Atoi(valueStr)
    if err != nil {
        return defaultValue
    }
    return value
}