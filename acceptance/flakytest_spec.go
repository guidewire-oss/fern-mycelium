package acceptance

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("FlakyTests Query", func() {
	It("should return flaky tests", func() {
		By("querying the GraphQL endpoint")

		query := `
			query {
				flakyTests(limit: 5, projectID: "Auth Suite") {
					testName
					passRate
					failureRate
					runCount
					lastFailure
				}
			}
		`

		reqBody, err := json.Marshal(map[string]string{
			"query": query,
		})
		Expect(err).ToNot(HaveOccurred())

		resp, err := http.Post(serverURL(), "application/json", bytes.NewBuffer(reqBody))
		Expect(err).ToNot(HaveOccurred())
		defer resp.Body.Close()

		Expect(resp.StatusCode).To(Equal(http.StatusOK))

		body, err := io.ReadAll(resp.Body)
		Expect(err).ToNot(HaveOccurred())

		var data struct {
			Data struct {
				FlakyTests []map[string]any `json:"flakyTests"`
			} `json:"data"`
		}

		err = json.Unmarshal(body, &data)
		Expect(err).ToNot(HaveOccurred())
		Expect(data.Data.FlakyTests).ToNot(BeEmpty())
		Expect(data.Data.FlakyTests[0]["runCount"]).Should(BeNumerically("==", 2))
		Expect(data.Data.FlakyTests[0]["passRate"]).Should(BeNumerically("==", 0))
		Expect(data.Data.FlakyTests[0]["failureRate"]).Should(BeNumerically("==", 1))
	})
})

func serverURL() string {
	url := os.Getenv("SERVER_URL")
	if url != "" {
		return url
	}
	// fallback if not running against external server
	return Server.URL + "/query"
}
