package checker

import (
	"net"
	"time"
)

func CheckHostByHostname(hostname string) time.Duration {
	nowatime := time.Now()
	conn, err := net.Dial("tcp", hostname + ":80")
	if err != nil {
		return -1
	} else {
		defer conn.Close()
		return time.Now().Sub(nowatime)
	}
}