package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"sync"
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
		InsecureSkipVerify:     true,
		SessionTicketsDisabled: false,
	}

	tlsConfig.MinVersion = tls.VersionTLS13
	// tlsConfig.CipherSuites = []uint16{
	// 	tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
	// 	tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
	// 	tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
	// 	tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
	// 	tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
	// 	tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
	// }

	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: httpUA.httpSetUserAgent(&http.Transport{
			DisableKeepAlives:   false,
			IdleConnTimeout:     300 * time.Second,
			MaxIdleConns:        128,
			MaxIdleConnsPerHost: 128,
			TLSClientConfig:     tlsConfig,
			DisableCompression:  false,
			ForceAttemptHTTP2:   true,
		}),
	}

	request, err := http.NewRequest("GET", "http://playmytime.com/", nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	responsePool := sync.Pool{
		New: func() interface{} {
			return new(http.Response)
		},
	}

	for i := 1; i < 100; i++ {
		response := responsePool.Get().(*http.Response)

		response, err = client.Do(request)
		if err != nil {
			fmt.Println(err)
			continue
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			fmt.Println(response.Status)
			continue
		}
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
