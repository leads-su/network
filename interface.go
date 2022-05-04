package network

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
)

// Interface represents structure of interface
type Interface struct {
	Index	int
	Address	string
	Name    string
	IPv4	string
	Subnet4	string
	IPv6	string
	Subnet6	string
}

// Interfaces represents structure of interfaces array
type Interfaces = []Interface

// SystemInterfaces returns list of system network interfaces
func SystemInterfaces() (Interfaces, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	var systemInterfaces Interfaces

	for _, iface := range interfaces {
		if iface.HardwareAddr == nil {
			continue
		}
		interfaceAddresses, err := iface.Addrs()
		if err != nil {
			return nil, err
		}

		version4 := ""
		version6 := ""

		for _, address := range interfaceAddresses {
			addressAsString := address.String()
			if strings.Count(addressAsString, ".") == 3 {
				version4 = addressAsString
			} else {
				version6 = addressAsString
			}
		}

		systemInterfaces = append(systemInterfaces, Interface{
			Index:		iface.Index,
			Address:	iface.HardwareAddr.String(),
			Name:		iface.Name,
			IPv4:		version4[:strings.Index(version4, "/")],
			Subnet4:	version4[strings.Index(version4, "/"):],
			IPv6:		version6[:strings.Index(version6, "/")],
			Subnet6:	version6[strings.Index(version6, "/"):],
		})
	}

	return systemInterfaces, nil
}

// InterfaceDetailsByName returns interface details while searching by Name
func InterfaceDetailsByName(name string) (*Interface, error) {
	interfaces, err := SystemInterfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range interfaces {
		if iface.Name == name {
			return &iface, nil
		}
	}
	return nil, fmt.Errorf("unable to find interface with given name - %s", name)
}

// InterfaceDetailsByMac returns interface details while searching by Mac address
func InterfaceDetailsByMac(address string) (*Interface, error) {
	interfaces, err := SystemInterfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range interfaces {
		if iface.Address == address {
			return &iface, nil
		}
	}
	return nil, fmt.Errorf("unable to find interface with given mac address - %s", address)
}

// GetIPv4ByName returns IPv4 address for interface by its Name
func GetIPv4ByName(name string) (string, error) {
	iface, err := InterfaceDetailsByName(name)
	if err != nil {
		return "", err
	}
	return iface.IPv4, nil
}

// GetIPv4ByMac returns IPv4 address for interface by its Mac Address
func GetIPv4ByMac(address string) (string, error) {
	iface, err := InterfaceDetailsByMac(address)
	if err != nil {
		return "", err
	}
	return iface.IPv4, nil
}

// GetIPv6ByName returns IPv6 address for interface by its Name
func GetIPv6ByName(name string) (string, error) {
	iface, err := InterfaceDetailsByName(name)
	if err != nil {
		return "", err
	}
	return iface.IPv6, nil
}

// GetIPv6ByMac returns IPv6 address for interface by its Mac Address
func GetIPv6ByMac(address string) (string, error) {
	iface, err := InterfaceDetailsByMac(address)
	if err != nil {
		return "", err
	}
	return iface.IPv6, nil
}

// ExternalIP represents response structure from IP resolver
type ExternalIP struct {
	Status          string  `json:"status"`
	Country         string  `json:"country"`
	CountryCode     string  `json:"countryCode"`
	Region          string  `json:"region"`
	RegionName      string  `json:"regionName"`
	City            string  `json:"city"`
	PostalCode      string  `json:"zip"`
	Latitude        float64 `json:"lat"`
	Longitude       float64 `json:"lon"`
	Timezone        string  `json:"timezone"`
	Provider        string  `json:"isp"`
	Organisation    string  `json:"org"`
	AS              string  `json:"as"`
	Address         string  `json:"query"`
}

// ExternalIPDetails get information for external IP
func ExternalIPDetails() (*ExternalIP, error) {
	request, err := http.Get("http://ip-api.com/json/")
	if err != nil {
		return nil, err
	}
	defer request.Body.Close()

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return nil, err
	}

	var externalIP *ExternalIP
	err = json.Unmarshal(body, &externalIP)
	if err != nil {
		return nil, err
	}
	return externalIP, nil
}

// ExternalIPAddress returns external IP address
func ExternalIPAddress() (string, error) {
	details, err := ExternalIPDetails()
	if err != nil {
		return "", err
	}
	return details.Address, nil
}