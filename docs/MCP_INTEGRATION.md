# Fern-Mycelium MCP Server Integration Guide

## Overview

Fern-Mycelium is an intelligent test context layer that provides AI agents with deep insights into test execution patterns, flaky test detection, and automated test analysis through the Model Context Protocol (MCP). This guide covers how to integrate the fern-mycelium MCP server with popular Large Language Models (LLMs) and AI platforms.

## What is Model Context Protocol (MCP)?

The Model Context Protocol is an open standard that enables AI assistants to securely access external data sources and tools. Fern-mycelium implements an MCP server that exposes test intelligence data through a standardized interface, allowing AI agents to:

- Query flaky test detection results
- Analyze test stability patterns
- Get automated recommendations for test improvements
- Access historical test execution data

## Quick Start

### 1. Deploy the MCP Server

```bash
# Deploy using KubeVela (recommended for production)
cd docs/kubevela
vela up -f vela.yaml

# Or run locally for development
export DB_URL="postgres://user:password@localhost:5432/fern_mycelium"
go run main.go serve --port 8081
```

### 2. Verify Server Health

```bash
curl http://localhost:8081/healthz
# Expected: {"status":"ok","message":"fern-mycelium is healthy 🍄"}
```

### 3. Test GraphQL API

```bash
curl -X POST http://localhost:8081/query \
  -H "Content-Type: application/json" \
  -d '{
    "query": "{ flakyTests(limit: 5, projectID: \"your-project\") { testName passRate failureRate runCount lastFailure } }"
  }'
```

## Integration Patterns

### Claude Desktop Integration

Claude Desktop natively supports MCP servers. Add fern-mycelium to your Claude Desktop configuration:

#### Configuration File Location
- **macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`
- **Windows**: `%APPDATA%\Claude\claude_desktop_config.json`
- **Linux**: `~/.config/Claude/claude_desktop_config.json`

#### Configuration Example

```json
{
  "mcpServers": {
    "fern-mycelium": {
      "command": "go",
      "args": ["run", "/path/to/fern-mycelium/main.go", "serve"],
      "env": {
        "DB_URL": "postgres://user:password@localhost:5432/fern_mycelium"
      }
    }
  }
}
```

#### Using with Claude Desktop

Once configured, you can ask Claude to analyze your tests:

```
"Can you analyze our flaky tests and provide recommendations for the most critical issues?"

"Show me test stability trends for the authentication module"

"What tests should we prioritize fixing based on their failure patterns?"
```

### OpenAI Integration

#### Custom GPT Integration

Create a custom GPT with fern-mycelium integration:

1. **Create a Custom GPT**
   - Go to OpenAI's GPT Builder
   - Configure the GPT with test analysis capabilities

2. **Add API Integration**
   ```json
   {
     "openapi": "3.0.0",
     "info": {
       "title": "Fern-Mycelium Test Intelligence API",
       "version": "1.0.0"
     },
     "servers": [
       {
         "url": "http://your-server:8081"
       }
     ],
     "paths": {
       "/query": {
         "post": {
           "operationId": "queryTests",
           "summary": "Query test data using GraphQL",
           "requestBody": {
             "required": true,
             "content": {
               "application/json": {
                 "schema": {
                   "type": "object",
                   "properties": {
                     "query": {"type": "string"},
                     "variables": {"type": "object"}
                   }
                 }
               }
             }
           }
         }
       }
     }
   }
   ```

#### Assistant API Integration

Use OpenAI's Assistant API with function calling:

```python
import openai
import requests

def query_flaky_tests(project_id, limit=10):
    """Function to query flaky tests from fern-mycelium"""
    query = """
    {
        flakyTests(limit: %d, projectID: "%s") {
            testName
            passRate
            failureRate
            runCount
            lastFailure
        }
    }
    """ % (limit, project_id)
    
    response = requests.post(
        "http://localhost:8081/query",
        json={"query": query},
        headers={"Content-Type": "application/json"}
    )
    return response.json()

