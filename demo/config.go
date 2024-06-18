package demo

import (
	"github.com/nolan/go-demo-server/flags"
	"github.com/urfave/cli/v2"
)

type CliConfig struct {
	Host         string
	Port         int
	EthRpc       string
	PprofEnabled bool
	PprofPort    int
	SyncEnabled  bool
	SyncInterval int
}

func NewConfig(cliCtx *cli.Context) *CliConfig {
	return &CliConfig{
		Host:         cliCtx.String(flags.HostFlag.Name),
		Port:         cliCtx.Int(flags.PortFlag.Name),
		EthRpc:       cliCtx.String(flags.EthRpcFlag.Name),
		PprofEnabled: cliCtx.Bool(flags.PprofEnabledFlag.Name),
		PprofPort:    cliCtx.Int(flags.PprofPortFlag.Name),
		SyncEnabled:  cliCtx.Bool(flags.SyncEnabledFlag.Name),
		SyncInterval: cliCtx.Int(flags.SyncIntervalFlag.Name),
	}
}
