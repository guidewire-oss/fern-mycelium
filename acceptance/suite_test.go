package acceptance

import (
	"context"
	"fmt"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/guidewire-oss/fern-mycelium/acceptance/fixtures"
	"github.com/guidewire-oss/fern-mycelium/internal/gql"
	"github.com/guidewire-oss/fern-mycelium/internal/gql/resolvers"
	"github.com/guidewire-oss/fern-mycelium/internal/server"
	"github.com/guidewire-oss/fern-mycelium/pkg/repo"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib" // registers the "pgx" driver with database/sql
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var Server *httptest.Server

func TestAcceptance(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Fern Mycelium Acceptance Suite")
}

var _ = BeforeSuite(func() {
	ctx := context.Background()
	container, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("fern"),
		postgres.WithUsername("user"),
		postgres.WithPassword("pass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	Expect(err).ToNot(HaveOccurred())

	host, err := container.Host(ctx)
	Expect(err).ToNot(HaveOccurred())

	port, err := container.MappedPort(ctx, "5432")
	Expect(err).ToNot(HaveOccurred())

	dsn := fmt.Sprintf("postgres://user:pass@%s:%s/fern?sslmode=disable", host, port.Port())
	fmt.Println("✅ Test DB DSN:", dsn)
	os.Setenv("DB_URL", dsn) //nolint:all

	dbpool, err := pgxpool.New(ctx, dsn)
	Expect(err).ToNot(HaveOccurred())

	expectSchema := fixtures.LoadSchema(ctx, dsn)
	Expect(expectSchema).To(Succeed())
	Expect(fixtures.SeedFlakyTests(ctx, dsn)).To(Succeed())

	repo := repo.NewFlakyTestRepo(dbpool)
	schema := gql.NewExecutableSchema(gql.Config{Resolvers: &resolvers.Resolver{FlakyRepo: repo}})
	handler := server.NewGraphQLServer(schema)

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.POST("/query", gin.WrapH(handler))
	Server = httptest.NewServer(r)
})

var _ = AfterSuite(func() {
	if Server != nil {
		Server.Close()
	}
})
