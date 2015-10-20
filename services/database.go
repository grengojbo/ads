package services

import (
	// "fmt"
	// "time"
	"os"
	// "strconv"

	"github.com/grengojbo/ads/config"
	"github.com/jackc/pgx"
	log "gopkg.in/inconshreveable/log15.v2"
)

type Database struct {
	// Bus    *bus.Bus
	Config  *config.Config
	DB      *pgx.ConnPool
	Release bool
	log     log.Logger
}

// Start PostgreSQL pool
func (self *Database) Start() {
	self.DB = self.GetPoll()
	defer self.DB.Close()
}

// Get PostgreSQL pool
func (self *Database) GetPoll() (pool *pgx.ConnPool) {
	pool, err := pgx.NewConnPool(self.getConfig())

	if err != nil {
		// log.Printf("Unable to create connection pool to database: %v\n", err)
		os.Exit(1)
	}
	return pool
}

func (self *Database) getConfig() (cfg pgx.ConnPoolConfig) {
	connConfig := pgx.ConnConfig{
		Host:     self.Config.DB.Host,
		User:     self.Config.DB.User,
		Password: self.Config.DB.Password,
		Database: self.Config.DB.Name,
	}

	// if self.Config.DB.Debug {
	//   lvl, err := log.LvlFromString("debug")
	//   if err != nil {
	//     fmt.Println(err)
	//   }
	//   logger.SetHandler(log.LvlFilterHandler(lvl, log.DiscardHandler()))
	//   connConfig.Logger = logger
	//   connConfig.LogLevel = pgx.LogLevelTrace
	// } else {
	//   lvl, err := log.LvlFromString("error")
	//   if err != nil {
	//     fmt.Println(err)
	//   }
	// logger.SetHandler(log.LvlFilterHandler(lvl, log.DiscardHandler()))
	// connConfig.Logger = logger
	// connConfig.LogLevel = pgx.LogLevelError
	// }
	//
	cfg = pgx.ConnPoolConfig{
		ConnConfig:     connConfig,
		MaxConnections: self.Config.DB.Connections,
		AfterConnect:   self.afterConnect,
	}
	return cfg
}

func (self *Database) afterConnect(conn *pgx.Conn) (err error) {
	_, err = conn.Prepare("getZoneById", `SELECT
  "id", "name"
FROM "stores"
WHERE id=$1
    `)
	if err != nil {
		// log.Printf("Prepare [getZoneById] %v\n", err)
		return
	}

	_, err = conn.Prepare("setShowBanner", `INSERT  INTO "banner_shows"
  (created_at, updated_at, show_date, show_time, store_id, is_bot, ses_uuid, user_mac, user_ip, ipv4, accept_language,
  ua_browser_family, ua_browser_version, ua_os_family, ua_device_family, is_mobile, user_agent, referrer)
VALUES ($1::timestamptz(0), NOW()::timestamptz(0), $1::date, $1::time(0), $3, $4, $2, $5, $6, $7, $8, $9, $10::int2, $11, $12, $13, $14, $15)
  `)
	if err != nil {
		log.Error("setShowBanner", err.Error())
		// log.Printf("Prepare [setShowBanner] %v\n", err)
		return
	}

	_, err = conn.Prepare("getZoneByName", `
SELECT
  "address",
  "city",
  "country",
  "created_at",
  "deleted_at",
  "email",
  "id",
  "latitude",
  "longitude",
  "name",
  "phone",
  "position",
  "region",
  "updated_at",
  "zip"
FROM "stores" WHERE lower(name)=lower($1)
    `)
	if err != nil {
		// log.Printf("Prepare [getZoneByName] %v\n", err)
		return
	}
	return
}
