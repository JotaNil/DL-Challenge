package ipdata

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

const (
	CountryCodeSwitzerland = "CH"
)

type IpData struct {
	IpFrom      int64    `json:"ip_from,omitempty"`
	IpTo        int64    `json:"ip_to,omitempty"`
	ProxyType   string   `json:"proxy_type,omitempty"`
	CountryCode string   `json:"country_code,omitempty"`
	CountyName  string   `json:"county_name,omitempty"`
	RegionName  string   `json:"region_name,omitempty"`
	CityName    string   `json:"city_name,omitempty"`
	ISP         string   `json:"isp,omitempty"`
	IpString    string   `json:"ip_string,omitempty"`
	IPList      []string `json:"ip_list,omitempty"`
}

type IspIpCount struct {
	Isp     string `json:"isp"`
	IpCount int64  `json:"ip_count"`
}

func (i IpData) GetIPCount() int64 {
	return i.IpTo - i.IpFrom + 1
}

func decimalIPToString(ip int64) string {
	ipFormat := "%d.%d.%d.%d"
	ipSegments := make([]int, 4)
	for i := 3; i >= 0; i-- {
		whole, decimal := math.Modf(float64(ip) / float64(ipConverterWeights[i]))
		ipSegments[i] = int(whole)
		ip = int64(decimal * float64(ipConverterWeights[i]))
	}
	return fmt.Sprintf(ipFormat, ipSegments[3], ipSegments[2], ipSegments[1], ipSegments[0])
}

func stringIPToDecimal(ip string) int64 {
	segments := strings.Split(ip, ".")
	var decimalValue int64
	for i := 3; i >= 0; i-- {
		segmentValue, _ := strconv.ParseInt(segments[3-i], 10, 64)
		segmentValueWeighted := segmentValue * ipConverterWeights[i]
		decimalValue += segmentValueWeighted
	}
	return decimalValue
}

var ipConverterWeights = make(map[int]int64)

func buildIPConverterWeights() {
	for i := 0; i <= 3; i++ {
		ipConverterWeights[i] = int64(math.Pow(256, float64(i)))
	}
}

func init() {
	buildIPConverterWeights()
}
