package config

import "github.com/kelseyhightower/envconfig"

func Config(pt any) error {
	return envconfig.Process("", pt)
}
