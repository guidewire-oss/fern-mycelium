package repo_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/guidewire-oss/fern-mycelium/pkg/repo"
	"github.com/guidewire-oss/fern-mycelium/pkg/repo/fakes"
	"github.com/jackc/pgx/v5"
)

type fakeRows struct {
	pgx.Rows
	index int
	data  [][]any
}

func (f *fakeRows) Next() bool {
	return f.index < len(f.data)
}

func (f *fakeRows) Scan(dest ...any) error {
	copy(dest, f.data[f.index])
	f.index++
	return nil
}

func (f *fakeRows) Close() {}

var _ = Describe("FlakyTestRepo", func() {
	var (
		ctx      context.Context
		fakeDB   *fakes.FakePgxQuerier
		repoInst repo.FlakyTestProvider
	)

	BeforeEach(func() {
		ctx = context.Background()
		fakeDB = &fakes.FakePgxQuerier{}
		repoInst = repo.NewFlakyTestRepo(fakeDB)
	})

	It("returns flaky test results from fake rows", func() {
		mockRows := &fakeRows{
			data: [][]any{
				{"LoginSpec", "auth_invalid_token", 40, 12, time.Date(2025, 4, 1, 10, 0, 0, 0, time.UTC)},
			},
		}

		fakeDB.QueryReturns(mockRows, nil)

		results, err := repoInst.GetFlakyTests(ctx, "policy-admin-ui", 1)
		Expect(err).To(BeNil())
		Expect(results).To(HaveLen(1))
		Expect(results[0].TestID).To(Equal("auth_invalid_token"))
		Expect(results[0].FailureRate).To(BeNumerically("~", 0.3, 0.01))
	})
})
