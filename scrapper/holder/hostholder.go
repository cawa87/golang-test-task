package holder

import (
	checker "../checker"
	"bufio"
	"errors"
	"log"
	"os"
	"sync"
	"time"
)

var CErrHostInvalid = errors.New("Host undefined")
var mu = &sync.Mutex{}

const CHostMaxTimeout = 500 * time.Millisecond
const CHostCheckInterval = 10 * time.Second

//HostList struct contains availability data of hosts. Key of any of Hosts will be a hostname, value by key defines last timeout or
// must be explained as "unavailable" in case of a negative value
type HostHolder struct {
	Hosts	map[string]time.Duration
	StopChan chan bool
}

//NewHostHolder returns new structure with already ran host checking
func NewHostHolder(hostSource string) *HostHolder {
	h := HostHolder{
		Hosts: make(map[string]time.Duration),
		StopChan: make(chan bool),
	}

	hostList := readHosts(hostSource)

	for i := range hostList {
		h.Hosts[hostList[i]] = 0
	}

	go h.Serve()
	return &h
}

func readHosts(filepath string) []string {
	var list []string
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		list = append(list, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return list
}

func (hl *HostHolder) GetDataByHostname(hostname string) (time.Duration, error) {

	if tm, ok := hl.Hosts[hostname]; ok {
		return tm, nil
	}

	return 0, CErrHostInvalid
}

func (hl *HostHolder) updateTimeout(hostname string, timeoutVal time.Duration) {
	mu.Lock()
	hl.Hosts[hostname] = timeoutVal
	mu.Unlock()
}

func (hl *HostHolder) update (){
	hosts := hl.Hosts
	wg := sync.WaitGroup{}
	for i := range hosts{
		wg.Add(1)
		go func(hostname string, wg *sync.WaitGroup) {
			doneChan := make(chan time.Duration)
			defer wg.Done()

			go func(doneCh chan time.Duration) {
				timeout := checker.CheckHostByHostname(hostname)
				doneChan <- timeout
			}(doneChan)

			ttlTimer := time.NewTimer(CHostMaxTimeout)

			select {
			case <-ttlTimer.C:
				go hl.updateTimeout(hostname, time.Duration(-1))
				return
			case dur:= <-doneChan:
				ttlTimer.Stop()
				go hl.updateTimeout(hostname, dur)
				return
			}

		}(i, &wg)
	}
	wg.Wait()
}

func (hl *HostHolder) Stop() {
	hl.StopChan <- true
}

//Serve run host checking
func (hl *HostHolder) Serve()  {
	eachMinuteUpdateTimer := time.NewTimer(CHostCheckInterval)
	select {
	case <-eachMinuteUpdateTimer.C:
		hl.update()

	case <-hl.StopChan:
		return
	}
}
