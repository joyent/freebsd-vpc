package agent

import "github.com/joyent/freebsd-vpc/db"

type Config struct {
	DBConfig db.Config `mapstructure:"db"`
	AgentConfig struct {
		Addresses struct {
			Internal string `mapstructure:"internal"`
		} `mapstructure:"addresses"`
	} `mapstructure:"agent"`
}
