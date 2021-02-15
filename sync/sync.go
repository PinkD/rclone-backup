package sync

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"

	"pinkd.moe/x/rclone-backup/config"
	"pinkd.moe/x/rclone-backup/stat"
)

func appendHomeEnv(cmd *exec.Cmd) error {
	home, ok := os.LookupEnv("HOME")
	if !ok {
		log.Println("HOME is not set, try to find it manually")
		var err error
		home, err = os.UserHomeDir()
		if err != nil {
			return err
		}
		return nil
	}
	log.Printf("HOME is %s\n", home)
	cmd.Env = append(cmd.Env, fmt.Sprintf("HOME=%s", home))
	return nil
}

func appendProxyEnv(cmd *exec.Cmd, p *config.Proxy) error {
	if p != nil {
		var env string
		switch p.Protocol {
		case config.ProxyHTTP:
			env = fmt.Sprintf("http_proxy=%s", p.URL)
		case config.ProxyHTTPS:
			env = fmt.Sprintf("https_proxy=%s", p.URL)
		default:
			return errors.New(fmt.Sprintf("unsupported proxy protocol: %s", p.Protocol))
		}
		cmd.Env = append(cmd.Env, env)
	}
	return nil
}

const CMDRClone = "rclone"

func backupWithConf(ctx context.Context, c *config.Conf) error {
	cmd := exec.CommandContext(ctx, CMDRClone, "sync", c.Path, c.Remote)
	err := appendProxyEnv(cmd, c.Proxy)
	if err != nil {
		return err
	}
	err = appendHomeEnv(cmd)
	if err != nil {
		return err
	}
	status := &stat.SyncStatus{
		Name:   c.Name,
		Status: stat.Syncing,
		Time:   time.Now(),
	}
	stat.Map.Store(c.Name, status)
	go func() {
		out, err := cmd.CombinedOutput()
		if err != nil {
			log.Println(string(out))
			log.Printf("Wait %s returns err: %s\n", cmd, err)
			status.Status = stat.Fail
			status.Time = time.Now()
			stat.Map.Store(c.Name, status)

		} else {
			status.Status = stat.Success
			now := time.Now()
			status.Time = now
			status.LastSuccessTime = now
			stat.Map.Store(c.Name, status)
			log.Printf("Backup %s finished at %s\n", c.Name, now)
		}
	}()
	return nil
}

func Backup(ctx context.Context, c *config.Conf) error {
	switch stat.Map.LoadStatus(c.Name) {
	case stat.Syncing:
		log.Printf("%s is syncing, skip next sync\n", c.Name)
		return nil
	}
	log.Printf("Start to sync %s\n", c.Name)
	err := backupWithConf(ctx, c)
	if err != nil {
		return errors.New(fmt.Sprintf("error when backup %s: %s", c.Name, err))
	}
	return nil
}

var cancelFunc context.CancelFunc
var wg sync.WaitGroup

func Start(confs map[string]*config.Conf) {
	ctx, cancel := context.WithCancel(context.Background())
	cancelFunc = cancel
	wg.Add(len(confs))
	for _, c := range confs {
		go func(c *config.Conf) {
			defer wg.Done()
			for {
				err := Backup(ctx, c)
				if err != nil {
					log.Println(err)
				}
				select {
				case <-ctx.Done():
					return
				case <-time.Tick(c.Interval.Duration):
				}
			}
		}(c)
	}
}

func Stop() {
	if cancelFunc != nil {
		cancelFunc()
		wg.Wait()
	}
}

func Restart(confs map[string]*config.Conf) {
	Stop()
	Start(confs)
}
