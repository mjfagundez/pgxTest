package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

func handleError(err error) {
	fmt.Fprintf(os.Stderr, "%v\n", err)
	os.Exit(1)
}

func run(ctx context.Context, dbURL string) error {
	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %w", err)
	}
	defer conn.Close(ctx)

	var name string
	var weight int64
	err = conn.QueryRow(ctx, "select name, weight from widgets where id=$1", 42).Scan(&name, &weight)
	if err != nil {
		return fmt.Errorf("query row failed: %w", err)
	}
	fmt.Println(name, weight)
	return nil
}

func main() {
	ctx := context.Background()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		handleError(fmt.Errorf("DATABASE_URL environment variable not set"))
	}

	err := run(ctx, dbURL)
	if err != nil {
		handleError(fmt.Errorf("error running query: %w", err))
	}
}
