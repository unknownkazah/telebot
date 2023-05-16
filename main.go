package main

import (
	"log"
	"os"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var numericKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Рассчитать нужное количество баллов на файнале"),
	),
)

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	go func() {
		for {
			sendPeriodicRequest(bot)
			time.Sleep(15 * time.Minute)
		}
	}()

	for update := range updates {
		if update.Message == nil {
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		switch update.Message.Text {
		case "/start":
			msg.Text = "Введите числовое значение для midterm"
			msg.ReplyMarkup = numericKeyboard
		case "/stop":
			msg.Text = "Keyboard closed"
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		case "Рассчитать нужное количество баллов на файнале":
			msg.Text = "Введите значение для midterm"
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		default:
			msg.Text = "Неверная команда. Используйте /start, /stop или Рассчитать нужное количество баллов на файнале"
			msg.ReplyMarkup = numericKeyboard
		}

		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}

		if update.Message.Text == "/stop" {
			continue
		}

		update = <-updates

		midterm, err := strconv.ParseFloat(update.Message.Text, 64)
		if err != nil {
			msg.Text = "Неверный ввод. Введите числовое значение для midterm"
			msg.ReplyMarkup = numericKeyboard
			bot.Send(msg)
			continue
		}

		msg.Text = "Введите значение для endterm"
		bot.Send(msg)

		update = <-updates

		endterm, err := strconv.ParseFloat(update.Message.Text, 64)
		if err != nil {
			msg.Text = "Неверный ввод. Введите числовое значение для endterm"
			msg.ReplyMarkup = numericKeyboard
			bot.Send(msg)
			continue
		}

		session := (70 - 0.6*((midterm+endterm)/2)) / 0.4

		result := strconv.FormatFloat(session, 'f', 2, 64)
		p := endterm + midterm
		total := ((p / 2) * 0.6) + (100 * 0.4)
		result1 := strconv.FormatFloat(total, 'f', 2, 64)

		if session >= 0 && session <= 50 {
			msg.Text = "Вам необходимо набрать 50 баллов"
		} else if p/2 < 50 {
			msg.Text = "У вас нет допуска к сессии"
		} else {
			msg.Text = "Требуемое количество баллов для сохранения стипендии: " + result
			msg.Text += "\nЕсли вы сдадите файнал на 100%, вы получите тотал: " + result1
		}

		msg.ReplyMarkup = numericKeyboard
		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}

func sendPeriodicRequest(bot *tgbotapi.BotAPI) {
	for {
		// Create the message request
		msg := tgbotapi.NewMessage(5970395353, "This is a periodic message!")

		// Send the message
		_, err := bot.Send(msg)
		if err != nil {
			log.Println("Error sending message:", err)
		}

		// Sleep for 15 minutes
		time.Sleep(15 * time.Minute)
	}
}
