package db_test

import (
	"log"
	"os"
	"testing"

	"github.com/guidewire-oss/fern-mycelium/internal/db"
	. "github.com/onsi/ginkgo/v2" //nolint:all
	. "github.com/onsi/gomega"    //nolint:all
)

func TestDB(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "DB Connection Suite")
}

var _ = Describe("Database Connection", func() {
	It("should connect to the database with a valid DB_URL", func() {
		if err := os.Setenv("DB_URL", "postgres://user:pass@localhost:5432/fern?sslmode=disable"); err != nil {
			log.Fatalf("Failed to set DB_URL: %v", err)
		}
		Expect(func() { db.Connect() }).ToNot(Panic()) //nolint:all
	})
})
