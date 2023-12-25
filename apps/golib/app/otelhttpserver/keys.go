package otelhttpserver

import (
	"go.opentelemetry.io/otel/attribute"
)

var (
	wroteBytesKey = attribute.Key("http.response.wrote_bytes") // if anything was written to the response writer, the total number of bytes written
)
