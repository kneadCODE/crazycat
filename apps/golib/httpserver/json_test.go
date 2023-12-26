package httpserver

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadJSON(t *testing.T) {
	type result struct {
		Key string `json:"key"`
	}

	// Given: no body
	r := httptest.NewRequest(http.MethodPost, "/abc", nil)

	// When:
	var res result
	err := ReadJSON(r, &res)

	// Then:
	require.Equal(t, &Error{Status: http.StatusBadRequest, Code: "json_parse_failed", Desc: io.EOF.Error()}, err)

	// Given: non-json body
	r = httptest.NewRequest(http.MethodPost, "/abc", bytes.NewReader([]byte("abc")))

	// When:
	err = ReadJSON(r, &res)

	// Then:
	require.Equal(t, &Error{Status: http.StatusBadRequest, Code: "json_parse_failed", Desc: "invalid character 'a' looking for beginning of value"}, err)

	// Given: invalid json
	r = httptest.NewRequest(http.MethodPost, "/abc", bytes.NewReader([]byte(`{"key":"value"`)))

	// When:
	err = ReadJSON(r, &res)

	// Then:
	require.Equal(t, &Error{Status: http.StatusBadRequest, Code: "json_parse_failed", Desc: io.ErrUnexpectedEOF.Error()}, err)

	// Given: valid json
	r = httptest.NewRequest(http.MethodPost, "/abc", bytes.NewReader([]byte(`{"key":"value"}`)))

	// When:
	err = ReadJSON(r, &res)

	// Then:
	require.NoError(t, err)
	require.Equal(t, "value", res.Key)
}

func TestWriteJSON(t *testing.T) {
	type result struct {
		Key string `json:"key"`
	}

	// Given: write json
	w := httptest.NewRecorder()

	// When:
	WriteJSON(context.Background(), w, result{Key: "val"}, nil)

	// Then:
	require.Equal(t, http.StatusOK, w.Result().StatusCode)
	v, err := io.ReadAll(w.Result().Body)
	require.NoError(t, err)
	require.Equal(t, "application/json", w.Header().Get("Content-Type"))
	require.Equal(t, "13", w.Header().Get("Content-Length"))
	require.Equal(t, `{"key":"val"}`, string(v))

	// Given: write json with header
	w = httptest.NewRecorder()

	// When:
	WriteJSON(context.Background(), w, result{Key: "val"}, map[string]string{"h1": "v1"})

	// Then:
	require.Equal(t, http.StatusOK, w.Result().StatusCode)
	v, err = io.ReadAll(w.Result().Body)
	require.NoError(t, err)
	require.Equal(t, "application/json", w.Header().Get("Content-Type"))
	require.Equal(t, "13", w.Header().Get("Content-Length"))
	require.Equal(t, `{"key":"val"}`, string(v))
	require.Equal(t, "v1", w.Result().Header.Get("h1"))

	// Given: write err
	w = httptest.NewRecorder()

	// When:
	WriteJSON(context.Background(), w, errors.New("err"), nil)

	// Then:
	require.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	v, err = io.ReadAll(w.Result().Body)
	require.NoError(t, err)
	require.Equal(t, "application/json", w.Header().Get("Content-Type"))
	require.Equal(t, strconv.Itoa(len(`{"code":"internal_server_error","description":"Internal Server Error"}`)), w.Header().Get("Content-Length"))
	require.Equal(t, `{"code":"internal_server_error","description":"Internal Server Error"}`, string(v))

	// Given: write err
	w = httptest.NewRecorder()

	// When:
	WriteJSON(context.Background(), w, &Error{Status: http.StatusBadRequest, Code: "code", Desc: "desc"}, nil)

	// Then:
	require.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	v, err = io.ReadAll(w.Result().Body)
	require.NoError(t, err)
	require.Equal(t, "application/json", w.Header().Get("Content-Type"))
	require.Equal(t, strconv.Itoa(len(`{"code":"code","description":"desc"}`)), w.Header().Get("Content-Length"))
	require.Equal(t, `{"code":"code","description":"desc"}`, string(v))
}
