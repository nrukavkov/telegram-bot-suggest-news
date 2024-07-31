package main

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/tucnak/telebot.v2"
)

// Структура для хранения данных о пользователях
type userMessageInfo struct {
	count     int
	firstTime time.Time
}

var userMessages = make(map[int64]*userMessageInfo)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	if err := godotenv.Load(".env.local"); err == nil {
		log.Info().Msg("Loaded vars from .env.local")
	} else {
		if err := godotenv.Load(".env"); err != nil {
			log.Warn().Msg("Cannot load .env file, env variables from system will be used")
		}
	}

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	adminID := os.Getenv("ADMIN_ID")
	rateLimit := os.Getenv("RATE_LIMIT")

	var maxMessages int
	var intervalDuration time.Duration

	// Parse RATE_LIMIT
	if rateLimit != "" && rateLimit != "0" && rateLimit != "0:0" {
		rlParts := strings.Split(rateLimit, ":")
		if len(rlParts) != 2 {
			log.Fatal().Msg("Wrong RATE_LIMIT param, expecting <number_of_messages>:<interval_in_seconds>")
		}
		var err error
		maxMessages, err = strconv.Atoi(rlParts[0])
		if err != nil {
			log.Fatal().Msgf("Wrong RATE_LIMIT param: %v", err)
		}
		interval, err := strconv.Atoi(rlParts[1])
		if err != nil {
			log.Fatal().Msgf("Wrong RATE_LIMIT param: %v", err)
		}
		intervalDuration = time.Duration(interval) * time.Second
	}

	// Init bot
	bot, err := telebot.NewBot(telebot.Settings{
		Token:  botToken,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal().Msgf("Error bot init: %v", err)
	}

	bot.Handle("/start", func(m *telebot.Message) {
		bot.Send(m.Sender, "Привет! Отправьте мне новость, которую хотите предложить для канала.")
		log.Info().Msgf("User %s (%d) sent command /start", m.Sender.Username, m.Sender.ID)
	})

	bot.Handle("/getid", func(m *telebot.Message) {
		bot.Send(m.Sender, "Ваш ID: "+strconv.FormatInt(m.Sender.ID, 10))
		log.Info().Msgf("User %s (%d) sent command /getid", m.Sender.Username, m.Sender.ID)
	})

	bot.Handle(telebot.OnText, func(m *telebot.Message) {
		if maxMessages > 0 && intervalDuration > 0 {
			now := time.Now()

			// Get info about messages
			userInfo, exists := userMessages[m.Sender.ID]
			if !exists {
				userInfo = &userMessageInfo{count: 0, firstTime: now}
				userMessages[m.Sender.ID] = userInfo
			}

			// Check limits
			if now.Sub(userInfo.firstTime) > intervalDuration {
				userInfo.count = 0
				userInfo.firstTime = now
			}

			if userInfo.count >= maxMessages {
				nextAllowedTime := userInfo.firstTime.Add(intervalDuration).Truncate(time.Minute)
				bot.Send(m.Sender, "Вы отправляете сообщения слишком часто. Пожалуйста, подождите до "+nextAllowedTime.Format("15:04")+" чтобы отправить следующее сообщение.")
				log.Warn().Msgf("Пользователь %s (%d) превысил лимит сообщений", m.Sender.Username, m.Sender.ID)
				return
			}

			userInfo.count++
		}

		admin, err := bot.ChatByID(adminID)
		if err != nil {
			log.Error().Msgf("Error get admin ID: %v", err)
			return
		}

		bot.Send(admin, "Новое предложение от @"+m.Sender.Username+":\n"+m.Text)
		bot.Send(m.Sender, "Спасибо! Ваше предложение отправлено на рассмотрение.")
		log.Info().Msgf("User %s (%d) send a suggestion: %s", m.Sender.Username, m.Sender.ID, m.Text)
	})

	log.Info().Msg("Bot is started")
	bot.Start()
}
