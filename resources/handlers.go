package resources

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/jaimegildesagredo/gomonitor/loads"
	"github.com/jaimegildesagredo/gomonitor/networks"
	"github.com/julienschmidt/httprouter"
)

const (
	BANDWIDTH_MONITOR_DELAY = 1 * time.Second
	LOAD_MONITOR_DELAY      = 1 * time.Second
)

func NewLoadHandler(loadService loads.LoadService) httprouter.Handle {
	var loads_ chan loads.Load
	var lastLoad loads.Load

	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		if loads_ == nil {
			loads_ = loadService.Monitor(LOAD_MONITOR_DELAY)

			go func() {
				for load := range loads_ {
					lastLoad = load
				}
			}()
		}

		serialized, err := serializeLoad(lastLoad)

		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println("Error marshalling load", err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(serialized)
	}
}

func serializeLoad(load loads.Load) ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"values":     load.Values,
		"created_at": load.CreatedAt,
	})
}

func NewNetworksHandler(interfacesService networks.InterfacesService) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		serialized, err := serializeInterfaces(interfacesService.FindAll())

		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println("Error marshalling network interfaces", err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write(serialized)
	}
}

func serializeInterfaces(interfaces []networks.Interface) ([]byte, error) {
	plain := []map[string]interface{}{}

	for _, interface_ := range interfaces {
		plain = append(plain, map[string]interface{}{
			"name":  interface_.Name,
			"state": interface_.State,
		})
	}

	return json.Marshal(plain)
}

func NewBandwidthHandler(interfacesService networks.InterfacesService) httprouter.Handle {
	lastBandwidthsByInterface := map[string]networks.Bandwidth{}

	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		interfaceName := params.ByName("name")
		bandwidth, found := lastBandwidthsByInterface[interfaceName]
		if !found {
			bandwidths, err := interfacesService.MonitorBandwidth(interfaceName, BANDWIDTH_MONITOR_DELAY)

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
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write(serialized)
	}
}
