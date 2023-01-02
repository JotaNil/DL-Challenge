package main

import (
	"DreamLabChallenge/cmd/api/ipdata"
	"DreamLabChallenge/cmd/services"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type Application struct {
	server *http.Server
}

// LoadAndRoute loads all the dependencies, prepare handlers and initialize routes
func (d *Application) LoadAndRoute() {
	// dataStorage
	ipv4ProxyDB, _ := services.ConnectToSQLDB(services.Ipv4ProxyDB)

	// ipData
	ipDataDao := ipdata.NewDao(ipv4ProxyDB)
	ipDataGateway := ipdata.NewGateway(ipDataDao)
	ipDataHandler := ipdata.NewHandler(ipDataGateway)

	// Routes --------------------------

	//ipData
	r := mux.NewRouter()

	r.HandleFunc("/ipdata/count/ip/{country_name}", ipDataHandler.GetIPCountByCountryName).Methods("GET")
	r.HandleFunc("/ipdata/{ip}", ipDataHandler.GetDataFromIP).Methods("GET")
	r.HandleFunc("/ipdata/top10/switzerland", ipDataHandler.GetTopISPsFromSwitzerland).Methods("GET")

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	d.server = srv
}
