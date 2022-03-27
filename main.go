package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	"github.com/shavkatjon/viktorina_bot_player/model"
	"github.com/shavkatjon/viktorina_bot_player/storage"
	"github.com/shavkatjon/viktorina_bot_player/utils"
)

var subject = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Matematika🧮"),
		tgbotapi.NewKeyboardButton("Tarix⏳"),
		tgbotapi.NewKeyboardButton("Ingliz tili🇬🇧"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Rus tili🇷🇺"),
		tgbotapi.NewKeyboardButton("Geografiya🗺"),
		tgbotapi.NewKeyboardButton("Adabiyot🪶"),
	),
)

func main() {

	var stop = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Viktorinani to'xtatish⏹️"),
		),
	)

	var groupMenuKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Profil📊"),
			tgbotapi.NewKeyboardButton("Reyting🔝"),
			tgbotapi.NewKeyboardButton("Viktorinani boshlash🧩"),
		),
	)

	e := godotenv.Load()
	utils.Check(e)

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_TOKEN2"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	// u := tgbotapi.NewUpdate(0)
	// u.Timeout = 60

	dbConn := storage.ConnectGameDb()
	if dbConn {
		fmt.Println("Connected to db")
	} else {
		fmt.Println("Not connected to db")
	}

	// Connect to history database
	dbConn = storage.ConnectHistoryDb()
	if dbConn {
		fmt.Println("Connected to db")
	} else {
		fmt.Println("Not connected to db")
	}

	go HistoryBot()

	// updates, _ := bot.GetUpdatesChan(u)
	_, err = bot.SetWebhook(tgbotapi.NewWebhook("https://tarixviktorinabot.herokuapp.com/" + bot.Token))
	utils.Check(err)
	updates := bot.ListenForWebhook("/" + bot.Token)
	go http.ListenAndServe(":"+os.Getenv("PORT"), nil)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		user := storage.GameGetUser(update.Message.Chat.ID)

		switch user.Step {
		case 1:
			user.FirstName = update.Message.Chat.FirstName
			user.LastName = update.Message.Chat.LastName
			user.UserName = update.Message.Chat.UserName
			user.Step++

			msg := tgbotapi.NewMessage(user.ChatId,
				"Assalomu alaykum "+user.FirstName+" "+user.LastName+", tarixdan viktorina botga xush kelibsiz!")
			msg.ReplyMarkup = groupMenuKeyboard
			m, _ := bot.Send(msg)

			user.MessageId = int64(m.MessageID)
			storage.GameUpdateUser(user)
		case 2:
			var m tgbotapi.Message
			msg := tgbotapi.NewMessage(user.ChatId, "")
			msg.ParseMode = "html"
			msg.ReplyMarkup = groupMenuKeyboard
			switch update.Message.Text {
			case "Reyting🔝":
				userList := storage.GameGetUserList(10)
				user.Step = 2

				fmt.Println(userList)
				msg.Text = utils.GameListToText(userList)
				m, _ = bot.Send(msg)
				msg.Text = "<i><b>Sizning profilingiz</b></i>:\n<b>Ism:</b> " + user.FirstName + "\n<b>Familiya:</b> " + user.LastName + "\n<b>Siz to'plagan ballar:</b> " + strconv.Itoa(user.Score)
			case "Viktorinani boshlash🧩":
				msg := tgbotapi.NewMessage(user.ChatId, "")
				msg.ParseMode = "html"
				msg.Text = "Qaysi fandan viktorinani boshlamoqchisiz?"
				msg.ReplyMarkup = subject
				user.Step++
				bot.Send((msg))
				user.MessageId = int64(m.MessageID)
				storage.GameUpdateUser(user)
			case "Profil📊":
				msg := tgbotapi.NewMessage(user.ChatId,
					"<i><b>Sizning profilingiz</b></i>:\n<b>Ism:</b> "+user.FirstName+"\n<b>Familiya:</b> "+user.LastName+"\n<b>Siz to'plagan ballar:</b> "+strconv.Itoa(user.Score))
				msg.ParseMode = "html"
				bot.Send((msg))
			default:
				msg := tgbotapi.NewMessage(user.ChatId, "Bunday buyruq yo'q🙃")
				msg.ReplyMarkup = groupMenuKeyboard
				user.Step = 2
				m, _ := bot.Send(msg)
				user.MessageId = int64(m.MessageID)
				storage.GameUpdateUser(user)
			}
			user.MessageId = int64(m.MessageID)
			storage.GameUpdateUser(user)
		case 3:
			switch update.Message.Text {
			case "Matematika🧮":
				user.Subject = "math"
				question := storage.GameGetQuestion(user.Subject)
				user.Question = question.Text
				user.Answer = question.Answer
				msg := tgbotapi.NewMessage(user.ChatId, "")
				msg.ParseMode = "html"
				if storage.GameGetNumberOfQuestions("math") == 0 {
					msg.Text = "Bu fan bo'yicha savollar hali qo'shilmagan."
					msg.ReplyMarkup = groupMenuKeyboard
					user.Step = 2
				} else {
					msg.Text = "<b>Viktorina boshlandi!!!</b>\n<b>Savol:</b> " + user.Question + "\n<b>Sizning javobingiz:</b>"
					msg.ReplyMarkup = stop
					user.Step++
				}
				bot.Send((msg))
			case "Tarix⏳":
				user.Subject = "history"
				question := storage.GameGetQuestion(user.Subject)
				user.Question = question.Text
				user.Answer = question.Answer
				msg := tgbotapi.NewMessage(user.ChatId, "")
				msg.ParseMode = "html"
				if storage.GameGetNumberOfQuestions("history") == 0 {
					msg.Text = "Bu fan bo'yicha savollar hali qo'shilmagan."
					msg.ReplyMarkup = groupMenuKeyboard
					user.Step = 2
				} else {
					msg.Text = "<b>Viktorina boshlandi!!!</b>\n<b>Savol:</b> " + user.Question + "\n<b>Sizning javobingiz:</b>"
					msg.ReplyMarkup = stop
					user.Step++
				}
				bot.Send((msg))
			case "Ingliz tili🇬🇧":
				user.Subject = "english"
				question := storage.GameGetQuestion(user.Subject)
				user.Question = question.Text
				user.Answer = question.Answer
				msg := tgbotapi.NewMessage(user.ChatId, "")
				msg.ParseMode = "html"
				if storage.GameGetNumberOfQuestions("english") == 0 {
					msg.Text = "Bu fan bo'yicha savollar hali qo'shilmagan."
					msg.ReplyMarkup = groupMenuKeyboard
					user.Step = 2
				} else {
					msg.Text = "<b>Viktorina boshlandi!!!</b>\n<b>Savol:</b> " + user.Question + "\n<b>Sizning javobingiz:</b>"
					msg.ReplyMarkup = stop
					user.Step++
				}
				bot.Send((msg))
			case "Rus tili🇷🇺":
				user.Subject = "russian"
				question := storage.GameGetQuestion(user.Subject)
				user.Question = question.Text
				user.Answer = question.Answer
				msg := tgbotapi.NewMessage(user.ChatId, "")
				msg.ParseMode = "html"
				if storage.GameGetNumberOfQuestions("russian") == 0 {
					msg.Text = "Bu fan bo'yicha savollar hali qo'shilmagan."
					msg.ReplyMarkup = groupMenuKeyboard
					user.Step = 2
				} else {
					msg.Text = "<b>Viktorina boshlandi!!!</b>\n<b>Savol:</b> " + user.Question + "\n<b>Sizning javobingiz:</b>"
					msg.ReplyMarkup = stop
					user.Step++
				}
				bot.Send((msg))
			case "Geografiya🗺":
				user.Subject = "geography"
				question := storage.GameGetQuestion(user.Subject)
				user.Question = question.Text
				user.Answer = question.Answer
				msg := tgbotapi.NewMessage(user.ChatId, "")
				msg.ParseMode = "html"
				if storage.GameGetNumberOfQuestions("geography") == 0 {
					msg.Text = "Bu fan bo'yicha savollar hali qo'shilmagan."
					msg.ReplyMarkup = groupMenuKeyboard
					user.Step = 2
				} else {
					msg.Text = "<b>Viktorina boshlandi!!!</b>\n<b>Savol:</b> " + user.Question + "\n<b>Sizning javobingiz:</b>"
					msg.ReplyMarkup = stop
					user.Step++
				}
				bot.Send((msg))
			case "Adabiyot🪶":
				user.Subject = "literature"
				question := storage.GameGetQuestion(user.Subject)
				user.Question = question.Text
				user.Answer = question.Answer
				msg := tgbotapi.NewMessage(user.ChatId, "")
				msg.ParseMode = "html"
				if storage.GameGetNumberOfQuestions("literature") == 0 {
					msg.Text = "Bu fan bo'yicha savollar hali qo'shilmagan."
					msg.ReplyMarkup = groupMenuKeyboard
					user.Step = 2
				} else {
					msg.Text = "<b>Viktorina boshlandi!!!</b>\n<b>Savol:</b> " + user.Question + "\n<b>Sizning javobingiz:</b>"
					msg.ReplyMarkup = stop
					user.Step++
				}
				bot.Send((msg))
			default:
				msg := tgbotapi.NewMessage(user.ChatId, "")
				msg.ParseMode = "html"
				msg.Text = "Bunday fan yo'q🙃"
				msg.ReplyMarkup = groupMenuKeyboard
				user.Step = 2
				bot.Send((msg))
			}
			var m tgbotapi.Message
			user.MessageId = int64(m.MessageID)
			storage.GameUpdateUser(user)
		case 4:
			switch strings.ToLower(update.Message.Text) {
			case "viktorinani to'xtatish⏹️":
				msg := tgbotapi.NewMessage(user.ChatId, "Viktorina to'xtatildi")
				msg.ReplyMarkup = groupMenuKeyboard
				user.Step = 2
				m, _ := bot.Send(msg)
				user.MessageId = int64(m.MessageID)
				storage.GameUpdateUser(user)
			default:
				question := storage.GameGetQuestion(user.Subject)
				user.Question = question.Text
				fmt.Println(question.Answer)
				fmt.Println(user.Answer)
				fmt.Println(update.Message.Text)
				if strings.EqualFold(strings.ToLower(update.Message.Text), strings.ToLower(user.Answer)) {
					msg := tgbotapi.NewMessage(user.ChatId, "")
					msg.ParseMode = "html"
					msg.ReplyMarkup = stop
					msg.Text = "Javobingiz to'g'ri😃"
					user.Score++
					m, _ := bot.Send(msg)
					msg.Text = "<b>Keyingi savol:\n</b>" + user.Question + "\n<b>Sizning javobingiz:\n</b>"
					bot.Send(msg)
					user.MessageId = int64(m.MessageID)
					storage.GameUpdateUser(user)
				} else {
					msg := tgbotapi.NewMessage(user.ChatId, "")
					msg.ParseMode = "html"
					msg.ReplyMarkup = stop
					msg.Text = "Javobingiz noto'g'ri🙁\nTo'g'ri javob: <b>" + user.Answer + "</b>"
					m, _ := bot.Send(msg)
					msg.Text = "<b>Keyingi savol:\n</b>" + user.Question + "\n<b>Sizning javobingiz:\n</b>"
					bot.Send(msg)
					user.MessageId = int64(m.MessageID)
					storage.GameUpdateUser(user)
				}
				user.Answer = question.Answer
			}
		default:
			msg := tgbotapi.NewMessage(user.ChatId, "Bunday buyruq yo'q🙃")
			msg.ReplyMarkup = groupMenuKeyboard
			user.Step = 2
			m, _ := bot.Send(msg)
			user.MessageId = int64(m.MessageID)
			storage.GameUpdateUser(user)
		}
		storage.GameUpdateUser(user)
	}
}

