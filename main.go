// Package go_rtsp_to_mp4
// Created by RTT.
// Author: teocci@yandex.com on 2021-Nov-16
package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type InitConf struct {
	Host      string
	Port      int64
	ModuleTag string
	CompanyID int64
	DroneID   int64
	FlightID  int64
	ServerURL string
}

const (
	prefixMP4File        = "record"
	formatMP4Filename    = "%s-%d-%s"
	formatMP4FilenameExt = "%s.mp4"
	formatMP4DirPath     = "%s/c-%d/d-%d"

	baseLogsPath = "/home/rtt/jinan/videos"
)

func main() {
	pid := os.Getpid()
	log.Println("PID:", pid)

	c := InitConf{
		FlightID:  134,
		CompanyID: 2,
		DroneID:   4,
		ServerURL: "rtsp://106.244.179.242:554/jinan_test",
	}

	t := time.Now()
	vn := fmt.Sprintf(formatMP4Filename, prefixMP4File, c.FlightID, t.Format("20060102-150405"))
	log.Println("RTT filename:", vn)

	videoDirPath := fmt.Sprintf(formatMP4DirPath, baseLogsPath, c.CompanyID, c.DroneID)

	vnExt := fmt.Sprintf(formatMP4FilenameExt, vn)
	videoPath := filepath.Join(videoDirPath, vnExt)
	log.Println("RTT file path:", videoPath)

	err := os.MkdirAll(filepath.Dir(videoPath), os.ModePerm)
	hasError(err, true)

	// use "-strftime", "1", to use "strftime names" like out%Y-%m-%d_%H-%M-%S.mp4
	cmdArguments := []string{
		"-rtsp_transport", "tcp",
		"-i", c.ServerURL,
		"-c", "copy",
		"-movflags", "faststart",
		videoPath,
	}

	var out bytes.Buffer
	cmd := exec.Command("ffmpeg", cmdArguments...)
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Command output: %q\n", out.String())
}
