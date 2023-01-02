package ipdata

import (
	"DreamLabChallenge/cmd/api/common"
	"encoding/json"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

/*
GetTopISPsFromSwitzerland(w http.ResponseWriter, r *http.Request)
	GetIPCountByCountryName(w http.ResponseWriter, r *http.Request)
	GetDataFromIP(w http.ResponseWriter, r *http.Request)
*/

func TestHandler_GetTopISPsFromSwitzerland(t *testing.T) {
	tests := []common.TestCase{
		{Scenario: "No error", TestFn: testHandlerGetTopISPsFromSwitzerlandNoErrors},
		{Scenario: "Gtw thrown error", TestFn: testHandlerGetTopISPsFromSwitzerlandDBError},
	}

	for _, testCase := range tests {
		t.Run(testCase.Scenario, testCase.TestFn)
	}

}

func TestHandler_GetIPCountByCountryName(t *testing.T) {
	tests := []common.TestCase{
		{Scenario: "No error", TestFn: testHandlerGetIpCountByCountryNameNoError},
		{Scenario: "No country param present error", TestFn: testHandlerGetIpCountByCountryNameNoCountryParamError},
		{Scenario: "Gateway thrown error", TestFn: testHandlerGetIpCountByCountryNameGtwError},
	}

	for _, testCase := range tests {
		t.Run(testCase.Scenario, testCase.TestFn)
	}

}

func TestHandler_GetDataFromIP(t *testing.T) {
	tests := []common.TestCase{
		{Scenario: "No error", TestFn: testHandlerGetDataFromIpNoError},
		{Scenario: "No ip param present error", TestFn: testHandlerGetDataFromIpNoIpParamError},
		{Scenario: "Invalid IP error", TestFn: testHandlerGetDataFromIpInvalidIpParamError},
		{Scenario: "Gateway thrown error", TestFn: testHandlerGetDataFromIpGtwError},
		{Scenario: "Gateway not found error", TestFn: testHandlerGetDataFromIpGtwNotFoundError},
	}

	for _, testCase := range tests {
		t.Run(testCase.Scenario, testCase.TestFn)
	}

}

func testHandlerGetTopISPsFromSwitzerlandNoErrors(t *testing.T) {
	type test struct {
		expectedCode int
		expectedBody string
		ipsCount     []IspIpCount
		err          error
		url          string
	}
	ispIpCount := utilGenerateIspIpCount(10)
	ispIpCountByes, _ := json.Marshal(ispIpCount)
	testCase := test{
		expectedCode: http.StatusOK,
		expectedBody: string(ispIpCountByes),
		url:          "/ipdata/top10/Switzerland",
		ipsCount:     ispIpCount,
		err:          nil,
	}

	ctrl := gomock.NewController(t)
	mockGtw := NewMockGateway(ctrl)
	testHandler := NewHandler(mockGtw)

	mockGtw.EXPECT().GetTopISPFromSwitzerland(gomock.Any()).
		Return(testCase.ipsCount, testCase.err)

	req, err := http.NewRequest("GET", testCase.url, nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(testHandler.GetTopISPsFromSwitzerland)

	httpHandler.ServeHTTP(rr, req)

	assert.Equal(t, testCase.expectedCode, rr.Code)
	assert.Equal(t, testCase.expectedBody, rr.Body.String())
}

func testHandlerGetTopISPsFromSwitzerlandDBError(t *testing.T) {
	type test struct {
		expectedCode int
		expectedBody string
		ipsCount     []IspIpCount
		err          error
		url          string
	}
	ispIpCount := utilGenerateIspIpCount(10)
	ispIpCountByes, _ := json.Marshal(ispIpCount)
	testCase := test{
		expectedCode: http.StatusInternalServerError,
		expectedBody: string(ispIpCountByes),
		url:          "/ipdata/top10/Switzerland",
		ipsCount:     ispIpCount,
		err:          errors.New("gtw error"),
	}

	ctrl := gomock.NewController(t)
	mockGtw := NewMockGateway(ctrl)
	testHandler := NewHandler(mockGtw)

	mockGtw.EXPECT().GetTopISPFromSwitzerland(gomock.Any()).
		Return(testCase.ipsCount, testCase.err)

	req, err := http.NewRequest("GET", testCase.url, nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(testHandler.GetTopISPsFromSwitzerland)

	httpHandler.ServeHTTP(rr, req)

	assert.Equal(t, testCase.expectedCode, rr.Code)
	assert.Equal(t, testCase.err.Error()+"\n", rr.Body.String())

}

func testHandlerGetIpCountByCountryNameNoError(t *testing.T) {
	type test struct {
		expectedCode int
		expectedBody string
		ipsCount     int64
		countryName  string
		err          error
		url          string
		muxVars      map[string]string
	}

	testCase := test{
		expectedCode: http.StatusOK,
		countryName:  "Ireland",
		url:          "/ipdata/count/ip/{country_name}",
		ipsCount:     42,
		err:          nil,
		muxVars:      map[string]string{"country_name": "Ireland"},
	}
	response := struct {
		CountryName string `json:"country_name"`
		IpCount     int64  `json:"ip_count"`
	}{testCase.countryName, testCase.ipsCount}
	responseByes, _ := json.Marshal(response)
	testCase.expectedBody = string(responseByes)

	ctrl := gomock.NewController(t)
	mockGtw := NewMockGateway(ctrl)
	testHandler := NewHandler(mockGtw)

	mockGtw.EXPECT().GetIpCountByCountryName(gomock.Any(), testCase.countryName).
		Return(testCase.ipsCount, testCase.err)

	req, err := http.NewRequest("GET", testCase.url, nil)
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, testCase.muxVars)

	rr := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(testHandler.GetIPCountByCountryName)

	httpHandler.ServeHTTP(rr, req)

	assert.Equal(t, testCase.expectedCode, rr.Code)
	assert.Equal(t, testCase.expectedBody, rr.Body.String())
}

func testHandlerGetIpCountByCountryNameNoCountryParamError(t *testing.T) {
	type test struct {
		expectedCode int
		expectedBody string
		ipsCount     int64
		countryName  string
		err          error
		url          string
		muxVars      map[string]string
	}

	testCase := test{
		expectedBody: "param: country_name requested parameter was not found bad request\n",
		expectedCode: http.StatusBadRequest,
		url:          "/ipdata/count/ip/{country_name}",
		muxVars:      map[string]string{},
	}

	ctrl := gomock.NewController(t)
	mockGtw := NewMockGateway(ctrl)
	testHandler := NewHandler(mockGtw)

	req, err := http.NewRequest("GET", testCase.url, nil)
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, testCase.muxVars)

	rr := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(testHandler.GetIPCountByCountryName)

	httpHandler.ServeHTTP(rr, req)

	assert.Equal(t, testCase.expectedCode, rr.Code)
	assert.Equal(t, testCase.expectedBody, rr.Body.String())
}

func testHandlerGetIpCountByCountryNameGtwError(t *testing.T) {
	type test struct {
		expectedCode int
		expectedBody string
		ipsCount     int64
		countryName  string
		err          error
		url          string
		muxVars      map[string]string
	}

	testCase := test{
		expectedBody: "internal server error\n",
		expectedCode: http.StatusInternalServerError,
		countryName:  "Ireland",
		url:          "/ipdata/count/ip/{country_name}",
		ipsCount:     42,
		err:          common.ErrorInternalServer,
		muxVars:      map[string]string{"country_name": "Ireland"},
	}

	ctrl := gomock.NewController(t)
	mockGtw := NewMockGateway(ctrl)
	testHandler := NewHandler(mockGtw)

	mockGtw.EXPECT().GetIpCountByCountryName(gomock.Any(), testCase.countryName).
		Return(testCase.ipsCount, testCase.err)

	req, err := http.NewRequest("GET", testCase.url, nil)
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, testCase.muxVars)

	rr := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(testHandler.GetIPCountByCountryName)

	httpHandler.ServeHTTP(rr, req)

	assert.Equal(t, testCase.expectedCode, rr.Code)
	assert.Equal(t, testCase.expectedBody, rr.Body.String())
}

