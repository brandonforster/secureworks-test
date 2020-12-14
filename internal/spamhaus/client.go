package spamhaus

import (
	"fmt"
	"net"
	"strings"
)

// Lookup will query the Spamhaus Domain name System BlockList and return the Return Code associated with given IP.
//
// IP is an IPv4 formatted address to be queried in the Spamhaus system.
//
// The string return value is return code returned back by the Spamhaus system.
func Lookup(IP string) ([]string, error) {
	if net.ParseIP(IP) == nil {
		return nil, fmt.Errorf("%s is not a valid IPv4 address", IP)
	}

	// we need to reverse the given IP to do the lookup
	octets := strings.Split(IP, ".")

	for i, j := 0, len(octets)-1; i < j; i, j = i+1, j-1 {
		octets[i], octets[j] = octets[j], octets[i]
	}

	reversed := strings.Join(octets, ".")

	queryAddress := reversed + ".zen.spamhaus.org"

	results, err := net.LookupHost(queryAddress)
	if err != nil {
		// unknown to Spamhaus and probably not a spammer
		if strings.Contains(err.Error(), "no such host") {
			return nil, nil
		}

		return nil, err
	}

	for _, code := range results {
		_, err = parseReturnCode(code)
		if err != nil {
			return nil, err
		}
	}

	return results, nil
}

func parseReturnCode(code string) (string, error) {
	knownCodes := map[string]string{
		"127.0.0.2":  "Direct UBE sources, spam operations & spam services",
		"127.0.0.3":  "Direct snowshoe spam sources detected via automation",
		"127.0.0.4":  "CBL (3rd party exploits such as proxies, trojans, etc.)",
		"127.0.0.6":  "CBL (3rd party exploits such as proxies, trojans, etc.)",
		"127.0.0.7":  "CBL (3rd party exploits such as proxies, trojans, etc.)",
		"127.0.0.9":  "Spamhaus DROP/EDROP Data (in addition to 127.0.0.2, since 01-Jun-2016)",
		"127.0.0.10": "End-user Non-MTA IP addresses set by ISP outbound mail policy",
		"127.0.0.11": "End-user Non-MTA IP addresses set by ISP outbound mail policy",
	}

	if _, exists := knownCodes[code]; exists {
		return code, nil
	}

	return "", fmt.Errorf("unknown return code %s", code)
}
