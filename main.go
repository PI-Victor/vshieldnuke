/* This program runs in the background on OSX Sierra and kills the McAffee vshield scanner

Installing it (requires root):

<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC -//Apple Computer//DTD PLIST 1.0//EN http://www.apple.com/DTDs/PropertyList-1.0.dtd >
<plist version="1.0">
  <dict>
    <key>Label</key>
    <string>com.user.fuckVShield</string> <- replace with your user
    <key>fuckVshield</key>
    <string>/full/path/to/fuckVShield</string> <- replace with full path.
    <key>KeepAlive</key>
    <true/>
  </dict>
</plist>

to enalbe:

launchctl load ~/Library/Launchagents/com.user.fuckVShield.plist

This will enable fuckVshield to run in the background every second,

But... why?
Because McAffee is a piece of shit antivirus, that's why.


MIT License

Copyright (c) 2017 Victor Palade

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

*/
package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func main() {
	var (
		ticker = time.Ticker{}
		c      = make(chan os.Signal, 1)
	)

	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	go func() {
		for _ = range c {
			ticker.Stop()
			os.Exit(0)
		}
	}()

	t := 5 * time.Second
	duration := time.NewTicker(t)

	for range duration.C {
		nukeShit()
	}

}

func nukeShit() {
	out, err := exec.Command("bash", "-c", "ps -ef | grep -i vshieldscanner | grep -v grep | awk '{print $2}'").Output()
	if err != nil {
		fmt.Printf("An error occured: %#v\n", err)
	}
	outputVShield := fmt.Sprintf("%s", out)
	processes := strings.Split(outputVShield, "\n")
	for _, proc := range processes {
		fmt.Printf("Killing process: %s\n", proc)
		if len(proc) == 0 {
			continue
		}
		proc, err := strconv.Atoi(proc)
		if err != nil {
			fmt.Printf("An error occured: %s\n", err)
		}

		go func(PID int) {
			rand.Seed(time.Now().Unix())
			d := time.Duration(rand.Intn(11-1) + 1)
			time.Sleep(d * time.Second)
			killProc, err := os.FindProcess(PID)
			if err != nil {
				fmt.Println("Error occured while finding process")
			}
			if err := killProc.Kill(); err != nil {
				fmt.Printf("Error killing process: %s\n", err)
			}
		}(proc)
	}
}
