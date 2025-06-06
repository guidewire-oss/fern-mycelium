# OpenAI Integration with Fern-Mycelium

## Overview

This guide covers multiple ways to integrate the fern-mycelium MCP server with OpenAI's ecosystem, including GPT models, the Assistant API, custom GPTs, and function calling patterns.

## Prerequisites

- OpenAI API key
- Fern-mycelium MCP server running
- Python or Node.js environment for API integration

## Integration Methods

### 1. Custom GPT Integration

Create a custom GPT that can query your test intelligence data.

#### Step 1: Create Custom GPT

1. Go to [OpenAI GPT Builder](https://chat.openai.com/gpts/editor)
2. Click "Create a GPT"
3. Configure your GPT with these settings:

**Name:** Test Intelligence Assistant  
**Description:** Analyzes test execution patterns and provides recommendations using fern-mycelium data.

**Instructions:**
```
You are a test intelligence assistant with access to fern-mycelium test data. You can:

1. Analyze flaky test patterns
2. Provide test stability recommendations  
3. Identify high-risk tests that need immediate attention
4. Generate test health reports
5. Suggest test improvement strategies

When users ask about tests, always query the latest data from fern-mycelium and provide specific, actionable insights based on actual test execution metrics.

Be concise but thorough in your analysis, and always prioritize the most critical issues first.
```

#### Step 2: Configure Actions

Add an action to connect to your fern-mycelium server:

```yaml
openapi: 3.0.0
info:
  title: Fern-Mycelium Test Intelligence API
  description: API for querying test execution data and flaky test patterns
  version: 1.0.0
servers:
  - url: http://your-server:8081
    description: Fern-Mycelium MCP Server

paths:
  /query:
    post:
      operationId: queryTestData
      summary: Execute GraphQL query for test data
      description: Query test intelligence data using GraphQL
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              properties:
                query:
                  type: string
                  description: GraphQL query string
                variables:
                  type: object
                  description: GraphQL variables
              required:
                - query
            examples:
              flaky_tests:
                summary: Query flaky tests
                value:
                  query: |
                    {
                      flakyTests(limit: 10, projectID: "web-app") {
                        testName
                        passRate
                        failureRate
                        runCount
                        lastFailure
                      }
                    }
      responses:
        '200':
          description: Successful response
          content:
            application/json:
              schema:
                type: object
                properties:
                  data:
                    type: object
                  errors:
                    type: array
                    items:
                      type: object

  /healthz:
    get:
      operationId: healthCheck
      summary: Check server health
      description: Verify that the fern-mycelium server is running
      responses:
        '200':
          description: Server is healthy
          content:
            application/json:
              schema:
                type: object
                properties:
                  status:
                    type: string
                  message:
                    type: string
```

#### Step 3: Test Your Custom GPT

```
User: "What are the most problematic tests in our web-app project?"

GPT: I'll query your test intelligence data to identify the most problematic tests.

[GPT calls the queryTestData action with appropriate GraphQL query]

Based on the latest data from fern-mycelium, here are the most problematic tests:

1. **TestPaymentFlow** - 67% failure rate (critical)
   - 23 total runs, 15 failures
   - Last failure: 2 hours ago
   - Recommendation: Immediate investigation needed

2. **TestUserAuthentication** - 34% failure rate (high priority)
   - 45 total runs, 15 failures  
   - Last failure: 1 day ago
   - Recommendation: Review auth token handling

3. **TestDatabaseConnection** - 12% failure rate (monitor)
   - 58 total runs, 7 failures
   - Last failure: 3 days ago
   - Recommendation: Check connection pooling

Priority Actions:
1. Fix TestPaymentFlow immediately
2. Schedule auth system review
3. Monitor database connection stability
```

### 2. Assistant API Integration

Use OpenAI's Assistant API with function calling to create programmable test intelligence.

#### Python Implementation

```python
import openai
import requests
import json
from typing import Dict, Any, List

class TestIntelligenceAssistant:
    def __init__(self, api_key: str, server_url: str = "http://localhost:8081"):
        self.client = openai.OpenAI(api_key=api_key)
        self.server_url = server_url
        self.assistant = None
        self._create_assistant()
    
    def _create_assistant(self):
        """Create the test intelligence assistant"""
        self.assistant = self.client.beta.assistants.create(
            name="Test Intelligence Assistant",
            instructions="""
            You are an expert test analysis assistant. Use the available functions to:
            1. Query test data from fern-mycelium
            2. Analyze test patterns and trends
            3. Provide actionable recommendations
            4. Generate test health reports
            
            Always provide specific, data-driven insights based on actual test execution metrics.
            Focus on the most critical issues first and provide clear action items.
            """,
            model="gpt-4-turbo-preview",
            tools=[
                {
                    "type": "function",
                    "function": {
                        "name": "query_flaky_tests",
                        "description": "Query flaky test data for a specific project",
                        "parameters": {
                            "type": "object",
                            "properties": {
                                "project_id": {
                                    "type": "string",
                                    "description": "The project ID to query tests for"
                                },
                                "limit": {
                                    "type": "integer",
                                    "description": "Maximum number of tests to return",
                                    "default": 10
                                }
                            },
                            "required": ["project_id"]
                        }
                    }
                },
                {
                    "type": "function", 
                    "function": {
                        "name": "check_server_health",
                        "description": "Check if the test intelligence server is healthy",
                        "parameters": {
                            "type": "object",
                            "properties": {}
                        }
                    }
                }
            ]
        )
    
    def query_flaky_tests(self, project_id: str, limit: int = 10) -> Dict[str, Any]:
        """Query flaky tests from fern-mycelium"""
        query = f"""
        {{
            flakyTests(limit: {limit}, projectID: "{project_id}") {{
                testName
                passRate
                failureRate
                runCount
                lastFailure
            }}
        }}
        """
        
        response = requests.post(
            f"{self.server_url}/query",
            json={"query": query},
            headers={"Content-Type": "application/json"}
        )
        
        if response.status_code == 200:
            return response.json()
        else:
            return {"error": f"HTTP {response.status_code}: {response.text}"}
    
    def check_server_health(self) -> Dict[str, Any]:
        """Check server health"""
        try:
            response = requests.get(f"{self.server_url}/healthz")
            return response.json() if response.status_code == 200 else {"error": "Server unhealthy"}
        except Exception as e:
            return {"error": str(e)}
    
    def handle_function_call(self, function_name: str, arguments: str) -> str:
        """Handle function calls from the assistant"""
        args = json.loads(arguments)
        
        if function_name == "query_flaky_tests":
            result = self.query_flaky_tests(**args)
            return json.dumps(result)
        elif function_name == "check_server_health":
            result = self.check_server_health()
            return json.dumps(result)
        else:
            return json.dumps({"error": f"Unknown function: {function_name}"})
    
    def chat(self, message: str) -> str:
        """Chat with the assistant"""
        # Create a thread
        thread = self.client.beta.threads.create()
        
        # Add user message
        self.client.beta.threads.messages.create(
            thread_id=thread.id,
            role="user",
            content=message
        )
        
        # Run the assistant
        run = self.client.beta.threads.runs.create(
            thread_id=thread.id,
            assistant_id=self.assistant.id
        )
        
        # Wait for completion and handle function calls
        while run.status in ["queued", "in_progress", "requires_action"]:
            if run.status == "requires_action":
                # Handle function calls
                tool_calls = run.required_action.submit_tool_outputs.tool_calls
                tool_outputs = []
                
                for tool_call in tool_calls:
                    function_name = tool_call.function.name
                    arguments = tool_call.function.arguments
                    output = self.handle_function_call(function_name, arguments)
                    
                    tool_outputs.append({
                        "tool_call_id": tool_call.id,
                        "output": output
                    })
                
                # Submit tool outputs
                run = self.client.beta.threads.runs.submit_tool_outputs(
                    thread_id=thread.id,
                    run_id=run.id,
                    tool_outputs=tool_outputs
                )
            else:
                # Wait a bit before checking again
                import time
                time.sleep(1)
                run = self.client.beta.threads.runs.retrieve(
                    thread_id=thread.id,
                    run_id=run.id
                )
        
        # Get the assistant's response
        messages = self.client.beta.threads.messages.list(thread_id=thread.id)
        return messages.data[0].content[0].text.value

# Usage example
def main():
    assistant = TestIntelligenceAssistant(
        api_key="your-openai-api-key",
        server_url="http://localhost:8081"
    )
    
    # Example conversations
    response = assistant.chat("What are the most problematic tests in our web-app project?")
    print(response)
    
    response = assistant.chat("Generate a test health report for all our projects")
    print(response)
    
    response = assistant.chat("Is the test intelligence server healthy?")
    print(response)

if __name__ == "__main__":
    main()
```

#### Node.js Implementation

```javascript
const OpenAI = require('openai');
const axios = require('axios');

class TestIntelligenceAssistant {
    constructor(apiKey, serverUrl = 'http://localhost:8081') {
        this.client = new OpenAI({ apiKey });
        this.serverUrl = serverUrl;
        this.assistant = null;
        this.init();
    }

    async init() {
        this.assistant = await this.client.beta.assistants.create({
            name: "Test Intelligence Assistant",
            instructions: `
                You are an expert test analysis assistant. Use the available functions to:
                1. Query test data from fern-mycelium
                2. Analyze test patterns and trends  
                3. Provide actionable recommendations
                4. Generate test health reports
            `,
            model: "gpt-4-turbo-preview",
            tools: [
                {
                    type: "function",
                    function: {
                        name: "query_flaky_tests",
                        description: "Query flaky test data for a specific project",
                        parameters: {
                            type: "object",
                            properties: {
                                project_id: { type: "string" },
                                limit: { type: "integer", default: 10 }
                            },
                            required: ["project_id"]
                        }
                    }
                }
            ]
        });
    }

    async queryFlakyTests(projectId, limit = 10) {
        const query = `
            {
                flakyTests(limit: ${limit}, projectID: "${projectId}") {
                    testName
                    passRate
                    failureRate
                    runCount
                    lastFailure
                }
            }
        `;

        try {
            const response = await axios.post(`${this.serverUrl}/query`, {
                query
            });
            return response.data;
        } catch (error) {
            return { error: error.message };
        }
    }

    async handleFunctionCall(functionName, args) {
        switch (functionName) {
            case 'query_flaky_tests':
                return await this.queryFlakyTests(args.project_id, args.limit);
            default:
                return { error: `Unknown function: ${functionName}` };
        }
    }

    async chat(message) {
        const thread = await this.client.beta.threads.create();
        
        await this.client.beta.threads.messages.create(
            thread.id,
            { role: "user", content: message }
        );

        let run = await this.client.beta.threads.runs.create(
            thread.id,
            { assistant_id: this.assistant.id }
        );

        while (['queued', 'in_progress', 'requires_action'].includes(run.status)) {
            if (run.status === 'requires_action') {
                const toolCalls = run.required_action.submit_tool_outputs.tool_calls;
                const toolOutputs = [];

                for (const toolCall of toolCalls) {
                    const args = JSON.parse(toolCall.function.arguments);
                    const output = await this.handleFunctionCall(
                        toolCall.function.name, 
                        args
                    );
                    
                    toolOutputs.push({
                        tool_call_id: toolCall.id,
                        output: JSON.stringify(output)
                    });
                }

                run = await this.client.beta.threads.runs.submitToolOutputs(
                    thread.id,
                    run.id,
                    { tool_outputs: toolOutputs }
                );
            } else {
                await new Promise(resolve => setTimeout(resolve, 1000));
                run = await this.client.beta.threads.runs.retrieve(thread.id, run.id);
            }
        }

        const messages = await this.client.beta.threads.messages.list(thread.id);
        return messages.data[0].content[0].text.value;
    }
}

// Usage
async function main() {
    const assistant = new TestIntelligenceAssistant('your-openai-api-key');
    
    const response = await assistant.chat(
        "Analyze test stability for the authentication project"
    );
    console.log(response);
}

main().catch(console.error);
```

### 3. Function Calling with ChatGPT

Direct integration using function calling in ChatGPT conversations.

#### Example Implementation

```python
import openai
import requests
import json

def get_flaky_tests(project_id: str, limit: int = 10):
    """Function to get flaky tests from fern-mycelium"""
    query = f"""
    {{
        flakyTests(limit: {limit}, projectID: "{project_id}") {{
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
    
    return response.json() if response.status_code == 200 else {"error": "Query failed"}

def chat_with_test_intelligence(user_message: str, api_key: str):
    """Chat with GPT using test intelligence functions"""
    client = openai.OpenAI(api_key=api_key)
    
    messages = [
        {
            "role": "system", 
            "content": "You are a test intelligence assistant. Use the available functions to query test data and provide insights."
        },
        {"role": "user", "content": user_message}
    ]
    
    tools = [
        {
            "type": "function",
            "function": {
                "name": "get_flaky_tests",
                "description": "Get flaky test data for a project",
                "parameters": {
                    "type": "object",
                    "properties": {
                        "project_id": {"type": "string", "description": "Project ID"},
                        "limit": {"type": "integer", "description": "Max results", "default": 10}
                    },
                    "required": ["project_id"]
                }
            }
        }
    ]
    
    response = client.chat.completions.create(
        model="gpt-4-turbo-preview",
        messages=messages,
        tools=tools,
        tool_choice="auto"
    )
    
    # Handle function calls
    if response.choices[0].message.tool_calls:
        for tool_call in response.choices[0].message.tool_calls:
            if tool_call.function.name == "get_flaky_tests":
                args = json.loads(tool_call.function.arguments)
                result = get_flaky_tests(**args)
                
                # Add function result to conversation
                messages.append(response.choices[0].message)
                messages.append({
                    "role": "tool",
                    "tool_call_id": tool_call.id,
                    "content": json.dumps(result)
                })
                
                # Get final response
                final_response = client.chat.completions.create(
                    model="gpt-4-turbo-preview",
                    messages=messages
                )
                
                return final_response.choices[0].message.content
    
    return response.choices[0].message.content

# Usage
result = chat_with_test_intelligence(
    "What are the most problematic tests in project 'web-app'?",
    "your-openai-api-key"
)
print(result)
```

### 4. Batch Processing

For analyzing multiple projects or large datasets:

```python
import asyncio
import aiohttp
import openai
from typing import List, Dict

class BatchTestAnalyzer:
    def __init__(self, api_key: str, server_url: str = "http://localhost:8081"):
        self.client = openai.OpenAI(api_key=api_key)
        self.server_url = server_url
    
    async def fetch_project_data(self, session: aiohttp.ClientSession, project_id: str):
        """Fetch test data for a single project"""
        query = f"""
        {{
            flakyTests(limit: 50, projectID: "{project_id}") {{
                testName
                passRate
                failureRate
                runCount
                lastFailure
            }}
        }}
        """
        
        async with session.post(
            f"{self.server_url}/query",
            json={"query": query}
        ) as response:
            data = await response.json()
            return {"project": project_id, "data": data}
    
    async def analyze_all_projects(self, project_ids: List[str]) -> Dict[str, str]:
        """Analyze test data for multiple projects"""
        async with aiohttp.ClientSession() as session:
            # Fetch data for all projects concurrently
            tasks = [
                self.fetch_project_data(session, project_id) 
                for project_id in project_ids
            ]
            results = await asyncio.gather(*tasks)
        
        # Analyze with GPT
        analysis_prompt = f"""
        Analyze the following test data for multiple projects and provide:
        1. Overall health assessment
        2. Most critical issues across all projects
        3. Recommendations for each project
        4. Resource allocation priorities
        
        Data: {json.dumps(results, indent=2)}
        """
        
        response = self.client.chat.completions.create(
            model="gpt-4-turbo-preview",
            messages=[{
                "role": "user",
                "content": analysis_prompt
            }],
            max_tokens=2000
        )
        
        return response.choices[0].message.content

# Usage
async def main():
    analyzer = BatchTestAnalyzer("your-openai-api-key")
    projects = ["web-app", "mobile-app", "api-service", "auth-service"]
    
    analysis = await analyzer.analyze_all_projects(projects)
    print(analysis)

asyncio.run(main())
```

## Best Practices

### 1. Error Handling

Always implement robust error handling:

```python
def safe_query_tests(project_id: str):
    """Safely query test data with error handling"""
    try:
        response = requests.post(
            "http://localhost:8081/query",
            json={"query": f'{{ flakyTests(projectID: "{project_id}") {{ testName }} }}'},
            timeout=30
        )
        response.raise_for_status()
        return response.json()
    except requests.exceptions.ConnectionError:
        return {"error": "Cannot connect to fern-mycelium server"}
    except requests.exceptions.Timeout:
        return {"error": "Query timeout - server may be overloaded"}
    except requests.exceptions.HTTPError as e:
        return {"error": f"HTTP error: {e.response.status_code}"}
    except Exception as e:
        return {"error": f"Unexpected error: {str(e)}"}
```

### 2. Rate Limiting

Implement rate limiting for production use:

```python
import time
from functools import wraps

def rate_limit(calls_per_minute: int):
    """Rate limiting decorator"""
    min_interval = 60.0 / calls_per_minute
    last_called = [0.0]
    
    def decorator(func):
        @wraps(func)
        def wrapper(*args, **kwargs):
            elapsed = time.time() - last_called[0]
            left_to_wait = min_interval - elapsed
            if left_to_wait > 0:
                time.sleep(left_to_wait)
            ret = func(*args, **kwargs)
            last_called[0] = time.time()
            return ret
        return wrapper
    return decorator

@rate_limit(30)  # 30 calls per minute
def query_test_data(project_id: str):
    # Your query implementation
    pass
```

### 3. Caching

Implement caching for better performance:

```python
import time
from typing import Optional

class CachedTestClient:
    def __init__(self, cache_ttl: int = 300):  # 5 minutes
        self.cache = {}
        self.cache_ttl = cache_ttl
    
    def get_cached_or_fetch(self, project_id: str) -> Optional[dict]:
        """Get data from cache or fetch from server"""
        cache_key = f"flaky_tests_{project_id}"
        current_time = time.time()
        
        if cache_key in self.cache:
            data, timestamp = self.cache[cache_key]
            if current_time - timestamp < self.cache_ttl:
                return data
        
        # Fetch fresh data
        data = self.fetch_test_data(project_id)
        self.cache[cache_key] = (data, current_time)
        return data
    
    def fetch_test_data(self, project_id: str) -> dict:
        # Your fetch implementation
        pass
```

## Monitoring and Observability

Track usage and performance:

```python
import logging
from datetime import datetime

class MonitoredTestClient:
    def __init__(self):
        self.logger = logging.getLogger(__name__)
        self.query_count = 0
        self.error_count = 0
    
    def log_query(self, project_id: str, success: bool, duration: float):
        """Log query metrics"""
        self.query_count += 1
        if not success:
            self.error_count += 1
        
        self.logger.info(f"Query: project={project_id}, success={success}, duration={duration:.2f}s")
    
    def get_metrics(self) -> dict:
        """Get client metrics"""
        return {
            "total_queries": self.query_count,
            "error_count": self.error_count,
            "error_rate": self.error_count / max(self.query_count, 1),
            "timestamp": datetime.now().isoformat()
        }
```

## Troubleshooting

### Common Issues

1. **Function calling not working**
   - Verify OpenAI API key has GPT-4 access
   - Check function schema matches exactly
   - Ensure proper JSON formatting in responses

2. **Server connection issues**
   - Verify fern-mycelium server is running
   - Check network connectivity and firewall rules
   - Test with curl before integrating

3. **Rate limiting**
   - Implement exponential backoff
   - Use async processing for batch operations
   - Cache frequently accessed data

4. **Large dataset handling**
   - Paginate GraphQL queries
   - Use streaming responses for large results
   - Implement proper timeout handling

For additional support, check the main [MCP Integration Guide](MCP_INTEGRATION.md) and [Usage Examples](USAGE_EXAMPLES.md).