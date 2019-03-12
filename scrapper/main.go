package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

const defaultProtocol = "http"
const fileName = "sites.txt"

// service is for checking service details
type service struct {
	name       string
	accessTime time.Duration
	alive      bool
	err        string
}

// metrics is for administrators to trace query statistics
type metrics struct {
	service service
	counter int
}

// add will add 1 to the service counter when the service is monitored
func (m *metrics) add(service service) {
	m.service = service
	m.counter++
}

// declare a global variable of type metrics to call when a service is monitored
var m metrics

// monitor is for monitoring service for health and access time
func (s *service) monitor() error {
	m.add(*s)
	log.Printf("Starting to monitor service: %s\n", s.name)
	start := time.Now()
	_, err := http.Get(fmt.Sprintf("%s://%s", defaultProtocol, s.name))

	if err == nil {
		accessTime := time.Since(start)
		s.accessTime = accessTime
		s.alive = true
		log.Printf("Finished monitoring service: %s in %v with success\n", s.name, accessTime)
		return nil
	}

	s.alive = false
	s.err = err.Error()
	log.Printf("Failed monitoring service: %s. %s\n", s.name, err)
	return err
}

// monitorServices monitors all services that serviceList holds
func monitorServices(serviceList []string) (services []service) {
	var wg sync.WaitGroup
	wg.Add(len(serviceList))
	for _, name := range serviceList {
		go func(name string) {
			defer wg.Done()
			service := new(service)
			service.name = name
			service.monitor()
			services = append(services, *service)
		}(name)
	}
	wg.Wait()

	return services
}

// readFile reads the file that contains list of services
func readFile(fileName string) (serviceList []string, err error) {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return serviceList, err
	}

	str := string(b)
	serviceList = strings.Split(str, "\n")

	if len(serviceList) == 0 {
		err := errors.New("No services found in the file")
		return serviceList, err
	}

	return serviceList, err
}

// getAliveServices gets all services that has been monitored without error and is alive
func getAliveServices() (aliveServices []service, err error) {
	serviceList, err := readFile(fileName)
	if err != nil {
		return aliveServices, err
	}

	if len(serviceList) == 0 {
		err := errors.New("No services found")
		return aliveServices, err
	}

	services := monitorServices(serviceList)

	for _, service := range services {
		if service.alive {
			aliveServices = append(aliveServices, service)
		}
	}

	if len(aliveServices) == 0 {
		err := errors.New("No alive services found")
		return aliveServices, err
	}

	return aliveServices, nil
}

// getServiceWithMinAccessTime gets the service with minimum access time
func getServiceWithMinAccessTime() (service service, err error) {
	aliveServices, err := getAliveServices()
	if err != nil {
		return service, err
	}

	service = aliveServices[0]
	for _, s := range aliveServices {
		if s.accessTime < service.accessTime {
			service = s
		}
	}

	return service, nil
}

// getServiceWithMaxAccessTime gets the service with maximum access time
func getServiceWithMaxAccessTime() (service service, err error) {
	aliveServices, err := getAliveServices()
	if err != nil {
		return service, err
	}

	service = aliveServices[0]
	for _, s := range aliveServices {
		if s.accessTime > service.accessTime {
			service = s
		}
	}

	return service, nil
}

func main() {
	max, _ := getServiceWithMaxAccessTime()
	min, _ := getServiceWithMinAccessTime()
	fmt.Println(min, max)
}
