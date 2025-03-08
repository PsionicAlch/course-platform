package sqlite_database

import (
	"github.com/PsionicAlch/course-platform/internal/database"
	"github.com/PsionicAlch/course-platform/internal/database/models"
	"github.com/PsionicAlch/course-platform/internal/database/sqlite_database/internal"
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

func (db *SQLiteDatabase) GetUserIpAddresses(userId string) ([]*models.WhitelistedIPModel, error) {
	query := `SELECT id, user_id, ip_address, created_at FROM whitelisted_ips WHERE user_id = ? ORDER BY created_at DESC;`

	var ipAddresses []*models.WhitelistedIPModel

	rows, err := db.connection.Query(query, userId)
	if err != nil {
		db.ErrorLog.Printf("Failed to query database for all user's (\"%s\") whitelisted IP addresses: %s\n", userId, err)
		return nil, err
	}

	for rows.Next() {
		var ipAddr models.WhitelistedIPModel

		err = rows.Scan(&ipAddr.ID, &ipAddr.UserID, &ipAddr.IPAddress, &ipAddr.CreatedAt)
		if err != nil {
			db.ErrorLog.Printf("Failed to query database row for user's (\"%s\") whitelisted IP address: %s\n", userId, err)
			return nil, err
		}

		ipAddresses = append(ipAddresses, &ipAddr)
	}

	if err := rows.Err(); err != nil {
		db.ErrorLog.Printf("Failed to query database for all user's (\"%s\") whitelisted IP addresses: %s\n", userId, err)
		return nil, err
	}

	return ipAddresses, nil
}

func (db *SQLiteDatabase) DeleteIPAddress(ipAddrId, userId string) error {
	query := `DELETE FROM whitelisted_ips WHERE id = ? AND user_id = ?;`

	_, err := db.connection.Exec(query, ipAddrId, userId)
	if err != nil {
		db.ErrorLog.Printf("Failed to delete IP address (\"%s\") from database (user id: \"%s\"): %s\n", ipAddrId, userId, err)
		return err
	}

	return nil
}
