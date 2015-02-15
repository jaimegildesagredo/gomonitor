package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/jaimegildesagredo/gomonitor/networks"
	"github.com/julienschmidt/httprouter"
)

const (
	BANDWIDTH_MONITOR_DELAY = 1 * time.Second
)

func main() {
	log.Println("Starting gomonitor")

	bandwidthService := networks.NewBandwidthServiceFactory()

	router := httprouter.New()
	router.GET("/networks/:name/bandwidth", newBandwidthHanler(bandwidthService))
	http.ListenAndServe(":3000", router)
}

func newBandwidthHanler(bandwidthService networks.BandwidthService) httprouter.Handle {
	lastBandwidthsByInterface := map[string]networks.Bandwidth{}

	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		interfaceName := params.ByName("name")
		bandwidth, found := lastBandwidthsByInterface[interfaceName]
		if !found {
			bandwidths, err := bandwidthService.MonitorBandwidth(interfaceName, BANDWIDTH_MONITOR_DELAY)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				log.Println("Error monitoring bandwidth for", interfaceName, err.Error())
				return
			}

			bandwidth = <-bandwidths

			go func() {
				for bandwidth := range bandwidths {
					lastBandwidthsByInterface[interfaceName] = bandwidth
				}
			}()
		}

		lastBandwidthsByInterface[interfaceName] = bandwidth

		serialized, err := json.Marshal(map[string]interface{}{
			"up":         bandwidth.Up,
			"down":       bandwidth.Down,
			"created_at": bandwidth.CreatedAt,
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
