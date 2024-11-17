package sqlite_database

import (
	"github.com/PsionicAlch/psionicalch-home/internal/database"
	"github.com/PsionicAlch/psionicalch-home/internal/database/sqlite_database/internal"
)

func (db *SQLiteDatabase) AddIPAddress(userId, ipAddr string) error {
	id, err := database.GenerateID()
	if err != nil {
		db.ErrorLog.Printf("Failed to generate new ID for IP address: %s\n", err)
		return err
	}

	err = internal.AddIPAddress(db.connection, id, userId, ipAddr)
	if err != nil {
		db.ErrorLog.Printf("Failed to save IP address %s to the database: %s\n", ipAddr, err)
		return err
	}

	return nil
}

func (db *SQLiteDatabase) GetUserIpAddresses(userId string) ([]string, error) {
	query := `SELECT ip_address FROM whitelisted_ips WHERE user_id = ?;`

	var ipAddresses []string

	rows, err := db.connection.Query(query, userId)
	if err != nil {
		db.ErrorLog.Printf("Failed to query database for all user's (\"%s\") whitelisted IP addresses: %s\n", userId, err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var ipAddr string

		err = rows.Scan(&ipAddr)
		if err != nil {
			db.ErrorLog.Printf("Failed to query database row for user's (\"%s\") whitelisted IP address: %s\n", userId, err)
			return nil, err
		}

		ipAddresses = append(ipAddresses, ipAddr)
	}

	if err := rows.Err(); err != nil {
		db.ErrorLog.Printf("Failed to query database for all user's (\"%s\") whitelisted IP addresses: %s\n", userId, err)
		return nil, err
	}

	return ipAddresses, nil
}
