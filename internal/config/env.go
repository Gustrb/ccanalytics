package config

type Environment string

const (
	EnvironmentDevelopment Environment = "development"
	EnvironmentProduction  Environment = "production"
)

type EnvironmentConfig struct {
	EnviromnmentName Environment `envconfig:"ENVIRONMENT_NAME" default:"development"`
}

var Environments EnvironmentConfig

func init() {
	if err := Config(&Environments); err != nil {
		panic(err)
	}
}
