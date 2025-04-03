package repo

import (
	"context"
	"time"

	"github.com/guidewire-oss/fern-mycelium/internal/gql"
	"github.com/jackc/pgx/v5"
)

//go:generate counterfeiter -o fakes/fake_flaky_test_provider.go . FlakyTestProvider
type FlakyTestProvider interface {
	GetFlakyTests(ctx context.Context, projectID string, limit int) ([]*gql.FlakyTest, error)
}

//go:generate counterfeiter -o fakes/fake_pgx_querier.go . PgxQuerier
type PgxQuerier interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

type FlakyTestRepo struct {
	db PgxQuerier
}

func NewFlakyTestRepo(db PgxQuerier) *FlakyTestRepo {
	return &FlakyTestRepo{db: db}
}

func (r *FlakyTestRepo) GetFlakyTests(ctx context.Context, projectID string, limit int) ([]*gql.FlakyTest, error) {
	query := `
		SELECT
			spec_runs.name AS test_name,
			spec_runs.file_name AS test_id,
			COUNT(*) AS total_runs,
			COUNT(*) FILTER (WHERE spec_runs.status != 'passed') AS failure_count,
			MAX(suite_runs.start_time) FILTER (WHERE spec_runs.status != 'passed') AS last_failure
		FROM spec_runs
		JOIN suite_runs ON spec_runs.suite_run_id = suite_runs.id
		WHERE suite_runs.project_id = $1
		GROUP BY spec_runs.name, spec_runs.file_name
		ORDER BY failure_count::float / COUNT(*) DESC
		LIMIT $2;
	`

	rows, err := r.db.Query(ctx, query, projectID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*gql.FlakyTest

	for rows.Next() {
		var testName, testID string
		var runCount, failureCount int
		var lastFailure time.Time

		if err := rows.Scan(&testName, &testID, &runCount, &failureCount, &lastFailure); err != nil {
			return nil, err
		}

		results = append(results, &gql.FlakyTest{
			TestID:      testID,
			TestName:    testName,
			PassRate:    float64(runCount-failureCount) / float64(runCount),
			FailureRate: float64(failureCount) / float64(runCount),
			LastFailure: lastFailure.Format(time.RFC3339),
			RunCount:    runCount,
		})
	}

	return results, nil
}
