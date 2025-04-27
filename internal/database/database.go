package database

import (
	"database/sql"
	"log/slog"

	_ "github.com/lib/pq"
)

type Database interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	Close() error
}

type PostgresDB struct {
	db *sql.DB
}

func NewPostgreSqlDB(connString string) (*PostgresDB, error) {
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresDB{db: db}, nil
}

func (p *PostgresDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return p.db.Query(query, args...)
}

func (p *PostgresDB) QueryRow(query string, args ...interface{}) *sql.Row {
	return p.db.QueryRow(query, args...)
}

func (p *PostgresDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return p.db.Exec(query, args...)
}

func (p *PostgresDB) Prepare(query string) (*sql.Stmt, error) {
	return p.db.Prepare(query)
}

func (p *PostgresDB) Close() error {
	return p.db.Close()
}

func NewSqlDatabase(host string, port string, user string, password string, dbName string, sslMode string) Database {
	connString := "host=" + host +
		" port=" + port +
		" user=" + user +
		" password=" + password +
		" dbname=" + dbName +
		" sslmode=" + sslMode
	db, err := NewPostgreSqlDB(connString)

	if err != nil {
		slog.Error("Failed to initialize SQL DB")
		panic(err)
	}
	return db
}
