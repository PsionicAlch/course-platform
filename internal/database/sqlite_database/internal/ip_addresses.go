package internal

import (
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

// AddIPAddress adds a new IP address to a user's whitelist of IP addresses. This function works with either a database
// connection or a database transaction. This function will NOT throw an error upon a unique constraint violation.
func AddIPAddress(dbFacade SqlDbFacade, id, userId, ipAddr string) error {
	query := `INSERT INTO whitelisted_ips (id, user_id, ip_address) VALUES (?, ?, ?);`

	_, err := dbFacade.Exec(query, id, userId, ipAddr)
	if err != nil {
		if sqliteErr, ok := err.(*sqlite.Error); ok && sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
			return nil
		}

		return err
	}

	return nil
}
