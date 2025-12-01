package handler

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"main/serp"
)

type handler struct {
	bot        *tgbotapi.BotAPI
	serpClient *serp.SerpClient
}

func New(bot *tgbotapi.BotAPI, serpClient *serp.SerpClient) *handler {
	return &handler{
		bot:        bot,
		serpClient: serpClient,
	}
}

func (h *handler) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := h.bot.GetUpdatesChan(u)

	for update := range updates {
		h.handleUpdate(update)
	}
}

func (h *handler) handleUpdate(update tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	if update.Message.IsCommand() {
		h.handleCommand(update)
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите команду. Например: /help")
	h.bot.Send(msg)
}

func (h *handler) handleCommand(update tgbotapi.Update) {
	cmd := update.Message.Command()
	args := update.Message.CommandArguments()
	chatID := update.Message.Chat.ID

	switch cmd {

	case "start":
		h.bot.Send(tgbotapi.NewMessage(chatID,
			"Привет! Я бот для проверки ссылок через Google Search API.\n"+
				"Используй /check <ссылка>"))

	case "help":
		h.bot.Send(tgbotapi.NewMessage(chatID,
			"Команды:\n/check https://site.com — проверить ссылку"))

	case "check":
		if args == "" {
			h.bot.Send(tgbotapi.NewMessage(chatID, "Использование: /check <ссылка>"))
			return
		}

		url := strings.TrimSpace(args)

		isPhishing, flag, err := h.serpClient.CheckURL(url)
		if err != nil {
			h.bot.Send(tgbotapi.NewMessage(chatID, err.Error()))
			return
		}

		if isPhishing {
			h.bot.Send(tgbotapi.NewMessage(chatID,
				fmt.Sprintf("Найден признак опасности: %s.\nСсылка может быть не безопасной!", flag)))
		} else {
			h.bot.Send(tgbotapi.NewMessage(chatID, "Ссылка выглядит безопасной."))
		}

	default:
		h.bot.Send(tgbotapi.NewMessage(chatID, "Неизвестная команда."))
	}
}
