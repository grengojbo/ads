package main

import (
	"fmt"
	"log"
	"os"

	"bitbucket.org/grengojbo/ads-core/config"
	// "bitbucket.org/grengojbo/ads-core/db"

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
	// cmdMigrate,
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

// var cmdMigrate = cli.Command{
// 	Name:  "migrate",
// 	Usage: "Perform database migrations",
// 	// Description: `Perform database migrations`,
// 	Action: runMigrate,
// }

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

	// pool := db.InitDB(&conf)
	// defer pool.Close()

	// release := false
	if c.Bool("release") {
		// release = true
		fmt.Printf("Listening on: %s:%d\n", c.String("host"), conf.Port)
	}
	// server := services.Server{Config: &conf, DB: pool, Release: release}
	// server.Start()

	// 	go func() {
	// 		var zoneName string
	// 		fmt.Printf("Query getZoneById ---> %v=%v\n", zoneId, zoneName)
	// 	}()
}

// func runMigrate(c *cli.Context) {
// 	fmt.Println("Start migration ...")
// 	// migrations.StartMigrate()
// 	fmt.Println("Finish migration.")
// }
