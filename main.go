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
	// tlsConfig.CipherSuites = []uint16{
	// 	tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
	// 	tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
	// 	tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
	// 	tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
	// 	tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
	// 	tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
	// }

	// defaultTransportDialContext := func(dialer *net.Dialer) func(context.Context, string, string) (net.Conn, error) {
	// 	return dialer.DialContext
	// }

	var wait sync.WaitGroup
	for k := 0; k < 10; k++ {
		wait.Add(1)
		go func() {
			// httpUA := &httpUserAgent{}

			tlsConfig := &tls.Config{
				InsecureSkipVerify:     true,
				SessionTicketsDisabled: false,
			}

			tlsConfig.MinVersion = tls.VersionTLS13

			client := &http.Client{
				Timeout: 10 * time.Second,
				// Transport: httpUA.httpSetUserAgent(&http.Transport{
				Transport: &http.Transport{
					DisableKeepAlives:   false,
					IdleConnTimeout:     300 * time.Second,
					MaxIdleConns:        128,
					MaxIdleConnsPerHost: 128,
					MaxConnsPerHost:     0,
					TLSClientConfig:     tlsConfig,
					DisableCompression:  false,
					ForceAttemptHTTP2:   true,
					// ReadBufferSize:      1,
					// DialContext: defaultTransportDialContext(&net.Dialer{
					// 	Timeout:   1 * time.Second,
					// 	KeepAlive: 300 * time.Second,
					// 	Resolver: &net.Resolver{
					// 		PreferGo: true,
					// 	},
					// }),
				},
			}

			request, err := http.NewRequest("GET", "https://playmytime.com/", nil)
			if err != nil {
				fmt.Println(err)
				return
			}

			request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:99.0) Gecko/20100101 Firefox/99.0")

			for i := 1; i < 100; i++ {
				response, err := client.Do(request)
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

			wait.Done()
		}()
	}

	wait.Wait()
}

func (m *httpUserAgent) RoundTrip(r *http.Request) (*http.Response, error) {
	// r.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:99.0) Gecko/20100101 Firefox/99.0")
	return m.inner.RoundTrip(r)
}

func (m *httpUserAgent) httpSetUserAgent(inner http.RoundTripper) http.RoundTripper {
	m.inner = inner
	return m
}
