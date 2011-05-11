// Copyright 2011 Petar Maymounkov. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"runtime"
	"time"
)

func MonitorMemProfile() {
	go func() {
		for {
			time.Sleep(10e9)
			mem, _ := runtime.MemProfile(nil, false)
			log.Printf("Mem# %d\n", mem)
		}
	}()
}
