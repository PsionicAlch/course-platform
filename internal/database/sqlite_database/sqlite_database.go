package sqlite_database

import (
	"database/sql"
	"fmt"

	"github.com/PsionicAlch/course-platform/internal/utils"
	_ "modernc.org/sqlite"
)

// TODO: Add database cleanup functions that run in the background.

type SQLiteDatabase struct {
	utils.Loggers
	fileName      string
	migrationsDir string
	connection    *sql.DB
}

func CreateSQLiteDatabase(fileName, migrationsDir string) (*SQLiteDatabase, error) {
	loggers := utils.CreateLoggers("SQLITE DATABASE")

	// Open a connection to the database.
	conn, err := sql.Open("sqlite", fmt.Sprintf(".%s", fileName))
	if err != nil {
		loggers.ErrorLog.Printf("Failed to connect to database: %s", err.Error())
		return nil, err
	}

	// Verify that the connection was successful.
	err = conn.Ping()
	if err != nil {
		loggers.ErrorLog.Printf("Failed to ping database: %s", err.Error())
		return nil, err
	}

	// Set maximum number of database connections to 1 to avoid database is locked error (or SQLITE_BUSY error).
	conn.SetMaxOpenConns(1)

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
		loggers.ErrorLog.Printf("Failed to run database startup query: %s", err.Error())
		return nil, err
	}

	sqliteDatabase := &SQLiteDatabase{
		Loggers:       loggers,
		fileName:      fileName,
		migrationsDir: migrationsDir,
		connection:    conn,
	}

	return sqliteDatabase, nil
}

func (db *SQLiteDatabase) Close() error {
	db.InfoLog.Println("Closing database connection")

	return db.connection.Close()
}