# Configure the assistant
assistant = openai.beta.assistants.create(
    name="Test Intelligence Assistant",
    instructions="You are a test analysis expert. Use the fern-mycelium API to analyze test patterns and provide actionable recommendations.",
    tools=[{
        "type": "function",
        "function": {
            "name": "query_flaky_tests",
            "description": "Query flaky test data from the test intelligence system",
            "parameters": {
                "type": "object",
                "properties": {
                    "project_id": {"type": "string"},
                    "limit": {"type": "integer", "default": 10}
                }
            }
        }
    }]
)
```

### Langchain Integration

Integrate fern-mycelium with Langchain applications:

```python
from langchain.tools import Tool
from langchain.agents import initialize_agent, AgentType
from langchain.llms import OpenAI
import requests

def query_test_intelligence(query_input):
    """Langchain tool for querying test intelligence"""
    project_id, limit = query_input.split(",")
    limit = int(limit.strip()) if limit.strip().isdigit() else 10
    
    graphql_query = f"""
    {{
        flakyTests(limit: {limit}, projectID: "{project_id.strip()}") {{
            testName
            passRate
            failureRate
            runCount
            lastFailure
        }}
    }}
    """
    
    response = requests.post(
        "http://localhost:8081/query",
        json={"query": graphql_query}
    )
    return response.json()

# Create Langchain tool
test_intelligence_tool = Tool(
    name="Test Intelligence",
    description="Query flaky test data. Input should be 'project_id, limit'",
    func=query_test_intelligence
)

# Initialize agent
llm = OpenAI(temperature=0)
agent = initialize_agent(
    [test_intelligence_tool],
    llm,
    agent=AgentType.ZERO_SHOT_REACT_DESCRIPTION,
    verbose=True
)

# Use the agent
result = agent.run("Analyze flaky tests for project 'web-app' and recommend fixes")
```

### LangGraph Integration

For more complex workflows with LangGraph:

```python
from langgraph.graph import StateGraph, END
from typing import TypedDict
import requests

class TestAnalysisState(TypedDict):
    project_id: str
    test_data: dict
    analysis_result: str
    recommendations: list

def fetch_test_data(state: TestAnalysisState):
    """Fetch test data from fern-mycelium"""
    query = f"""
    {{
        flakyTests(limit: 20, projectID: "{state['project_id']}") {{
            testName
            passRate
            failureRate
            runCount
            lastFailure
        }}
    }}
    """
    
    response = requests.post(
        "http://localhost:8081/query",
        json={"query": query}
    )
    
    state["test_data"] = response.json()
    return state

def analyze_test_patterns(state: TestAnalysisState):
    """Analyze test patterns using LLM"""
    # Your LLM analysis logic here
    pass

def generate_recommendations(state: TestAnalysisState):
    """Generate actionable recommendations"""
    # Your recommendation logic here
    pass

# Build the graph
workflow = StateGraph(TestAnalysisState)
workflow.add_node("fetch_data", fetch_test_data)
workflow.add_node("analyze", analyze_test_patterns)
workflow.add_node("recommend", generate_recommendations)

workflow.set_entry_point("fetch_data")
workflow.add_edge("fetch_data", "analyze")
workflow.add_edge("analyze", "recommend")
workflow.add_edge("recommend", END)

app = workflow.compile()
```

## Custom LLM Integration

### REST API Integration

For any LLM platform that supports HTTP requests:

```javascript
// Example: Integrating with Anthropic's API
const anthropic = new Anthropic({
  apiKey: process.env.ANTHROPIC_API_KEY,
});

