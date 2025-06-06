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
        spec_runs.spec_description AS test_name,
        COUNT(*) AS total_runs,
        COUNT(*) FILTER (WHERE spec_runs.status <> 'passed') AS failure_count,
        MAX(spec_runs.end_time) FILTER (WHERE spec_runs.status <> 'passed') AS last_failure
    FROM spec_runs
    JOIN suite_runs ON spec_runs.suite_id = suite_runs.id
    WHERE suite_runs.suite_name = $1
    GROUP BY spec_runs.spec_description
    ORDER BY (COUNT(*) FILTER (WHERE spec_runs.status <> 'passed'))::float / COUNT(*) DESC
    LIMIT $2;
	`
	rows, err := r.db.Query(ctx, query, projectID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*gql.FlakyTest

	for rows.Next() {
		var testName string
		var runCount, failureCount int
		var lastFailure *time.Time

		if err := rows.Scan(&testName, &runCount, &failureCount, &lastFailure); err != nil {
			return nil, err
		}

		test := &gql.FlakyTest{
			TestID:      testName, // Use test name as ID for now
			TestName:    testName,
			PassRate:    float64(runCount-failureCount) / float64(runCount),
			FailureRate: float64(failureCount) / float64(runCount),
			RunCount:    runCount,
		}

		if lastFailure != nil {
			formattedTime := lastFailure.Format(time.RFC3339)
			test.LastFailure = &formattedTime
		}

		results = append(results, test)
	}

	return results, nil
}
