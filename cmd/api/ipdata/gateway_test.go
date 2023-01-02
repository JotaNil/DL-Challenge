package ipdata

import (
	"DreamLabChallenge/cmd/api/common"
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGateway_GetIspIpsByCountryCode(t *testing.T) {
	tests := []common.TestCase{
		{Scenario: "No error", TestFn: testGtwGetIspIpsByCountryCodeNoError},
		{Scenario: "Dao thrown error", TestFn: testGtwGetIspIpsByCountryCodeDBError},
		{Scenario: "Invalid country code error", TestFn: testGtwGetIspIpsByCountryInvalidCodeError},
	}

	for _, testCase := range tests {
		t.Run(testCase.Scenario, testCase.TestFn)
	}
}

func TestGateway_GetIpCountByCountryName(t *testing.T) {
	tests := []common.TestCase{
		{Scenario: "No error", TestFn: testGtwGetIpCountByCountryNameNoErrors},
		{Scenario: "Dao thrown error", TestFn: testGtwGetIpCountByCountryNameDBError},
		{Scenario: "Invalid country error", TestFn: testGtwGetIpCountByCountryNameInvalidCountryNameError},
	}

	for _, testCase := range tests {
		t.Run(testCase.Scenario, testCase.TestFn)
	}
}

func TestGateway_GetTopIspFromSwitzerland(t *testing.T) {
	tests := []common.TestCase{
		{Scenario: "No error", TestFn: testGtwGetTopIspFromSwitzerlandNoError},
		{Scenario: "Dao thrown error", TestFn: testGtwGetTopIspFromSwitzerlandDBError},
	}

	for _, testCase := range tests {
		t.Run(testCase.Scenario, testCase.TestFn)
	}
}

func TestGetDataFromIP(t *testing.T) {
	tests := []common.TestCase{
		{Scenario: "No error", TestFn: testGtwGetDataFromIPNoError},
		{Scenario: "Dao thrown error", TestFn: testGtwGetDataFromIPDbError},
	}

	for _, testCase := range tests {
		t.Run(testCase.Scenario, testCase.TestFn)
	}
}

// GetIspIpsByCountryCode

func testGtwGetIspIpsByCountryCodeNoError(t *testing.T) {
	type test struct {
		code   string
		limit  int
		output []IspIpCount
		err    error
	}
	testData := test{code: "AR", limit: 1, output: []IspIpCount{{Isp: "mainIsp", IpCount: 1234}}, err: nil}

	ctrl := gomock.NewController(t)
	mockDao := NewMockDao(ctrl)
	gtw := NewGateway(mockDao)

	mockDao.EXPECT().
		GetTopIspByCountryCode(gomock.Any(), testData.code, testData.limit).
		Return(testData.output, testData.err)

	output, err := gtw.GetIspIpsByCountryCode(context.Background(), testData.code, testData.limit)

	assert.Equal(t, testData.output, output)
	assert.True(t, errors.Is(err, testData.err))
}

func testGtwGetIspIpsByCountryCodeDBError(t *testing.T) {
	type test struct {
		code   string
		limit  int
		output []IspIpCount
		err    error
	}
	testData := test{code: "AR", limit: 1, output: []IspIpCount{}, err: errors.New("connection error")}

	ctrl := gomock.NewController(t)
	mockDao := NewMockDao(ctrl)
	gtw := NewGateway(mockDao)

	mockDao.EXPECT().
		GetTopIspByCountryCode(gomock.Any(), testData.code, testData.limit).
		Return(testData.output, testData.err)

	output, err := gtw.GetIspIpsByCountryCode(context.Background(), testData.code, testData.limit)

	assert.Equal(t, testData.output, output)
	assert.True(t, errors.Is(err, testData.err))
}

func testGtwGetIspIpsByCountryInvalidCodeError(t *testing.T) {
	type test struct {
		code   string
		limit  int
		output []IspIpCount
		err    error
	}
	testData := test{code: "badCode", limit: 1, output: []IspIpCount{}, err: common.ErrorBadRequest}

	ctrl := gomock.NewController(t)
	mockDao := NewMockDao(ctrl)
	gtw := NewGateway(mockDao)

	output, err := gtw.GetIspIpsByCountryCode(context.Background(), testData.code, testData.limit)

	assert.Equal(t, testData.output, output)
	assert.True(t, errors.Is(err, testData.err))
}

// GetIpCountByCountryName

func testGtwGetIpCountByCountryNameNoErrors(t *testing.T) {
	type test struct {
		name        string
		countryName string
		output      int64
		err         error
	}
	testData := test{name: "No error", countryName: "Ireland", output: 1526, err: nil}

	ctrl := gomock.NewController(t)
	mockDao := NewMockDao(ctrl)
	gtw := NewGateway(mockDao)

	mockDao.EXPECT().
		GetIpSumByCountry(gomock.Any(), testData.countryName).
		Return(testData.output, testData.err)

	output, err := gtw.GetIpCountByCountryName(context.Background(), testData.countryName)

	assert.Equal(t, testData.output, output)
	assert.True(t, errors.Is(err, testData.err))
}

func testGtwGetIpCountByCountryNameDBError(t *testing.T) {
	type test struct {
		countryName string
		output      int64
		err         error
	}
	testData := test{countryName: "Ireland", output: 0, err: errors.New("connection error")}

	ctrl := gomock.NewController(t)
	mockDao := NewMockDao(ctrl)
	gtw := NewGateway(mockDao)

	mockDao.EXPECT().
		GetIpSumByCountry(gomock.Any(), testData.countryName).
		Return(testData.output, testData.err)

	output, err := gtw.GetIpCountByCountryName(context.Background(), testData.countryName)

	assert.Equal(t, testData.output, output)
	assert.True(t, errors.Is(err, testData.err))
}

func testGtwGetIpCountByCountryNameInvalidCountryNameError(t *testing.T) {
	type test struct {
		countryName string
		output      int64
		err         error
	}
	testData := test{countryName: "invalid country", output: 0, err: common.ErrorBadRequest}

	ctrl := gomock.NewController(t)
	mockDao := NewMockDao(ctrl)
	gtw := NewGateway(mockDao)

	output, err := gtw.GetIpCountByCountryName(context.Background(), testData.countryName)

	assert.Equal(t, testData.output, output)
	assert.True(t, errors.Is(err, testData.err))
}

// GetTopIspFromSwitzerland

func testGtwGetTopIspFromSwitzerlandNoError(t *testing.T) {
	type test struct {
		code   string
		limit  int
		output []IspIpCount
		err    error
	}
	testData := test{code: "CH", limit: 10, output: utilGenerateIspIpCount(10), err: nil}

	ctrl := gomock.NewController(t)
	mockDao := NewMockDao(ctrl)
	gtw := NewGateway(mockDao)

	mockDao.EXPECT().
		GetTopIspByCountryCode(gomock.Any(), testData.code, testData.limit).
		Return(testData.output, testData.err)

	output, err := gtw.GetTopISPFromSwitzerland(context.Background())

	assert.Equal(t, testData.output, output)
	assert.True(t, errors.Is(err, testData.err))
}

func testGtwGetTopIspFromSwitzerlandDBError(t *testing.T) {
	type test struct {
		code   string
		limit  int
		output []IspIpCount
		err    error
	}
	testData := test{code: "CH", limit: 10, output: []IspIpCount{}, err: errors.New("connection error")}

	ctrl := gomock.NewController(t)
	mockDao := NewMockDao(ctrl)
	gtw := NewGateway(mockDao)

	mockDao.EXPECT().
		GetTopIspByCountryCode(gomock.Any(), testData.code, testData.limit).
		Return(testData.output, testData.err)

	output, err := gtw.GetTopISPFromSwitzerland(context.Background())

	assert.Equal(t, testData.output, output)
	assert.True(t, errors.Is(err, testData.err))
}

// GetDataFromIP

func testGtwGetDataFromIPNoError(t *testing.T) {
	type test struct {
		ipString  string
		ipDecimal int64
		output    IpData
		err       error
	}
	testData := test{ipString: "127.0.0.1", ipDecimal: 2130706433, output: mockIpDataGateway, err: nil}

	ctrl := gomock.NewController(t)
	mockDao := NewMockDao(ctrl)
	gtw := NewGateway(mockDao)

	mockDao.EXPECT().
		GetByIp(gomock.Any(), testData.ipDecimal).
		Return(testData.output, testData.err)

	output, err := gtw.GetDataFromIP(context.Background(), testData.ipString)

	assert.Equal(t, testData.output, output)
	assert.True(t, errors.Is(err, testData.err))
}

func testGtwGetDataFromIPDbError(t *testing.T) {
	type test struct {
		ipString  string
		ipDecimal int64
		output    IpData
		err       error
	}
	testData := test{ipString: "127.0.0.1", ipDecimal: 2130706433, output: IpData{}, err: errors.New("db error")}

	ctrl := gomock.NewController(t)
	mockDao := NewMockDao(ctrl)
	gtw := NewGateway(mockDao)

	mockDao.EXPECT().
		GetByIp(gomock.Any(), testData.ipDecimal).
		Return(testData.output, testData.err)

	output, err := gtw.GetDataFromIP(context.Background(), testData.ipString)

	assert.Equal(t, testData.output, output)
	assert.True(t, errors.Is(err, testData.err))
}

// mock utils

var mockIpDataGateway = IpData{
	IpFrom:      2130706433,
	IpTo:        2130706433,
	ProxyType:   "PUB",
	CountryCode: "ES",
	CountryName: "Spain",
	RegionName:  "Spain",
	CityName:    "Barcelona",
	ISP:         "IPS",
	IpString:    "127.0.0.1",
}

func utilGenerateIspIpCount(count int) []IspIpCount {
	mockData := make([]IspIpCount, 0)
	for i := 0; i < count; i++ {
		mockData = append(mockData, IspIpCount{Isp: string(rune('A' - 1 + i)), IpCount: int64(i * 10)})
	}

	return mockData
}
