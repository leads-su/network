package network

import (
	"errors"
	"net"
	"strconv"
	"time"
)

// Available represents structure of available configuration
type Available struct {
	host		string
	port		uint
	timeout		int
}

// AvailableOptions represents structure of available options
type AvailableOptions struct {
	Host		string
	Port		uint
	Timeout		int
}

// AvailablePingOptions represents structure of available options with ping
type AvailablePingOptions struct {
	Scheme				string
	Host				string
	Port				uint
	WithDNSResolveTime	bool
	Timeout				int
}

// IsServiceAvailable checks if service specified in options is available
func IsServiceAvailable(options *AvailableOptions) (bool, error) {
	available, err := initializeAvailableStruct(options)
	if err != nil {
		return false, err
	}
	return available.isAvailable(), nil
}

// IsServiceAvailableWithPing check if service specified in options is available and also retrieve ping
func IsServiceAvailableWithPing(options *AvailablePingOptions) (bool, int64, error) {
	available, err := initializeAvailableStruct(&AvailableOptions{
		Host:    	options.Host,
		Port:    	options.Port,
		Timeout:	options.Timeout,
	})
	if err != nil {
		return false, -1, err
	}
	if !available.isAvailable() {
		return false, -1, nil
	}
	responseTime, err := ResponseTime(&PingOptions{
		Scheme:				options.Scheme,
		Host:				options.Host,
		Port:				options.Port,
		WithDNSResolveTime:	options.WithDNSResolveTime,
	})
	return true, responseTime, nil
}

// initializeAvailableStruct initializes Available struct for later usage
func initializeAvailableStruct(options *AvailableOptions) (*Available, error){
	available := &Available{
		host:    options.Host,
		port:    options.Port,
		timeout: options.Timeout,
	}

	if available.host == "" {
		return nil, errors.New("please specify target host")
	}

	if available.port == 0 {
		return nil, errors.New("please specify target port")
	}

	if available.timeout == 0 {
		available.timeout = 5
	}

	return available, nil
}

// isAvailable check if specified target is available
func (available *Available) isAvailable() bool {
	timeout := time.Duration(available.timeout) * time.Second
	connector, err := net.DialTimeout("tcp", net.JoinHostPort(available.host, strconv.Itoa(int(available.port))), timeout)
	if err != nil {
		return false
	}
	if connector != nil {
		defer connector.Close()
		return true
	}
	return false
}