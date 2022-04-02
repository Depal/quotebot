package bot

import (
	"github.com/tucnak/telebot"
)

func (s *Service) announceCommand(command string, message *telebot.Message) {
	sender := message.Sender
	s.log.Infof("[%s]: Called by %s (ID: %v, @%s)...", command, sender.FirstName, sender.ID, sender.Username)
}

func (s *Service) finishCommand(command string) {
	s.log.Debugf("[%s]: Ok", command)
}
