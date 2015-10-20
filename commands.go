package main

import (
	"fmt"
	"log"
	"os"

	"github.com/grengojbo/ads/config"
	"github.com/grengojbo/ads/services"

	"github.com/codegangsta/cli"
	"github.com/jinzhu/configor"
	// "github.com/fatih/color"
)

// Show debug message
func debug(v ...interface{}) {
	if os.Getenv("DEBUG") != "" {
		log.Println(v...)
	}
}

// Error assert
func assert(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// Loading configuration from yaml file
func getConfig(c *cli.Context) (config.Config, error) {
	yamlPath := c.GlobalString("config")
	conf := config.Config{}

	err := configor.Load(&conf, yamlPath)
	return conf, err
}

var Commands = []cli.Command{
	cmdServer,
}

var cmdServer = cli.Command{
	Name:   "server",
	Usage:  "Start web server",
	Action: runWeb,
	Flags: []cli.Flag{
		cli.IntFlag{
			Name:  "port, p",
			Usage: "port number to start web server",
		},
		cli.StringFlag{"host", "", "Host to start web server", ""},
		cli.BoolFlag{
			Name:   "release",
			Usage:  "Release mode in production.",
			EnvVar: "GIN_MODE",
		},
	},
}

func runWeb(c *cli.Context) {
	ConfigRuntime()
	conf, err := getConfig(c)
	assert(err)

	conf.Host = c.String("host")
	if c.Int("port") > 0 {
		conf.Port = c.Int("port")
	}

	fmt.Printf("App Version: %s\n", Version)
	fmt.Printf("Build Time: %s\n", BuildTime)
	fmt.Printf("Git Commit Hash: %s\n", GitHash)

	release := false
	if c.Bool("release") {
		release = true
		fmt.Printf("Listening on: %s:%d\n", c.String("host"), conf.Port)
	}
	l := services.Logger{Config: &conf, Release: release}
	l.Start()

	db := services.Database{Config: &conf, Log: &l, Release: release}
	// db.Start()
	pool := db.GetPoll()
	defer pool.Close()
	server := services.Server{Config: &conf, Log: &l, DB: pool, Release: release}
	server.Start()
}
