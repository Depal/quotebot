package bot

import (
	"fmt"
	"github.com/Depal/quotebot/internal/bot/static"
	"github.com/tucnak/telebot"
)

func (s *Service) handleQuote(message *telebot.Message) {
	s.announceCommand(static.CommandQuote, message)

	quoted := message.ReplyTo
	if quoted == nil {
		_, err := s.bot.Reply(message, "Please quote a message first")
		if err != nil {
			s.log.Error(err)
			return
		}
		return
	}

	response := fmt.Sprintf("Will create a sticker of message \"%s\" from %s", quoted.Text, quoted.Sender.FirstName)

	_, err := s.bot.Reply(message, response)
	if err != nil {
		s.log.Error(err)
		return
	}

	s.finishCommand(static.CommandQuote)
}
