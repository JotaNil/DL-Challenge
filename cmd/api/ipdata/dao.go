package ipdata

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

const (
	selectQuery            = "SELECT ip_from,ip_to,country_code,country_name,isp FROM proxydata.ip2location WHERE 1=1"
	getIPsPerCountryQuery  = "SELECT SUM(ip_to - ip_from + 1) FROM proxydata.ip2location WHERE country_name = $1 "
	getTopIspByCountryCode = "SELECT isp, sum(ip_to-ip_from+1) as difference FROM proxydata.ip2location WHERE country_code = $1 GROUP BY isp order by difference DESC LIMIT $2"
	selectByIPQuery        = "SELECT ip_from,ip_to,country_code,country_name,isp,region_name,city_name,proxy_type FROM proxydata.ip2location WHERE $1 BETWEEN ip_from AND ip_to"
)

//go:generate mockgen -destination=mock_dao.go -package=ipdata -source=dao.go Dao

type Dao interface {
	Get(ctx context.Context, filters *IpData) ([]IpData, error)
	GetByIp(ctx context.Context, ip int64) (IpData, error)
	GetIpSumByCountry(ctx context.Context, countryName string) (int64, error)
	GetTopIspByCountryCode(ctx context.Context, countryCode string, limit int) ([]IspIpCount, error)
}

func NewDao(dbConnection *sql.DB) Dao {
	return dao{db: dbConnection}
}

type dao struct {
	db *sql.DB
}

// Get uses IpData to filter by IpData.ISP, IpData.CountryName and IpData.CountryCode
func (d dao) Get(ctx context.Context, filters *IpData) ([]IpData, error) {
	rows, err := d.db.Query(processFilters(selectQuery, filters))
	if err != nil {
		err = fmt.Errorf("error with get query with DB.  %w", err)
		return []IpData{}, err
	}
	defer rows.Close()

	ipData := make([]IpData, 0)
	for rows.Next() {
		data := IpData{}
		err = rows.Scan(&data.IpFrom,
			&data.IpTo,
			&data.CountryCode,
			&data.CountryName,
			&data.ISP)
		if err != nil {
			return ipData, err
		}
		ipData = append(ipData, data)
	}
	err = rows.Err()
	if err != nil {
		err = fmt.Errorf("error with get query while scanning rows.  %w", err)
		return ipData, err
	}

	return ipData, nil
}

// GetTopIspByCountryCode get the top (limit) ISPs from the given countryCode
func (d dao) GetTopIspByCountryCode(ctx context.Context, countryCode string, limit int) ([]IspIpCount, error) {
	rows, err := d.db.QueryContext(ctx, getTopIspByCountryCode, countryCode, limit)
	if err != nil {
		err = fmt.Errorf("error with get query with DB.  %w", err)
		return []IspIpCount{}, err
	}
	defer rows.Close()

	ipData := make([]IspIpCount, 0)
	for rows.Next() {
		data := IspIpCount{}
		err = rows.Scan(&data.Isp, &data.IpCount)
		if err != nil {
			return ipData, err
		}
		ipData = append(ipData, data)
	}
	err = rows.Err()
	if err != nil {
		err = fmt.Errorf("error with get query while scanning rows.  %w", err)
		return ipData, err
	}

	return ipData, nil
}

// GetByIp gets all the data of the given ip in decimal format
func (d dao) GetByIp(ctx context.Context, ip int64) (IpData, error) {
	row := d.db.QueryRowContext(ctx, selectByIPQuery, ip)

	ipData := IpData{}
	err := row.Scan(&ipData.IpFrom,
		&ipData.IpTo,
		&ipData.CountryCode,
		&ipData.CountryName,
		&ipData.ISP,
		&ipData.RegionName,
		&ipData.CityName,
		&ipData.ProxyType)
	if err != nil {
		return IpData{}, err
	}

	err = row.Err()
	if err != nil {
		err = fmt.Errorf("error with get query while scanning rows.  %w", err)
		return IpData{}, err
	}

	return ipData, nil
}

// GetIpSumByCountry gets the number of Ips of the given countryName
func (d dao) GetIpSumByCountry(ctx context.Context, countryName string) (int64, error) {
	row := d.db.QueryRowContext(ctx, getIPsPerCountryQuery, countryName)

	var ipSum int64
	err := row.Scan(&ipSum)
	if err != nil {
		return 0, err
	}

	err = row.Err()
	if err != nil {
		err = fmt.Errorf("error with get query while scanning rows.  %w", err)
		return 0, err
	}

	return ipSum, nil
}

func processFilters(query string, filters *IpData) string {
	if filters != nil {
		if filters.ISP != "" {
			query = query + " AND isp = '" + filters.ISP + "'"
		}
		if filters.CountryCode != "" {
			query = query + " AND country_code = '" + filters.CountryCode + "'"
		}
		if filters.CountryName != "" {
			query = query + " AND country_name = '" + filters.CountryName + "'"
		}
	}

	return query
}
