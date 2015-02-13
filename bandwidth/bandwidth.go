package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/jaimegildesagredo/gomonitor/network"
)

func main() {
	interfaceName, delay := parseArgs()

	log.Println("Monitor interface", interfaceName)
	log.Println("Monitor delay", delay, "seconds")

	bandwidthService := network.NewBandwidthServiceFactory(interfaceName)
	bandwidths := bandwidthService.MonitorBandwidth(time.Duration(delay) * time.Second)

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
