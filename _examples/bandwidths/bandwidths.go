package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/jaimegildesagredo/gomonitor/networks"
)

func main() {
	interfaceName, delay := parseArgs()

	log.Println("Monitor interface", interfaceName)
	log.Println("Monitor delay", delay, "seconds")

	interfacesService := networks.NewInterfacesServiceFactory()
	bandwidths, err := interfacesService.MonitorBandwidth(interfaceName, time.Duration(delay)*time.Second)

	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case bandwidth := <-bandwidths:
			fmt.Println("D:", bandwidth.Down/1000, "KB/s")
			fmt.Println("U:", bandwidth.Up/1000, "KB/s")
		}
	}
}

func parseArgs() (string, int) {
	interfaceName := flag.String("interface", "", "The interface to monitor")
	delay := flag.Int("delay", 1, "The monitor seconds delay")
	flag.Parse()

	if *interfaceName == "" {
		log.Fatal("'-interface' argument required")
	}
	return *interfaceName, *delay
}
