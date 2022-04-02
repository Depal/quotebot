package bot

import (
	"github.com/Depal/quotebot/internal/bot/static"
	"github.com/Depal/quotebot/internal/rating"
	"github.com/Depal/quotebot/pkg/logger"
	"github.com/tucnak/telebot"
	"os"
	"time"
)

type Service struct {
	log    logger.ILogger
	Rating *rating.Service
	bot    *telebot.Bot
}

func Initialize(log logger.ILogger, rating *rating.Service) *Service {
	return &Service{
		log:    log,
		Rating: rating,
	}
}

func (s *Service) Start() (err error) {
	s.bot, err = telebot.NewBot(telebot.Settings{
		Token:  os.Getenv(static.EnvBotToken),
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		return err
	}

	s.initializeHandlers()

	go s.bot.Start()
	s.log.Info("Bot started")

	return nil
}
