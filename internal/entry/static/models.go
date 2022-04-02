package static

import (
	"github.com/Depal/quotebot/internal/bot"
	"github.com/Depal/quotebot/internal/rating"
)

type Services struct {
	Rating *rating.Service
	Bot    *bot.Service
}
