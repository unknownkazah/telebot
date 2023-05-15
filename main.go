package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"strconv"
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

	msg := tgbotapi.NewMessage(0, "Введите числовое значение для term1")
	msg.ReplyMarkup = numericKeyboard
	bot.Send(msg)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore non-Message updates
			continue
		}

		msg = tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

		switch update.Message.Text {
		case "open":
			msg.ReplyMarkup = numericKeyboard
		case "close":
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		case "Рассчитать нужное количество баллов на файнале":
			msg.Text = "Введите значение для midterm"
			//msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			bot.Send(msg)
		default:
			msg.Text = "default"

		}

		//if _, err := bot.Send(msg); err != nil {
		//	log.Panic(err)
		//}
		// Wait for the next message from the user
		update = <-updates

		// Parse the value for term1 from the user's message
		midterm, err := strconv.ParseFloat(update.Message.Text, 64)
		if err != nil {
			msg.Text = "Неверный ввод. Введите числовое значение для midterm"
			bot.Send(msg)
			continue
		}

		// Prompt the user for the value of term2
		msg.Text = "Введите значение для endterm"
		bot.Send(msg)
		update = <-updates

		// Parse the value for term2 from the user's message
		endterm, err := strconv.ParseFloat(update.Message.Text, 64)
		if err != nil {
			msg.Text = "Неверный ввод. Введите числовое значение для endterm"
			bot.Send(msg)
			continue
		}

		// Calculate the required number of points for the final exam using the given formula
		//final := (term1+term2)/2*0.6 + session*0.4
		session := (70 - 0.6*((midterm+endterm)/2)) / 0.4
		result := strconv.FormatFloat(session, 'f', 2, 64)
		p := endterm + midterm

		if session >= 0 && session <= 50 {
			msg.Text = "Вам необходимо набрать 50 баллов"
		} else if p/2 < 50 {
			msg.Text = "У вас нет допуска к сесии"

		} else {
			msg.Text = "Требуемое количество баллов для сохранения стипендии: " + result
		}

		// Send the result back to the user

		msg.ReplyMarkup = numericKeyboard
		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}
