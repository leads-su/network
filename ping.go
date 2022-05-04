package network

import (
	"fmt"
	"net/http"
	"net/http/httptrace"
	"time"
)

// Ping represents structure of ping configuration
type Ping struct {
	scheme			string
	host			string
	port			uint
	withDNSResolve	bool
}

// PingOptions represents structure of ping options
type PingOptions struct {
	Scheme				string
	Host				string
	Port				uint
	WithDNSResolveTime	bool
}

// ResponseTime returns response time
func ResponseTime(options *PingOptions) (int64, error) {
	ping := &Ping{
		scheme:			options.Scheme,
		host:			options.Host,
		port:			options.Port,
		withDNSResolve:	options.WithDNSResolveTime,
	}

	if ping.host == "" {
		return -1, fmt.Errorf("please provide host")
	}

	if ping.port == 0 {
		if ping.scheme == "http" {
			ping.port = 80
		} else if ping.scheme == "https" {
			ping.port = 443
		} else {
			return -1, fmt.Errorf("unable to guess port from scheme as it is empty")
		}
	}

	if ping.scheme == "" {
		if ping.port == 80 {
			ping.scheme = "http"
		} else if ping.port == 443 {
			ping.scheme = "https"
		} else {
			return -1, fmt.Errorf("unable to guess scheme from port")
		}
	}

	return ping.roundTrip(), nil
}

// fullPath returns full path to host
func (ping *Ping) fullPath() string {
	return fmt.Sprintf("%s://%s:%d", ping.scheme, ping.host, ping.port)
}

// roundTrip calculates round trip between local and remote hosts
func (ping *Ping) roundTrip() int64 {
	request, _ := http.NewRequest("GET", ping.fullPath(), nil)
	var connectStart, dnsStart time.Time
	var connectEnd, dnsEnd time.Duration

	trace := &httptrace.ClientTrace{
		DNSStart: func (dsi httptrace.DNSStartInfo) {
			dnsStart = time.Now()
		},
		DNSDone: func (ddi httptrace.DNSDoneInfo) {
			dnsEnd = time.Since(dnsStart)
		},
		ConnectStart: func (network, addr string) {
			connectStart = time.Now()
		},
		ConnectDone: func (network, addr string, err error) {
			connectEnd = time.Since(connectStart)
		},
	}

	request = request.WithContext(httptrace.WithClientTrace(request.Context(), trace))
	if _, err := http.DefaultTransport.RoundTrip(request); err != nil {
		return -1
	}

	if ping.withDNSResolve {
		return connectEnd.Milliseconds() + dnsEnd.Milliseconds()
	}

	return connectEnd.Milliseconds()
}