package resolvers_test

import (
	"context"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/guidewire-oss/fern-mycelium/internal/gql"
	"github.com/guidewire-oss/fern-mycelium/internal/gql/resolvers"
	"github.com/guidewire-oss/fern-mycelium/pkg/repo/fakes"
)

func TestResolvers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Resolver Suite")
}

var _ = Describe("FlakyTests Resolver", func() {
	var (
		fakeRepo *fakes.FakeFlakyTestProvider
		resolver *resolvers.Resolver
		ctx      context.Context
	)

	BeforeEach(func() {
		fakeRepo = &fakes.FakeFlakyTestProvider{}
		resolver = &resolvers.Resolver{FlakyRepo: fakeRepo}
		ctx = context.Background()
	})

	It("should return flaky test data from the fake repository", func() {
		expected := []*gql.FlakyTest{
			{
				TestID:      "test-123",
				TestName:    "Login should timeout on invalid credentials",
				PassRate:    0.7,
				FailureRate: 0.3,
				LastFailure: "2025-04-01T18:00:00Z",
				RunCount:    42,
			},
		}

		fakeRepo.GetFlakyTestsReturns(expected, nil)

		result, err := resolver.Query().FlakyTests(ctx, 1, "policy-admin-ui")

		Expect(err).To(BeNil())
		Expect(result).To(Equal(expected))
		Expect(fakeRepo.GetFlakyTestsCallCount()).To(Equal(1))

		_, projID, limit := fakeRepo.GetFlakyTestsArgsForCall(0)
		Expect(projID).To(Equal("policy-admin-ui"))
		Expect(limit).To(Equal(1))
	})
})
