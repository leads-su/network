# Network Package for Go Lang
This package provides necessary information for network interfaces in the system.  
It also allows to retrieve RTT information for external hosts and check their availability.

## Availability and RTT

### Retrieve RTT
```go
roundTrip, err := network.ResponseTime(&network.PingOptions{
	Host: "google.com",
	Port: 443,
})

if err != nil {
	logger.Fatalf("main", "failed to calculate round trip - %s", err.Error())
}

fmt.Printf("%dms\n", roundTrip)
```

### Check Availability
```go
available, err := network.IsServiceAvailable(&network.AvailableOptions{
	Host:               "google.com",
	Port:               443,
})

if err != nil {
	logger.Fatalf("main", "%s", err.Error())
}

if available {
    fmt.Println("service is available")
} else {
	fmt.Println("service is unavailable")
}
```

### Check Availability And Return RTT
```go
available, ping, err := network.IsServiceAvailableWithPing(&network.AvailablePingOptions{
	Host:               "google.com",
	Port:               443,
})

if err != nil {
	logger.Fatalf("main", "%s", err.Error())
}

if available {
	fmt.Printf("service is available and it's ping is %dms", ping)
} else {
	fmt.Println("service is unavailable")
}
```

## Network Interfaces

### External IP Details
```go
address, err := network.ExternalIPDetails()
if err != nil {
	logger.Fatalf("main", "%s", err.Error())
}

fmt.Printf("%v\n", address)
```

### External IP
```go
address, err := network.ExternalIPAddress()
if err != nil {
	logger.Fatalf("main", "%s", err.Error())
}

fmt.Println(address)
```