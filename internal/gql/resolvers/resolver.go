package resolvers

import "github.com/guidewire-oss/fern-mycelium/pkg/repo"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	FlakyRepo repo.FlakyTestProvider
}
