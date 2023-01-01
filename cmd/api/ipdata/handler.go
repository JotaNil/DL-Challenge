package ipdata

import (
	"DreamLabChallenge/cmd/api/common"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Handler interface {
	GetTopISPsFromSwitzerland(w http.ResponseWriter, r *http.Request) //
	GetIPCountByCountryName(w http.ResponseWriter, r *http.Request)   //
	GetDataFromIP(w http.ResponseWriter, r *http.Request)             //
}
type handler struct {
	gtw Gateway
}

func NewHandler(gtw Gateway) Handler {
	return handler{gtw: gtw}
}

func (h handler) GetTopISPsFromSwitzerland(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	topTenCH, err := h.gtw.GetTopISPFromSwitzerland(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	response, err := json.Marshal(topTenCH)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(response)
}

func (h handler) GetIPCountByCountryName(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	timer := time.Now()
	countyName, err := common.GetParamFromRequest(r, "country_name")
	if err != nil {
		err = fmt.Errorf("param: country_name %w", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ipCount, err := h.gtw.GetIpCountByCountyName(ctx, countyName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	response, err := json.Marshal(struct {
		CountryName        string `json:"country_name"`
		IpCount            int64  `json:"ip_count"`
		ElapsedTimeInMilis int64  `json:"elapsed_time_in_milis"`
	}{countyName, ipCount, time.Now().Sub(timer).Milliseconds()})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(response)

}

func (h handler) GetDataFromIP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ip, err := common.GetParamFromRequest(r, "ip")
	if err != nil {
		err = fmt.Errorf("param: ip %w", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ipData, err := h.gtw.GetDataFromIP(ctx, ip)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(ipData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(response)
}

func (h handler) SelectTopISPByCountryCode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	countryCode, err := common.GetParamFromRequest(r, "country_code")
	if err != nil {
		err = fmt.Errorf("param: country_code %w", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	topISPs, err := h.gtw.GetIspIpsByCountryCode(ctx, countryCode, 100)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(topISPs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(response)
}
