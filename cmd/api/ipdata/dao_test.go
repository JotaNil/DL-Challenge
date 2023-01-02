package ipdata

import (
	"DreamLabChallenge/cmd/api/common"
	"DreamLabChallenge/cmd/services"
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

/*
	GetByIp(ctx context.Context, ip int64) (IpData, error)
	GetIpSumByCountry(ctx context.Context, countryName string) (int64, error)
	GetTopIspByCountryCode(ctx context.Context, countryCode string, limit int) ([]IspIpCount, error)
*/

func TestDao_GetByIp(t *testing.T) {
	tests := []common.TestCase{
		{Scenario: "No error", TestFn: testDaoGetByIpNoError},
		{Scenario: "No rows error", TestFn: testDaoGetByIpRowNotFoundError},
		{Scenario: "Connection error", TestFn: testDaoGetByIpRowConnectionError},
	}

	for _, testCase := range tests {
		t.Run(testCase.Scenario, testCase.TestFn)
	}
}

func TestDao_GetIpSumByCountry(t *testing.T) {
	tests := []common.TestCase{
		{Scenario: "No error", TestFn: testDaoGetIpSumByCountryNoError},
		{Scenario: "No rows error", TestFn: testDaoGetIpSumByCountryNotFoundError},
		{Scenario: "Connection error", TestFn: testDaoGetIpSumByCountryConnectionError},
	}

	for _, testCase := range tests {
		t.Run(testCase.Scenario, testCase.TestFn)
	}

}

func TestDao_GetTopIspByCountryCode(t *testing.T) {
	tests := []common.TestCase{
		{Scenario: "No error", TestFn: testDaoGetTopIspByCountryCodeNoError},
		{Scenario: "No rows error", TestFn: testDaoGetTopIspByCountryCodeNotFoundError},
		{Scenario: "Connection error", TestFn: testDaoGetTopIspByCountryConnectionError},
	}

	for _, testCase := range tests {
		t.Run(testCase.Scenario, testCase.TestFn)
	}

}

// GetByIp

func testDaoGetByIpNoError(t *testing.T) {
	type test struct {
		ip     int64
		rows   *sqlmock.Rows
		output IpData
		err    error
	}

	testData := test{ip: 2130706433, rows: getIpDataRows(), output: mockIpDataDao, err: nil}
	mockDB, mockHandler := services.ConnectToSQLDB(services.MockDB)
	mockDao := NewDao(mockDB)

	mockHandler.ExpectQuery(regexp.QuoteMeta(selectByIPQuery)).WithArgs(testData.ip).WillReturnRows(testData.rows)

	output, err := mockDao.GetByIp(context.Background(), testData.ip)
	assert.Equal(t, testData.output, output)
	assert.True(t, errors.Is(err, testData.err))
}

func testDaoGetByIpRowNotFoundError(t *testing.T) {
	type test struct {
		ip     int64
		rows   *sqlmock.Rows
		output IpData
		err    error
	}

	rowsWithError := getIpDataRows().RowError(0, sql.ErrNoRows)
	testData := test{ip: 2130706433, rows: rowsWithError, output: IpData{}, err: common.ErrorNotFound}
	mockDB, mockHandler := services.ConnectToSQLDB(services.MockDB)
	mockDao := NewDao(mockDB)

	mockHandler.ExpectQuery(regexp.QuoteMeta(selectByIPQuery)).WithArgs(testData.ip).WillReturnRows(testData.rows)

	output, err := mockDao.GetByIp(context.Background(), testData.ip)
	assert.Equal(t, testData.output, output)
	assert.True(t, errors.Is(err, testData.err))

}

func testDaoGetByIpRowConnectionError(t *testing.T) {
	type test struct {
		ip     int64
		rows   *sqlmock.Rows
		output IpData
		err    error
	}

	rowsWithError := getIpDataRows().RowError(0, errors.New("connectionError"))
	testData := test{ip: 2130706433, rows: rowsWithError, output: IpData{}, err: common.ErrorInternalServer}
	mockDB, mockHandler := services.ConnectToSQLDB(services.MockDB)
	mockDao := NewDao(mockDB)

	mockHandler.ExpectQuery(regexp.QuoteMeta(selectByIPQuery)).WithArgs(testData.ip).WillReturnRows(testData.rows)

	output, err := mockDao.GetByIp(context.Background(), testData.ip)
	assert.Equal(t, testData.output, output)
	assert.True(t, errors.Is(err, testData.err))

}

// GetIpSumByCountry

func testDaoGetIpSumByCountryNoError(t *testing.T) {
	type test struct {
		countryName string
		rows        *sqlmock.Rows
		output      int64
		err         error
	}

	testData := test{countryName: "Ireland", rows: getIpSumByCountryRows(), output: 42, err: nil}
	mockDB, mockHandler := services.ConnectToSQLDB(services.MockDB)
	mockDao := NewDao(mockDB)

	mockHandler.ExpectQuery(regexp.QuoteMeta(getIPsPerCountryQuery)).WithArgs(testData.countryName).WillReturnRows(testData.rows)

	output, err := mockDao.GetIpSumByCountry(context.Background(), testData.countryName)
	assert.Equal(t, testData.output, output)
	assert.True(t, errors.Is(err, testData.err))

}

func testDaoGetIpSumByCountryNotFoundError(t *testing.T) {
	type test struct {
		countryName string
		rows        *sqlmock.Rows
		output      int64
		err         error
	}

	rowsWithError := getIpSumByCountryRows().RowError(0, sql.ErrNoRows)
	testData := test{countryName: "Ireland", rows: rowsWithError, output: 0, err: common.ErrorNotFound}
	mockDB, mockHandler := services.ConnectToSQLDB(services.MockDB)
	mockDao := NewDao(mockDB)

	mockHandler.ExpectQuery(regexp.QuoteMeta(getIPsPerCountryQuery)).WithArgs(testData.countryName).WillReturnRows(testData.rows)

	output, err := mockDao.GetIpSumByCountry(context.Background(), testData.countryName)
	assert.Equal(t, testData.output, output)
	assert.True(t, errors.Is(err, testData.err))

}

func testDaoGetIpSumByCountryConnectionError(t *testing.T) {
	type test struct {
		countryName string
		rows        *sqlmock.Rows
		output      int64
		err         error
	}

	rowsWithError := getIpSumByCountryRows().RowError(0, errors.New("connection error"))
	testData := test{countryName: "Ireland", rows: rowsWithError, output: 0, err: common.ErrorInternalServer}
	mockDB, mockHandler := services.ConnectToSQLDB(services.MockDB)
	mockDao := NewDao(mockDB)

	mockHandler.ExpectQuery(regexp.QuoteMeta(getIPsPerCountryQuery)).WithArgs(testData.countryName).WillReturnRows(testData.rows)

	output, err := mockDao.GetIpSumByCountry(context.Background(), testData.countryName)
	assert.Equal(t, testData.output, output)
	assert.True(t, errors.Is(err, testData.err))

}

// GetTopIspByCountryCode

func testDaoGetTopIspByCountryCodeNoError(t *testing.T) {
	type test struct {
		countryCode string
		limit       int
		rows        *sqlmock.Rows
		output      []IspIpCount
		err         error
	}

	testData := test{countryCode: "AR", limit: 10, rows: getTopIspByCountryCodeRows(), output: mockIspIpCountDao, err: nil}
	mockDB, mockHandler := services.ConnectToSQLDB(services.MockDB)
	mockDao := NewDao(mockDB)

	mockHandler.ExpectQuery(regexp.QuoteMeta(getTopIspByCountryCode)).WithArgs(testData.countryCode, testData.limit).WillReturnRows(testData.rows)

	output, err := mockDao.GetTopIspByCountryCode(context.Background(), testData.countryCode, testData.limit)
	assert.Equal(t, testData.output, output)
	assert.True(t, errors.Is(err, testData.err))

}

func testDaoGetTopIspByCountryCodeNotFoundError(t *testing.T) {
	type test struct {
		countryCode string
		limit       int
		rows        *sqlmock.Rows
		output      []IspIpCount
		err         error
	}

	rowsWithError := getTopIspByCountryCodeRows().RowError(2, sql.ErrNoRows)
	testData := test{countryCode: "AR", limit: 10, rows: rowsWithError, output: []IspIpCount{}, err: common.ErrorNotFound}
	mockDB, mockHandler := services.ConnectToSQLDB(services.MockDB)
	mockDao := NewDao(mockDB)

	mockHandler.ExpectQuery(regexp.QuoteMeta(getTopIspByCountryCode)).WithArgs(testData.countryCode, testData.limit).WillReturnRows(testData.rows)

	output, err := mockDao.GetTopIspByCountryCode(context.Background(), testData.countryCode, testData.limit)
	assert.Equal(t, testData.output, output)
	assert.True(t, errors.Is(err, testData.err))

}
func testDaoGetTopIspByCountryConnectionError(t *testing.T) {
	type test struct {
		countryCode string
		limit       int
		rows        *sqlmock.Rows
		output      []IspIpCount
		err         error
	}

	rowsWithError := getTopIspByCountryCodeRows().RowError(0, errors.New("connection error"))
	testData := test{countryCode: "AR", limit: 10, rows: rowsWithError, output: []IspIpCount{}, err: common.ErrorInternalServer}
	mockDB, mockHandler := services.ConnectToSQLDB(services.MockDB)
	mockDao := NewDao(mockDB)

	mockHandler.ExpectQuery(regexp.QuoteMeta(getTopIspByCountryCode)).WithArgs(testData.countryCode, testData.limit).WillReturnRows(testData.rows)

	output, err := mockDao.GetTopIspByCountryCode(context.Background(), testData.countryCode, testData.limit)
	assert.Equal(t, testData.output, output)
	assert.True(t, errors.Is(err, testData.err))

}

// mock utils

var mockIpDataDao = IpData{
	IpFrom:      2130706433,
	IpTo:        2130706433,
	ProxyType:   "PUB",
	CountryCode: "ES",
	CountryName: "Spain",
	RegionName:  "Spain",
	CityName:    "Barcelona",
	ISP:         "IPS",
}

var mockIspIpCountDao = utilGenerateIspIpCount(10)

func getIpDataRows() *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{"ip_from", "ip_to", "country_code", "country_name", "isp", "region_name", "city_name", "proxy_type"})
	rows.AddRow(mockIpDataGateway.IpFrom, mockIpDataGateway.IpTo, mockIpDataGateway.CountryCode, mockIpDataGateway.CountryName, mockIpDataGateway.ISP, mockIpDataGateway.RegionName, mockIpDataGateway.CityName, mockIpDataGateway.ProxyType)

	return rows
}

func getTopIspByCountryCodeRows() *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{"isp", "difference"})
	for _, count := range mockIspIpCountDao {
		rows.AddRow(count.Isp, count.IpCount)
	}

	return rows
}

func getIpSumByCountryRows() *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{"sum"})
	rows.AddRow(42)

	return rows
}
