package ipdata

import (
	"context"
	"fmt"
)

type Gateway interface {
	SelectTopHundred(ctx context.Context) ([]IpData, error)
	GetIspIpsByCountryCode(ctx context.Context, countryCode string, limit int) ([]IspIpCount, error)
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
	sortedISPs, err := g.GetIspIpsByCountryCode(ctx, CountryCodeSwitzerland, 10)
	if err != nil {
		err = fmt.Errorf("error getting Ips Ip count.  %w", err)
		return []IspIpCount{}, err
	}

	return sortedISPs[:10], err
}

func (g gateway) GetIspIpsByCountryCode(ctx context.Context, countryCode string, limit int) ([]IspIpCount, error) {
	countryIPData, err := g.dao.GetTopIspByCountryCode(ctx, countryCode, limit)
	if err != nil {
		return []IspIpCount{}, err
	}

	return countryIPData, nil
}

func (g gateway) GetIpCountByCountyName(ctx context.Context, countryName string) (int64, error) {
	countryIPCount, err := g.dao.GetIpSumByCountry(ctx, countryName)
	if err != nil {
		return 0, err
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
