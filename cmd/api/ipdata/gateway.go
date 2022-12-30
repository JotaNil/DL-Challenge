package ipdata

import (
	"context"
	"fmt"
	"sort"
)

type Gateway interface {
	SelectTopHundred(ctx context.Context) ([]IpData, error)
	GetIspIpsByCountryCode(ctx context.Context, countryCode string) ([]IspIpCount, error)
	GetTopISPFromSwitzerland(ctx context.Context) ([]IspIpCount, error)
	GetIpCountByCountyName(ctx context.Context, countryName string) (int64, error)
	GetDataFromIP(ctx context.Context, ip string) (IpData, error)
}

type gateway struct {
	dao Dao
}

func NewGateway(dao Dao) Gateway {
	return gateway{dao: dao}
}

func (g gateway) SelectTopHundred(ctx context.Context) ([]IpData, error) {
	data, err := g.dao.Get(ctx, nil)
	if err != nil {
		return []IpData{}, err
	}

	return data, nil
}
func (g gateway) GetTopISPFromSwitzerland(ctx context.Context) ([]IspIpCount, error) {
	sortedISPs, err := g.GetIspIpsByCountryCode(ctx, CountryCodeSwitzerland)
	if err != nil {
		err = fmt.Errorf("error getting Ips Ip count.  %w", err)
		return []IspIpCount{}, err
	}

	return sortedISPs[:10], err
}

func (g gateway) GetIspIpsByCountryCode(ctx context.Context, countryCode string) ([]IspIpCount, error) {
	countryIPData, err := g.dao.Get(ctx, &IpData{CountryCode: countryCode})
	if err != nil {
		return []IspIpCount{}, err
	}

	countMap := make(map[string]int64)
	for i := range countryIPData {
		countMap[countryIPData[i].ISP] = countMap[countryIPData[i].ISP] + countryIPData[i].GetIPCount()
	}

	sortedISPs := make([]IspIpCount, 0)
	for isp, ipcount := range countMap {
		sortedISPs = append(sortedISPs, IspIpCount{Isp: isp, IpCount: ipcount})
	}
	sort.Slice(sortedISPs, func(i, j int) bool {
		return sortedISPs[i].IpCount > sortedISPs[j].IpCount
	})

	return sortedISPs, nil
}

func (g gateway) GetIpCountByCountyName(ctx context.Context, countryName string) (int64, error) {
	countryIPData, err := g.dao.Get(ctx, &IpData{CountyName: countryName})
	if err != nil {
		return 0, err
	}

	var countryIPCount int64
	for i := range countryIPData {
		countryIPCount = countryIPCount + countryIPData[i].GetIPCount()
	}

	return countryIPCount, nil
}
func (g gateway) GetDataFromIP(ctx context.Context, ip string) (IpData, error) {
	ipData, err := g.dao.GetByIp(ctx, stringIPToDecimal(ip))
	if err != nil {
		err = fmt.Errorf("error getting Ips Ip count.  %w", err) //TODO:better error
		return IpData{}, err
	}
	ipData.IpString = ip

	return ipData, nil
}
