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
		return fmt.Errorf("failed to create table, %w", err)
	}
	return nil
}

func insertData(ctx context.Context, conn *pgx.Conn, name string, weight int64) error {
	_, err := conn.Exec(ctx, "INSERT INTO widgets (name, weight) VALUES ($1, $2)", name, weight)
	if err != nil {
		return fmt.Errorf("failed to insert widget, %w", err)
	}
	return nil
}

func run(ctx context.Context) error {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		handleError("DATABASE_URL environment variable not set", fmt.Errorf("missing variable"), true)
	}

	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		return fmt.Errorf("unable to connect to database, %w", err)
	}
	defer conn.Close(ctx)

	err = createTable(ctx, conn)
	if err != nil {
		return err
	}

	err = insertData(ctx, conn, "Freeman", 100)
	if err != nil {
		return err
	}

	var name string
	var weight int64
	err = conn.QueryRow(ctx, "select name, weight from widgets where name=$1", "Freeman").Scan(&name, &weight)
	if err != nil {
		return err
	}
	fmt.Println(name, weight)
	return nil
}

func main() {
	ctx := context.Background()

	if err := run(ctx); err != nil {
		handleError("run failed", err, true)
	}
}
