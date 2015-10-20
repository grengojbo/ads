package services

import (
	"fmt"
	"net"
	"time"

	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/grengojbo/ads/config"
	"github.com/jackc/pgx"
	"github.com/mssola/user_agent"
	"github.com/nu7hatch/gouuid"
)

type Server struct {
	Config  *config.Config
	DB      *pgx.ConnPool
	Release bool
	r, s    *gin.Engine
	Log     *Logger
}

// Start Web Server
func (self *Server) Start() {
	self.Log.Info("starting server service...")

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

	self.r.Run(fmt.Sprintf("%s:%d", self.Config.Host, self.Config.Port))
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
	go self.saveShow(t, sesUuid, storeID, c.Param("umac"), c.ClientIP(), c.Request.Header.Get("Accept-Language")[0:2], c.Request.Referer(), c.Request.UserAgent())

	c.Header("cache-control", "priviate, max-age=0, no-cache")
	c.Header("pragma", "no-cache")
	c.Header("expires", "-1")
	c.Header("Last-Modified", fmt.Sprintf("%v", t))
	c.Header("Content-Type", "text/javascript")
	c.String(200, fmt.Sprintf("var uatv_me_uuid='%s';", sesUuid))
}

func (self *Server) saveShow(t time.Time, sesUuid string, storeID int, userMac string, remoteIp string, acceptLanguage string, refererSrc string, userAgent string) {
	ipv4 := false
	var zoneId pgx.NullInt32
	var uaBrowserVersion pgx.NullInt16
	var zoneName pgx.NullString
	var mac pgx.NullString
	var referer pgx.NullString

	umac, err := net.ParseMAC(userMac)
	if err == nil {
		mac = pgx.NullString{String: umac.String(), Valid: true}
	}

	if err := self.DB.QueryRow("getZoneById", storeID).Scan(&zoneId, &zoneName); err != nil {
		self.Log.Error("Is not", "zoneId", storeID)
	}
	ua := user_agent.New(userAgent)

	uaBrowserFamily, version := ua.Browser()
	browserVersion, err := strconv.Atoi(version[0:strings.Index(version, ".")])
	if err == nil {
		uaBrowserVersion = pgx.NullInt16{Int16: int16(browserVersion), Valid: true}
	}

	ip, _, err := net.SplitHostPort(remoteIp)
	if err != nil {
		ip = "127.0.0.1"
	}
	if net.ParseIP(ip).To4() != nil {
		ipv4 = true
	}

	if len(refererSrc) > 1 {
		referer = pgx.NullString{String: refererSrc, Valid: true}
	}

	if _, err := self.DB.Exec("setShowBanner", t, sesUuid, zoneId, ua.Bot(), mac, ip, ipv4, acceptLanguage, uaBrowserFamily, uaBrowserVersion, ua.OS(), ua.Platform(), ua.Mobile(), userAgent, referer); err != nil {
		self.Log.Error("Exec", "setShowBanner", err)
	}

}
