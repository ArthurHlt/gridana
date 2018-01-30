package utils

import (
	"crypto/tls"
	"github.com/ArthurHlt/gridana/model"
	"net"
	"net/http"
	"time"
)

type BasicAuthRoundTripper struct {
	user     string
	password string
	wrap     http.RoundTripper
}

func (t BasicAuthRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.SetBasicAuth(t.user, t.password)
	return t.wrap.RoundTrip(req)
}

func CreateClient(config model.DriverConfig) *http.Client {
	var transport http.RoundTripper
	transport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: config.InsecureSkipVerify,
		},
	}
	if config.User != "" {
		transport = &BasicAuthRoundTripper{
			user:     config.User,
			password: config.Password,
			wrap:     transport,
		}
	}
	return &http.Client{
		Transport: transport,
	}
}
