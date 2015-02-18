package main

import (
	"log"
	"net/http"
	"time"

	"github.com/jaimegildesagredo/gomonitor/dashboards"
	"github.com/jaimegildesagredo/gomonitor/networks"
	"github.com/jaimegildesagredo/gomonitor/resources"
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
	router.GET("/networks/:name/bandwidth", resources.NewBandwidthHandler(interfacesService))
	router.GET("/networks", resources.NewNetworksHandler(interfacesService))
	router.GET("/dashboards/network", dashboards.NewNetworkDashboardHandler())
	router.GET("/gomonitor.js", dashboards.NewGomonitorJsHandler())
	http.ListenAndServe(":3000", router)
}
