package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"reflect"
	"time"

	holder "./holder"
	stats "./stats"
)

var hostHolder *holder.HostHolder
var statistics *stats.CounterService

func main()  {
	hostHolder = holder.NewHostHolder("sites.txt")
	statistics = stats.NewCounterService()
	http.HandleFunc("/getTiming/random", getStatsByRandom)
	http.HandleFunc("/getTiming/min", getStatsByMin)
	http.HandleFunc("/getTiming/max", getStatsByMax)
	http.HandleFunc("/stats", getStatsByReqested)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))

}


func getStatsByReqested(w http.ResponseWriter, r *http.Request) {
	counters := statistics.GetCounters()
	fmt.Fprintf(w, "Reqested URI\t\t | Count\n")
	for i := range counters {
		fmt.Fprintf(w, "%s\t\t | %d\n", i, counters[i])
	}

}

func getStatsByRandom(w http.ResponseWriter, r *http.Request) {
	statistics.IncreaseURIStat(r.RequestURI)
	availableHosts := reflect.ValueOf(hostHolder.Hosts).MapKeys()
	randHostname := availableHosts[rand.Intn(len(availableHosts))]

	timingOfHost, err := hostHolder.GetDataByHostname(fmt.Sprintf("%s", randHostname))
	if err != nil {
		fmt.Fprintf(w, "Undefined host.")
		return
	}
	if timingOfHost < 0 {
		fmt.Fprintf(w, "Host unavailable by timeout")
		return
	}
	fmt.Fprintf(w, "Timiout to host %s is %v", randHostname, timingOfHost)
}

func getStatsByMin(w http.ResponseWriter, r *http.Request) {
	statistics.IncreaseURIStat(r.RequestURI)
	availableHosts := hostHolder.Hosts

	keys := reflect.ValueOf(availableHosts).MapKeys()
	var minHostName string
	minHostTimeout := availableHosts[fmt.Sprintf("%s", keys[0])]
	for i := range availableHosts {
		if availableHosts[i] < minHostTimeout && availableHosts[i] > time.Duration(0) {
			minHostTimeout = availableHosts[i]
			minHostName = i
		}
	}

	if "" == minHostName {
		fmt.Fprintf(w, "No one of hosts are alive.")
		return
	}

	fmt.Fprintf(w, "Host %s is most avaialable of others and has timeout %v", minHostName, minHostTimeout)
}

func getStatsByMax(w http.ResponseWriter, r *http.Request) {
	statistics.IncreaseURIStat(r.RequestURI)
	availableHosts := hostHolder.Hosts
	var maxHostName string
	maxHostTimeout := time.Duration(0)
	for i := range availableHosts {
		if availableHosts[i] > maxHostTimeout && availableHosts[i] > time.Duration(0) {
			maxHostTimeout = availableHosts[i]
			maxHostName = i
		}
	}

	if "" == maxHostName {
		fmt.Fprintf(w, "No one of hosts are alive.")
		return
	}

	fmt.Fprintf(w, "Host %s is the worst by avaialability with timeout %v", maxHostName, maxHostTimeout)
}