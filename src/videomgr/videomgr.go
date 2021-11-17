// Package videomgr
// Created by RTT.
// Author: teocci@yandex.com on 2021-Nov-17
package videomgr

import (
	"bytes"
	"fmt"
	"github.com/teocci/go-rtsp-to-mp4/src/datamgr"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

const (
	cmdFFMPEG = "ffmpeg"

	prefixMP4File        = "record"
	formatMP4Filename    = "%s-%d-%s"
	formatMP4FilenameExt = "%s.mp4"
	formatMP4DirPath     = "%s/c-%d/d-%d"

	baseLogsPath = "/home/rtt/jinan/videos"
)

type VideoExec struct {
	Executor *exec.Cmd
	Out      *bytes.Buffer
	// Buffered channel of outbound messages.
	Done    chan struct{}
	Signals chan os.Signal
}

func New(c datamgr.InitConf) (videoExec *VideoExec) {
	t := time.Now()
	vn := fmt.Sprintf(formatMP4Filename, prefixMP4File, c.FlightID, t.Format("20060102-150405"))
	log.Println("RTT filename:", vn)

	videoDirPath := fmt.Sprintf(formatMP4DirPath, baseLogsPath, c.CompanyID, c.DroneID)

	vnExt := fmt.Sprintf(formatMP4FilenameExt, vn)
	videoPath := filepath.Join(videoDirPath, vnExt)
	log.Println("RTT file path:", videoPath)

	err := os.MkdirAll(filepath.Dir(videoPath), os.ModePerm)
	hasError(err, true)

	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)

	// use "-strftime", "1", to use "strftime names" like out%Y-%m-%d_%H-%M-%S.mp4
	cmdArguments := []string{
		"-rtsp_transport", "tcp",
		"-i", c.ServerURL,
		"-c", "copy",
		"-movflags", "faststart",
		videoPath,
	}

	cmd := exec.Command(cmdFFMPEG, cmdArguments...)

	videoExec = &VideoExec{
		Executor: cmd,
		Out:      &bytes.Buffer{},
		Done:     make(chan struct{}),
		Signals:  signals,
	}

	go videoExec.onExecution()

	return
}

func (v *VideoExec) onExecution() {
	defer v.Close()

	for {
		select {
		case <-v.Done:
			log.Printf("Done.")
			return
		case s := <-v.Signals:
			switch s {
			case syscall.SIGHUP:
				log.Println("onExecution-> SIGHUP")
			case os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT:
				log.Printf("onExecution-> %s\n", s.String())
				v.endProcess()
			case os.Kill:
				log.Printf("onExecution-> %s\n", s.String())
				v.endProcess()
			}

			// Close file
			select {
			case <-v.Done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

func (v *VideoExec) endProcess() {
	err := v.Executor.Process.Signal(syscall.SIGTERM)
	hasError(err, false)
}

func (v *VideoExec) Start() {
	v.Executor.Stdout = v.Out
	err := v.Executor.Start()
	hasError(err, true)
	log.Printf("started process %v\n", v.PID())

	err = v.Executor.Wait()
	log.Printf("Command output: %q\n", err)
}

func (v *VideoExec) PID() int {
	return v.Executor.Process.Pid
}

func (v *VideoExec) Close() {
	v.Close()
}
