// Package main
// Created by RTT.
// Author: teocci@yandex.com on 2021-Nov-16
package main

import (
	"log"
	"os"

	"github.com/teocci/go-samples-rtsp-mp4/src/datamgr"
	"github.com/teocci/go-samples-rtsp-mp4/src/videomgr"
)

func main() {
	pid := os.Getpid()
	log.Println("PID:", pid)

	initConf := datamgr.InitConf{
		FlightID:  134,
		CompanyID: 2,
		DroneID:   4,
		ServerURL: "rtsp://106.244.179.242:554/jinan_test",
	}

	stream := videomgr.New(initConf)
	stream.Start()
}
