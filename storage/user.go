package storage

import (
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/shavkatjon/viktorina-bot/model"
	"github.com/shavkatjon/viktorina-bot/utils"
)

// game.db

func GameInsertUser(user model.GameUser) int64 {
	query := `
		INSERT INTO users 
		(
			"chat_id",
			"first_name",
			"last_name", 
			"username",
			"step", 
			"message_id",
			"score",
			"subject",
		 	"question",
		 	"answer"
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := gameDb.Exec(
		query,
		user.ChatId,
		user.FirstName,
		user.LastName,
		user.UserName,
		user.Step,
		user.MessageId,
		user.Score,
		user.Subject,
		user.Question,
		user.Answer,
	)
	utils.Check(err)

	var id int64

	query = `
		SELECT 
			"id" 
		FROM users 
		WHERE "chat_id" = $1
	`
	row := gameDb.QueryRow(query, user.ChatId)
	err = row.Scan(&id)

	utils.Check(err)

	return id
}

func GameUpdateUser(user model.GameUser) {
	query := `
		UPDATE users SET 
			"first_name" = $1,
			"last_name" = $2, 
			"username" = $3,
			"step" = $4, 
			"message_id" = $5,
			"score" = $6,
			"subject" = $7,
			"question" = $8,
			"answer" = $9
		WHERE "chat_id" = $10
	`

	_, err := gameDb.Exec(
		query,
		user.FirstName,
		user.LastName,
		user.UserName,
		user.Step,
		user.MessageId,
		user.Score,
		user.Subject,
		user.Question,
		user.Answer,
		user.ChatId,
	)

	utils.Check(err)
}

func GameGetUser(chatId int64) model.GameUser {
	query := `
		SELECT 
			"chat_id",
			"first_name",
			"last_name",
			"username",
			"step",
			"message_id",
			"score",
			"subject",
			"question",
			"answer"
		FROM users WHERE "chat_id" = $1
	`

	row := gameDb.QueryRow(query, chatId)

	var user model.GameUser

	err := row.Scan(
		&user.ChatId,
		&user.FirstName,
		&user.LastName,
		&user.UserName,
		&user.Step,
		&user.MessageId,
		&user.Score,
		&user.Subject,
		&user.Question,
		&user.Answer,
	)

	if err != nil {
		fmt.Println(err)

		user.ChatId = chatId
		user.Step = 1
		user.Id = GameInsertUser(user)
	}

	return user
}

func GameGetUserList(Limit int64) []model.GameUser {
	var uList []model.GameUser
	var user model.GameUser

	query := `
		SELECT
			"first_name",
			"last_name", 
			"score"
		FROM users ORDER BY "score" DESC LIMIT $1
	`
	rows, err := gameDb.Query(query, Limit)
	utils.Check(err)
	if err == nil {
		for rows.Next() {
			err := rows.Scan(
				&user.FirstName,
				&user.LastName,
				&user.Score,
			)
			utils.Check(err)
			uList = append(uList, user)
		}
	}

	return uList
}

// qa.db

func InsertUser(user model.User) int64 {
	query := `
		INSERT INTO users 
		(
			"chat_id",
			"first_name",
			"last_name",
			"username",
			"step",
			"message_id",
			"question_id",
		 	"subject"
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := questionDb.Exec(
		query,
		user.ChatId,
		user.FirstName,
		user.LastName,
		user.UserName,
		user.Step,
		user.MessageId,
		user.QuestionId,
		user.Subject,
	)
	utils.Check(err)

	var id int64

	query = `
		SELECT 
			"id" 
		FROM users 
		WHERE "chat_id" = $1
	`
	row := questionDb.QueryRow(query, user.ChatId)
	err = row.Scan(&id)

	utils.Check(err)

	return id
}

func UpdateUser(user model.User) {
	query := `
		UPDATE users SET 
			"first_name" = $1,
			"last_name" = $2, 
			"username" = $3,
			"step" = $4, 
			"message_id" = $5,
			"question_id" = $6,
			"subject" = $7
		WHERE "chat_id" = $8
	`

	_, err := questionDb.Exec(
		query,
		user.FirstName,
		user.LastName,
		user.UserName,
		user.Step,
		user.MessageId,
		user.QuestionId,
		user.Subject,
		user.ChatId,
	)

	utils.Check(err)
}

func GetUser(chatId int64) model.User {
	query := `
		SELECT 
			"id",
			"chat_id",
			"first_name",
			"last_name",
			"username",
			"step",
			"message_id",
			"question_id",
			"subject"
		FROM users WHERE "chat_id" = $1`

	row := questionDb.QueryRow(query, chatId)

	var user model.User

	err := row.Scan(
		&user.Id,
		&user.ChatId,
		&user.FirstName,
		&user.LastName,
		&user.UserName,
		&user.Step,
		&user.MessageId,
		&user.QuestionId,
		&user.Subject,
	)

	if err != nil {
		fmt.Println(err)

		user.ChatId = chatId
		user.Step = 1
		user.Id = InsertUser(user)
	}

	return user
}
