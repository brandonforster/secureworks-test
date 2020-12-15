package interfaces

// BlocklistClient defines the minimal behavior needed to perform lookups on an external blocklist
type BlocklistClient interface {
	// Lookup will query an externally maintained blocklist and return some return status code associated with given IP.
	//
	// IP is an IPv4 formatted address to be queried in the externally maintained blocklist.
	//
	// The string return value is return code returned back by the externally maintained blocklist.
	Lookup(IP string) ([]string, error)
}
