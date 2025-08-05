package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
)

// DB holds the database connection pool.
type DB struct {
	*sql.DB
}

// Config holds the database connection details.
type Config struct {
	User     string
	Password string
	Host     string
	Port     string
	DBName   string
}

// New creates and returns a new DB connection pool.
func New(cfg Config) (*DB, error) {
	// DSN: Data Source Name
	// [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)

	// Open a connection pool.
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("could not open sql connection: %w", err)
	}

	// Set connection pool settings for better performance.
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Ping the database to verify the connection is alive.
	if err := db.Ping(); err != nil {
		db.Close() // Close the connection if ping fails.
		return nil, fmt.Errorf("could not ping database: %w", err)
	}

	return &DB{db}, nil
}

// QueryRow executes a query that is expected to return at most one row.
func (d *DB) QueryRow(query string, args ...interface{}) *sql.Row {
	log.Printf("Executing query: %s with args: %v", query, args)
	return d.DB.QueryRow(query, args...)
}

// Query executes a query that returns rows, typically a SELECT.
func (d *DB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	log.Printf("Executing query: %s with args: %v", query, args)
	return d.DB.Query(query, args...)
}

// Exec executes a query without returning any rows.
// The args are for any placeholder parameters in the query.
func (d *DB) Exec(query string, args ...interface{}) (sql.Result, error) {
	log.Printf("Executing exec query: %s with args: %v", query, args)
	return d.DB.Exec(query, args...)
}

// GetTables returns a list of all tables in the database.
func (d *DB) GetTables() ([]string, error) {
	rows, err := d.Query("SHOW TABLES")
	if err != nil {
		return nil, fmt.Errorf("could not show tables: %w", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, fmt.Errorf("could not scan table name: %w", err)
		}
		tables = append(tables, tableName)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return tables, nil
}

// Transaction executes a function within a database transaction.
// If the function returns an error, the transaction is rolled back. Otherwise, it's committed.
func (d *DB) Transaction(fn func(*sql.Tx) error) error {
	tx, err := d.Begin()
	if err != nil {
		return fmt.Errorf("could not begin transaction: %w", err)
	}

	if err := fn(tx); err != nil {
		// If an error occurs, roll back the transaction.
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("transaction error: %v, rollback error: %v", err, rbErr)
		}
		return fmt.Errorf("transaction error: %w", err)
	}

	// If everything is fine, commit the transaction.
	return tx.Commit()
}
