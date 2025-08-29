package adapter

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/http/httptrace"
	"strings"
	"time"
)

type DefaultHttpClient struct {
	client *http.Client
}

func NewHttpClient() *DefaultHttpClient {
	return &DefaultHttpClient{client: &http.Client{}}
}

func (c *DefaultHttpClient) DoRequest(method, url string, body any, headers map[string]string) (int, []byte, error) {
	var (
		start, dnsStart, connStart, tlsStart, gotConn, wroteRequest, firstByte        time.Time
		dnsDur, tcpDur, tlsDur, connDur, ttfbDur, downloadDur, prepareDur, processDur time.Duration
		reused, wasIdle                                                               bool
		idleTime                                                                      time.Duration
	)

	prepareStart := time.Now()
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return 0, nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	prepareDur = time.Since(prepareStart)

	trace := &httptrace.ClientTrace{
		GetConn:           func(_ string) { start = time.Now() },
		DNSStart:          func(_ httptrace.DNSStartInfo) { dnsStart = time.Now() },
		DNSDone:           func(_ httptrace.DNSDoneInfo) { dnsDur = time.Since(dnsStart) },
		ConnectStart:      func(_, _ string) { connStart = time.Now() },
		ConnectDone:       func(_, _ string, _ error) { tcpDur = time.Since(connStart) },
		TLSHandshakeStart: func() { tlsStart = time.Now() },
		TLSHandshakeDone: func(_ tls.ConnectionState, _ error) {
			tlsDur = time.Since(tlsStart)
		},
		GotConn: func(info httptrace.GotConnInfo) {
			gotConn = time.Now()
			reused, wasIdle, idleTime = info.Reused, info.WasIdle, info.IdleTime
		},
		WroteRequest:         func(_ httptrace.WroteRequestInfo) { wroteRequest = time.Now() },
		GotFirstResponseByte: func() { firstByte = time.Now(); ttfbDur = firstByte.Sub(wroteRequest); connDur = gotConn.Sub(start) },
	}

	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))

	resp, err := c.client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()

	downloadStart := time.Now()
	data, err := io.ReadAll(resp.Body)
	downloadDur = time.Since(downloadStart)

	processStart := time.Now()
	processDur = time.Since(processStart)

	totalDur := time.Since(start)

	fmt.Printf("\n--- HTTP Diagnostics ---\n")
	fmt.Printf("Status: %d, Proto: %s, Content-Length: %d\n", resp.StatusCode, resp.Proto, len(data))
	fmt.Printf("Response Time: %.2f ms\n", ms(totalDur))
	fmt.Printf("Prepare: %.2f ms\n", ms(prepareDur))
	fmt.Printf("DNS Lookup: %.2f ms\n", ms(dnsDur))
	fmt.Printf("TCP Handshake: %.2f ms\n", ms(tcpDur))
	fmt.Printf("TLS Handshake: %.2f ms\n", ms(tlsDur))
	fmt.Printf("Connection: %.2f ms (Reused=%v, WasIdle=%v, Idle=%.2f ms)\n",
		ms(connDur), reused, wasIdle, ms(idleTime))
	fmt.Printf("TTFB: %.2f ms\n", ms(ttfbDur))
	fmt.Printf("Download: %.2f ms\n", ms(downloadDur))
	fmt.Printf("Process: %.2f ms\n", ms(processDur))

	if resp.TLS != nil {
		fmt.Printf("TLS Version: %s, Cipher: %s\n",
			tlsVersionToStr(resp.TLS.Version),
			tls.CipherSuiteName(resp.TLS.CipherSuite))
		if len(resp.TLS.PeerCertificates) > 0 {
			cert := resp.TLS.PeerCertificates[0]
			fmt.Printf("TLS Cert CN: %s, Issuer: %s, ValidUntil: %s\n",
				cert.Subject.CommonName, cert.Issuer.CommonName, cert.NotAfter)
		}
	}

	fmt.Println("\n--- Headers Analysis ---")
	seenHeaders := map[string]bool{}
	for k, v := range resp.Header {
		key := strings.ToLower(k)
		if seenHeaders[key] {
			fmt.Printf("[!] Duplicate header: %s\n", k)
		}
		seenHeaders[key] = true

		// Подозрительные
		if key == "server" || key == "x-powered-by" {
			fmt.Printf("[!] Suspicious header (%s): %v\n", k, v)
		}
		if key == "content-type" && len(v) == 0 {
			fmt.Printf("[!] Missing Content-Type!\n")
		}
	}

	// 🔎 Проверка cookies
	fmt.Println("\n--- Cookies Analysis ---")
	seenCookies := map[string]bool{}
	for _, c := range resp.Cookies() {
		if seenCookies[c.Name] {
			fmt.Printf("[!] Duplicate cookie: %s\n", c.Name)
		}
		seenCookies[c.Name] = true

		// Подозрительные
		if !c.HttpOnly {
			fmt.Printf("[!] Cookie %s is not HttpOnly\n", c.Name)
		}
		if !c.Secure {
			fmt.Printf("[!] Cookie %s is not Secure\n", c.Name)
		}
		if c.MaxAge > 86400*30 {
			fmt.Printf("[!] Cookie %s has very long Max-Age: %d seconds\n", c.Name, c.MaxAge)
		}
	}

	return resp.StatusCode, data, err
}

// Утилиты
func ms(d time.Duration) float64 {
	return float64(d.Microseconds()) / 1000
}

func tlsVersionToStr(v uint16) string {
	switch v {
	case tls.VersionTLS10:
		return "TLS1.0"
	case tls.VersionTLS11:
		return "TLS1.1"
	case tls.VersionTLS12:
		return "TLS1.2"
	case tls.VersionTLS13:
		return "TLS1.3"
	default:
		return "Unknown"
	}
}
