package otelhttpserver

import (
	"go.opentelemetry.io/otel/attribute"
)

var (
	readBytesKey  = attribute.Key("http.request.read_bytes")   // if anything was read from the request body, the total number of bytes read
	readErrorKey  = attribute.Key("http.request.read_error")   // If an error occurred while reading a request, the string of the error (io.EOF is not recorded)
	wroteBytesKey = attribute.Key("http.response.wrote_bytes") // if anything was written to the response writer, the total number of bytes written
	writeErrorKey = attribute.Key("http.response.write_error") // if an error occurred while writing a reply, the string of the error (io.EOF is not recorded)
)
