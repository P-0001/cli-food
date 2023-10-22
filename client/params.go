package client

import (
	"io"

	http "github.com/bogdanfinn/fhttp"
)

type TLSHeaders = http.Header

type TLSParams struct {
	Client           HttpClient
	Method           string
	Url              string
	Headers          TLSHeaders
	Body             io.Reader
	ExpectedResponse int
}
