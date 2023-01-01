package ipdata

import (
	"context"
	"fmt"
)

//go:generate mockgen -destination=mock_gateway.go -package=ipdata -source=gateway.go Gateway

type Gateway interface {
	// GetIspIpsByCountryCode returns the top (limit) ISPs of the (countryCode) given
	GetIspIpsByCountryCode(ctx context.Context, countryCode string, limit int) ([]IspIpCount, error)
	// GetTopISPFromSwitzerland returns a list of the top 10 ISPs based on how many IPs does it have
	GetTopISPFromSwitzerland(ctx context.Context) ([]IspIpCount, error)
	// GetIpCountByCountyName returns the number of Ips of the given countryName
	GetIpCountByCountyName(ctx context.Context, countryName string) (int64, error)
	// GetDataFromIP gets the data associated from the given IP in string(xxx.xxx.xxx.xxx) format
	GetDataFromIP(ctx context.Context, ip string) (IpData, error)
}

type gateway struct {
	dao Dao
}

func NewGateway(dao Dao) Gateway {
	return gateway{dao: dao}
}

func (g gateway) GetTopISPFromSwitzerland(ctx context.Context) ([]IspIpCount, error) {
	sortedISPs, err := g.GetIspIpsByCountryCode(ctx, CountryCodeSwitzerland, 10)
	if err != nil {
		err = fmt.Errorf("error getting Ips Ip count.  %w", err)
		return []IspIpCount{}, err
	}

	return sortedISPs[:10], err
}

// GetIspIpsByCountryCode returns the top (limit) ISPs of the (countryCode) given
func (g gateway) GetIspIpsByCountryCode(ctx context.Context, countryCode string, limit int) ([]IspIpCount, error) {
	countryIPData, err := g.dao.GetTopIspByCountryCode(ctx, countryCode, limit)
	if err != nil {
		return []IspIpCount{}, err
	}

	return countryIPData, nil
}

// GetIpCountByCountyName returns the number of Ips of the given countryName
func (g gateway) GetIpCountByCountyName(ctx context.Context, countryName string) (int64, error) {
	countryIPCount, err := g.dao.GetIpSumByCountry(ctx, countryName)
	if err != nil {
		return 0, err
	}

	return countryIPCount, nil
}

// GetDataFromIP gets the data associated from the given IP in string(xxx.xxx.xxx.xxx) format
func (g gateway) GetDataFromIP(ctx context.Context, ip string) (IpData, error) {
	ipData, err := g.dao.GetByIp(ctx, stringIPToDecimal(ip))
	if err != nil {
		err = fmt.Errorf("error getting Ips Ip count.  %w", err)
		return IpData{}, err
	}
	ipData.IpString = ip

	return ipData, nil
}
