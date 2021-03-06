package utils

import (
	"github.com/Dreamacro/clash/log"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/thank243/StairUnlocker-Go/config"
)

type TgBot struct {
	SendMessage string
	Check       bool
	Bot         *tg.BotAPI
}

var Cfg *config.SuConfig

func (tb *TgBot) NewBot(cfg *config.SuConfig) {
	Cfg = cfg
	bot, err := tg.NewBotAPI(cfg.Telegram.TelegramToken)
	if err != nil {
		panic(err)
	}
	if cfg.LogLevel == 0 {
		bot.Debug = true
	}
	tb.Bot = bot
	log.Infoln("Authorized on account %s", bot.Self.UserName)
}

func (tb *TgBot) TelegramUpdates(buf *chan bool) {
	bot := tb.Bot
	u := tg.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}
		if update.Message.Chat.ID == Cfg.Telegram.ID {
			switch update.Message.Text {
			case "/start":
				_, _ = bot.Send(tg.NewMessage(update.Message.Chat.ID, "/check Check all node.\n/stat Show last status."))
			case "/check":
				if !tb.Check {
					tb.Check = true
					*buf <- true
					_, _ = bot.Send(tg.NewMessage(update.Message.Chat.ID, "Checking all nodes..."))
				} else {
					_, _ = bot.Send(tg.NewMessage(update.Message.Chat.ID, "Duplication, Checking all nodes..."))
				}
			case "/stat":
				_, _ = bot.Send(tg.NewMessage(update.Message.Chat.ID, tb.SendMessage))
			default:
				_, _ = bot.Send(tg.NewMessage(update.Message.Chat.ID, "Invalid command"))
			}
		}
		log.Infoln("TelegramUpdates Bot: [ID: %d], Text: %s", update.Message.From.ID, update.Message.Text)
	}
}
