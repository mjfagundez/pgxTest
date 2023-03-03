package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

func handleError(msg string, err error, exit bool) {
	fmt.Fprintf(os.Stderr, "%s: %v\n", msg, err)
	if exit {
		os.Exit(1)
	}
}

func run(ctx context.Context, dbURL string) error {
	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		handleError("unable to connect to database", err, true)
	}
	defer conn.Close(ctx)

	var name string
	var weight int64
	err = conn.QueryRow(ctx, "select name, weight from widgets where id=$1", 42).Scan(&name, &weight)
	if err != nil {
		handleError("query row failed", err, false)
		return fmt.Errorf("query row failed: %w", err)
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
