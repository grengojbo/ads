package services

import (
	"fmt"
	// "net"
	"time"

	"strconv"

	"github.com/grengojbo/ads/config"
	// "bitbucket.org/grengojbo/ads-core/db"
	"github.com/gin-gonic/gin"
	// "github.com/jackc/pgx"
	"github.com/nu7hatch/gouuid"
	// log "gopkg.in/inconshreveable/log15.v2"
)

type Server struct {
	// Bus    *bus.Bus
	Config  *config.Config
	DB      *Database
	Release bool
	r, s    *gin.Engine
	Log     *Logger
	// log     log.Logger
}

// Start Web Server
func (self *Server) Start() {
	self.Log.Info("starting server service...")
	// logger := log.New()
	// self.log = log.New()
	// mlog.Info("starting server service")

	if self.Release {
		gin.SetMode(gin.ReleaseMode)
		self.r = gin.New()
	} else {
		self.r = gin.New()
		self.r.Use(gin.Logger())
	}
	self.r.Use(gin.Recovery())

	self.r.NoRoute(self.redirect)
	self.r.GET("/ping", self.ping)

	show := self.r.Group("show")
	show.GET(":region_id/:umac/ping.js", self.showPing)

	// go self.r.Run(fmt.Sprintf("%s:%d", self.Config.Host, self.Config.Port))
	self.r.Run(fmt.Sprintf("%s:%d", self.Config.Host, self.Config.Port))
}

// Stop Web Server
func (self *Server) Stop() {
	// mlog.Info("server service stopped")
	// nothing here
}

// Redirect no route
func (self *Server) redirect(c *gin.Context) {
	c.Redirect(301, "/")
}

// crossdomain.xml

// ping pong :)
func (self *Server) ping(c *gin.Context) {
	c.String(200, "pong")
}

func (self *Server) showPing(c *gin.Context) {
	u4, err := uuid.NewV4()
	if err != nil {
		self.Log.Error("Is not generate uuid", err.Error())
	}
	sesUuid := fmt.Sprintf("%v", u4)

	storeID, err := strconv.Atoi(c.Param("region_id"))
	if err != nil {
		self.Log.Error("Region ID is not integer", storeID)
	}

	t := time.Now().UTC()
	// core.SaveShow(self.DB, t, sesUuid, storeID, c.Param("umac"), c.ClientIP(), c.Request.Header.Get("Accept-Language")[0:2], c.Request.Referer(), c.Request.UserAgent())

	c.Header("cache-control", "priviate, max-age=0, no-cache")
	c.Header("pragma", "no-cache")
	c.Header("expires", "-1")
	c.Header("Last-Modified", fmt.Sprintf("%v", t))
	// c.Header("Date", t.Format(time.RFC1123))
	c.Header("Content-Type", "text/javascript")
	c.String(200, fmt.Sprintf("var uatv_me_uuid='%s';", sesUuid))
}
