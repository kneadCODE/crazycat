package internal

import (
	"errors"
	"os"
	"testing"

	"github.com/kneadCODE/crazycat/apps/golib/app/config"
	"github.com/stretchr/testify/require"
)

func Test_NewSentryHub(t *testing.T) {
	type testCase struct {
		givenDSNKey string
		expErr      error
		expHub      bool
	}
	tcs := map[string]testCase{
		"no dsn": {},
		"invalid dsn": {
			givenDSNKey: "abcd",
			expErr:      errors.New("init sentry failed: [Sentry] DsnParseError: invalid scheme"),
		},
		"enabled": {
			givenDSNKey: "https://something@somethingelse.ingest.sentry.io/1",
			expHub:      true,
		},
	}
	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			// Given:
			defer os.Unsetenv("SENTRY_DSN")
			require.NoError(t, os.Setenv("SENTRY_DSN", tc.givenDSNKey))

			// When:
			hub, err := NewSentryHub(config.Config{})

			// Then:
			if tc.expErr != nil {
				require.EqualError(t, err, tc.expErr.Error())
				require.Nil(t, hub)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expHub, hub != nil)
			}
		})
	}
}
