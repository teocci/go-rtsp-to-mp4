// Package go_rtsp_to_mp4
// Created by RTT.
// Author: teocci@yandex.com on 2021-Nov-12
package main

import "log"

func hasError(e error, fatal bool, format ...string) {
	hasFormat := false
	if len(format) > 0 {
		hasFormat = true
	}
	if e != nil && !fatal && !hasFormat {
		log.Println(e)
		return
	}

	if e != nil && !fatal && hasFormat {
		log.Printf(format[0], e)
		return
	}

	if e != nil && fatal {
		log.Fatal(e)
	}
}
