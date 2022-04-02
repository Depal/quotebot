package bot

import (
	"github.com/Depal/quotebot/internal/bot/static"
	"github.com/Depal/quotebot/internal/rating"
	"github.com/Depal/quotebot/pkg/logger"
	"github.com/tucnak/telebot"
	"time"
)

type Service struct {
	log    logger.ILogger
	Rating *rating.Service
	bot    telebot.Bot
}

func Initialize(log logger.ILogger, rating *rating.Service) *Service {
	return &Service{
		log:    log,
		Rating: rating,
	}
}

func (s *Service) Start() {
	s.bot = telebot.Bot{
		Token:  static.EnvBotToken,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	s.initializeHandlers()

	go s.bot.Start()
	s.log.Info("Bot started")
}
