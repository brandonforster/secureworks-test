package sqlite

import (
	"fmt"

	"github.com/jmoiron/sqlx"

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
		return nil, fmt.Errorf("could not open sqlite db: %s", err)
	}

	err = sqliteDb.Ping()
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %s", err)
	}

	return &Client{db: sqliteDb}, nil
}

func (c *Client) AddIPDetails(newIP model.IPDetails) (model.IPDetails, error) {

}
