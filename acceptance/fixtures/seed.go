package fixtures

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/guidewire/fern-reporter/pkg/db/migrations"
	_ "github.com/lib/pq"
)

// LoadSchema runs database migrations using fern-reporter's embedded migration files
func LoadSchema(ctx context.Context, dsn string) error {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to open db: %w", err)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create postgres driver: %w", err)
	}

	sourceDriver, err := iofs.New(migrations.EmbeddedMigrations, ".")
	if err != nil {
		return fmt.Errorf("failed to init embedded migration source: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", sourceDriver, "postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to init migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration failed: %w", err)
	}

	return nil
}

//	func LoadSchema(ctx context.Context, db *pgxpool.Pool) error {
//		driver, _ := postgres.WithInstance(db, &postgres.Config{})
//		source, _ := iofs.New(fernmigrations.Migrations, ".")
//		m, _ := migrate.NewWithInstance("iofs", source, "postgres", driver)
//		return m.Up()
//		// path, _ := filepath.Abs("fixtures/schema.sql")
//		// schemaBytes, err := os.ReadFile(path)
//		// if err != nil {
//		// 	return err
//		// }
//		// _, err = db.Exec(ctx, string(schemaBytes))
//		// return err
//	}
func SeedFlakyTests(ctx context.Context, dsn string) error {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to open db: %w", err)
	}
	defer db.Close()

	statements := []string{
		`INSERT INTO test_runs (id,  start_time, end_time, git_branch, git_sha, build_trigger_actor, build_url, test_seed)
     VALUES (1,  CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'main', 'abc123', 'tester', 'https://ci.example.com/build/1', 100);`,

		`INSERT INTO suite_runs (id, test_run_id, suite_name, start_time, end_time)
		 VALUES (1, 1, 'Auth Suite', NOW(), NOW());`,

		`INSERT INTO spec_runs (id, suite_id, spec_description,  status, message, start_time, end_time)
		 VALUES
		 (1, 1, 'LoginService handles expired tokens',  'failed', 'message1', NOW(), NOW()),
		 (2, 1, 'LoginService handles expired tokens',  'failed', 'message2', NOW(), NOW());`,

		`INSERT INTO tags (id, name)
		 VALUES (1, 'flaky');`,

		`INSERT INTO spec_run_tags (spec_run_id, tag_id)
		 VALUES (1, 1);`,

		`INSERT INTO project_details (id, name, team_name,comment, created_at, updated_at)
		 VALUES (1, 'demo', 'team-a', 'comment-1', NOW(), NOW());`,
	}

	for _, stmt := range statements {
		if _, err := db.ExecContext(ctx, stmt); err != nil {
			return fmt.Errorf("seed statement failed: %w", err)
		}
	}

	return nil
}

// func SeedFlakyTests(ctx context.Context, dsn string) error {
// 	db, err := sql.Open("postgres", dsn)
// 	if err != nil {
// 		return fmt.Errorf("failed to open db: %w", err)
// 	}
// 	defer db.Close()
//
// 	// INSERT into projects
// 	_, err = db.ExecContext(ctx,
// 		`INSERT INTO test_runs VALUES (1, 'suite-1', 'demo project', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);`,
// 	)
// 	if err != nil {
// 		return fmt.Errorf("insert into test_runs failed: %w", err)
// 	}
//
// 	// INSERT into suite_runs
// 	_, err = db.ExecContext(ctx,
// 		`INSERT INTO suite_runs (id, project_id, git_branch, git_commit, start_time)
// 		 VALUES ($1, $2, $3, $4, $5)`,
// 		"suite-1", "demo", "main", "abc123", time.Now(),
// 	)
// 	if err != nil {
// 		return fmt.Errorf("insert into suite_runs failed: %w", err)
// 	}
//
// 	// INSERT into spec_runs
// 	_, err = db.ExecContext(ctx,
// 		`INSERT INTO spec_runs (id, suite_run_id, name, file_name, status, duration_ms)
// 		 VALUES
// 		  ($1, $2, $3, $4, $5, $6),
// 		  ($7, $8, $9, $10, $11, $12),
// 		  ($13, $14, $15, $16, $17, $18)`,
// 		"run-1", "suite-1", "LoginService handles expired tokens", "auth-invalid-token", "failed", 120,
// 		"run-2", "suite-1", "LoginService handles expired tokens", "auth-invalid-token", "failed", 115,
// 		"run-3", "suite-1", "LoginService handles valid tokens", "auth-valid-token", "passed", 90,
// 	)
// 	if err != nil {
// 		return fmt.Errorf("insert into spec_runs failed: %w", err)
// 	}
//
// 	return nil
// }
