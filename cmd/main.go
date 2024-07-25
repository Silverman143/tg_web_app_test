package main

import (
	back "project-2x"
	database "project-2x/pkg/database"
	handler "project-2x/pkg/handlers"
	"project-2x/pkg/service"
	"project-2x/pkg/telegramBot"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"os"
)



func main() {

	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := initConfig(); err != nil{
		logrus.Fatalf("Init configuration fale with error: %s", err.Error())
	}

	if err := godotenv.Load(".env"); err != nil {
		logrus.Fatalf("Dotenv loading fail with error : %s", err.Error())
	}

	db, err := database.NewPostgresDB( &database.Config{
		// Host: viper.GetString("postgres.host"),
		// Port: viper.GetString("postgres.port"),
		// Username: os.Getenv("POSTGRES_USER"),
		// Password: os.Getenv("POSTGRES_PASSWORD"),
		// DBname: viper.GetString("postgres.dbName"),
		// SSLmode: "disable",
		URL: os.Getenv("POSTGRES_URL"),
		
	})

	if err != nil {
		logrus.Fatalf("Cant connect postgres db with error: %s", err.Error())
	}

	defer db.Close()

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		logrus.Fatalf("TELEGRAM_BOT_TOKEN is not set in the environment")
	}

	bot, err := telegramBot.NewBot(botToken)
	if err != nil {
		logrus.Fatalf("Failed to create Telegram bot: %s", err.Error())
	}

	datadase := database.NewDatabase(db)
	service := service.NewService(datadase, bot)
	handler := handler.NewHandler(service, bot)
	srv := new(back.Server)

	err = datadase.System.UpdateDailyBonuses(viper.GetStringMapString("dailyBonuses"))

	if err != nil {
		logrus.Fatalf("failed to update daily bonuses: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Значение по умолчанию, если переменная окружения не установлена
	}

	// Инициализация роутов и добавление обработчика вебхуков Telegram
	router := handler.InitRouts(bot)

	err = bot.Start()

	if err != nil {
		logrus.Fatalf("Failed to start Telegram bot: %s", err.Error())
	}

	if err := srv.Run(port, router); err != nil {
		logrus.Fatalf( "error while server is running, %s", err.Error())
	}
	
}



func initConfig() error {
	viper.AddConfigPath("./configs")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	return viper.ReadInConfig()
}