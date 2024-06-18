package flags

import (
	"fmt"
	"github.com/urfave/cli/v2"
)

const EnvVarPrefix = "DEMO"

func prefixEnvVars(name string) []string {
	return []string{EnvVarPrefix + "_" + name}
}

var (
	HostFlag = &cli.StringFlag{
		Name:    "host",
		Usage:   "Hostname or IP address to bind to",
		EnvVars: prefixEnvVars("HOST"),
		Value:   "0.0.0.0",
	}
	PortFlag = &cli.IntFlag{
		Name:    "port",
		Usage:   "Port to bind to",
		EnvVars: prefixEnvVars("PORT"),
		Value:   9090,
	}
	PprofEnabledFlag = &cli.BoolFlag{
		Name:    "pprof-enabled",
		Usage:   "Enable pprof",
		EnvVars: prefixEnvVars("PPROF_ENABLED"),
		Value:   false,
	}
	PprofPortFlag = &cli.IntFlag{
		Name:    "pprof-port",
		Usage:   "pprof port",
		EnvVars: prefixEnvVars("PPROF_PORT"),
		Value:   9091,
	}
	SyncEnabledFlag = &cli.BoolFlag{
		Name:    "sync-enabled",
		Usage:   "Enable to sync block number",
		EnvVars: prefixEnvVars("SYNC_ENABLED"),
		Value:   false,
	}
	SyncIntervalFlag = &cli.IntFlag{
		Name:    "sync-interval",
		Usage:   "Sync interval in seconds",
		EnvVars: prefixEnvVars("SYNC_INTERVAL"),
		Value:   3,
	}
	EthRpcFlag = &cli.StringFlag{
		Name:    "eth-rpc",
		Usage:   "eth rpc endpoint",
		EnvVars: prefixEnvVars("ETH_RPC"),
	}
)

// 必须参数
var requiredFlags = []cli.Flag{
	EthRpcFlag,
}

// 可选参数
var optionalFlags = []cli.Flag{
	HostFlag, PortFlag, PprofEnabledFlag, PprofPortFlag, SyncEnabledFlag, SyncIntervalFlag,
}

func init() {
	Flags = append(requiredFlags, optionalFlags...)
}

var Flags = []cli.Flag{}

func CheckRequired(ctx *cli.Context) error {
	for _, f := range requiredFlags {
		if !ctx.IsSet(f.Names()[0]) {
			return fmt.Errorf("flag %s is required", f.Names()[0])
		}
	}
	return nil
}
