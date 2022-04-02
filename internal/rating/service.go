package rating

import (
	"github.com/Depal/quotebot/pkg/logger"
	"github.com/jmoiron/sqlx"
)

type Service struct {
	log logger.ILogger
	db  *sqlx.DB
}

func Initialize(log logger.ILogger, db *sqlx.DB) *Service {
	return &Service{
		log: log,
		db:  db,
	}
}
