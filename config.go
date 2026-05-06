package redis

import (
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

type Config struct {
	Addr     string
	Host     string
	Port     string
	Username string
	Password string
	DB       int
	Tls      bool
	Timeout  time.Duration
}

func (c *Config) Initialize() {
	if c.Addr == "" {
		if strings.Index(c.Host, ":") > 0 {
			// port already defined in host
			c.Addr = c.Host
		} else {
			c.Addr = c.Host + ":" + c.Port
		}
	}
	if c.Timeout == 0 {
		c.Timeout = time.Second * 2
	}
}

var Forever = time.Duration(redis.KeepTTL)