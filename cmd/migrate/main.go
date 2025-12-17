package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/xernobyl/formbricks_worktrial/internal/config"
	"github.com/xernobyl/formbricks_worktrial/pkg/database"
)

func main() {
	ctx := context.Background()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Get Postgres pool
	db, err := database.NewPostgresPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := runMigrations(ctx, db); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("All migrations completed successfully")
}

// runMigrations runs the migrations on the given pool
// Gets all .sql files on the migration folder and runs them
func runMigrations(ctx context.Context, db *pgxpool.Pool) error {
	// TODO Add DB versioning somewhere

	migrationsDir := "migrations"

	// Get all files on the migrations folder
	files, err := filepath.Glob(filepath.Join(migrationsDir, "*.sql"))
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	// Sort migrations alphanumerically
	sort.Strings(files)

	for _, file := range files {
		log.Printf("Running migration: %s", filepath.Base(file))

		// Get file content
		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", file, err)
		}

		// Run SQL
		if _, err := db.Exec(ctx, string(content)); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", file, err)
		}

		log.Printf("Completed migration: %s", filepath.Base(file))
	}

	return nil
}
