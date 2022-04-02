package bot

import "github.com/Depal/quotebot/internal/bot/static"

func (s *Service) initializeHandlers() {
	s.addEndpoint(static.CommandPing, s.handlePing)
	s.addEndpoint(static.CommandQuote, s.handleQuote)
	s.addEndpoint(static.CommandAlert, s.handleAlert)
	// ...
}

func (s *Service) addEndpoint(endpoint string, handler interface{}) {
	s.bot.Handle(endpoint, handler)
	s.log.Infof("Command initialized: %s", endpoint)
}
