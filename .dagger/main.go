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
		From("golang:1.24.3").
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

func (m *FernMycelium) Publish(
	ctx context.Context,
	src *dagger.Directory,
	version string,
	githubToken dagger.Secret,
) error {
	container, err := m.Build(ctx, src)
	if err != nil {
		return err
	}

	imageTag := fmt.Sprintf("ghcr.io/guidewire-oss/fern-mycelium:%s", version)

	_, err = container.
		WithRegistryAuth("ghcr.io", "guidewire-oss", &githubToken).
		Publish(ctx, imageTag)

	return err
}

// Publish pushes the built image to ttl.sh for temporary sharing
// func (f *FernMycelium) Publish(
// 	ctx context.Context,
// 	// +defaultPath="."
// 	src *dagger.Directory,
// ) (string, error) {
// 	container, err := f.Build(ctx, src)
// 	if err != nil {
// 		return "", err
// 	}
//
// 	timestamp := time.Now().Unix()
// 	tag := fmt.Sprintf("ttl.sh/fern-mycelium-%d:1h", timestamp)
// 	ref, err := container.Publish(ctx, tag)
// 	if err != nil {
// 		return "", err
// 	}

// 	return fmt.Sprintf("✅ Published to %s", ref), nil
// }

// Test runs unit tests using Ginkgo
func (f *FernMycelium) Test(
	ctx context.Context,
	// +defaultPath="."
	src *dagger.Directory,
) (string, error) {
	log.Println("✅ Running Ginkgo tests...")
	output, err := dag.Container().
		From("golang:1.24.3").
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

func (f *FernMycelium) Acceptance(
	ctx context.Context,
	// +defaultPath="."
	src *dagger.Directory,
) (string, error) {
	log.Println("✅ Running Ginkgo tests...")

	output, err := dag.Container().
		From("golang:1.24.3").
		WithMountedDirectory("/src", src).
		WithWorkdir("/src").
		WithServiceBinding("docker", dag.Docker().Cli().Engine()).
		WithEnvVariable("DOCKER_HOST", "tcp://docker:2375").
		WithExec([]string{"go", "install", "github.com/onsi/ginkgo/v2/ginkgo@latest"}).
		WithExec([]string{"ginkgo", "-r", "--vv", "-p", "acceptance/"}).
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
		From("golang:1.24.3").
		WithSecretVariable("GITHUB_AUTH_TOKEN", githubToken).
		WithExec([]string{"go", "install", "github.com/ossf/scorecard/v4@latest"}).
		WithExec([]string{"scorecard", fmt.Sprintf("--repo=%s", repo)}).
		Stdout(ctx)
	if err != nil {
		return "", err
	}
	return output, nil
}

func (m *FernMycelium) Pipeline(
	ctx context.Context,
	src *dagger.Directory,
) error {
	var err error

	// Step 1: Lint
	_, err = m.Lint(ctx, src)
	if err != nil {
		return err
	}

	// Step 2: Test
	_, err = m.Test(ctx, src)
	if err != nil {
		return err
	}

	// Step 3: Scan (with built image name)
	_, err = m.Scan(ctx, src)
	if err != nil {
		return err
	}

	// Step 4: Publish the image
	// _, err = m.Publish(ctx, src)
	// if err != nil {
	// 	return err
	// }

	return nil
}

func (m *FernMycelium) SBOM(ctx context.Context, container *dagger.Container) (*dagger.File, error) {
	syft := dag.Container().
		From("anchore/syft:latest").
		WithMountedDirectory("/input", container.Rootfs()).
		WithWorkdir("/input").
		WithExec([]string{"syft", ".", "-o", "spdx-json", "-q", "--file", "/sbom.json"})

	return syft.File("/sbom.json"), nil
}

func (m *FernMycelium) Cosign(ctx context.Context, image string) error {
	cosign := dag.Container().
		From("gcr.io/projectsigstore/cosign:v2.2.3").
		// WithMountedSecret("/cosign/creds", dag.EnvVariable("GITHUB_TOKEN")).
		WithEnvVariable("COSIGN_EXPERIMENTAL", "1").
		WithEnvVariable("GITHUB_TOKEN", "env://GITHUB_TOKEN").
		WithExec([]string{"cosign", "sign", "--yes", image})

	_, err := cosign.Sync(ctx)
	return err
}

func (m *FernMycelium) Release(ctx context.Context, src *dagger.Directory, version string, githubToken dagger.Secret) error {
	container, err := m.Build(ctx, src)
	if err != nil {
		return err
	}

	// if err := m.Publish(ctx, container, version, githubToken); err != nil {
	// 	return err
	// }

	if err := m.Cosign(ctx, fmt.Sprintf("ghcr.io/guidewire-oss/fern-mycelium:%s", version)); err != nil {
		return err
	}

	sbomFile, err := m.SBOM(ctx, container)
	if err != nil {
		return err
	}

	// Optionally export SBOM file to local or GitHub release asset
	_, err = sbomFile.Export(ctx, "fern-mycelium-sbom.json")
	return err
}

// Deploy deploys the application to k3d cluster using KubeVela
func (m *FernMycelium) Deploy(
	ctx context.Context,
	// +defaultPath="."
	src *dagger.Directory,
) (string, error) {
	log.Println("🚀 Deploying to k3d cluster using KubeVela...")

	// Build the container first
	container, err := m.Build(ctx, src)
	if err != nil {
		return "", fmt.Errorf("failed to build container: %w", err)
	}

	// Load the image into k3d
	imageRef := "fern-mycelium:latest"
	_, err = container.Publish(ctx, imageRef)
	if err != nil {
		return "", fmt.Errorf("failed to publish image: %w", err)
	}

	// Apply KubeVela component definitions and application
	output, err := dag.Container().
		From("oamdev/vela-cli:latest").
		WithMountedDirectory("/manifests", src.Directory("docs/kubevela")).
		WithWorkdir("/manifests").
		WithExec([]string{"kubectl", "create", "namespace", "fern", "--dry-run=client", "-o", "yaml"}).
		WithExec([]string{"kubectl", "apply", "-f", "-"}).
		WithExec([]string{"vela", "def", "apply", "cnpg.cue"}).
		WithExec([]string{"vela", "def", "apply", "gateway.cue"}).
		WithExec([]string{"kubectl", "apply", "-f", "vela.yaml"}).
		Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to apply KubeVela manifests: %w", err)
	}

	return fmt.Sprintf("✅ Deployed successfully using KubeVela:\n%s", output), nil
}

// A coding agent for developing new features
func (m *FernMycelium) Develop(
	ctx context.Context,
	// Assignment to complete
	assignment string,
	// +defaultPath="/"
	source *dagger.Directory,
) (*dagger.Directory, error) {
	// Environment with agent inputs and outputs
	environment := dag.Env(dagger.EnvOpts{Privileged: true}).
		WithStringInput("assignment", assignment, "the assignment to complete").
		WithWorkspaceInput(
			"workspace",
			dag.Workspace(source),
			"the workspace with tools to edit code").
		WithWorkspaceOutput(
			"completed",
			"the workspace with the completed assignment")

	// Detailed prompt stored in markdown file
	promptFile := dag.CurrentModule().Source().File("develop_prompt.md")

	// Put it all together to form the agent
	work := dag.LLM().
		WithEnv(environment).
		WithPromptFile(promptFile)

	// Get the output from the agent
	completed := work.
		Env().
		Output("completed").
		AsWorkspace()
	completedDirectory := completed.GetSource().WithoutDirectory("node_modules")

	// Make sure the tests really pass
	_, err := m.Test(ctx, completedDirectory)
	if err != nil {
		return nil, err
	}

	// Return the Directory with the assignment completed
	return completedDirectory, nil
}
