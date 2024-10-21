package sqlite_database

import (
	"database/sql"

	"github.com/PsionicAlch/psionicalch-home/internal/utils"
)

type SQLiteDatabase struct {
	utils.Loggers
	connection *sql.DB
}

func CreateSQLiteDatabase() *SQLiteDatabase {
	// Create database loggers.
	dbLoggers := utils.CreateLoggers("DATABASE")

	// Open a connection to the database.
	conn, err := sql.Open("sqlite", "./db/db.sqlite")
	if err != nil {
		dbLoggers.ErrorLog.Fatal("Failed to connect to the database: ", err)
	}

	// Verify that the connection was successful.
	err = conn.Ping()
	if err != nil {
		dbLoggers.ErrorLog.Fatal("Failed to ping the database: ", err)
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
		dbLoggers.ErrorLog.Fatalln("Failed to run performance tuning commands: ", err)
	}

	dbLoggers.InfoLog.Println("Successfully created database connection.")

	return &SQLiteDatabase{
		Loggers:    dbLoggers,
		connection: conn,
	}
}

func (db *SQLiteDatabase) Close() {
	db.connection.Close()
}
