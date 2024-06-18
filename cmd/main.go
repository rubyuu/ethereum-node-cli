package main

import (
	"github.com/nolan/go-demo-server/common"
	"github.com/nolan/go-demo-server/demo"
	"github.com/nolan/go-demo-server/flags"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "demo"
	app.Usage = "startup a json rpc server"
	app.Description = "Learn and practice golang"
	app.Flags = flags.Flags
	app.Action = common.LifecycleCmd(demo.DemoMain())

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
