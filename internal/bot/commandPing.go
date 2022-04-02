package bot

import (
	"github.com/Depal/quotebot/internal/bot/static"
	"github.com/tucnak/telebot"
)

func (s *Service) handlePing(message *telebot.Message) {
	s.announceCommand(static.CommandPing, message)

	_, err := s.bot.Reply(message, "pong")
	if err != nil {
		s.log.Error(err)
		return
	}

	s.finishCommand(static.CommandPing)
}
