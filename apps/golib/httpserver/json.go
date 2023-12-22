package httpserver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/kneadCODE/crazycat/apps/golib/app"
)

// ReadJSON reads the http.Request body and attempts to parse it into the desired type v. If read fails or parsing
// fails, it returns an error. v here should be a pointer to the actual type.
func ReadJSON(r *http.Request, v any) error {
	// We don't need to close the req body after reading because:
	// The Server will close the request body. The ServeHTTP Handler does not need to.
	// Ref - https://pkg.go.dev/net/http#Request

	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return &Error{Status: http.StatusBadRequest, Code: "json_parse_failed", Desc: err.Error()}
	}

	return nil
}

// WriteJSON parses the given v to JSON equivalent and writes it to the http.ResponseWriter. It also sets the relevant
// headers such as status, Content-Type and Content-Length.
func WriteJSON(ctx context.Context, w http.ResponseWriter, v interface{}, headers map[string]string) {
	for k, v := range headers {
		w.Header().Add(k, v)
	}

	switch parsed := v.(type) {
	case *Error:
		w.WriteHeader(parsed.Status)
	case error:
		w.WriteHeader(http.StatusInternalServerError)
		v = ErrInternalServer // We don't want to return internal err details, so we transform it.
	}

	vBytes, err := json.Marshal(v) // Need to do this way instead of Encode because we need to log it
	if err != nil {
		app.TrackErrorEvent(ctx, fmt.Errorf("httpserver:WriteJSON: %w", err)) // TODO: Add any additional fields if needed.
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(vBytes))) // TODO: Check if this causes problems or not.

	// TODO: Log any additional fields if needed.
	app.TrackInfoEvent(ctx, `Writing JSON Body: 
		BODY: [%s],
		HEADERS: [%v]
		`, vBytes, w.Header(),
	)

	if _, err = w.Write(vBytes); err != nil {
		app.TrackErrorEvent(ctx, fmt.Errorf("httpserver:WriteJSON: %w", err)) // TODO: Add any additional fields if needed.
	}
}
