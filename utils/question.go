package utils

import (
	"math/rand"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/shavkatjon/viktorina-bot/model"
)

func ListToText(list model.QuestionList) string {
	var text string

	for i, question := range list.List {
		text += "<b>" + strconv.Itoa(int(list.Page-1)*10+i+1) + "</b>." + " <b>Indeks:</b> " + strconv.Itoa(int(question.Id)) + "\n"
		text += " <b>Fan:</b> " + question.Subject + "\n"
		text += " <b>Savol:</b> " + question.Text + "\n"
		text += " <b>Javob:</b> " + question.Answer + "\n"
	}

	text += "\n<i>Umumiy savollar soni:</i> " + strconv.Itoa(int(list.Count))

	return text
}

func BuildInline(user *model.GameUser, variants []string) tgbotapi.InlineKeyboardMarkup {

	variants = append(variants, user.Answer)

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(variants), func(i, j int) { variants[i], variants[j] = variants[j], variants[i] })

	var List [][]tgbotapi.InlineKeyboardButton
	for _, variant := range variants {
		if variant == user.Answer {
			List = append(List, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(variant, "correct")))
		} else {
			List = append(List, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(variant, "wrong")))
		}
	}

	return tgbotapi.InlineKeyboardMarkup{InlineKeyboard: List}
}
