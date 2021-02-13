package config

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func Equals(e, a interface{}) {
	if e == a {
		return
	}
	if reflect.DeepEqual(e, a) {
		return
	}
	panic(fmt.Sprintf("%#v does not equal to expected: %#v", a, e))
}

func TestReadConfig(t *testing.T) {
	confStr := `[
  {
    "name": "-",
    "remote": "gd:/backup",
    "interval": "1h",
    "proxy": "http://127.0.0.1:1080"
  },
  {
    "name": "path1",
    "path": "/path1"
  },
  {
    "name": "path2",
    "path": "/path2",
    "remote": "gd:/project",
    "proxy": "socks5://127.0.0.1:1080"
  }
]`
	confs, err := parseConfig([]byte(confStr))
	if err != nil {
		panic(err)
	}
	d, err := time.ParseDuration("1h")
	if err != nil {
		panic(err)
	}
	interval := Duration{d}
	for _, c := range confs {
		switch c.Name {
		case "-":
			panic("shouldn't contain global conf")
		case "path1":
			Equals(c.Path, "/path1")
			Equals(c.Interval, interval)
			Equals(c.Remote, "gd:/backup")
			Equals(c.Proxy.URL, "http://127.0.0.1:1080")
			Equals(c.Proxy.Protocol, "http")
		case "path2":
			Equals(c.Path, "/path2")
			Equals(c.Interval, interval)
			Equals(c.Remote, "gd:/project")
			Equals(c.Proxy.URL, "socks5://127.0.0.1:1080")
			Equals(c.Proxy.Protocol, "socks5")
		}
	}
}

func TestReadConfigNoInterval(t *testing.T) {
	confStr := `[
  {
    "name": "-",
    "remote": "gd:/backup"
  },
  {
    "name": "path1",
    "path": "/path1",
    "interval": "1h"
 },
  {
    "name": "path2",
    "path": "/path2",
    "remote": "gd:/project"
  }
]`
	_, err := parseConfig([]byte(confStr))
	if err == nil {
		panic("There should be a no interval error")
	}
	t.Logf("test no interval: %s", err)
}

func TestReadConfigNoRemote(t *testing.T) {
	confStr := `[
  {
    "name": "-",
    "interval": "1h"
  },
  {
    "name": "path1",
    "path": "/path1"
 },
  {
    "name": "path2",
    "path": "/path2",
    "remote": "gd:/project"
  }
]`
	_, err := parseConfig([]byte(confStr))
	if err == nil {
		panic("There should be a no remote error")
	}
	t.Logf("test no remote: %s", err)
}
