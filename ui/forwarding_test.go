package ui

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_refreshRemoteHost(t *testing.T) {

	type tc struct {
		remoteHost string
		expected   []string
	}

	testCases := []tc{
		{
			remoteHost: "this_is_the_testing-127.0.0.4:9000",
			expected: []string{
				"this_is_the_testing-127.0.0.4:9000",
				"his_is_the_testing-127.0.0.4:9000",
				"is_is_the_testing-127.0.0.4:9000",
				"s_is_the_testing-127.0.0.4:9000",
				"_is_the_testing-127.0.0.4:9000",
				"is_the_testing-127.0.0.4:9000",
				"s_the_testing-127.0.0.4:9000",
				"_the_testing-127.0.0.4:9000",
				"the_testing-127.0.0.4:9000",
				"he_testing-127.0.0.4:9000  thi",
				"e_testing-127.0.0.4:9000  this",
				"_testing-127.0.0.4:9000  this_",
				"testing-127.0.0.4:9000  this_i",
				"esting-127.0.0.4:9000  this_is",
				"sting-127.0.0.4:9000  this_is_",
				"ting-127.0.0.4:9000  this_is_t",
				"ing-127.0.0.4:9000  this_is_th",
				"ng-127.0.0.4:9000  this_is_the",
				"g-127.0.0.4:9000  this_is_the_",
				"-127.0.0.4:9000  this_is_the_t",
				"127.0.0.4:9000  this_is_the_te",
				"27.0.0.4:9000  this_is_the_tes",
				"7.0.0.4:9000  this_is_the_test",
				".0.0.4:9000  this_is_the_testi",
				"0.0.4:9000  this_is_the_testin",
				".0.4:9000  this_is_the_testing",
				"0.4:9000  this_is_the_testing-",
				".4:9000  this_is_the_testing-1",
				"4:9000  this_is_the_testing-12",
				":9000  this_is_the_testing-127",
				"9000  this_is_the_testing-127.",
				"000  this_is_the_testing-127.0",
				"00  this_is_the_testing-127.0.",
				"0  this_is_the_testing-127.0.0",
				"this_is_the_testing-127.0.0.4:9000",
				"his_is_the_testing-127.0.0.4:9000",
				"is_is_the_testing-127.0.0.4:9000",
				"s_is_the_testing-127.0.0.4:9000",
				"_is_the_testing-127.0.0.4:9000",
				"is_the_testing-127.0.0.4:9000",
				"s_the_testing-127.0.0.4:9000",
				"_the_testing-127.0.0.4:9000",
				"the_testing-127.0.0.4:9000",
				"he_testing-127.0.0.4:9000  thi",
				"e_testing-127.0.0.4:9000  this",
				"_testing-127.0.0.4:9000  this_",
				"testing-127.0.0.4:9000  this_i",
				"esting-127.0.0.4:9000  this_is",
				"sting-127.0.0.4:9000  this_is_",
				"ting-127.0.0.4:9000  this_is_t",
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.remoteHost, func(t *testing.T) {
			for i, s := range testCase.expected {
				t.Run(fmt.Sprintf(`#%d_%s`, i, s), func(t *testing.T) {
					assert.Equal(t, s, refreshRemoteHost(testCase.remoteHost))
				})
			}
		})
	}
}
