package db

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

func ConnectDB() (*sql.DB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	if host == "" || port == "" || user == "" || password == "" || dbname == "" {
		return nil, fmt.Errorf("database environment variables are required")
	}
	portNumber, err := strconv.Atoi(port)
	if err != nil || portNumber < 1 || portNumber > 65535 {
		return nil, fmt.Errorf("DB_PORT must be a valid port")
	}
	maxOpen, err := connectionLimit("DB_MAX_OPEN_CONNS", 10)
	if err != nil {
		return nil, err
	}
	maxIdle, err := connectionLimit("DB_MAX_IDLE_CONNS", 5)
	if err != nil {
		return nil, err
	}
	if maxOpen == 0 || maxIdle > maxOpen {
		return nil, fmt.Errorf("pool limits must satisfy 0 <= idle <= open and open > 0")
	}
	maxLifetime, err := connectionDuration("DB_CONN_MAX_LIFETIME", 30*time.Minute)
	if err != nil {
		return nil, err
	}
	maxIdleTime, err := connectionDuration("DB_CONN_MAX_IDLE_TIME", 5*time.Minute)
	if err != nil {
		return nil, err
	}
	sslmode := os.Getenv("DB_SSLMODE")
	if sslmode == "" {
		sslmode = "disable"
	}
	databaseURL := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(user, password),
		Host:   host + ":" + port,
		Path:   "/" + dbname,
	}
	values := url.Values{}
	values.Set("sslmode", sslmode)
	databaseURL.RawQuery = values.Encode()
	db, err := sql.Open("postgres", databaseURL.String())
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(maxOpen)
	db.SetMaxIdleConns(maxIdle)
	db.SetConnMaxLifetime(maxLifetime)
	db.SetConnMaxIdleTime(maxIdleTime)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("could not connect to database")
	}
	return db, nil
}

func connectionLimit(name string, fallback int) (int, error) {
	value := os.Getenv(name)
	if value == "" {
		return fallback, nil
	}
	limit, err := strconv.Atoi(value)
	if err != nil || limit < 0 {
		return 0, fmt.Errorf("%s must be a non-negative integer", name)
	}
	return limit, nil
}

func connectionDuration(name string, fallback time.Duration) (time.Duration, error) {
	value := os.Getenv(name)
	if value == "" {
		return fallback, nil
	}
	duration, err := time.ParseDuration(value)
	if err != nil || duration <= 0 {
		return 0, fmt.Errorf("%s must be a positive duration", name)
	}
	return duration, nil
}
