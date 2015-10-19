package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/codegangsta/cli"
)

var (
	Version   = "0.1.0"
	BuildTime = "2015-09-20 UTC"
	GitHash   = "c00"
)

func ConfigRuntime() {
	nuCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(nuCPU)
	fmt.Printf("Running with %d CPUs\n", nuCPU)
}

func main() {
	app := cli.NewApp()
	app.Name = "ads"
	app.Version = Version
	app.Usage = "Advertising System"
	app.Author = "Oleg Dolya"
	app.Email = "oleg.dolya@gmail.com"
	app.EnableBashCompletion = true
	app.Commands = Commands
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "verbose",
			Usage: "Verbose mode",
		},
		cli.BoolFlag{
			Name:  "debug",
			Usage: "Debug mode",
		},
		cli.StringFlag{
			Name:   "config, c",
			Value:  "config/config.yml",
			Usage:  "config file to use (config/config.yml)",
			EnvVar: "APP_CONFIG",
		},
	}

	app.Run(os.Args)
}
