package render

import (
	"github.com/Depal/quotebot/pkg/logger"
	"github.com/tucnak/telebot"
)

type Service struct {
	log logger.ILogger
}

func Initialize(log logger.ILogger) *Service {
	return &Service{log: log}
}

func (s *Service) CreateSingleSticker(image telebot.Photo, username string, text string) (err error) {
	panic("implement me") // TODO
}
