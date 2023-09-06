package redis

import (
	"strings"
	"time"
)

type Config struct {
	Addr     string
	Host     string
	Port     string
	Password string
	DB       int
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
