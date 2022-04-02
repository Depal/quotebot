package bot

import (
	"github.com/Depal/quotebot/internal/bot/static"
	"github.com/tucnak/telebot"
)

func (s *Service) handleAlert(message *telebot.Message) {
	s.announceCommand(static.CommandAlert, message)

	_, err := s.bot.Send(message.Chat, "@Depal, сделай меня")
	if err != nil {
		s.log.Error(err)
		return
	}

	s.finishCommand(static.CommandAlert)
}
