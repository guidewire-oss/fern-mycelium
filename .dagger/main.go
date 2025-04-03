// A generated module for FernMycelium functions
//
// This module has been generated via dagger init and serves as a reference to
// basic module structure as you get started with Dagger.
//
// Two functions have been pre-created. You can modify, delete, or add to them,
// as needed. They demonstrate usage of arguments and return types using simple
// echo and grep commands. The functions can be called from the dagger CLI or
// from one of the SDKs.
//
// The first line in this comment block is a short description line and the
// rest is a long description with more detail on the module's purpose or usage,
// if appropriate. All modules should have a short description.

package main

import (
	"context"
	"dagger/fern-mycelium/internal/dagger"
	"fmt"
	"log"
	"time"
)

// FernMycelium defines the reusable Dagger pipeline components
type FernMycelium struct{}

func (f *FernMycelium) Build(
	ctx context.Context,
	// +defaultPath="."
	src *dagger.Directory,
) (*dagger.Container, error) {
	log.Println("🔨 Building slim Alpine image with counterfeiter")

	builder := dag.Container().
		From("golang:1.24").
		WithMountedDirectory("/src", src).
		WithWorkdir("/src").
		WithExec([]string{"go", "mod", "tidy"}).
		WithExec([]string{"go", "install", "github.com/maxbrunsfeld/counterfeiter/v6@latest"}).
		WithEnvVariable("PATH", "/go/bin:/usr/local/go/bin:$PATH").
		WithExec([]string{
			"counterfeiter",
			"-generate",
			"-o", "pkg/repo/fakes/fake_flaky_test_provider.go",
			"github.com/guidewire-oss/fern-mycelium/pkg/repo.FlakyTestProvider",
		}).
		WithExec([]string{
			"counterfeiter",
			"-generate",
			"-o", "pkg/repo/fakes/fake_pgx_querier.go",
			"github.com/guidewire-oss/fern-mycelium/pkg/repo.PgxQuerier",
		}).
		WithExec([]string{"go", "build", "-o", "/app/fern-mycelium"})

	runtime := dag.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "--no-cache", "add", "ca-certificates"}).
		WithFile("/fern-mycelium", builder.File("/app/fern-mycelium")).
		WithEntrypoint([]string{"/fern-mycelium"})

	return runtime, nil
}

// Scan runs Trivy scan on the built container image
func (f *FernMycelium) Scan(
	ctx context.Context,
	// +defaultPath="."
	src *dagger.Directory,
) (string, error) {
	log.Println("🔍 Running Trivy filesystem scan on built container...")
	container, err := f.Build(ctx, src)
	if err != nil {
		return "", err
	}

	output, err := dag.Container().
		From("aquasec/trivy:latest").
		WithMountedDirectory("/scan", container.Rootfs()).
		WithExec([]string{"trivy", "fs", "--exit-code", "1", "--severity", "CRITICAL,HIGH", "/scan"}).
		Stdout(ctx)
	if err != nil {
		return "", err
	}

	return output, nil
}

// Publish pushes the built image to ttl.sh for temporary sharing
func (f *FernMycelium) Publish(
	ctx context.Context,
	// +defaultPath="."
	src *dagger.Directory,
) (string, error) {
	container, err := f.Build(ctx, src)
	if err != nil {
		return "", err
	}

	timestamp := time.Now().Unix()
	tag := fmt.Sprintf("ttl.sh/fern-mycelium-%d:1h", timestamp)
	ref, err := container.Publish(ctx, tag)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("✅ Published to %s", ref), nil
}

// Test runs unit tests using Ginkgo
func (f *FernMycelium) Test(
	ctx context.Context,
	// +defaultPath="."
	src *dagger.Directory,
) (string, error) {
	log.Println("✅ Running Ginkgo tests...")
	output, err := dag.Container().
		From("golang:1.24").
		WithMountedDirectory("/src", src).
		WithWorkdir("/src").
		WithExec([]string{"go", "install", "github.com/onsi/ginkgo/v2/ginkgo@latest"}).
		WithExec([]string{"ginkgo", "-r", "-p", "--skip-package", "acceptance"}).
		Stdout(ctx)
	if err != nil {
		return "", err
	}
	return output, nil
}

// Lint runs static analysis with golangci-lint
func (f *FernMycelium) Lint(
	ctx context.Context,
	// +defaultPath="."
	src *dagger.Directory,
) (string, error) {
	log.Println("🧼 Linting with golangci-lint...")
	output, err := dag.Container().
		From("golangci/golangci-lint:latest").
		WithMountedDirectory("/src", src).
		WithWorkdir("/src").
		WithExec([]string{"golangci-lint", "run", "--timeout=3m"}).
		Stdout(ctx)
	if err != nil {
		return "", err
	}
	return output, nil
}

// CheckOpenSSF runs OpenSSF Scorecard analysis using CLI with GitHub token
func (f *FernMycelium) CheckOpenSSF(
	ctx context.Context,
	repo string,
	// +optional
	// +secret
	githubToken *dagger.Secret,
) (string, error) {
	log.Println("🛡 Running OpenSSF Scorecard with GitHub token...")

	output, err := dag.Container().
		From("golang:1.24").
		WithSecretVariable("GITHUB_AUTH_TOKEN", githubToken).
		WithExec([]string{"go", "install", "github.com/ossf/scorecard/v4@latest"}).
		WithExec([]string{"scorecard", fmt.Sprintf("--repo=%s", repo)}).
		Stdout(ctx)
	if err != nil {
		return "", err
	}
	return output, nil
}