async function analyzeTestsWithClaude(projectId) {
  // First, fetch test data
  const testData = await fetch('http://localhost:8081/query', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      query: `{
        flakyTests(limit: 10, projectID: "${projectId}") {
          testName
          passRate
          failureRate
          runCount
          lastFailure
        }
      }`
    })
  }).then(r => r.json());

  // Then analyze with Claude
  const message = await anthropic.messages.create({
    model: "claude-3-sonnet-20240229",
    max_tokens: 1000,
    messages: [{
      role: "user",
      content: `Analyze this test data and provide recommendations: ${JSON.stringify(testData)}`
    }]
  });

  return message.content;
}
```

### WebSocket Integration

For real-time test analysis:

```python
import asyncio
import websockets
import json

async def real_time_test_monitor():
    """Real-time test monitoring with LLM analysis"""
    uri = "ws://localhost:8081/ws"  # If you add WebSocket support
    
    async with websockets.connect(uri) as websocket:
        while True:
            # Receive test events
            data = await websocket.recv()
            test_event = json.loads(data)
            
            # Analyze with your LLM
            if test_event['type'] == 'test_failure':
                analysis = await analyze_with_llm(test_event)
                print(f"LLM Analysis: {analysis}")

async def analyze_with_llm(test_event):
    """Analyze test event with your preferred LLM"""
    # Your LLM analysis logic here
    pass
```

## Advanced Usage Patterns

### Batch Analysis

For analyzing large test suites:

```python
def batch_analyze_tests(project_ids, batch_size=50):
    """Analyze multiple projects in batches"""
    results = []
    
    for project_id in project_ids:
        query = f"""
        {{
            flakyTests(limit: {batch_size}, projectID: "{project_id}") {{
                testName
                passRate
                failureRate
                runCount
                lastFailure
            }}
        }}
        """
        
        response = requests.post(
            "http://localhost:8081/query",
            json={"query": query}
        )
        
        if response.status_code == 200:
            results.append({
                "project": project_id,
                "data": response.json()
            })
    
    return results
```

### Trend Analysis

For historical trend analysis:

```python
def analyze_test_trends(project_id, days=30):
    """Analyze test stability trends over time"""
    # This would require extending the GraphQL schema
    # to include time-based queries
    query = f"""
    {{
        testTrends(
            projectID: "{project_id}",
            days: {days}
        ) {{
            date
            totalTests
            flakyTests
            averageFailureRate
        }}
    }}
    """
    
    # Implementation would depend on your specific needs
    pass
```

## Security Considerations

### Authentication

Add authentication to your MCP server:

```go
// Add to your server setup
func authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if !validateToken(token) {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        next.ServeHTTP(w, r)
    })
}
```

### Rate Limiting

Implement rate limiting for production use:

```go
import "golang.org/x/time/rate"

func rateLimitMiddleware(limiter *rate.Limiter) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if !limiter.Allow() {
                http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}
```

## Troubleshooting

### Common Issues

1. **Connection Errors**
   ```bash
   # Check if server is running
   curl http://localhost:8081/healthz
   
   # Check database connectivity
   psql $DB_URL -c "SELECT 1"
   ```

2. **GraphQL Errors**
   ```bash
   # Validate GraphQL query syntax
   curl -X POST http://localhost:8081/query \
     -H "Content-Type: application/json" \
     -d '{"query": "{ __schema { types { name } } }"}'
   ```

3. **Performance Issues**
   ```bash
   # Monitor server performance
   curl http://localhost:8081/metrics  # If you add metrics endpoint
   ```

### Debugging

Enable debug logging:

```bash
export LOG_LEVEL=debug
go run main.go serve --port 8081
```

## Contributing

To extend the MCP server with additional LLM integrations:

1. Fork the repository
2. Add your integration in a new package
3. Update this documentation
4. Submit a pull request

## Support

For issues and questions:
- GitHub Issues: [Report bugs and feature requests](https://github.com/guidewire-oss/fern-mycelium/issues)
- Documentation: Check the `/docs` folder for additional guides
- Examples: See the `/examples` directory for more integration patterns