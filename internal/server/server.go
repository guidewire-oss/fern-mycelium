package server

import (
	"log"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/guidewire-oss/fern-mycelium/internal/db"
	"github.com/guidewire-oss/fern-mycelium/internal/gql"
	"github.com/guidewire-oss/fern-mycelium/internal/gql/resolvers"
	"github.com/guidewire-oss/fern-mycelium/pkg/repo"
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
	router.POST("/query", gin.WrapH(handler.NewDefaultServer(schema)))

	log.Println("🚀 GraphQL Playground available at http://localhost:8080/graphql")
	log.Println("✅ Health check available at http://localhost:8080/healthz")

	// Start server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
