package app

import (
	"context"
	"errors"
	"os"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	defer resetStubs()

	var sigTerm os.Signal = syscall.SIGTERM

	type testCase struct {
		s1Err   bool
		s2Err   bool
		exitSig *os.Signal
	}
	tcs := map[string]testCase{
		"all ok, sigkill sent": {
			exitSig: &os.Kill,
		},
		"all ok, siginterrupt sent": {
			exitSig: &os.Interrupt,
		},
		"all ok, sigterm sent": {
			exitSig: &sigTerm,
		},
		"s1 error, should terminate s2 as well": {
			s1Err: true,
		},
		"s2 error, should terminate s1 as well": {
			s2Err: true,
		},
	}
	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Given:
			exitChan := make(chan os.Signal, 1)
			exitSignalStub = func() <-chan os.Signal {
				return exitChan
			}

			var s1Executed bool
			s1 := func(ctx context.Context) error {
				s1Executed = true
				if tc.s1Err {
					time.Sleep(2 * time.Second)
					return errors.New("s1 error")
				}

				for {
					select {
					case <-ctx.Done():
						return nil
					}
				}
			}
			var s2Executed bool
			s2 := func(ctx context.Context) error {
				s2Executed = true
				if tc.s2Err {
					time.Sleep(2 * time.Second)
					return errors.New("s2 error")
				}

				for {
					select {
					case <-ctx.Done():
						return nil
					}
				}
			}

			// When:
			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				defer wg.Done()
				Run(context.Background(), s1, s2)
			}()

			time.Sleep(500 * time.Millisecond)

			// Then:
			if tc.exitSig != nil {
				exitChan <- *tc.exitSig
			}

			wg.Wait()

			require.True(t, s1Executed)
			require.True(t, s2Executed)
		})
	}
}