func testHandlerGetDataFromIpNoError(t *testing.T) {
	type test struct {
		ip           string
		expectedCode int
		expectedBody string
		ipData       IpData
		err          error
		url          string
		muxVars      map[string]string
	}

	upDataBytes, _ := json.Marshal(mockIpDataGateway)
	testCase := test{
		ip:           "127.0.0.1",
		expectedCode: http.StatusOK,
		expectedBody: string(upDataBytes),
		url:          "/ipdata/{ip}",
		ipData:       mockIpDataGateway,
		err:          nil,
	}
	testCase.muxVars = map[string]string{"ip": testCase.ip}
	ctrl := gomock.NewController(t)
	mockGtw := NewMockGateway(ctrl)
	testHandler := NewHandler(mockGtw)

	mockGtw.EXPECT().GetDataFromIP(gomock.Any(), testCase.ip).
		Return(testCase.ipData, testCase.err)

	req, err := http.NewRequest("GET", testCase.url, nil)
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, testCase.muxVars)

	rr := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(testHandler.GetDataFromIP)

	httpHandler.ServeHTTP(rr, req)

	assert.Equal(t, testCase.expectedCode, rr.Code)
	assert.Equal(t, testCase.expectedBody, rr.Body.String())

}

func testHandlerGetDataFromIpNoIpParamError(t *testing.T) {
	type test struct {
		ip           string
		expectedCode int
		expectedBody string
		ipData       IpData
		err          error
		url          string
		muxVars      map[string]string
	}

	testCase := test{
		expectedCode: http.StatusBadRequest,
		expectedBody: "param: ip requested parameter was not found bad request\n",
		url:          "/ipdata/{ip}",
		muxVars:      map[string]string{"ip": ""},
	}
	ctrl := gomock.NewController(t)
	mockGtw := NewMockGateway(ctrl)
	testHandler := NewHandler(mockGtw)

	req, err := http.NewRequest("GET", testCase.url, nil)
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, testCase.muxVars)

	rr := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(testHandler.GetDataFromIP)

	httpHandler.ServeHTTP(rr, req)

	assert.Equal(t, testCase.expectedCode, rr.Code)
	assert.Equal(t, testCase.expectedBody, rr.Body.String())
}

