package static

const (
	EnvFile        string = ".env"
	EnvDatabaseUrl string = "DATABASE_URL"
)

type Environment string

const (
	EnvironmentLocal      Environment = "LOCAL"
	EnvironmentProduction Environment = "PRODUCTION"
)
