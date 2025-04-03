package server

import (
	"log"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/guidewire-oss/fern-mycelium/internal/db"
	"github.com/guidewire-oss/fern-mycelium/internal/gql"
	"github.com/guidewire-oss/fern-mycelium/internal/gql/resolvers"
)

func Start() {
	db.Connect()
	router := gin.Default()

	// Health check
	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "fern-mycelium is healthy 🍄",
		})
	})

	// GraphQL schema
	schema := gql.NewExecutableSchema(gql.Config{Resolvers: &resolvers.Resolver{}})
	router.GET("/graphql", gin.WrapH(playground.Handler("Mycel GraphQL Playground", "/query")))
	router.POST("/query", gin.WrapH(handler.New(schema)))

	log.Println("🚀 GraphQL Playground available at http://localhost:8080/graphql")
	log.Println("✅ Health check available at http://localhost:8080/healthz")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
