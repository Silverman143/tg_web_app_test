package main

import (
	"log"
	"strings"

	tgbotapi "project-2x/pkg/telegramBot/all"

	"github.com/sirupsen/logrus"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("631028305:AAEGxnNsXDWSHDm-k8fUKbhF-62aLmN--3I")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	// Установка уровня логирования на Debug
	logrus.SetLevel(logrus.DebugLevel)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		logrus.Panic(err)
	}

	for update := range updates {
		if update.Message != nil && update.Message.IsCommand() {
			if strings.HasPrefix(update.Message.Text, "/start") {
				logrus.Debug("--------------------- Start bot command ------------------")

				// Извлечение реферального ключа из сообщения
				var refKey string
				parts := strings.Fields(update.Message.Text)
				if len(parts) > 1 {
					refKey = parts[1]
				}

				logrus.Debugf("ref key = %s \nfrom user %s", refKey, update.Message.From.UserName)
			} else {
				logrus.Debugf("Command = %s \nfrom user %s", update.Message.Text, update.Message.From.UserName)
			}
		}
	}
}

// "INSERT INTO Referrals (referrer, referral, bonus) VALUES ((SELECT id FROM Users WHERE telegram_id = $1), (SELECT id FROM Users WHERE telegram_id = $2), $3) ON CONFLICT (referrer, referral) DO NOTHING"

// "SELECT u.telegram_id, u.user_name, u.avatar_url, u.global_rank, u.stars_balance
// 	FROM Referrals r
// 	JOIN Users u ON r.referral = u.id
// 	WHERE r.referrer = (SELECT id FROM Users WHERE telegram_id = $1)
// 	LIMIT $2 OFFSET $3;"

// 	"SELECT u.telegram_id, u.user_name, u.avatar_url, u.global_rank, u.stars_balance
// 	FROM Referrals r
// 	JOIN Users u ON r.referral = u.id
// 	WHERE r.referrer = (SELECT id FROM Users WHERE telegram_id = 22222)
// 	LIMIT 50 OFFSET 0;"