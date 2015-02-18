package dashboards

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const (
	NETWORK_DASHBOARD_HTML_PATH = "dashboards/static/html/network_dashboard.html"
	LOAD_DASHBOARD_HTML_PATH    = "dashboards/static/html/load_dashboard.html"
	GOMONITOR_JS_PATH           = "dashboards/static/js/src/gomonitor.js"
)

var NETWORK_DASHBOARD_HTML []byte
var LOAD_DASHBOARD_HTML []byte
var GOMONITOR_JS []byte

func NewNetworkDashboardHandler() httprouter.Handle {
	contentType := "text/html"
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		if len(NETWORK_DASHBOARD_HTML) == 0 {
			html, err := ioutil.ReadFile(NETWORK_DASHBOARD_HTML_PATH)

			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				log.Println("Error reading", NETWORK_DASHBOARD_HTML_PATH, err.Error())
				return
			}

			w.Header().Set("Content-Type", contentType)
			w.Write(html)
			return
		}

		w.Header().Set("Content-Type", contentType)
		w.Write(NETWORK_DASHBOARD_HTML)
	}
}

func NewLoadDashboardHandler() httprouter.Handle {
	contentType := "text/html"
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		if len(LOAD_DASHBOARD_HTML) == 0 {
			html, err := ioutil.ReadFile(LOAD_DASHBOARD_HTML_PATH)

			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				log.Println("Error reading", LOAD_DASHBOARD_HTML_PATH, err.Error())
				return
			}

			w.Header().Set("Content-Type", contentType)
			w.Write(html)
			return
		}

		w.Header().Set("Content-Type", contentType)
		w.Write(LOAD_DASHBOARD_HTML)
	}
}

func NewGomonitorJsHandler() httprouter.Handle {
	contentType := "text/javascript"

	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		if len(GOMONITOR_JS) == 0 {
			html, err := ioutil.ReadFile(GOMONITOR_JS_PATH)

			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				log.Println("Error reading", GOMONITOR_JS_PATH, err.Error())
				return
			}

			w.Header().Set("Content-Type", contentType)
			w.Write(html)
			return
		}

		w.Header().Set("Content-Type", contentType)
		w.Write(GOMONITOR_JS)
	}
}
