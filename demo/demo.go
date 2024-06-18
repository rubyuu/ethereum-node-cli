package demo

import (
	"github.com/nolan/go-demo-server/common"
	"github.com/nolan/go-demo-server/flags"
	"github.com/urfave/cli/v2"
)

func DemoMain() common.LifecycleAction {
	return func(cliCtx *cli.Context) (common.Lifecycle, error) {
		// check flags
		if err := flags.CheckRequired(cliCtx); err != nil {
			return nil, err
		}

		// create config
		cfg := NewConfig(cliCtx)
		return DemoServiceFromCLIConfig(cliCtx.Context, cfg)
	}
}
