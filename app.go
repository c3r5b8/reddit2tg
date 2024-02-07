package main

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/c3r5b8/r2g/sqlite"
	"github.com/joho/godotenv"

	_ "github.com/mattn/go-sqlite3"
)

var (
	//go:embed schema.sql
	ddl     string
	queries *sqlite.Queries
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	initScarper()
	// go scroblePage(NextUrl)
	db, err := sql.Open("sqlite3", "./posts.db")
	if err != nil {
		fmt.Println(err)
	}

	// create tables
	if _, err := db.ExecContext(context.Background(), ddl); err != nil {
		fmt.Println(err)
	}

	queries = sqlite.New(db)

	ctx := context.Background()

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill, syscall.SIGTERM)
	defer cancel()

	if err := runBot(ctx); err != nil {
		fmt.Println(err)
		fmt.Println("Bot stopped")
		defer os.Exit(1)
	}
}
