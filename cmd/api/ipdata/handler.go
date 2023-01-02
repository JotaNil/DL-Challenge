package ipdata

import (
	"DreamLabChallenge/cmd/api/common"
	"encoding/json"
	"fmt"
	"net/http"
)

type Handler interface {
	GetTopISPsFromSwitzerland(w http.ResponseWriter, r *http.Request)
	GetIPCountByCountryName(w http.ResponseWriter, r *http.Request)
	GetDataFromIP(w http.ResponseWriter, r *http.Request)
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
		common.HandlerErrorResponse(w, err)
		return
	}

	response, err := json.Marshal(topTenCH)
	if err != nil {
		common.HandlerErrorResponse(w, err)
		return
	}
	w.Write(response)
}

func (h handler) GetIPCountByCountryName(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	countyName, err := common.GetParamFromRequest(r, "country_name")
	if err != nil {
		err = fmt.Errorf("param: country_name %w", err)
		common.HandlerErrorResponse(w, err)
		return
	}
	ipCount, err := h.gtw.GetIpCountByCountryName(ctx, countyName)
	if err != nil {
		common.HandlerErrorResponse(w, err)
		return
	}

	response, err := json.Marshal(struct {
		CountryName string `json:"country_name"`
		IpCount     int64  `json:"ip_count"`
	}{countyName, ipCount})

	if err != nil {
		common.HandlerErrorResponse(w, err)
		return
	}

	w.Write(response)

}

func (h handler) GetDataFromIP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ip, err := common.GetParamFromRequest(r, "ip")
	if err != nil {
		err = fmt.Errorf("param: ip %w", err)
		common.HandlerErrorResponse(w, err)
		return
	}

	isValid := isValidIp(ip)
	if !isValid {
		err = fmt.Errorf("ip is not a valid Ipv4 format %w", common.ErrorBadRequest)
		common.HandlerErrorResponse(w, err)
		return
	}

	ipData, err := h.gtw.GetDataFromIP(ctx, ip)
	if err != nil {
		common.HandlerErrorResponse(w, err)
		return
	}

	response, err := json.Marshal(ipData)
	if err != nil {
		common.HandlerErrorResponse(w, err)
		return
	}

	w.Write(response)
}
