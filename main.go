package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"pinkd.moe/x/rclone-backup/config"
	"pinkd.moe/x/rclone-backup/sync"
)

func main() {
	var file string
	flag.StringVar(&file, "config,c", "config.json", "config file path")
	flag.Parse()

	confs, err := config.ReadConfig(file)
	if err != nil {
		panic(err)
	}
	sync.Start(confs)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT)
	signal.Notify(c, syscall.SIGTERM)
	<-c
}