func HistoryBot() {
	// Buttons
	defaultMenuKeyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Savollar to'plami📋"),
			tgbotapi.NewKeyboardButton("Savol qo'shish📥"),
			tgbotapi.NewKeyboardButton("Savolni tahrirlash📝"),
		),
	)

	editMenuKeyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Savolni tahrirlash♻️"),
			tgbotapi.NewKeyboardButton("Savolni o'chirish🗑"),
		),
	)

	e := godotenv.Load()
	utils.Check(e)

	// Initialize bot
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_TOKEN"))
	utils.Check(err)

	bot.Debug = true
	fmt.Println("Authorized on account", bot.Self.UserName)

	// Localniy run qilish uchun pastdagi 4 ta qatorni kommmentdan chiqarish kerak

	// u := tgbotapi.NewUpdate(0)
	// u.Timeout = 60
	// updates, _ := bot.GetUpdatesChan(u)

	// Herokuda run qilish uchun pastdagi 4 ta qatorni kommentdan chiqarish kerak

	_, err = bot.SetWebhook(tgbotapi.NewWebhook("https://tarixviktorinabot.herokuapp.com/" + bot.Token))
	utils.Check(err)
	updates := bot.ListenForWebhook("/" + bot.Token)
	go http.ListenAndServe(":"+os.Getenv("PORT"), nil)

	for update := range updates {
		if update.CallbackQuery != nil {
			user := storage.GetUser(int64(update.CallbackQuery.From.ID))

			callbackQueryDataSplit := strings.Split(update.CallbackQuery.Data, "#")

			switch user.Step {
			case 2:
				switch callbackQueryDataSplit[0] {
				case "next":
					{
						limit, err := strconv.Atoi(callbackQueryDataSplit[1])
						utils.Check(err)
						page, err := strconv.Atoi(callbackQueryDataSplit[2])
						utils.Check(err)

						msg := tgbotapi.NewEditMessageText(user.ChatId, int(user.MessageId), "")
						msg.ParseMode = "html"

						var inlineKeyboard tgbotapi.InlineKeyboardMarkup
						user.Step = 2
						questionList := storage.GetQuestionList(int64(limit), int64(page))
						if questionList.Page < questionList.PageCount && questionList.Page > 1 {
							inlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(
								tgbotapi.NewInlineKeyboardRow(
									tgbotapi.NewInlineKeyboardButtonData("Prev", "prev#"+strconv.Itoa(int(questionList.Limit))+"#"+strconv.Itoa(int(questionList.Page-1))),
									tgbotapi.NewInlineKeyboardButtonData(strconv.Itoa(int(questionList.Page))+"/"+strconv.Itoa(int(questionList.PageCount)), "page"),
									tgbotapi.NewInlineKeyboardButtonData("Next", "next#"+strconv.Itoa(int(questionList.Limit))+"#"+strconv.Itoa(int(questionList.Page+1))),
								),
							)
							msg.ReplyMarkup = &inlineKeyboard
						} else if questionList.Page < questionList.PageCount {
							inlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(
								tgbotapi.NewInlineKeyboardRow(
									tgbotapi.NewInlineKeyboardButtonData(strconv.Itoa(int(questionList.Page))+"/"+strconv.Itoa(int(questionList.PageCount)), "page"),
									tgbotapi.NewInlineKeyboardButtonData("Next", "next#"+strconv.Itoa(int(questionList.Limit))+"#"+strconv.Itoa(int(questionList.Page+1))),
								),
							)
							msg.ReplyMarkup = &inlineKeyboard
						} else if questionList.Page > 1 {
							inlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(
								tgbotapi.NewInlineKeyboardRow(
									tgbotapi.NewInlineKeyboardButtonData("Prev", "prev#"+strconv.Itoa(int(questionList.Limit))+"#"+strconv.Itoa(int(questionList.Page-1))),
									tgbotapi.NewInlineKeyboardButtonData(strconv.Itoa(int(questionList.Page))+"/"+strconv.Itoa(int(questionList.PageCount)), "page"),
								),
							)
							msg.ReplyMarkup = &inlineKeyboard
						}

						msg.Text = utils.ListToText(questionList)
						bot.Send(msg)
						storage.UpdateUser(user)
					}
				case "prev":
					{
						limit, err := strconv.Atoi(callbackQueryDataSplit[1])
						utils.Check(err)
						page, err := strconv.Atoi(callbackQueryDataSplit[2])
						utils.Check(err)

						msg := tgbotapi.NewEditMessageText(user.ChatId, int(user.MessageId), "")
						msg.ParseMode = "html"

						var inlineKeyboard tgbotapi.InlineKeyboardMarkup
						user.Step = 2
						questionList := storage.GetQuestionList(int64(limit), int64(page))
						if questionList.Page < questionList.PageCount && questionList.Page > 1 {
							inlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(
								tgbotapi.NewInlineKeyboardRow(
									tgbotapi.NewInlineKeyboardButtonData("Prev", "prev#"+strconv.Itoa(int(questionList.Limit))+"#"+strconv.Itoa(int(questionList.Page-1))),
									tgbotapi.NewInlineKeyboardButtonData(strconv.Itoa(int(questionList.Page))+"/"+strconv.Itoa(int(questionList.PageCount)), "page"),
									tgbotapi.NewInlineKeyboardButtonData("Next", "next#"+strconv.Itoa(int(questionList.Limit))+"#"+strconv.Itoa(int(questionList.Page+1))),
								),
							)
							msg.ReplyMarkup = &inlineKeyboard
						} else if questionList.Page < questionList.PageCount {
							inlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(
								tgbotapi.NewInlineKeyboardRow(
									tgbotapi.NewInlineKeyboardButtonData(strconv.Itoa(int(questionList.Page))+"/"+strconv.Itoa(int(questionList.PageCount)), "page"),
									tgbotapi.NewInlineKeyboardButtonData("Next", "next#"+strconv.Itoa(int(questionList.Limit))+"#"+strconv.Itoa(int(questionList.Page+1))),
								),
							)
							msg.ReplyMarkup = &inlineKeyboard
						} else if questionList.Page > 1 {
							inlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(
								tgbotapi.NewInlineKeyboardRow(
									tgbotapi.NewInlineKeyboardButtonData("Prev", "prev#"+strconv.Itoa(int(questionList.Limit))+"#"+strconv.Itoa(int(questionList.Page-1))),
									tgbotapi.NewInlineKeyboardButtonData(strconv.Itoa(int(questionList.Page))+"/"+strconv.Itoa(int(questionList.PageCount)), "page"),
								),
							)
							msg.ReplyMarkup = &inlineKeyboard
						}

						msg.Text = utils.ListToText(questionList)
						bot.Send(msg)
						storage.UpdateUser(user)
					}
				}
			}

		}

		if update.Message == nil {
			continue
		}

		user := storage.GetUser(update.Message.Chat.ID)

		switch user.Step {
		case 1:
			user.FirstName = update.Message.Chat.FirstName
			user.LastName = update.Message.Chat.LastName
			user.UserName = update.Message.Chat.UserName
			user.Step++

			msg := tgbotapi.NewMessage(user.ChatId,
				"Assalomu alaykum!\n"+user.FirstName+" "+user.LastName+" tarixdan viktorina boti( @tarixviktorina_bot ) uchun savollar tuzadigan botga xush kelibsiz!")

			msg.ParseMode = "html"
			msg.ReplyMarkup = defaultMenuKeyboard
			m, _ := bot.Send(msg)

			user.MessageId = int64(m.MessageID)
			storage.UpdateUser(user)
		case 2:
			var m tgbotapi.Message
			msg := tgbotapi.NewMessage(user.ChatId, "")
			msg.ParseMode = "html"

			switch update.Message.Text {
			case "Savollar to'plami📋":
				{
					var inlineKeyboard tgbotapi.InlineKeyboardMarkup
					user.Step = 2
					questionList := storage.GetQuestionList(10, 1)
					if questionList.Page < questionList.PageCount {
						inlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(
							tgbotapi.NewInlineKeyboardRow(
								tgbotapi.NewInlineKeyboardButtonData(strconv.Itoa(int(questionList.Page))+"/"+strconv.Itoa(int(questionList.PageCount)), "page"),
								tgbotapi.NewInlineKeyboardButtonData("Next", "next#"+strconv.Itoa(int(questionList.Limit))+"#"+strconv.Itoa(int(questionList.Page+1))),
							),
						)
						msg.ReplyMarkup = inlineKeyboard
					}

					msg.Text = utils.ListToText(questionList)
					m, _ = bot.Send(msg)
				}
			case "Savol qo'shish📥":
				{
					user.Step++
					msg.Text = "Kiritmoqchi bo'lgan savolingizning turini tanlang:"
					msg.ReplyMarkup = subject
					m, _ = bot.Send(msg)
				}
			case "Savolni tahrirlash📝":
				{
					user.Step = 6
					msg.Text = "Tahrirlamoqchi yoki o'chirmoqchi bo'lgan savolingizni indeksini kiriting:"
					msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
					m, _ = bot.Send(msg)
				}
			default:
				{
					msg = tgbotapi.NewMessage(user.ChatId, "Men bunday buyruqni bilmayman😕")
					msg.ReplyMarkup = defaultMenuKeyboard
					m, _ = bot.Send(msg)
				}
			}

			user.MessageId = int64(m.MessageID)
			storage.UpdateUser(user)
		case 3:
			user.Step++
			msg := tgbotapi.NewMessage(user.ChatId, "Yangi savolni kiriting:")
			msg.ParseMode = "html"
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			var question model.Question
			switch update.Message.Text {
			case "Matematika🧮":
				question.Subject = "math"
			case "Tarix⏳":
				question.Subject = "history"
			case "Ingliz tili🇬🇧":
				question.Subject = "english"
			case "Rus tili🇷🇺":
				question.Subject = "russian"
			case "Geografiya🗺":
				question.Subject = "geography"
			case "Adabiyot🪶":
				question.Subject = "literature"
			default:
				msg := tgbotapi.NewMessage(user.ChatId, "")
				msg.ParseMode = "html"
				msg.Text = "Bunday fan yo'q🙃"
				msg.ReplyMarkup = defaultMenuKeyboard
				user.Step = 2
				bot.Send((msg))
			}
			question.UserId = user.Id
			user.QuestionId = storage.InsertQuestion(question)
			m, _ := bot.Send(msg)
			user.MessageId = int64(m.MessageID)
			storage.UpdateUser(user)
		case 4:
			user.Step++
			msg := tgbotapi.NewMessage(user.ChatId, "Savolingizning javobini kiriting:")
			msg.ParseMode = "html"
			var question model.Question
			question = storage.GetQuestion(user.QuestionId)
			question.Text = update.Message.Text
			question.UserId = user.Id
			storage.UpdateQuestion(question)
			m, _ := bot.Send(msg)
			user.MessageId = int64(m.MessageID)
			storage.UpdateUser(user)
		case 5:
			var question model.Question
			question = storage.GetQuestion(user.QuestionId)
			question.Answer = update.Message.Text
			question.Status = 1
			storage.UpdateQuestion(question)

			msg := tgbotapi.NewMessage(user.ChatId, "Savol muvaffaqiyatli qo'shildi✅")
			msg.ParseMode = "html"
			msg.ReplyMarkup = defaultMenuKeyboard
			m, _ := bot.Send(msg)

			user.Step = 2
			user.MessageId = int64(m.MessageID)
			storage.UpdateUser(user)
		case 6:
			msg := tgbotapi.NewMessage(user.ChatId, "")
			msg.ParseMode = "html"

			questionId, err := strconv.Atoi(update.Message.Text)
			isQuestionExists := storage.IsQuestionExists(int64(questionId))

			if questionId == 0 || err != nil || !isQuestionExists {
				msg.Text = "Siz notog'ri indeks kiritdingiz. Iltimos qayta harakat qilib ko'ring."
				msg.ReplyMarkup = defaultMenuKeyboard
				user.Step = 2
			} else {
				user.QuestionId = int64(questionId)
				user.Step++
				msg.Text = "Endi esa savolni o'chirish yoki tahrirlashni tanlang."
				msg.ReplyMarkup = editMenuKeyboard
			}

			m, _ := bot.Send(msg)
			user.MessageId = int64(m.MessageID)
			storage.UpdateUser(user)
		case 7:
			msg := tgbotapi.NewMessage(user.ChatId, "")
			msg.ParseMode = "html"

			switch update.Message.Text {
			case "Savolni tahrirlash♻️":
				{
					user.Step++
					msg.Text = "Yangi savolni kiriting:"
					msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
				}
			case "Savolni o'chirish🗑":
				{
					user.Step = 2
					isExists := storage.IsExists(user.QuestionId, user.Id)
					status := storage.DeleteQuestion(user.QuestionId, user.Id)
					if status && isExists {
						msg.Text = strconv.Itoa(int(user.QuestionId)) + " indeksli savol muvaffaqiyatli o'chirildi!"
					} else {
						msg.Text = strconv.Itoa(int(user.QuestionId)) + " indeksli savolni o'chirishning imkoni bo'lmadi, qayta urinib ko'ring! (yuzaga kelishi mumkin bo'lgan muammolar: bunday indeksli savol mavjud bo'lmasligi mumkin yoki shu savolnining egasi boshqa.)"
					}
					msg.ReplyMarkup = defaultMenuKeyboard
				}
			}

			m, _ := bot.Send(msg)
			user.MessageId = int64(m.MessageID)

			storage.UpdateUser(user)
		case 8:
			user.Step++
			msg := tgbotapi.NewMessage(user.ChatId, "Savolingizning javobini kiriting:")
			msg.ParseMode = "html"
			var question model.Question
			question = storage.GetQuestion(user.QuestionId)
			question.Text = update.Message.Text
			question.Answer = ""
			question.Status = 0
			storage.UpdateQuestion(question)
			m, _ := bot.Send(msg)
			user.MessageId = int64(m.MessageID)
			storage.UpdateUser(user)
		case 9:
			var question model.Question
			question = storage.GetQuestion(user.QuestionId)
			question.Answer = update.Message.Text
			question.Status = 1
			storage.UpdateQuestion(question)

			msg := tgbotapi.NewMessage(user.ChatId, "Savol muvaffaqiyatli tahrirlandi✅")
			msg.ParseMode = "html"
			msg.ReplyMarkup = defaultMenuKeyboard
			m, _ := bot.Send(msg)

			user.Step = 2
			user.MessageId = int64(m.MessageID)
			storage.UpdateUser(user)
		default:
			msg := tgbotapi.NewMessage(user.ChatId, "Bunday buyruq yo'q🙃")
			msg.ParseMode = "html"
			m, _ := bot.Send(msg)
			user.Step = 2
			user.MessageId = int64(m.MessageID)
			storage.UpdateUser(user)
		}
	}
}
