package spamhaus

import (
	"net"
	"strings"
)

// Lookup will query the Spamhaus Domain name System BlockList and return the Return Code associated with given IP.
//
// IP is an IPv4 formatted address to be queried in the Spamhaus system.
//
// The string return value is return code returned back by the Spamhaus system.
func Lookup(IP string) (string, error) {
	// we need to reverse the given IP to do the lookup
	splitAddress := strings.Split(IP, ".")

	for i, j := 0, len(splitAddress)-1; i < len(splitAddress)/2; i, j = i+1, j-1 {
		splitAddress[i], splitAddress[j] = splitAddress[j], splitAddress[i]
	}

	reversed := strings.Join(splitAddress, ".")

	queryAddress := reversed + ".zen.spamhaus.org"

	result, err := net.LookupHost(queryAddress)
	if err != nil {
		return "", err
	}

	print(result)
	return "", nil
}

