package main

//чисто для тестирования, всё скопировано.
import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/TanyaMasharova/hockey-team/internal/api/http/handlers"
	"github.com/TanyaMasharova/hockey-team/internal/repository/postgres"
	"github.com/TanyaMasharova/hockey-team/internal/service/admin"
	"github.com/TanyaMasharova/hockey-team/internal/service/auth"
	"github.com/TanyaMasharova/hockey-team/internal/service/matches"
	"github.com/TanyaMasharova/hockey-team/internal/service/seat"
	"github.com/TanyaMasharova/hockey-team/internal/service/sector"
	"github.com/TanyaMasharova/hockey-team/internal/service/ticket"
	"github.com/TanyaMasharova/hockey-team/pkg/database"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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

	logger.Info("Запуск приложения v0.3.0")

	// 2. Загружаем .env
	logger.Info("Загрузка .env файла")
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
	//4. Подключение  к БД
	logger.WithFields(logrus.Fields{
		"host": cfg.Host,
		"port": cfg.Port,
		"db":   cfg.DBName,
		"user": cfg.User,
	}).Info("Подключение к базе данных")

	db, err := database.NewPostgresDB(cfg, logger)
	if err != nil {
		logger.WithError(err).Fatal("Не удалось подключиться к БД")
	}
	defer db.Close()

	logger.Info("База данных успешно подключена!")

	//5. Инициализация репозитория
	userRepo := postgres.NewUserRepository(db)
	matchRepo := postgres.NewMatchRepository(db)
	ticketRepo := postgres.NewTicketRepository(db)
	sectorRepo := postgres.NewSectorRepository(db)
	seatRepo := postgres.NewSeatRepository(db)
	adminRepo := postgres.NewAdminRepository(db)

	//6. Инициализация сервиса
	authService := auth.NewService(userRepo)
	matchService := matches.NewService(matchRepo)
	ticketService := ticket.NewService(ticketRepo)
	sectorService := sector.NewService(sectorRepo)
	seatService := seat.NewService(seatRepo)
	adminService := admin.NewService(adminRepo)

	//7. Инициализация хэндлера
	userHandler := handlers.NewUserHandler(authService, logger)
	matchesHandler := handlers.NewMatchesHandler(matchService, logger)
	ticketHandler := handlers.NewTicketHandler(ticketService, logger)
	sectorHandler := handlers.NewSectorHandler(sectorService, logger)
	seatHandler := handlers.NewSeatHandler(seatService, logger)
	adminHandler := handlers.NewAdminHandler(adminService, logger)

	//8. Настройка маршрутизаора gin
	//!!!!! Почитать про моды
	ginMode := getEnv("GIN_MODE", "debug")
	gin.SetMode(ginMode)

	router := gin.Default()
	corsAllowedOrigins := getEnvAsSlice("CORS_ALLOWED_ORIGINS", []string{"http://localhost:3000"})

	router.Use(cors.New(cors.Config{ //настройка корсов
		AllowOrigins:     corsAllowedOrigins,                                           //с каких доменов разрешены запросы
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}, //какие методы разрешены
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},          //какие заголовки можно добавлять
		ExposeHeaders:    []string{"Content-Length"},                                   //какие заголовки сможет читать фронтенд
		AllowCredentials: true,                                                         //передавать куки и авторизацию
		//MaxAge:           12 * time.Hour, //можно кэшировать по времени результат CORS-проверки
	}))

	api := router.Group("/api")
	{
		api.POST("/register", userHandler.Register)
		api.POST("/login", userHandler.Login)

		api.GET("/matches", matchesHandler.GetMatches)
		api.GET("/matches/:id", matchesHandler.GetMatchByID)
		api.GET("/matchesStats", matchesHandler.GetStatsMatches)
		api.GET("/user/:id", userHandler.GetUserByID)
		api.GET("/tickets/user/:user_id", ticketHandler.GetUserTickets)
		api.PATCH("/user/:user_id/profile/field", userHandler.UpdateField)
		api.POST("/tickets", ticketHandler.CreateTicket)

		api.GET("/stadium/sectors", sectorHandler.GetAllSectors)
		api.GET("/stadium/sectors/:sectorId/seats", seatHandler.GetSeatsBySector)

		api.GET("/admin/stats-summary", adminHandler.GetStatsSummary)
		api.GET("/admin/stats", adminHandler.GetAllStats)
	}
	//9. Настройка http-Сервера
	port := getEnv("HTTP_PORT", "8080")
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	//10. Создаём контекст, который удем отменять при получении сигналов

	//Сигналы:
	//- **SIGTERM** — стандартный сигнал на корректное завершение;
	//- **SIGINT** — прерывание (например, Ctrl+C);
	//- **SIGKILL** — немедленное завершение (не перехватывается).

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM) //фоновый налюдатель, который ловит сигналы о завершении (например, пользователь нажал на Ctrl+c)

	defer stop() //выполнится вторым. Отключаем наблюдателя сигналов

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.WithError(err).Fatal("Server is not running, an error occurred")
		}
	}()
	<-ctx.Done()
	logger.Info("Termination signal received, shutdown the server")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel() //выполнится первым.
	//освобождает следующие ресурсы:
	//таймер
	//дочерние контексты (в коде их нет)

	//server.Shutdown() перестает принимать новые запросы, ждет их завершения 10 секунд
	//Если 10 секунд прошло, а запросы ещё висят — закрывается принудительно
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.WithError(err).Error("Error stopping the server")
	} else {
		logger.Info("The server was stopped correctly")
	}

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

func getEnvAsSlice(key string, defaultValue []string) []string {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}

	parts := strings.Split(valueStr, ",")
	values := make([]string, 0, len(parts))
	for _, part := range parts {
		value := strings.TrimSpace(part)
		if value != "" {
			values = append(values, value)
		}
	}

	if len(values) == 0 {
		return defaultValue
	}

	return values
}
