package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// GraphQL request and response structures
type GraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

type GraphQLResponse struct {
	Data   map[string]interface{} `json:"data,omitempty"`
	Errors []GraphQLError         `json:"errors,omitempty"`
}

type GraphQLError struct {
	Message string        `json:"message"`
	Path    []interface{} `json:"path,omitempty"`
}

type HealthResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type FlakyTest struct {
	TestID      string  `json:"testID"`
	TestName    string  `json:"testName"`
	PassRate    float64 `json:"passRate"`
	FailureRate float64 `json:"failureRate"`
	RunCount    int     `json:"runCount"`
	LastFailure *string `json:"lastFailure"`
}

const (
	baseURL = "http://localhost:8081"
)

func main() {
	fmt.Println("🧪 Testing fern-mycelium MCP Server")
	fmt.Println("==================================================")

	// Test 1: Health check
	fmt.Println("\n1. Testing Health Check...")
	if !testHealthCheck() {
		os.Exit(1)
	}

	// Test 2: Query flaky tests
	fmt.Println("\n2. Testing Flaky Test Detection...")
	flakyTests, ok := testFlakyTestQuery()
	if !ok {
		os.Exit(1)
	}

	// Test 3: Demonstrate AI Agent Use Case
	fmt.Println("\n3. AI Agent Analysis Simulation...")
	simulateAgentAnalysis(flakyTests)

	fmt.Println("\n✅ MCP Server Test Complete!")
	fmt.Println("🎯 The fern-mycelium MCP server successfully provides test intelligence context")
	fmt.Println("   that AI agents can use for automated test analysis and recommendations!")
}

func testHealthCheck() bool {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(baseURL + "/healthz")
	if err != nil {
		fmt.Printf("❌ Health check error: %v\n", err)
		return false
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("❌ Health check failed: %d\n", resp.StatusCode)
		return false
	}

	var health HealthResponse
	if err := json.NewDecoder(resp.Body).Decode(&health); err != nil {
		fmt.Printf("❌ Failed to decode health response: %v\n", err)
		return false
	}

	fmt.Printf("✅ Health: %s - %s\n", health.Status, health.Message)
	return true
}

func testFlakyTestQuery() ([]FlakyTest, bool) {
	query := `{
		flakyTests(limit: 10, projectID: "MCP Server Tests") {
			testID
			testName
			passRate
			failureRate
			runCount
			lastFailure
		}
	}`

	result, err := executeGraphQLQuery(query, nil)
	if err != nil {
		fmt.Printf("❌ Failed to query flaky tests: %v\n", err)
		return nil, false
	}

	if len(result.Errors) > 0 {
		fmt.Printf("❌ GraphQL errors: %v\n", result.Errors)
		return nil, false
	}

	// Parse flaky tests from response
	flakyTestsData, ok := result.Data["flakyTests"].([]interface{})
	if !ok {
		fmt.Printf("❌ Invalid flaky tests data format\n")
		return nil, false
	}

	var flakyTests []FlakyTest
	for _, item := range flakyTestsData {
		testMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		test := FlakyTest{
			TestID:      getString(testMap, "testID"),
			TestName:    getString(testMap, "testName"),
			PassRate:    getFloat64(testMap, "passRate"),
			FailureRate: getFloat64(testMap, "failureRate"),
			RunCount:    int(getFloat64(testMap, "runCount")),
		}

		if lastFailure, exists := testMap["lastFailure"]; exists && lastFailure != nil {
			if failureStr, ok := lastFailure.(string); ok {
				test.LastFailure = &failureStr
			}
		}

		flakyTests = append(flakyTests, test)
	}

	fmt.Printf("✅ Found %d tests with flaky behavior:\n", len(flakyTests))
	for _, test := range flakyTests {
		status := "🔴 FLAKY"
		if test.FailureRate == 0 {
			status = "🟢 STABLE"
		}

		fmt.Printf("   %s %s\n", status, test.TestName)
		fmt.Printf("      - Pass Rate: %.1f%%\n", test.PassRate*100)
		fmt.Printf("      - Failure Rate: %.1f%%\n", test.FailureRate*100)
		fmt.Printf("      - Total Runs: %d\n", test.RunCount)
		if test.LastFailure != nil {
			fmt.Printf("      - Last Failure: %s\n", *test.LastFailure)
		}
		fmt.Println()
	}

	return flakyTests, true
}

func simulateAgentAnalysis(flakyTests []FlakyTest) {
	fmt.Println("🤖 Agent: Analyzing test stability patterns...")

	var highRisk, moderateRisk, stable []FlakyTest

	for _, test := range flakyTests {
		if test.FailureRate > 0.3 {
			highRisk = append(highRisk, test)
		} else if test.FailureRate > 0 {
			moderateRisk = append(moderateRisk, test)
		} else {
			stable = append(stable, test)
		}
	}

	fmt.Printf("\n📊 Test Stability Analysis:\n")
	fmt.Printf("   🔴 High Risk (>30%% failure): %d tests\n", len(highRisk))
	fmt.Printf("   🟡 Moderate Risk (1-30%% failure): %d tests\n", len(moderateRisk))
	fmt.Printf("   🟢 Stable (0%% failure): %d tests\n", len(stable))

	if len(highRisk) > 0 {
		fmt.Printf("\n🚨 Agent Recommendation: Immediate attention needed for:\n")
		for _, test := range highRisk {
			fmt.Printf("   - %s (%.1f%% failure rate)\n", test.TestName, test.FailureRate*100)
		}
	}

	if len(moderateRisk) > 0 {
		fmt.Printf("\n⚠️  Agent Recommendation: Monitor closely:\n")
		for _, test := range moderateRisk {
			fmt.Printf("   - %s (%.1f%% failure rate)\n", test.TestName, test.FailureRate*100)
		}
	}
}

func executeGraphQLQuery(query string, variables map[string]interface{}) (*GraphQLResponse, error) {
	reqBody := GraphQLRequest{
		Query:     query,
		Variables: variables,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Post(baseURL+"/query", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error %d: %s", resp.StatusCode, string(body))
	}

	var result GraphQLResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// Helper functions for type conversion
func getString(m map[string]interface{}, key string) string {
	if val, ok := m[key].(string); ok {
		return val
	}
	return ""
}

func getFloat64(m map[string]interface{}, key string) float64 {
	if val, ok := m[key].(float64); ok {
		return val
	}
	return 0
}