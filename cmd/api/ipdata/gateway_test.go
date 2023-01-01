package ipdata

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetIspIpsByCountryCode(t *testing.T) {
	type test struct {
		name   string
		code   string
		limit  int
		output []IspIpCount
		err    error
	}

	tests := []test{
		{name: "test no error", code: "AR", limit: 1, output: []IspIpCount{{Isp: "mainIsp", IpCount: 1234}}, err: nil},
		{name: "test db error", code: "AR", limit: 1, output: []IspIpCount{}, err: errors.New("db error")},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockDao := NewMockDao(ctrl)
			gtw := NewGateway(mockDao)

			mockDao.EXPECT().
				GetTopIspByCountryCode(gomock.Any(), testCase.code, testCase.limit).
				Return(testCase.output, testCase.err)

			output, err := gtw.GetIspIpsByCountryCode(context.Background(), testCase.code, testCase.limit)

			assert.Equal(t, testCase.output, output)
			assert.Equal(t, testCase.err, err)
		})
	}
}
