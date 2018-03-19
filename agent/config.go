package agent

import "github.com/joyent/freebsd-vpc/db"

type Config struct {
	DBConfig db.Config `mapstructure:"db"`
	AgentConfig struct {
		RPCAddress string `mapstructure:"rpc_server_address"`
	} `mapstructure:"agent"`
}
