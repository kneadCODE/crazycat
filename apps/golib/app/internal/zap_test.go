package internal

import (
	"testing"

	"github.com/kneadCODE/crazycat/apps/golib/app"
	"github.com/stretchr/testify/require"
)

func Test_NewZap(t *testing.T) {
	// Given: && When:
	l, err := NewZap(app.Config{Name: "name", Env: app.EnvDev})

	// Then:
	require.NoError(t, err)
	require.NotNil(t, l)

	// Given: && When:
	l, err = NewZap(app.Config{Name: "name", Env: app.EnvStaging})

	// Then:
	require.NoError(t, err)
	require.NotNil(t, l)

	// Given: && When:
	l, err = NewZap(app.Config{Name: "name", Env: app.EnvProd})

	// Then:
	require.NoError(t, err)
	require.NotNil(t, l)
}

// func TestSomething(t *testing.T) {
// 	require.NoError(t, os.Setenv(string(semconv.ServiceNameKey), "name"))
// 	require.NoError(t, os.Setenv(string(semconv.ServiceNamespaceKey), "project"))
// 	require.NoError(t, os.Setenv(string(semconv.ServiceVersionKey), "v1.0.0"))
// 	require.NoError(t, os.Setenv(string(semconv.ServiceInstanceIDKey), "instance"))
// 	require.NoError(t, os.Setenv(string(semconv.DeploymentEnvironmentKey), "development"))
// 	defer os.Unsetenv(string(semconv.ServiceNameKey))
// 	defer os.Unsetenv(string(semconv.ServiceNamespaceKey))
// 	defer os.Unsetenv(string(semconv.ServiceVersionKey))
// 	defer os.Unsetenv(string(semconv.ServiceInstanceIDKey))
// 	defer os.Unsetenv(string(semconv.DeploymentEnvironmentKey))
//
// 	ctx, finish, err := Init()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer finish()
//
// 	tracer := otelTracerFromContext(ctx)
// 	ctx, _ = tracer.Start(ctx, "some_span")
//
// 	TrackDebugEvent(ctx, "hello world")
// 	TrackInfoEvent(ctx, "hello world")
// 	TrackWarnEvent(ctx, "hello world")
// 	TrackErrorEvent(ctx, errors.New("some err"))
//
// 	t.Fail()
// }
