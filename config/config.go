package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"path/filepath"
	"regexp"
	"time"
)

const GlobalConfName = "-"

type Duration struct {
	time.Duration
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		d.Duration = time.Duration(value)
		return nil
	case string:
		var err error
		d.Duration, err = time.ParseDuration(value)
		if err != nil {
			return err
		}
		return nil
	default:
		return errors.New("invalid duration")
	}
}

func (d Duration) Zero() bool {
	return d.Duration == 0
}

const (
	ProxySocks5 = "socks5"
	ProxyHTTP   = "http"
	ProxyHTTPS  = "https"
)

type Proxy struct {
	Protocol string
	URL      string
}

func (p Proxy) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.URL)
}

func (p *Proxy) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	r, err := regexp.Compile("([0-9a-z]+)://.+")
	if err != nil {
		return err
	}
	switch value := v.(type) {
	case string:
		matches := r.FindAllStringSubmatch(value, -1)
		p.Protocol = matches[0][1]
		p.URL = value
		return nil
	default:
		return errors.New("invalid proxy format")
	}
}

type Conf struct {
	Name     string   `json:"name"`
	Remote   string   `json:"remote"`
	Path     string   `json:"path"`
	Interval Duration `json:"interval"`
	Proxy    *Proxy   `json:"proxy,omitempty"`
}

func mergeGlobalConf(c, g *Conf) error {
	if len(c.Path) == 0 {
		return errors.New(fmt.Sprintf("no path defined for %s", c.Name))
	}
	if len(c.Remote) == 0 {
		if g == nil || len(g.Remote) == 0 {
			return errors.New(fmt.Sprintf("no remote defined for %s", c.Name))
		}
		c.Remote = g.Remote
	}
	// perform directory copy manually because rclone only support file copy
	if c.Remote[len(c.Remote)-1] == '/' {
		// split with \ on windows
		// split with / on linux
		_, filename := filepath.Split(c.Path)
		// always join with /
		c.Remote = path.Join(c.Remote, filename)
		log.Printf("Perform directory copy for %s to %s because remote end with /\n", c.Path, c.Remote)
	}
	if c.Interval.Zero() {
		if g == nil || g.Interval.Zero() {
			return errors.New(fmt.Sprintf("no interval defined for %s", c.Name))
		}
		c.Interval = g.Interval
	}
	if c.Proxy == nil {
		if g != nil {
			c.Proxy = g.Proxy
		}
	}
	return nil
}

func parseConfig(data []byte) (map[string]*Conf, error) {
	var cs []*Conf
	err := json.Unmarshal(data, &cs)
	if err != nil {
		return nil, err
	}
	confs := make(map[string]*Conf)
	var globalConf *Conf
	for _, c := range cs {
		if c.Name == GlobalConfName {
			globalConf = c
			break
		}
	}
	for _, c := range cs {
		if c.Name == GlobalConfName {
			continue
		}
		err := mergeGlobalConf(c, globalConf)
		if err != nil {
			return nil, err
		}
		confs[c.Name] = c
	}
	return confs, nil
}

func ReadConfig(file string) (map[string]*Conf, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return parseConfig(data)
}
