package sqlite

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	"github.com/brandonforster/resolver/graph/model"
)

type Client struct {
	db *sqlx.DB
}

func (c *Client) Close() error {
	return c.db.Close()
}

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

func (c *Client) GetIPDetail(id string) (*model.IPDetails, error) {
	var modelIP IPDetails
	err := c.db.Get(&modelIP, "SELECT * FROM ip WHERE id = ?", id)
	if err != nil {
		return nil, err
	}

	retval := modelIP.toContract()
	return &retval, nil
}

func (c *Client) GetIPDetailByAddress(address string) (*model.IPDetails, error) {
	var modelIP IPDetails
	err := c.db.Get(&modelIP, "SELECT * FROM ip WHERE ip_address = ?", address)
	if err != nil {
		return nil, err
	}

	retval := modelIP.toContract()
	return &retval, nil
}

func (c *Client) GetIPDetails() ([]*model.IPDetails, error) {
	var models []IPDetails
	err := c.db.Get(&models, "SELECT * FROM ip;")
	if err != nil {
		return nil, err
	}

	contracts := make([]*model.IPDetails, len(models))
	for i := range contracts {
		thisContract := models[i].toContract()
		contracts[i] = &thisContract
	}

	return contracts, nil
}
