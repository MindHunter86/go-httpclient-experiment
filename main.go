package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/profile"
)

type httpUserAgent struct {
	inner http.RoundTripper
}

func main() {
	defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()

	httpUA := &httpUserAgent{}

	tlsConfig := &tls.Config{
		InsecureSkipVerify:     false,
		SessionTicketsDisabled: false,
	}

	tlsConfig.MinVersion = tls.VersionTLS12
	tlsConfig.CipherSuites = []uint16{
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
		tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: httpUA.httpSetUserAgent(&http.Transport{
			DisableKeepAlives:   false,
			IdleConnTimeout:     300 * time.Second,
			MaxIdleConnsPerHost: 128,
			TLSClientConfig:     tlsConfig,
			DisableCompression:  false,
		}),
	}

	request, err := http.NewRequest("GET", "https://playmytime.com/", nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		fmt.Println(response.Status)
		return
	}
}

func (m *httpUserAgent) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:99.0) Gecko/20100101 Firefox/99.0")
	return m.inner.RoundTrip(r)
}

func (m *httpUserAgent) httpSetUserAgent(inner http.RoundTripper) http.RoundTripper {
	m.inner = inner
	return m
}
