package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

func handleError(msg string, err error, exit bool) error {
	fmt.Fprintf(os.Stderr, "%s: %v\n", msg, err)
	if exit {
		os.Exit(1)
	}
	return fmt.Errorf("%s: %w", msg, err)
}

func createTable(ctx context.Context, conn *pgx.Conn) error {
	_, err := conn.Exec(ctx, "CREATE TABLE IF NOT EXISTS widgets (name TEXT PRIMARY KEY, weight BIGINT)")
	if err != nil {
		handleError("unable to create table", err, false)
	}
	return err
}

func insertData(ctx context.Context, conn *pgx.Conn, name string, weight int64) error {
	_, err := conn.Exec(ctx, "INSERT INTO widgets (name, weight) VALUES ($1, $2)", "Freeman", 80)
	if err != nil {
		handleError("unable to insert data", err, false)
	}
	return err
}

func run(ctx context.Context, dbURL string) error {
	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		handleError("unable to connect to database", err, true)
	}
	defer conn.Close(ctx)

	err = createTable(ctx, conn)
	if err != nil {
		handleError("unable to create table", err, false)
	}

	err = insertData(ctx, conn, "widget1", 100)
	if err != nil {
		handleError("unable to insert data", err, false)
	}

	var name string
	var weight int64
	err = conn.QueryRow(ctx, "select name, weight from widgets where =$1", 42).Scan(&name, &weight)
	if err != nil {
		handleError("query row failed", err, false)
	}
	fmt.Println(name, weight)
	return nil
}

func main() {
	ctx := context.Background()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		handleError("DATABASE_URL environment variable not set", fmt.Errorf("missing variable"), true)
	}

	if err := run(ctx, dbURL); err != nil {
		handleError("run failed", err, true)
	}
}
