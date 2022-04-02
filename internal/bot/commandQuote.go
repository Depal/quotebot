package bot

import (
	"fmt"
	"github.com/Depal/quotebot/internal/bot/static"
	"github.com/tucnak/telebot"
)

func (s *Service) handleQuote(message *telebot.Message) {
	s.announceCommand(static.CommandQuote, message)

	quoted := message.ReplyTo

	response := fmt.Sprintf("Will create sticker of message \"%s\" from %s", quoted.Text, quoted.Sender.FirstName)

	_, err := s.bot.Reply(message, response)
	if err != nil {
		s.log.Error(err)
		return
	}

	s.finishCommand(static.CommandQuote)
}
