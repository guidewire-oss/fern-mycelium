package server

import (
	"context"
	"log"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/guidewire-oss/fern-mycelium/internal/db"
	"github.com/guidewire-oss/fern-mycelium/internal/gql"
	"github.com/guidewire-oss/fern-mycelium/internal/gql/resolvers"
	"github.com/guidewire-oss/fern-mycelium/pkg/repo"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func Start() {
	// Connect to the fern-reporter DB
	pool, err := db.Connect()
	if err != nil {
		log.Fatalf("❌ Failed to get db connection: %v", err)
	}

	// Inject your flaky test provider
	flakyRepo := repo.NewFlakyTestRepo(pool)

	// Create GraphQL schema with real dependencies
	resolver := &resolvers.Resolver{
		FlakyRepo: flakyRepo,
	}
	schema := gql.NewExecutableSchema(gql.Config{Resolvers: resolver})

	// Setup router
	router := gin.Default()

	// Health check endpoint
	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "fern-mycelium is healthy 🍄",
		})
	})

	// GraphQL endpoints
	router.GET("/graphql", gin.WrapH(playground.Handler("Mycel GraphQL Playground", "/query")))
	router.POST("/query", gin.WrapH(NewGraphQLServer(schema)))

	log.Println("🚀 GraphQL Playground available at http://localhost:8080/graphql")
	log.Println("✅ Health check available at http://localhost:8080/healthz")

	// Start server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}

func NewGraphQLServer(schema graphql.ExecutableSchema) *handler.Server {
	srv := handler.New(schema)

	// Add transports (e.g., POST only for production)
	srv.AddTransport(transport.POST{})

	// Optional: configure caching and introspection
	// srv.SetQueryCache(lru.New(1000))
	srv.Use(extension.Introspection{})

	// Optional: error presenter
	srv.SetErrorPresenter(func(ctx context.Context, err error) *gqlerror.Error {
		return graphql.DefaultErrorPresenter(ctx, err)
	})

	return srv
}
