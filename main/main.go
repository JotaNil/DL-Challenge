package main

import (
	"DreamLabChallenge/cmd/api/ipdata"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

func main() {
	r := mux.NewRouter()
	dao := ipdata.NewDao()
	gtw := ipdata.NewGateway(dao)
	ipDataHandler := ipdata.NewHandler(gtw)
	r.HandleFunc("/ipdata", ipDataHandler.GetData).Methods("GET")
	r.HandleFunc("/ipdata/topISP/{country_code}", ipDataHandler.SelectTopISPByCountryCode).Methods("GET")
	r.HandleFunc("/ipdata/count/ip/{country_name}", ipDataHandler.GetIPCountByCountryName).Methods("GET")
	r.HandleFunc("/ipdata/{ip}", ipDataHandler.GetDataFromIP).Methods("GET")
	r.HandleFunc("/ipdata/top10/switzerland", ipDataHandler.GetTopISPsFromSwitzerland).Methods("GET")

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
