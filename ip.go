package cfdyndns

import (
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

// InvalidIPError is returned when an value returned is invalid.
// This error should be returned by the source itself.
type InvalidIPError string

// Error implements error.Error
func (err InvalidIPError) Error() string {
	return "Invalid IP: " + string(err)
}

func GetOutboundIP(url string, timeout time.Duration) (net.IP, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)

	client := &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			MaxIdleConns:      1,
			IdleConnTimeout:   3 * time.Second,
			DisableKeepAlives: true,
		},
	}

	// Do the request and read the body for non-error results.
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[ERROR] could not GET %q: %v\n", url, err)
		return nil, err
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[ERROR] could not read response from %q: %v\n", url, err)
		return nil, err
	}

	raw := string(bytes)

	// validate the IP
	externalIP := net.ParseIP(strings.TrimSpace(raw))
	if externalIP == nil {
		log.Printf("[ERROR] %q returned an invalid IP: %v\n", url, err)
		return nil, InvalidIPError(raw)
	}

	// returned the parsed IP
	return externalIP, nil
}
