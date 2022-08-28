package clients

import (
	"sync"
	"time"

	"github.com/valyala/fasthttp"
)

var (
	onceFastHTTP sync.Once
	fastHTTPConn *fasthttp.Client
)

func initializeFHHTP() {
	readTimeout, _ := time.ParseDuration("5s")
	writeTimeout, _ := time.ParseDuration("s")
	maxIdleConnDuration, _ := time.ParseDuration("1h")

	client := &fasthttp.Client{
		ReadTimeout:                   readTimeout,
		WriteTimeout:                  writeTimeout,
		MaxIdleConnDuration:           maxIdleConnDuration,
		NoDefaultUserAgentHeader:      true, // Don't send: User-Agent: fasthttp
		DisableHeaderNamesNormalizing: true, // If you set the case on your headers correctly you can enable this
		DisablePathNormalizing:        true,
		Dial: (&fasthttp.TCPDialer{
			Concurrency:      4096,
			DNSCacheDuration: time.Hour,
		}).Dial,
	}

	fastHTTPConn = client
}

func FHTTP() *fasthttp.Client {
	onceFastHTTP.Do(initializeFHHTP)
	return fastHTTPConn
}