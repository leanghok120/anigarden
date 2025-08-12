package main

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func getDbPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	appDir := filepath.Join(dir, "anigarden")
	if err := os.MkdirAll(appDir, 0775); err != nil {
		return "", err
	}

	return filepath.Join(appDir, "watchlist.db"), nil
}

func initDB() {
	var err error
	var dbPath string

	dbPath, err = getDbPath()
	if err != nil {
		log.Fatalf("failed to get db path: %v\n", err)
	}

	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("failed to init db: %v\n", err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS watchlist (
		anime_id TEXT NOT NULL UNIQUE,
		added_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)
	`)
	if err != nil {
		log.Fatalf("failed to create watchlist table: %v\n", err)
	}
}

func getWatchlist() []string {
	rows, err := db.Query(`SELECT anime_id FROM watchlist ORDER BY added_at DESC`)
	if err != nil {
		log.Fatalf("failed to get watchlist: %v\n", err)
	}
	defer rows.Close()

	var animeIds []string
	for rows.Next() {
		var animeId string
		if err := rows.Scan(&animeId); err != nil {
			log.Fatalf("failed to scan rows from watchlist: %v\n", err)
		}
		animeIds = append(animeIds, animeId)
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("error iterating watchlist rows: %v\n", err)
	}

	return animeIds
}

func addAnimeToWatchlist(animeId string) {
	_, err := db.Exec(`INSERT INTO watchlist (anime_id) VALUES (?)`, animeId)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			return
		}
		log.Fatalf("failed to add %s to watchlist: %v\n", animeId, err)
	}
}

func removeAnimeFromWatchlist(animeId string) {
	_, err := db.Exec(`DELETE FROM watchlist WHERE anime_id = (?)`, animeId)
	if err != nil {
		log.Fatalf("failed to remove %s from watchlist: %v\n", animeId, err)
	}
}
