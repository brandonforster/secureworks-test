// interfaces defines a set of contracts that implementations can be built to.
// The application should use contracts and not implementations whenever possible.
package interfaces

import "github.com/brandonforster/resolver/graph/model"

// DBClient describes the principal methods needed to operate a database with the system.
type DBClient interface {
	// Close will attempt to close the database connection.
	// This should usually be called with a defer statement when opening any connection to the database.
	Close() error

	// AddIPDetails is the creation operation for adding new data to the database.
	// It returns a representation of the object added to the database if successful; an error otherwise.
	//
	// contract is the data model used by most of the system for IPDetails and represents the object to be stored.
	//
	// It returns a representation of the object added to the database if successful; an error otherwise.
	AddIPDetails(contract model.IPDetails) (*model.IPDetails, error)

	// UpdateIPDetails is the update operation for editing data to the database.
	// It returns a representation of the object updated in the database if successful; an error otherwise.
	//
	// contract is the data model used by most of the system for IPDetails and represents the object to be stored.
	//
	// It returns a representation of the object updated it the database if successful; an error otherwise.
	UpdateIPDetails(contract model.IPDetails) (*model.IPDetails, error)

	// GetIPDetailByAddress is the retrieval operation for getting data from the database.
	// It returns a representation of the object requested in the database if successful; an error otherwise.
	//
	// address is the IP address that we want to query the database for.
	//
	// It returns a representation of the object requested in the database if successful; an error otherwise.
	GetIPDetailByAddress(address string) (*model.IPDetails, error)
}
