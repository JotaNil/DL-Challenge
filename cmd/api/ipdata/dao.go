package ipdata

import (
	"DreamLabChallenge/cmd/api/common"
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
)

const (
	ipdataSchemaTableName = "proxydata.ip2location"

	getIPsPerCountryQuery  = "SELECT SUM(ip_to - ip_from + 1) FROM " + ipdataSchemaTableName + " WHERE country_name = $1 "
	getTopIspByCountryCode = "SELECT isp, sum(ip_to-ip_from+1) as difference FROM " + ipdataSchemaTableName + " WHERE country_code = $1 GROUP BY isp order by difference DESC LIMIT $2"
	selectByIPQuery        = "SELECT ip_from,ip_to,country_code,country_name,isp,region_name,city_name,proxy_type FROM " + ipdataSchemaTableName + " WHERE $1 BETWEEN ip_from AND ip_to"
)

//go:generate mockgen -destination=mock_dao.go -package=ipdata -source=dao.go Dao

type Dao interface {
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

// GetTopIspByCountryCode get the top (limit) ISPs from the given countryCode
func (d dao) GetTopIspByCountryCode(ctx context.Context, countryCode string, limit int) ([]IspIpCount, error) {
	rows, err := d.db.QueryContext(ctx, getTopIspByCountryCode, countryCode, limit)
	if err != nil {
		return []IspIpCount{}, err
	}
	defer rows.Close()

	ipData := make([]IspIpCount, 0)
	for rows.Next() {
		data := IspIpCount{}
		err := rows.Scan(&data.Isp, &data.IpCount)
		if err != nil {
			return ipData, err
		}
		ipData = append(ipData, data)
	}
	err = rows.Err()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = fmt.Errorf("error with get query with DB.  %w", common.ErrorNotFound)
			return []IspIpCount{}, err
		}
		err = fmt.Errorf("error with get query while scanning rows. %s %w", err.Error(), common.ErrorInternalServer)
		return []IspIpCount{}, err
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
		if errors.Is(err, sql.ErrNoRows) {
			err = fmt.Errorf("error with get query with DB.  %w", common.ErrorNotFound)
			return IpData{}, err
		}
		err = fmt.Errorf("error with get query while scanning rows. %s %w", err.Error(), common.ErrorInternalServer)
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
		if errors.Is(err, sql.ErrNoRows) {
			err = fmt.Errorf("error with get query with DB.  %w", common.ErrorNotFound)
			return 0, err
		}
		err = fmt.Errorf("error with get query while scanning rows. %s %w", err.Error(), common.ErrorInternalServer)
		return 0, err
	}

	return ipSum, nil
}
