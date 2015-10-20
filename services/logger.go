package services

import (
	"fmt"

	"github.com/grengojbo/ads/config"
	log "gopkg.in/inconshreveable/log15.v2"
)

type Logger struct {
	Config  *config.Config
	Release bool
	log.Logger
	Level log.Lvl
}

// Start Logger service
func (self *Logger) Start() {
	var err error
	self.Logger = log.New()
	if self.Release {
		self.Level, err = log.LvlFromString("error")
		if err != nil {
			fmt.Println(err)
		}
	} else {
		self.Level, err = log.LvlFromString("debug")
		if err != nil {
			fmt.Println(err)
		}
	}
	self.Logger.SetHandler(log.LvlFilterHandler(self.Level, log.StdoutHandler))
	// self.SetHandler(log.LvlFilterHandler(lvl, log.DiscardHandler()))
	self.Info("Start Logger service...")
}
