package main

import (
	"fmt"
	"log"
	// "net"
	"os"
	// "time"

	"bitbucket.org/grengojbo/ads-core/config"
	// "bitbucket.org/grengojbo/ads-core/core"
	"bitbucket.org/grengojbo/ads-core/db"
	"bitbucket.org/grengojbo/ads-core/services"
	"github.com/codegangsta/cli"
	// "github.com/gin-gonic/gin"
	"github.com/jinzhu/configor"
	// "github.com/fatih/color"
	// "github.com/qor/qor-example/config/admin"
	// "github.com/qor/qor-example/db/migrations"
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

	// if _, err := os.Stat(yamlPath); err != nil {
	// 	return config, errors.New("config path not valid")
	// }

	// ymlData, err := ioutil.ReadFile(yamlPath)
	// if err != nil {
	// 	return config, err
	// }

	// err = yaml.Unmarshal([]byte(ymlData), &config)
	err := configor.Load(&conf, yamlPath)
	return conf, err
}

var Commands = []cli.Command{
	cmdServer,
	// cmdMigrate,
}

var cmdServer = cli.Command{
	Name:        "server",
	Usage:       "Start web server",
	Description: `Start QOR Admin web server`,
	Action:      runWeb,
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

	pool := db.InitDB(&conf)
	defer pool.Close()

	release := false
	if c.Bool("release") {
		release = true
		fmt.Printf("Listening on: %s:%d\n", c.String("host"), conf.Port)
	}
	server := services.Server{Config: &conf, DB: pool, Release: release}
	server.Start()

	// 	go func() {

	// 		var zoneName string
	// 		fmt.Printf("Query getZoneById ---> %v=%v\n", zoneId, zoneName)

	// 	}()
	// 	// INSERT INTO "banner_shows" ( "ses_uuid", "store_id", "user_ip", "user_mac", "created_at", "updated_at", "show_year", "show_month", "show_day", "show_hour", "show_minute", "user_agent", "accept_language" )
	// 	// VALUES ('df9e9ecf-b17c-4fe6-502e-aaddb55b961c', 8, '127.0.0.1:61526', '01:23:32:bb:63:12', NOW(), NOW(), Extract(YEAR from Now()), Extract(month from Now()), Extract(DAY from Now()), Extract(hour from Now()), Extract(minute from Now()), 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/45.0.2454.101 Safari/537.36', 'ru');
}

// func runMigrate(c *cli.Context) {
// 	fmt.Println("Start migration ...")
// 	// migrations.StartMigrate()
// 	fmt.Println("Finish migration.")
// }
