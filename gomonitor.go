package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/jaimegildesagredo/gomonitor/network"
	"github.com/julienschmidt/httprouter"
)

const (
	BANDWIDTH_MONITOR_DELAY = 1 * time.Second
)

func main() {
	log.Println("Starting gomonitor")

	bandwidthService := network.NewBandwidthServiceFactory()

	router := httprouter.New()
	router.GET("/networks/:name/bandwidth", newBandwidthHanler(bandwidthService))
	http.ListenAndServe(":3000", router)
}

func newBandwidthHanler(bandwidthService network.BandwidthService) httprouter.Handle {
	bandwidthsByInterface := map[string]chan network.Bandwidth{}

	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		interfaceName := params.ByName("name")

		bandwidths, found := bandwidthsByInterface[interfaceName]
		if !found {
			bandwidths = bandwidthService.MonitorBandwidth(interfaceName, BANDWIDTH_MONITOR_DELAY)
			bandwidthsByInterface[interfaceName] = bandwidths
		}

		bandwidth := <-bandwidths

		serialized, err := json.Marshal(map[string]int{
			"up":   bandwidth.Up,
			"down": bandwidth.Down,
		})

		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println("Error marshalling bandwidth", err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(serialized)
	}
}
