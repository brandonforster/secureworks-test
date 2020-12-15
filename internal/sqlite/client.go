// sqlite is a package designed to interact specifically with the SQLite3 database engine.
// It is configured by default to use SQLite's "write to disk" functionality to store data in a flat file on the OS.
package sqlite

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	"github.com/brandonforster/resolver/graph/model"
)

// DBClient is the atomic unit of interacting with SQLite.
type Client struct {
	db *sqlx.DB
}

// Close will attempt to close the database connection.
// This should usually be called with a defer statement when opening any connection to the database.
func (c *Client) Close() error {
	return c.db.Close()
}

// NewClient creates a client and opens a connection to a database.
// It allows for connection to arbitrarily named databases.
// It returns a pointer to a DBClient object if successful; an error otherwise.
//
// filename is the name of the flat file where the data will be stored.
//
// It returns a pointer to a DBClient object if successful; an error otherwise.
func NewClient(filename string) (*Client, error) {
	sqliteDb, err := sqlx.Open("sqlite3", fmt.Sprintf("file:%s?_fk=true&_busy_timeout=5000&_journal_mode=WAL", filename))
	if err != nil {
		return nil, err
	}

	err = sqliteDb.Ping()
	if err != nil {
		return nil, err
	}

	return &Client{db: sqliteDb}, nil
}

// AddIPDetails is the creation operation for adding new data to the database.
// It returns a representation of the object added to the database if successful; an error otherwise.
//
// contract is the data model used by most of the system for IPDetails and represents the object to be stored.
//
// It returns a representation of the object added to the database if successful; an error otherwise.
func (c *Client) AddIPDetails(contract model.IPDetails) (*model.IPDetails, error) {
	var modelIP IPDetails
	modelIP.fromContract(contract)

	// since we aren't using the int64 ID, we don't actually care about the result
	_, err := c.db.Exec(
		"INSERT INTO ip(id, created_at, updated_at, response_code, ip_address) VALUES ($1, $2, $3, $4, $5)",
		modelIP.ID,
		modelIP.CreatedAt,
		modelIP.UpdatedAt,
		modelIP.ResponseCode,
		modelIP.IPAddress,
	)
	if err != nil {
		return nil, err
	}

	retval := modelIP.toContract()
	return &retval, nil
}

// UpdateIPDetails is the update operation for editing data to the database.
// It returns a representation of the object updated in the database if successful; an error otherwise.
//
// contract is the data model used by most of the system for IPDetails and represents the object to be stored.
//
// It returns a representation of the object updated it the database if successful; an error otherwise.
func (c *Client) UpdateIPDetails(contract model.IPDetails) (*model.IPDetails, error) {
	var modelIP IPDetails
	modelIP.fromContract(contract)

	// since we aren't using the int64 ID, we don't actually care about the result
	_, err := c.db.Exec(
		"UPDATE ip SET updated_at=$1, response_code=$2 WHERE id = ?",
		modelIP.UpdatedAt,
		modelIP.ResponseCode,
		modelIP.ID,
	)
	if err != nil {
		return nil, err
	}

	retval := modelIP.toContract()
	return &retval, nil
}

// GetIPDetailByAddress is the retrieval operation for getting data from the database.
// It returns a representation of the object requested in the database if successful; an error otherwise.
//
// address is the IP address that we want to query the database for.
//
// It returns a representation of the object requested in the database if successful; an error otherwise.
func (c *Client) GetIPDetailByAddress(address string) (*model.IPDetails, error) {
	var modelIP IPDetails
	err := c.db.Get(&modelIP, "SELECT * FROM ip WHERE ip_address = ?", address)
	if err != nil {
		return nil, err
	}

	retval := modelIP.toContract()
	return &retval, nil
}
