#!/usr/bin/env python3
"""
Simple MCP client to test fern-mycelium MCP server functionality.
This script demonstrates how an AI agent would query the test context.
"""

import json
import requests
import sys

def test_graphql_query(query, variables=None):
    """Test a GraphQL query against the fern-mycelium server."""
    url = "http://localhost:8081/query"
    payload = {
        "query": query
    }
    if variables:
        payload["variables"] = variables
    
    try:
        response = requests.post(url, json=payload, headers={'Content-Type': 'application/json'})
        response.raise_for_status()
        return response.json()
    except requests.exceptions.RequestException as e:
        print(f"❌ Error querying GraphQL: {e}")
        return None

def main():
    print("🧪 Testing fern-mycelium MCP Server")
    print("=" * 50)
    
    # Test 1: Health check
    print("\n1. Testing Health Check...")
    try:
        response = requests.get("http://localhost:8081/healthz")
        if response.status_code == 200:
            health_data = response.json()
            print(f"✅ Health: {health_data['status']} - {health_data['message']}")
        else:
            print(f"❌ Health check failed: {response.status_code}")
            return
    except Exception as e:
        print(f"❌ Health check error: {e}")
        return
    
    # Test 2: Query flaky tests
    print("\n2. Testing Flaky Test Detection...")
    flaky_query = """
    {
        flakyTests(limit: 10, projectID: "MCP Server Tests") {
            testID
            testName
            passRate
            failureRate
            runCount
            lastFailure
        }
    }
    """
    
    result = test_graphql_query(flaky_query)
    if result and 'data' in result and 'flakyTests' in result['data']:
        flaky_tests = result['data']['flakyTests']
        print(f"✅ Found {len(flaky_tests)} tests with flaky behavior:")
        
        for test in flaky_tests:
            status = "🔴 FLAKY" if test['failureRate'] > 0 else "🟢 STABLE"
            print(f"   {status} {test['testName']}")
            print(f"      - Pass Rate: {test['passRate']:.1%}")
            print(f"      - Failure Rate: {test['failureRate']:.1%}")
            print(f"      - Total Runs: {test['runCount']}")
            if test['lastFailure']:
                print(f"      - Last Failure: {test['lastFailure']}")
            print()
    else:
        print(f"❌ Failed to query flaky tests: {result}")
        return
    
    # Test 3: Demonstrate AI Agent Use Case
    print("3. AI Agent Analysis Simulation...")
    print("🤖 Agent: Analyzing test stability patterns...")
    
    high_risk_tests = [test for test in flaky_tests if test['failureRate'] > 0.3]
    moderate_risk_tests = [test for test in flaky_tests if 0 < test['failureRate'] <= 0.3]
    stable_tests = [test for test in flaky_tests if test['failureRate'] == 0]
    
    print(f"\n📊 Test Stability Analysis:")
    print(f"   🔴 High Risk (>30% failure): {len(high_risk_tests)} tests")
    print(f"   🟡 Moderate Risk (1-30% failure): {len(moderate_risk_tests)} tests")
    print(f"   🟢 Stable (0% failure): {len(stable_tests)} tests")
    
    if high_risk_tests:
        print(f"\n🚨 Agent Recommendation: Immediate attention needed for:")
        for test in high_risk_tests:
            print(f"   - {test['testName']} ({test['failureRate']:.1%} failure rate)")
    
    if moderate_risk_tests:
        print(f"\n⚠️  Agent Recommendation: Monitor closely:")
        for test in moderate_risk_tests:
            print(f"   - {test['testName']} ({test['failureRate']:.1%} failure rate)")
    
    print("\n✅ MCP Server Test Complete!")
    print("🎯 The fern-mycelium MCP server successfully provides test intelligence context")
    print("   that AI agents can use for automated test analysis and recommendations!")

if __name__ == "__main__":
    main()