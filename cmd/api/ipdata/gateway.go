package ipdata

import (
	"DreamLabChallenge/cmd/api/common"
	"context"
	"fmt"
)

//go:generate mockgen -destination=mock_gateway.go -package=ipdata -source=gateway.go Gateway

type Gateway interface {
	// GetIspIpsByCountryCode returns the top (limit) ISPs of the (countryCode) given
	GetIspIpsByCountryCode(ctx context.Context, countryCode string, limit int) ([]IspIpCount, error)
	// GetIpCountByCountryName returns the number of Ips of the given countryName
	GetIpCountByCountryName(ctx context.Context, countryName string) (int64, error)
	// GetTopISPFromSwitzerland returns a list of the top 10 ISPs based on how many IPs does it have
	GetTopISPFromSwitzerland(ctx context.Context) ([]IspIpCount, error)
	// GetDataFromIP gets the data associated from the given IP in string(xxx.xxx.xxx.xxx) format
	GetDataFromIP(ctx context.Context, ip string) (IpData, error)
}

type gateway struct {
	dao Dao
}

func NewGateway(dao Dao) Gateway {
	return gateway{dao: dao}
}

// GetIspIpsByCountryCode returns the top (limit) ISPs of the (countryCode) given
func (g gateway) GetIspIpsByCountryCode(ctx context.Context, countryCode string, limit int) ([]IspIpCount, error) {
	isValid := isValidCountryCode(countryCode)
	if !isValid {
		return []IspIpCount{}, fmt.Errorf("invalid country_code  %w", common.ErrorBadRequest)

	}
	countryIPData, err := g.dao.GetTopIspByCountryCode(ctx, countryCode, limit)
	if err != nil {
		return []IspIpCount{}, err
	}

	return countryIPData, nil
}

// GetIpCountByCountryName returns the number of Ips of the given countryName
func (g gateway) GetIpCountByCountryName(ctx context.Context, countryName string) (int64, error) {
	isValid := isValidCountryName(countryName)
	if !isValid {
		return 0, fmt.Errorf("invalid country_code  %w", common.ErrorBadRequest)
	}
	countryIPCount, err := g.dao.GetIpSumByCountry(ctx, countryName)
	if err != nil {
		return 0, err
	}

	return countryIPCount, nil
}

// GetTopISPFromSwitzerland returns a list of the top 10 ISPs based on how many IPs does it have
func (g gateway) GetTopISPFromSwitzerland(ctx context.Context) ([]IspIpCount, error) {
	sortedISPs, err := g.GetIspIpsByCountryCode(ctx, CountryCodeSwitzerland, 10)
	if err != nil {
		err = fmt.Errorf("error getting Ips Ip count.  %w", err)
		return []IspIpCount{}, err
	}

	return sortedISPs, err
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
