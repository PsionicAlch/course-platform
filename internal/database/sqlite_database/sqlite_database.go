package sqlite_database

import (
	"database/sql"
	"fmt"

	"github.com/PsionicAlch/psionicalch-home/internal/database/errors"
	_ "modernc.org/sqlite"
)

type SQLiteDatabase struct {
	fileName      string
	migrationsDir string
	connection    *sql.DB
}

func CreateSQLiteDatabase(fileName, migrationsDir string) (*SQLiteDatabase, error) {
	// Open a connection to the database.
	conn, err := sql.Open("sqlite", fmt.Sprintf(".%s", fileName))
	if err != nil {
		return nil, errors.CreateFailedToConnectToDatabase(err.Error())
	}

	// Verify that the connection was successful.
	err = conn.Ping()
	if err != nil {
		return nil, errors.CreateFailedToConnectToDatabase(err.Error())
	}

	// Set maximum number of database connections to 1 to avoid database is locked error (or SQLITE_BUSY error).
	// conn.SetMaxOpenConns(1)

	// SQLite performance tuning according to https://phiresky.github.io/blog/2020/sqlite-performance-tuning/.
	_, err = conn.Exec(`pragma journal_mode = WAL;
	pragma synchronous = normal;
	pragma temp_store = memory;
	pragma mmap_size = 30000000000;
	pragma page_size = 32768;
	pragma vacuum;
	pragma optimize;
	`)
	if err != nil {
		return nil, errors.CreateFailedToConnectToDatabase(err.Error())
	}

	sqliteDatabase := &SQLiteDatabase{
		fileName:      fileName,
		migrationsDir: migrationsDir,
		connection:    conn,
	}

	return sqliteDatabase, nil
}

func (db *SQLiteDatabase) Close() error {
	return db.connection.Close()
}
