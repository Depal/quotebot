package entry

import (
	"fmt"
	"github.com/Depal/quotebot/internal/bot"
	botstatic "github.com/Depal/quotebot/internal/bot/static"
	"github.com/Depal/quotebot/internal/entry/static"
	"github.com/Depal/quotebot/internal/rating"
	"github.com/Depal/quotebot/internal/render"
	"github.com/Depal/quotebot/pkg/logger"
	"github.com/Depal/quotebot/pkg/migration"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

type App struct {
	Environment static.Environment
	log         logger.ILogger
	db          *sqlx.DB
	services    *static.Services
}

func Initialize() *App {
	return &App{}
}

func (a *App) Setup() (err error) {
	err = a.determineEnvironment()
	if err != nil {
		return err
	}

	err = a.initializeLogger()
	if err != nil {
		return err
	}

	err = a.initializeGodotenv()
	if err != nil {
		return err
	}

	err = a.checkEnvironmentVariables()
	if err != nil {
		return err
	}

	err = a.initializeDatabase()
	if err != nil {
		return err
	}

	err = a.db.Ping()
	if err != nil {
		return err
	}

	err = migration.Apply(a.db, a.log)
	if err != nil {
		return err
	}

	err = a.initializeServices()
	if err != nil {
		return err
	}

	return nil
}

func (a *App) determineEnvironment() (err error) {
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	if strings.Contains(hostname, "local") {
		a.Environment = static.EnvironmentLocal
	} else {
		a.Environment = static.EnvironmentProduction
	}

	return nil
}

func (a *App) initializeLogger() (err error) {
	var zapLogger *zap.Logger

	zapLogger, err = zap.NewDevelopment()
	if err != nil {
		return err
	}

	a.log = zapLogger.Sugar()
	a.log.Infof("Environment: %s", a.Environment)

	return nil
}

func (a *App) initializeGodotenv() (err error) {
	if a.Environment == static.EnvironmentProduction {
		return nil
	}

	err = godotenv.Load(static.EnvFile)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) checkEnvironmentVariables() (err error) {
	_, ok := os.LookupEnv(static.EnvDatabaseUrl)
	if !ok {
		return fmt.Errorf("environment variable not set: %s", static.EnvDatabaseUrl)
	}

	_, ok = os.LookupEnv(botstatic.EnvBotToken)
	if !ok {
		return fmt.Errorf("environment variable not set: %s", botstatic.EnvBotToken)
	}

	return nil
}

func (a *App) initializeDatabase() (err error) {
	databaseUrl := os.Getenv(static.EnvDatabaseUrl)

	db, err := sqlx.Open("pgx", databaseUrl)
	if err != nil {
		return err
	}

	a.db = db

	return err
}

func (a *App) initializeServices() (err error) {
	a.services = &static.Services{}

	ratingService := rating.Initialize(a.log, a.db)
	renderService := render.Initialize(a.log)
	botService := bot.Initialize(a.log, ratingService, renderService)

	a.services = &static.Services{
		Rating: ratingService,
		Bot:    botService,
	}

	return nil
}

func (a *App) Start() (err error) {
	err = a.services.Bot.Start()
	if err != nil {
		return err
	}

	a.awaitQuitSignal()

	return nil
}

func (a *App) awaitQuitSignal() {
	a.log.Info("Working until a quit signal is received...")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
}

func (a *App) Teardown() (err error) {
	err = a.db.Close()
	if err != nil {
		return err
	}

	return nil
}
