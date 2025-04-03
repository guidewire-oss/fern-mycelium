package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.70

import (
	"context"
	"fmt"

	"github.com/guidewire-oss/fern-mycelium/internal/gql"
)

// Health is the resolver for the health field.
func (r *queryResolver) Health(ctx context.Context) (string, error) {
	panic(fmt.Errorf("not implemented: Health - health"))
}

// FlakyTests is the resolver for the flakyTests field.
func (r *queryResolver) FlakyTests(ctx context.Context, limit int, projectID string) ([]*gql.FlakyTest, error) {
	// mock := []*gql.FlakyTest{
	// 	{
	// 		TestID:      "auth-invalid-token",
	// 		TestName:    "LoginService handles expired tokens",
	// 		PassRate:    0.72,
	// 		FailureRate: 0.28,
	// 		LastFailure: "2025-03-30T18:44:10Z",
	// 		RunCount:    50,
	// 	},
	// }

	return r.FlakyRepo.GetFlakyTests(ctx, projectID, limit)
	// Eventually: fetch by projectID from DB
	// return mock, nil
}

// Query returns gql.QueryResolver implementation.
func (r *Resolver) Query() gql.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
