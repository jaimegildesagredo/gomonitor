package dashboards

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const (
	NETWORK_DASHBOARD_HTML_PATH = "static/html/network_dashboard.html"
)

var NETWORK_DASHBOARD_HTML []byte

func NewNetworkDashboardHandler() httprouter.Handle {
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
