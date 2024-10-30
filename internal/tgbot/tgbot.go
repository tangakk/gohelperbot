package tgbot

import (
	"bytes"
	"context"
	"gohelperbot/pkg/questions"
	"os"
	"os/signal"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Tgbot struct {
	Config
	Q questions.Questions
}

func New(config Config) *Tgbot {
	return &Tgbot{Config: config}
}

func (tg *Tgbot) Run() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(tg.defaultHandler),
		bot.WithCallbackQueryDataHandler("button", bot.MatchTypePrefix, tg.callbackHandler),
	}

	b, err := bot.New(tg.Key, opts...)
	if nil != err {
		panic(err)
	}

	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, tg.startHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/help", bot.MatchTypeExact, tg.startHandler)

	b.Start(ctx)
}

func (tg *Tgbot) defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	txt := update.Message.Text
	q := tg.Q.Ask(txt)
	var kb models.InlineKeyboardMarkup
	kb.InlineKeyboard = make([][]models.InlineKeyboardButton, len(q.Subquestions))

	for i, v := range q.Subquestions {
		kb.InlineKeyboard[i] = []models.InlineKeyboardButton{{Text: v, CallbackData: "button " + v}}
	}

	text_message := "*" + q.Text + "*\n\n" + q.Answer

	ok := false
	if q.Extra != "" {
		err := tg.sendAttachment(ctx, b, q.Extra, update.Message.Chat.ID, text_message, kb)
		if err == nil {
			ok = true
		}
	}

	if !ok {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.Message.Chat.ID,
			Text:        text_message,
			ReplyMarkup: kb,
			ParseMode:   models.ParseModeMarkdown,
		})
	}
}

func (tg *Tgbot) callbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})

	q := tg.Q.Ask(update.CallbackQuery.Data)
	var kb models.InlineKeyboardMarkup
	kb.InlineKeyboard = make([][]models.InlineKeyboardButton, len(q.Subquestions))

	for i, v := range q.Subquestions {
		kb.InlineKeyboard[i] = []models.InlineKeyboardButton{{Text: v, CallbackData: "button " + v}}
	}

	text_message := "*" + q.Text + "*\n\n" + q.Answer

	ok := false
	if q.Extra != "" {
		err := tg.sendAttachment(ctx, b, q.Extra, update.CallbackQuery.Message.Message.Chat.ID,
			text_message, kb)
		if err == nil {
			ok = true
		}
	}

	if !ok {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
			Text:        text_message,
			ReplyMarkup: kb,
			ParseMode:   models.ParseModeMarkdown,
		})
	}

	b.EditMessageReplyMarkup(ctx, &bot.EditMessageReplyMarkupParams{
		ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
		MessageID:   update.CallbackQuery.Message.Message.ID,
		ReplyMarkup: nil,
	})
}

func (tg *Tgbot) sendAttachment(ctx context.Context, b *bot.Bot, path string, chatID int64, text string, kb models.InlineKeyboardMarkup) error {
	fileData, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	if strings.Contains(path, "mp4") {
		_, err := b.SendVideo(ctx, &bot.SendVideoParams{
			ChatID:      chatID,
			Video:       &models.InputFileUpload{Filename: "video.mp4", Data: bytes.NewReader(fileData)},
			Caption:     text,
			ReplyMarkup: kb,
			ParseMode:   models.ParseModeMarkdown,
		})
		if err != nil {
			return err
		}
	} else {
		_, err := b.SendPhoto(ctx, &bot.SendPhotoParams{
			ChatID:      chatID,
			Photo:       &models.InputFileUpload{Filename: "photo.png", Data: bytes.NewReader(fileData)},
			Caption:     text,
			ReplyMarkup: kb,
			ParseMode:   models.ParseModeMarkdown,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (tg *Tgbot) startHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	var kb models.InlineKeyboardMarkup
	kb.InlineKeyboard = make([][]models.InlineKeyboardButton, len(tg.StartButtons))

	for i, v := range tg.StartButtons {
		kb.InlineKeyboard[i] = []models.InlineKeyboardButton{{Text: v, CallbackData: "button " + v}}
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        "Бот может отвечать на вопросы\\. Задайте вопрос или нажмите на кнопки",
		ParseMode:   models.ParseModeMarkdown,
		ReplyMarkup: kb,
	})
}
