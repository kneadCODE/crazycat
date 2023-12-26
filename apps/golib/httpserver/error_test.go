package httpserver

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestError_Error(t *testing.T) {
	require.Equal(t,
		"httpserver:Error: Status:[200],Code:[code],Desc:[desc]",
		Error{Status: http.StatusOK, Code: "code", Desc: "desc"}.Error(),
	)
	require.Equal(t,
		"httpserver:Error: Status:[400],Code:[c],Desc:[d]",
		Error{Status: http.StatusBadRequest, Code: "c", Desc: "d"}.Error(),
	)
	require.Equal(t,
		"httpserver:Error: Status:[0],Code:[],Desc:[]",
		Error{}.Error(),
	)
}
