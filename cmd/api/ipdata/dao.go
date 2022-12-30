package ipdata

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
)

const (
	host   = "localhost"
	port   = 5432
	user   = "postgres"
	dbname = "ipv4-proxy-dreamlab"

	selectQuery     = "SELECT ip_from,ip_to,country_code,country_name,isp FROM proxydata.ip2location WHERE 1=1"
	selectByIPQuery = "SELECT ip_from,ip_to,country_code,country_name,isp,region_name,city_name,proxy_type FROM proxydata.ip2location WHERE $1 BETWEEN ip_from AND ip_to"
)

var password = os.Getenv("bdpassword")

type Dao interface {
	Get(ctx context.Context, filters *IpData) ([]IpData, error)
	GetByIp(ctx context.Context, ip int64) (IpData, error)
}

func NewDao() Dao {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	return dao{db: db}
}

type dao struct {
	db *sql.DB
}

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
			&data.CountyName,
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

func (d dao) GetByIp(ctx context.Context, ip int64) (IpData, error) {
	rows := d.db.QueryRow(selectByIPQuery, ip)

	ipData := IpData{}
	err := rows.Scan(&ipData.IpFrom,
		&ipData.IpTo,
		&ipData.CountryCode,
		&ipData.CountyName,
		&ipData.ISP,
		&ipData.RegionName,
		&ipData.CityName,
		&ipData.ProxyType)
	if err != nil {
		return IpData{}, err
	}

	err = rows.Err()
	if err != nil {
		err = fmt.Errorf("error with get query while scanning rows.  %w", err)
		return IpData{}, err
	}

	return ipData, nil
}

func processFilters(query string, filters *IpData) string {
	if filters != nil {
		if filters.ISP != "" {
			query = query + " AND isp = '" + filters.ISP + "'"
		}
		if filters.CountryCode != "" {
			query = query + " AND country_code = '" + filters.CountryCode + "'"
		}
		if filters.CountyName != "" {
			query = query + " AND country_name = '" + filters.CountyName + "'"
		}
	}

	return query
}