func testHandlerGetDataFromIpInvalidIpParamError(t *testing.T) {
	type test struct {
		ip           string
		expectedCode int
		expectedBody string
		ipData       IpData
		err          error
		url          string
		muxVars      map[string]string
	}

	testCase := test{
		expectedCode: http.StatusBadRequest,
		expectedBody: "ip is not a valid Ipv4 format bad request\n",
		url:          "/ipdata/{ip}",
		muxVars:      map[string]string{"ip": "badIP"},
	}
	ctrl := gomock.NewController(t)
	mockGtw := NewMockGateway(ctrl)
	testHandler := NewHandler(mockGtw)

	req, err := http.NewRequest("GET", testCase.url, nil)
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, testCase.muxVars)

	rr := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(testHandler.GetDataFromIP)

	httpHandler.ServeHTTP(rr, req)

	assert.Equal(t, testCase.expectedCode, rr.Code)
	assert.Equal(t, testCase.expectedBody, rr.Body.String())
}

func testHandlerGetDataFromIpGtwError(t *testing.T) {
	type test struct {
		ip           string
		expectedCode int
		expectedBody string
		ipData       IpData
		err          error
		url          string
		muxVars      map[string]string
	}

	testCase := test{
		ip:           "127.0.0.1",
		expectedCode: http.StatusInternalServerError,
		expectedBody: "internal server error\n",
		url:          "/ipdata/{ip}",
		ipData:       IpData{},
		err:          common.ErrorInternalServer,
	}
	testCase.muxVars = map[string]string{"ip": testCase.ip}
	ctrl := gomock.NewController(t)
	mockGtw := NewMockGateway(ctrl)
	testHandler := NewHandler(mockGtw)

	mockGtw.EXPECT().GetDataFromIP(gomock.Any(), testCase.ip).
		Return(testCase.ipData, testCase.err)

	req, err := http.NewRequest("GET", testCase.url, nil)
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, testCase.muxVars)

	rr := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(testHandler.GetDataFromIP)

	httpHandler.ServeHTTP(rr, req)

	assert.Equal(t, testCase.expectedCode, rr.Code)
	assert.Equal(t, testCase.expectedBody, rr.Body.String())
}

func testHandlerGetDataFromIpGtwNotFoundError(t *testing.T) {
	type test struct {
		ip           string
		expectedCode int
		expectedBody string
		ipData       IpData
		err          error
		url          string
		muxVars      map[string]string
	}

	testCase := test{
		ip:           "127.0.0.1",
		expectedCode: http.StatusNotFound,
		expectedBody: "not found\n",
		url:          "/ipdata/{ip}",
		ipData:       IpData{},
		err:          common.ErrorNotFound,
	}
	testCase.muxVars = map[string]string{"ip": testCase.ip}
	ctrl := gomock.NewController(t)
	mockGtw := NewMockGateway(ctrl)
	testHandler := NewHandler(mockGtw)

	mockGtw.EXPECT().GetDataFromIP(gomock.Any(), testCase.ip).
		Return(testCase.ipData, testCase.err)

	req, err := http.NewRequest("GET", testCase.url, nil)
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, testCase.muxVars)

	rr := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(testHandler.GetDataFromIP)

	httpHandler.ServeHTTP(rr, req)

	assert.Equal(t, testCase.expectedCode, rr.Code)
	assert.Equal(t, testCase.expectedBody, rr.Body.String())
}
