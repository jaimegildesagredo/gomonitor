package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/jaimegildesagredo/gomonitor/networks"
	"github.com/julienschmidt/httprouter"
)

const (
	BANDWIDTH_MONITOR_DELAY     = 1 * time.Second
	NETWORK_DASHBOARD_HTML_PATH = "dashboards/network/index.html"
)

var NETWORK_DASHBOARD_HTML []byte

func main() {
	log.Println("Starting gomonitor")

	interfacesService := networks.NewInterfacesServiceFactory()

	router := httprouter.New()
	router.GET("/networks/:name/bandwidth", newBandwidthHanler(interfacesService))
	router.GET("/networks", newNetworksHandler(interfacesService))
	router.GET("/dashboards/network", newNetworkDashboardHandler())
	http.ListenAndServe(":3000", router)
}

func newNetworkDashboardHandler() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		if len(NETWORK_DASHBOARD_HTML) == 0 {
			html, err := ioutil.ReadFile(NETWORK_DASHBOARD_HTML_PATH)

			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				log.Println("Error reading", NETWORK_DASHBOARD_HTML_PATH, err.Error())
				return
			}

			w.Header().Set("Content-Type", "text/html")
			w.Write(html)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		w.Write(NETWORK_DASHBOARD_HTML)
	}
}

func newNetworksHandler(interfacesService networks.InterfacesService) httprouter.Handle {
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

func newBandwidthHanler(interfacesService networks.InterfacesService) httprouter.Handle {
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
