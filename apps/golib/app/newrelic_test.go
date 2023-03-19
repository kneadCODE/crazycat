package app

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_newNewRelic(t *testing.T) {
	type testCase struct {
		givenLicense string
		givenCfg     Config
		expErr       error
		expApp       bool
	}
	tcs := map[string]testCase{
		"no license": {},
		"invalid license": {
			givenLicense: "abcd",
			expErr:       errors.New("init new relic failed: license length is not 40"),
		},
		"err: cfg without name": {
			givenLicense: "1234567890123456789012345678901234567890",
			expErr:       errors.New("init new relic failed: string AppName required"),
		},
		"enabled": {
			givenLicense: "1234567890123456789012345678901234567890",
			givenCfg:     Config{Name: "name"},
			expApp:       true,
		},
	}
	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			// Given:
			defer os.Unsetenv("NEW_RELIC_LICENSE_KEY")
			require.NoError(t, os.Setenv("NEW_RELIC_LICENSE_KEY", tc.givenLicense))

			// When:
			nrApp, err := newNewRelic(tc.givenCfg)

			// Then:
			if tc.expErr != nil {
				require.EqualError(t, err, tc.expErr.Error())
				require.Nil(t, nrApp)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expApp, nrApp != nil)
			}
		})
	}
}
