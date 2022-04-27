package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	"github.com/shavkatjon/viktorina-bot/model"
	"github.com/shavkatjon/viktorina-bot/storage"
	"github.com/shavkatjon/viktorina-bot/utils"
)

var subjectKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Matematika"),
		tgbotapi.NewKeyboardButton("Fizika"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Kimyo"),
		tgbotapi.NewKeyboardButton("Biologiya"),
	),
)

func main() {
	// Buttons
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

	// Connect to game database
	dbConn := storage.ConnectGameDb()
	if dbConn {
		fmt.Println("Connected to game.db")
	} else {
		fmt.Println("Not connected to game.db")
	}

	// Connect to history database
	dbConn = storage.ConnectHistoryDb()
	if dbConn {
		fmt.Println("Connected to history.db")
	} else {
		fmt.Println("Not connected to history.db")
	}

	go HistoryBot()

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_TOKEN2"))
	utils.Check(err)

	bot.Debug = true

	log.Printf("Authorized  on account %s", bot.Self.UserName)

	// Localniy run qilish uchun pastdagi 4 ta qatorni kommmentdan chiqarish kerak

	//_, err = bot.RemoveWebhook()
	//u := tgbotapi.NewUpdate(0)
	//u.Timeout = 60
	//updates, _ := bot.GetUpdatesChan(u)

	// Herokuda run qilish uchun pastdagi 8 ta qatorni kommentdan chiqarish kerak bo'ladi

	_, err = bot.RemoveWebhook()
	_, err = bot.SetWebhook(tgbotapi.NewWebhook("https://pmviktorinabot.herokuapp.com/" + bot.Token))
	utils.Check(err)
	updates := bot.ListenForWebhook("/" + bot.Token)
	go func() {
		err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
		utils.Check(err)
	}()

	for update := range updates {
		if update.CallbackQuery != nil {
			user := storage.GameGetUser(update.CallbackQuery.Message.Chat.ID)
			callbackQueryDataSplit := strings.Split(update.CallbackQuery.Data, "#")

			switch user.Step {
			case 4:
				_, err := bot.DeleteMessage(tgbotapi.DeleteMessageConfig{ChatID: user.ChatId, MessageID: int(user.MessageId)})
				utils.Check(err)

				switch callbackQueryDataSplit[0] {
				case "correct":
					congratsMsg := tgbotapi.NewMessage(user.ChatId, "")
					congratsMsg.ParseMode = "html"
					congratsMsg.ReplyMarkup = stop
					congratsMsg.Text = "Javobingiz to'g'ri😃"
					_, err = bot.Send(congratsMsg)
					utils.Check(err)
					user.Score++

					question := storage.GameGetQuestion(user.Subject)
					user.Question = question.Text
					user.Answer = question.Answer

					if question.Image != "" {
						response, _ := http.Get(question.Image)

						b, _ := io.ReadAll(response.Body)
						file := tgbotapi.FileBytes{
							Name:  "picture",
							Bytes: b,
						}

						msg := tgbotapi.NewPhotoUpload(user.ChatId, file)
						msg.ParseMode = "html"
						rand.Seed(time.Now().UnixNano())
						randomType := rand.Intn(2)

						if randomType == 1 {
							msg.Caption = "<b>Savol:</b> " + user.Question + "\n\n<b>Javobni kiriting:</b>"
						} else {
							variants := []string{question.Variant1, question.Variant2, question.Variant3}

							msg.Caption = "<b>Savol:</b> " + user.Question + "\n\n<b>Javobni tanlang:</b>"
							msg.ReplyMarkup = utils.BuildInline(&user, variants)
						}

						m, _ := bot.Send(msg)
						user.MessageId = int64(m.MessageID)

						err := response.Body.Close()
						utils.Check(err)
					} else {
						msg := tgbotapi.NewMessage(user.ChatId, "")
						msg.ParseMode = "html"
						rand.Seed(time.Now().UnixNano())
						randomType := rand.Intn(2)

						if randomType == 1 {
							msg.Text = "<b>Savol:</b> " + user.Question + "\n\n<b>Javobni kiriting:</b>"
						} else {
							variants := []string{question.Variant1, question.Variant2, question.Variant3}

							msg.Text = "<b>Savol:</b> " + user.Question + "\n\n<b>Javobni tanlang:</b>"
							msg.ReplyMarkup = utils.BuildInline(&user, variants)
						}

						m, _ := bot.Send(msg)
						user.MessageId = int64(m.MessageID)
					}
				case "wrong":
					congratsMsg := tgbotapi.NewMessage(user.ChatId, "")
					congratsMsg.ParseMode = "html"
					congratsMsg.ReplyMarkup = stop
					congratsMsg.Text = "Javobingiz noto'g'ri🙁\n\nTo'g'ri javob: <b>" + user.Answer + "</b>"
					_, err = bot.Send(congratsMsg)
					utils.Check(err)

					question := storage.GameGetQuestion(user.Subject)
					user.Question = question.Text
					user.Answer = question.Answer

					if question.Image != "" {
						response, _ := http.Get(question.Image)

						b, _ := io.ReadAll(response.Body)
						file := tgbotapi.FileBytes{
							Name:  "picture",
							Bytes: b,
						}

						msg := tgbotapi.NewPhotoUpload(user.ChatId, file)
						msg.ParseMode = "html"
						rand.Seed(time.Now().UnixNano())
						randomType := rand.Intn(2)

						if randomType == 1 {
							msg.Caption = "<b>Savol:</b> " + user.Question + "\n\n<b>Javobni kiriting:</b>"
						} else {
							variants := []string{question.Variant1, question.Variant2, question.Variant3}

							msg.Caption = "<b>Savol:</b> " + user.Question + "\n\n<b>Javobni tanlang:</b>"
							msg.ReplyMarkup = utils.BuildInline(&user, variants)
						}

						m, _ := bot.Send(msg)
						user.MessageId = int64(m.MessageID)

						err := response.Body.Close()
						utils.Check(err)
					} else {
						msg := tgbotapi.NewMessage(user.ChatId, "")
						msg.ParseMode = "html"
						rand.Seed(time.Now().UnixNano())
						randomType := rand.Intn(2)

						if randomType == 1 {
							msg.Text = "<b>Savol:</b> " + user.Question + "\n\n<b>Javobni kiriting:</b>"
						} else {
							variants := []string{question.Variant1, question.Variant2, question.Variant3}

							msg.Text = "<b>Savol:</b> " + user.Question + "\n\n<b>Javobni tanlang:</b>"
							msg.ReplyMarkup = utils.BuildInline(&user, variants)
						}

						m, _ := bot.Send(msg)
						user.MessageId = int64(m.MessageID)
					}
				}
			}
			storage.GameUpdateUser(user)
		}

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

			msgToAdmin := tgbotapi.NewMessage(738151092, "")
			msgToAdmin.ParseMode = "html"
			msgToAdmin.Text = "Viktorina botga kirgan foydalanuvchi: " + user.FirstName + " " + user.LastName + "\nusername: @" + user.UserName
			_, err := bot.Send(msgToAdmin)
			utils.Check(err)

			msg := tgbotapi.NewMessage(user.ChatId,
				"Assalomu Alaykum, "+user.FirstName+" "+user.LastName+"!\nViktorina botga xush kelibsiz!")

			msg.ParseMode = "html"
			msg.ReplyMarkup = groupMenuKeyboard

			m, _ := bot.Send(msg)
			user.MessageId = int64(m.MessageID)
		case 2:
			switch update.Message.Text {
			case "Profil📊":
				msg := tgbotapi.NewMessage(user.ChatId,
					"<i><b>Sizning profilingiz</b></i>:\n<b>Ism:</b> "+user.FirstName+"\n<b>Familiya:</b> "+user.LastName+"\n<b>Siz to'plagan ballar:</b> "+strconv.Itoa(user.Score))
				msg.ParseMode = "html"
				msg.ReplyMarkup = groupMenuKeyboard

				m, _ := bot.Send(msg)
				user.MessageId = int64(m.MessageID)
			case "Reyting🔝":
				userList := storage.GameGetUserList(10)

				msg := tgbotapi.NewMessage(user.ChatId, utils.GameListToText(userList))
				msg.ParseMode = "html"
				msg.ReplyMarkup = groupMenuKeyboard

				m, _ := bot.Send(msg)
				user.MessageId = int64(m.MessageID)
			case "Viktorinani boshlash🧩":
				user.Step++

				msg := tgbotapi.NewMessage(user.ChatId, "Qaysi fandan viktorinani boshlamoqchisiz?")
				msg.ParseMode = "html"
				msg.ReplyMarkup = subjectKeyboard

				m, _ := bot.Send(msg)
				user.MessageId = int64(m.MessageID)
			default:
				msg := tgbotapi.NewMessage(user.ChatId, "Bunday buyruq yo'q🙃")
				msg.ParseMode = "html"
				msg.ReplyMarkup = groupMenuKeyboard

				m, _ := bot.Send(msg)
				user.MessageId = int64(m.MessageID)
			}
		case 3:
			switch update.Message.Text {
			case "Matematika":
				user.Subject = "math"
			case "Fizika":
				user.Subject = "physics"
			case "Kimyo":
				user.Subject = "chemistry"
			case "Biologiya":
				user.Subject = "biology"
			default:
				user.Step = 2
				msg := tgbotapi.NewMessage(user.ChatId, "Bunday fan yo'q🙃")
				msg.ParseMode = "html"
				msg.ReplyMarkup = groupMenuKeyboard
				m, _ := bot.Send(msg)
				user.MessageId = int64(m.MessageID)
				storage.GameUpdateUser(user)
				continue
			}

			question := storage.GameGetQuestion(user.Subject)

			if storage.GameGetNumberOfQuestions(user.Subject) == 0 {
				user.Step = 2
				msg := tgbotapi.NewMessage(user.ChatId, "Bu fan bo'yicha savollar hali qo'shilmagan.")
				msg.ParseMode = "html"
				msg.ReplyMarkup = groupMenuKeyboard
				m, _ := bot.Send(msg)
				user.MessageId = int64(m.MessageID)
				storage.GameUpdateUser(user)
				continue
			}

			user.Question = question.Text
			user.Answer = question.Answer

			if question.Image != "" {
				response, _ := http.Get(question.Image)

				b, _ := io.ReadAll(response.Body)
				file := tgbotapi.FileBytes{
					Name:  "picture",
					Bytes: b,
				}

				msg := tgbotapi.NewPhotoUpload(user.ChatId, file)
				msg.ParseMode = "html"
				startMsg := tgbotapi.NewMessage(user.ChatId, "<b>Viktorina boshlandi!!!</b>")
				startMsg.ParseMode = "html"
				startMsg.ReplyMarkup = stop
				_, err = bot.Send(startMsg)
				utils.Check(err)

				rand.Seed(time.Now().UnixNano())
				randomType := rand.Intn(2)

				if randomType == 1 {
					msg.Caption = "<b>Savol:</b> " + user.Question + "\n\n<b>Javobni kiriting:</b>"
				} else {
					variants := []string{question.Variant1, question.Variant2, question.Variant3}

					msg.Caption = "<b>Savol:</b> " + user.Question + "\n\n<b>Javobni tanlang:</b>"
					msg.ReplyMarkup = utils.BuildInline(&user, variants)
				}

				user.Step++

				m, _ := bot.Send(msg)
				user.MessageId = int64(m.MessageID)

				err := response.Body.Close()
				utils.Check(err)
			} else {
				msg := tgbotapi.NewMessage(user.ChatId, "")
				msg.ParseMode = "html"

				startMsg := tgbotapi.NewMessage(user.ChatId, "<b>Viktorina boshlandi!!!</b>")
				startMsg.ParseMode = "html"
				startMsg.ReplyMarkup = stop
				_, err = bot.Send(startMsg)
				utils.Check(err)

				rand.Seed(time.Now().UnixNano())
				randomType := rand.Intn(2)

				if randomType == 1 {
					msg.Text = "<b>Savol:</b> " + user.Question + "\n\n<b>Javobni kiriting:</b>"
				} else {
					variants := []string{question.Variant1, question.Variant2, question.Variant3}

					msg.Text = "<b>Savol:</b> " + user.Question + "\n\n<b>Javobni tanlang:</b>"
					msg.ReplyMarkup = utils.BuildInline(&user, variants)
				}

				user.Step++

				m, _ := bot.Send(msg)
				user.MessageId = int64(m.MessageID)
			}
		case 4:
			switch strings.ToLower(update.Message.Text) {
			case "viktorinani to'xtatish⏹️":
				user.Step = 2

				msg := tgbotapi.NewMessage(user.ChatId, "Viktorina to'xtatildi")
				msg.ReplyMarkup = groupMenuKeyboard

				m, _ := bot.Send(msg)
				user.MessageId = int64(m.MessageID)
			default:
				_, err := bot.DeleteMessage(tgbotapi.DeleteMessageConfig{ChatID: user.ChatId, MessageID: int(user.MessageId)})
				utils.Check(err)

				if strings.EqualFold(strings.ToLower(update.Message.Text), strings.ToLower(user.Answer)) {
					user.Score++
					congratsMsg := tgbotapi.NewMessage(user.ChatId, "Javobingiz to'g'ri😃")
					congratsMsg.ParseMode = "html"
					congratsMsg.ReplyMarkup = stop
					_, err = bot.Send(congratsMsg)
					utils.Check(err)
				} else {
					congratsMsg := tgbotapi.NewMessage(user.ChatId, "Javobingiz noto'g'ri🙁\n\nTo'g'ri javob: <b>"+user.Answer+"</b>")
					congratsMsg.ParseMode = "html"
					congratsMsg.ReplyMarkup = stop
					_, err = bot.Send(congratsMsg)
					utils.Check(err)
				}

				question := storage.GameGetQuestion(user.Subject)
				user.Question = question.Text
				user.Answer = question.Answer

				if question.Image != "" {
					response, _ := http.Get(question.Image)

					b, _ := io.ReadAll(response.Body)
					file := tgbotapi.FileBytes{
						Name:  "picture",
						Bytes: b,
					}

					msg := tgbotapi.NewPhotoUpload(user.ChatId, file)
					msg.ParseMode = "html"
					rand.Seed(time.Now().UnixNano())
					randomType := rand.Intn(2)

					if randomType == 1 {
						msg.Caption = "<b>Savol:</b> " + user.Question + "\n\n<b>Javobni kiriting:</b>"
					} else {
						variants := []string{question.Variant1, question.Variant2, question.Variant3}

						msg.Caption = "<b>Savol:</b> " + user.Question + "\n\n<b>Javobni tanlang:</b>"
						msg.ReplyMarkup = utils.BuildInline(&user, variants)
					}

					m, _ := bot.Send(msg)
					user.MessageId = int64(m.MessageID)

					err := response.Body.Close()
					utils.Check(err)
				} else {
					msg := tgbotapi.NewMessage(user.ChatId, "")
					msg.ParseMode = "html"
					rand.Seed(time.Now().UnixNano())
					randomType := rand.Intn(2)

					if randomType == 1 {
						msg.Text = "<b>Savol:</b> " + user.Question + "\n\n<b>Javobni kiriting:</b>"
					} else {
						variants := []string{question.Variant1, question.Variant2, question.Variant3}

						msg.Text = "<b>Savol:</b> " + user.Question + "\n\n<b>Javobni tanlang:</b>"
						msg.ReplyMarkup = utils.BuildInline(&user, variants)
					}

					m, _ := bot.Send(msg)
					user.MessageId = int64(m.MessageID)
				}
			}
		default:
			user.Step = 2

			msg := tgbotapi.NewMessage(user.ChatId, "Bunday buyruq yo'q🙃")
			msg.ReplyMarkup = groupMenuKeyboard

			m, _ := bot.Send(msg)
			user.MessageId = int64(m.MessageID)
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
			tgbotapi.NewKeyboardButton("Savolni tahrirlash♻"),
			tgbotapi.NewKeyboardButton("Savolni o'chirish🗑"),
		),
	)

	e := godotenv.Load()
	utils.Check(e)

	// Initialize bot
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_TOKEN"))
	utils.Check(err)

	bot.Debug = true
	log.Printf("Authorized  on account %s", bot.Self.UserName)

	// Localniy run qilish uchun pastdagi 4 ta qatorni kommmentdan chiqarish kerak

	//_, err = bot.RemoveWebhook()
	//u := tgbotapi.NewUpdate(0)
	//u.Timeout = 60
	//updates, _ := bot.GetUpdatesChan(u)

	// Herokuda run qilish uchun pastdagi 8 ta qatorni kommentdan chiqarish kerak

	_, err = bot.RemoveWebhook()
	_, err = bot.SetWebhook(tgbotapi.NewWebhook("https://pmviktorinabot.herokuapp.com/" + bot.Token))
	utils.Check(err)
	updates := bot.ListenForWebhook("/" + bot.Token)
	go func() {
		err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
		utils.Check(err)
	}()

	for update := range updates {
		if update.CallbackQuery != nil {
			user := storage.GetUser(int64(update.CallbackQuery.From.ID))
			callbackQueryDataSplit := strings.Split(update.CallbackQuery.Data, "#")

			switch user.Step {
			case 2:
				switch callbackQueryDataSplit[0] {
				case "next":
					{
						var inlineKeyboard tgbotapi.InlineKeyboardMarkup

						msg := tgbotapi.NewEditMessageText(user.ChatId, int(user.MessageId), "")
						msg.ParseMode = "html"

						limit, err := strconv.Atoi(callbackQueryDataSplit[1])
						utils.Check(err)
						page, err := strconv.Atoi(callbackQueryDataSplit[2])
						utils.Check(err)
						questionList := storage.GetQuestionList(user.Subject, int64(limit), int64(page))

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
						_, err = bot.Send(msg)
						utils.Check(err)
					}
				case "prev":
					{
						var inlineKeyboard tgbotapi.InlineKeyboardMarkup

						msg := tgbotapi.NewEditMessageText(user.ChatId, int(user.MessageId), "")
						msg.ParseMode = "html"

						limit, err := strconv.Atoi(callbackQueryDataSplit[1])
						utils.Check(err)
						page, err := strconv.Atoi(callbackQueryDataSplit[2])
						utils.Check(err)
						questionList := storage.GetQuestionList(user.Subject, int64(limit), int64(page))

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
						_, err = bot.Send(msg)
						utils.Check(err)
					}
				}
			case 5:
				_, err := bot.DeleteMessage(tgbotapi.DeleteMessageConfig{ChatID: user.ChatId, MessageID: int(user.MessageId)})
				utils.Check(err)

				user.Step++
				msg := tgbotapi.NewMessage(user.ChatId, "")
				msg.ParseMode = "html"
				msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)

				switch callbackQueryDataSplit[0] {
				case "yes":
					msg.Text = "Savolingiz uchun rasm jo'nating:"
				default:
					user.Step = 7
					msg.Text = "Savolning tog'ri javobini kiriting:"
				}

				m, _ := bot.Send(msg)
				user.MessageId = int64(m.MessageID)
			}
			storage.UpdateUser(user)
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

			msgToAdmin := tgbotapi.NewMessage(738151092, "")
			msgToAdmin.ParseMode = "html"
			msgToAdmin.Text = "Admin botga kirgan foydalanuvchi:  " + user.FirstName + " " + user.LastName + "\nusername: @" + user.UserName
			_, err := bot.Send(msgToAdmin)
			utils.Check(err)

			msg := tgbotapi.NewMessage(user.ChatId,
				"Assalomu Alaykum, "+user.FirstName+" "+user.LastName+"!\nViktorina uchun savollar tuzadigan admin botga xush kelibsiz!")

			msg.ParseMode = "html"
			msg.ReplyMarkup = defaultMenuKeyboard
			m, _ := bot.Send(msg)

			user.MessageId = int64(m.MessageID)
		case 2:
			var m tgbotapi.Message
			msg := tgbotapi.NewMessage(user.ChatId, "")
			msg.ParseMode = "html"

			switch update.Message.Text {
			case "Savollar to'plami📋":
				{
					user.Step = 100
					msg.Text = "Fanni tanlang:"
					msg.ReplyMarkup = subjectKeyboard
					m, _ = bot.Send(msg)
				}
			case "Savol qo'shish📥":
				{
					user.Step++
					msg.Text = "Kiritmoqchi bo'lgan savolning fanini tanlang:"
					msg.ReplyMarkup = subjectKeyboard
					m, _ = bot.Send(msg)
				}
			case "Savolni tahrirlash📝":
				{
					user.Step = 11
					msg.Text = "Tahrirlamoqchi yoki o'chirmoqchi bo'lgan savolning indeksini kiriting:"
					msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
					m, _ = bot.Send(msg)
				}
			case "/start":
				{
					user.FirstName = update.Message.Chat.FirstName
					user.LastName = update.Message.Chat.LastName
					user.UserName = update.Message.Chat.UserName
					user.Step = 2

					msg = tgbotapi.NewMessage(user.ChatId,
						"Assalomu Alaykum, "+user.FirstName+" "+user.LastName+"!\nViktorina uchun savollar tuzadigan admin botga xush kelibsiz!")
					msg.ReplyMarkup = defaultMenuKeyboard
					m, _ = bot.Send(msg)
				}
			default:
				{
					msg = tgbotapi.NewMessage(user.ChatId, "Bunday buyruq yo'q🙃")
					msg.ReplyMarkup = defaultMenuKeyboard
					m, _ = bot.Send(msg)
				}
			}

			user.MessageId = int64(m.MessageID)
		case 3:
			user.Step++
			msg := tgbotapi.NewMessage(user.ChatId, "Yangi savolni kiriting:")
			msg.ParseMode = "html"
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			var question model.Question
			switch update.Message.Text {
			case "Matematika":
				question.Subject = "math"
			case "Fizika":
				question.Subject = "physics"
			case "Kimyo":
				question.Subject = "chemistry"
			case "Biologiya":
				question.Subject = "biology"
			default:
				msg.Text = "Bunday fan yo'q🙃"
				msg.ReplyMarkup = defaultMenuKeyboard
				m, _ := bot.Send(msg)
				user.Step = 2
				user.MessageId = int64(m.MessageID)
				storage.UpdateUser(user)
				continue
			}
			question.UserId = user.Id
			user.QuestionId = storage.InsertQuestion(question)

			m, _ := bot.Send(msg)
			user.MessageId = int64(m.MessageID)
		case 4:
			user.Step++
			var question model.Question
			question = storage.GetQuestion(user.QuestionId)
			question.Text = update.Message.Text
			question.UserId = user.Id
			storage.UpdateQuestion(question)

			msg := tgbotapi.NewMessage(user.ChatId, "Savolda rasm qatnashganmi?")
			inlineBTN := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("Ha", "yes"),
					tgbotapi.NewInlineKeyboardButtonData("Yo'q", "no"),
				),
			)
			msg.ReplyMarkup = inlineBTN
			msg.ParseMode = "html"
			m, _ := bot.Send(msg)
			user.MessageId = int64(m.MessageID)
		case 5:
			_, err := bot.DeleteMessage(tgbotapi.DeleteMessageConfig{ChatID: user.ChatId, MessageID: int(user.MessageId)})
			utils.Check(err)

			msg := tgbotapi.NewMessage(user.ChatId, "Savolda rasm qatnashganmi?")
			inlineBTN := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("Ha", "yes"),
					tgbotapi.NewInlineKeyboardButtonData("Yo'q", "no"),
				),
			)
			msg.ReplyMarkup = inlineBTN
			msg.ParseMode = "html"
			m, _ := bot.Send(msg)
			user.MessageId = int64(m.MessageID)
		case 6:
			user.Step++
			var question model.Question
			question = storage.GetQuestion(user.QuestionId)
			if update.Message.Photo != nil {
				url, err := bot.GetFileDirectURL((*update.Message.Photo)[len(*update.Message.Photo)-1].FileID)
				if err == nil {
					question.Image = url
				} else {
					panic(err)
				}
			}
			storage.UpdateQuestion(question)

			msg := tgbotapi.NewMessage(user.ChatId, "Savolning tog'ri javobini kiriting:")
			msg.ParseMode = "html"
			m, _ := bot.Send(msg)
			user.MessageId = int64(m.MessageID)
		case 7:
			user.Step++
			var question model.Question
			question = storage.GetQuestion(user.QuestionId)
			question.Answer = update.Message.Text
			storage.UpdateQuestion(question)

			msg := tgbotapi.NewMessage(user.ChatId, "Javobning birinchi noto'g'ri variantini kitiring:")
			msg.ParseMode = "html"
			m, _ := bot.Send(msg)
			user.MessageId = int64(m.MessageID)
		case 8:
			user.Step++
			var question model.Question
			question = storage.GetQuestion(user.QuestionId)
			question.Variant1 = update.Message.Text
			storage.UpdateQuestion(question)

			msg := tgbotapi.NewMessage(user.ChatId, "Javobning ikkinchi noto'g'ri variantini kitiring:")
			msg.ParseMode = "html"
			m, _ := bot.Send(msg)
			user.MessageId = int64(m.MessageID)
		case 9:
			user.Step++
			var question model.Question
			question = storage.GetQuestion(user.QuestionId)
			question.Variant2 = update.Message.Text
			storage.UpdateQuestion(question)

			msg := tgbotapi.NewMessage(user.ChatId, "Javobning uchinchi noto'g'ri variantini kitiring:")
			msg.ParseMode = "html"
			m, _ := bot.Send(msg)
			user.MessageId = int64(m.MessageID)
		case 10:
			var question model.Question
			question = storage.GetQuestion(user.QuestionId)
			question.Variant3 = update.Message.Text
			question.Status = 1
			storage.UpdateQuestion(question)

			msg := tgbotapi.NewMessage(user.ChatId, "Savol muvaffaqiyatli qo'shildi✅")
			msg.ParseMode = "html"
			msg.ReplyMarkup = defaultMenuKeyboard
			m, _ := bot.Send(msg)
			user.MessageId = int64(m.MessageID)
			user.Step = 2
		case 11:
			questionId, err := strconv.Atoi(update.Message.Text)
			isQuestionExists := storage.IsQuestionExists(int64(questionId))

			msg := tgbotapi.NewMessage(user.ChatId, "")
			msg.ParseMode = "html"

			if questionId == 0 || err != nil || !isQuestionExists {
				user.Step = 2
				msg.Text = "Siz notog'ri indeks kiritdingiz. Iltimos qayta harakat qilib ko'ring."
				msg.ReplyMarkup = defaultMenuKeyboard
			} else {
				user.Step++
				user.QuestionId = int64(questionId)
				msg.Text = "Endi esa savolni o'chirish yoki tahrirlashni tanlang."
				msg.ReplyMarkup = editMenuKeyboard
			}

			m, _ := bot.Send(msg)
			user.MessageId = int64(m.MessageID)
		case 12:
			msg := tgbotapi.NewMessage(user.ChatId, "")
			msg.ParseMode = "html"

			switch update.Message.Text {
			case "Savolni tahrirlash♻":
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
						msg.Text = strconv.Itoa(int(user.QuestionId)) + " indeksli savolni o'chirishning imkoni bo'lmadi, qayta urinib ko'ring! (yuzaga kelishi mumkin bo'lgan muammolar: bunday indeksli savol mavjud bo'lmasligi mumkin yoki shu savolnining muallifi siz emassiz.)"
					}
					msg.ReplyMarkup = defaultMenuKeyboard
				}
			default:
				user.Step = 2
				msg.Text = "Bunday buyruq yo'q🙃"
				msg.ReplyMarkup = defaultMenuKeyboard
				m, _ := bot.Send(msg)
				user.MessageId = int64(m.MessageID)
				storage.UpdateUser(user)
				continue
			}

			m, _ := bot.Send(msg)
			user.MessageId = int64(m.MessageID)
		case 13:
			user.Step++
			var question model.Question
			question = storage.GetQuestion(user.QuestionId)
			question.Text = update.Message.Text
			storage.UpdateQuestion(question)

			msg := tgbotapi.NewMessage(user.ChatId, "Savolning to'g'ri javobini kiriting:")
			msg.ParseMode = "html"
			msg.ReplyMarkup = defaultMenuKeyboard
			m, _ := bot.Send(msg)
			user.MessageId = int64(m.MessageID)
		case 14:
			user.Step = 2
			var question model.Question
			question = storage.GetQuestion(user.QuestionId)
			question.Answer = update.Message.Text
			question.Status = 1
			storage.UpdateQuestion(question)

			msg := tgbotapi.NewMessage(user.ChatId, "Savol muvaffaqiyatli tahrirlandi✅")
			msg.ParseMode = "html"
			msg.ReplyMarkup = defaultMenuKeyboard
			m, _ := bot.Send(msg)
			user.MessageId = int64(m.MessageID)
		case 100:
			user.Step = 2
			msg := tgbotapi.NewMessage(user.ChatId, "")
			msg.ParseMode = "html"
			msg.ReplyMarkup = defaultMenuKeyboard

			switch update.Message.Text {
			case "Matematika":
				user.Subject = "math"
			case "Fizika":
				user.Subject = "physics"
			case "Kimyo":
				user.Subject = "chemistry"
			case "Biologiya":
				user.Subject = "biology"
			default:
				msg.Text = "Bunday fan yo'q🙃"
				msg.ReplyMarkup = defaultMenuKeyboard
				m, _ := bot.Send(msg)
				user.MessageId = int64(m.MessageID)
				storage.UpdateUser(user)
				continue
			}

			var inlineKeyboard tgbotapi.InlineKeyboardMarkup

			questionList := storage.GetQuestionList(user.Subject, 10, 1)

			if len(questionList.List) == 0 {
				msg.Text = "Bu fan bo'yicha savollar hali qo'shilmagan!"
				msg.ReplyMarkup = defaultMenuKeyboard
				m, _ := bot.Send(msg)
				user.MessageId = int64(m.MessageID)
				storage.UpdateUser(user)
				continue
			}

			// this is for default keyborad
			{
				msg.Text = user.Subject
				msg.ReplyMarkup = defaultMenuKeyboard
				_, err = bot.Send(msg)
				utils.Check(err)
			}

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
			m, _ := bot.Send(msg)
			user.MessageId = int64(m.MessageID)
		default:
			user.Step = 2
			msg := tgbotapi.NewMessage(user.ChatId, "Bunday buyruq yo'q🙃")
			msg.ParseMode = "html"
			m, _ := bot.Send(msg)
			user.MessageId = int64(m.MessageID)
		}
		storage.UpdateUser(user)
	}
}
