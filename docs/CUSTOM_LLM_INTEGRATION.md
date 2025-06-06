# Custom LLM Integration with Fern-Mycelium

## Overview

This guide covers how to integrate fern-mycelium with custom LLMs, self-hosted models, and various AI frameworks beyond OpenAI and Claude. Whether you're using local models, enterprise AI platforms, or building custom agents, this documentation provides integration patterns and examples.

## Integration Architectures

### 1. Direct HTTP API Integration

The simplest integration pattern for any LLM that can make HTTP requests.

#### Basic HTTP Client Pattern

```python
import requests
import json
from typing import Dict, Any, Optional

class FernMyceliumClient:
    def __init__(self, server_url: str = "http://localhost:8081"):
        self.server_url = server_url
        self.session = requests.Session()
        self.session.headers.update({
            'Content-Type': 'application/json',
            'User-Agent': 'FernMycelium-Client/1.0'
        })
    
    def health_check(self) -> Dict[str, Any]:
        """Check server health"""
        try:
            response = self.session.get(f"{self.server_url}/healthz", timeout=10)
            response.raise_for_status()
            return response.json()
        except Exception as e:
            return {"status": "error", "message": str(e)}
    
    def query_flaky_tests(self, project_id: str, limit: int = 10) -> Dict[str, Any]:
        """Query flaky tests for a project"""
        query = f"""
        {{
            flakyTests(limit: {limit}, projectID: "{project_id}") {{
                testID
                testName
                passRate
                failureRate
                runCount
                lastFailure
            }}
        }}
        """
        
        payload = {"query": query}
        
        try:
            response = self.session.post(
                f"{self.server_url}/query",
                json=payload,
                timeout=30
            )
            response.raise_for_status()
            return response.json()
        except Exception as e:
            return {"errors": [{"message": str(e)}]}
    
    def get_test_summary(self, project_id: str) -> str:
        """Get a formatted test summary for LLM consumption"""
        data = self.query_flaky_tests(project_id, limit=50)
        
        if "errors" in data:
            return f"Error fetching test data: {data['errors'][0]['message']}"
        
        if not data.get("data", {}).get("flakyTests"):
            return f"No test data found for project '{project_id}'"
        
        tests = data["data"]["flakyTests"]
        
        # Format summary for LLM
        summary = f"Test Intelligence Summary for '{project_id}':\n\n"
        summary += f"Total tests analyzed: {len(tests)}\n\n"
        
        # Categorize tests
        high_risk = [t for t in tests if t["failureRate"] > 0.3]
        moderate_risk = [t for t in tests if 0 < t["failureRate"] <= 0.3]
        stable = [t for t in tests if t["failureRate"] == 0]
        
        summary += f"🔴 High Risk Tests (>30% failure): {len(high_risk)}\n"
        summary += f"🟡 Moderate Risk Tests (1-30% failure): {len(moderate_risk)}\n"
        summary += f"🟢 Stable Tests (0% failure): {len(stable)}\n\n"
        
        if high_risk:
            summary += "CRITICAL ISSUES:\n"
            for test in high_risk[:5]:  # Top 5 most critical
                summary += f"- {test['testName']}: {test['failureRate']:.1%} failure rate\n"
            summary += "\n"
        
        return summary

# Usage with any LLM that can execute Python
client = FernMyceliumClient()
test_summary = client.get_test_summary("web-app")
print(test_summary)
```

### 2. Hugging Face Integration

For models hosted on Hugging Face or using the Transformers library.

#### Local Model Integration

```python
from transformers import AutoTokenizer, AutoModelForCausalLM, pipeline
import torch
from typing import List, Dict

class HuggingFaceTestAgent:
    def __init__(self, model_name: str = "microsoft/DialoGPT-large"):
        self.fern_client = FernMyceliumClient()
        self.model_name = model_name
        self.tokenizer = AutoTokenizer.from_pretrained(model_name)
        self.model = AutoModelForCausalLM.from_pretrained(model_name)
        
        # Add padding token if not present
        if self.tokenizer.pad_token is None:
            self.tokenizer.pad_token = self.tokenizer.eos_token
        
        # Create text generation pipeline
        self.generator = pipeline(
            "text-generation",
            model=self.model,
            tokenizer=self.tokenizer,
            device=0 if torch.cuda.is_available() else -1
        )
    
    def analyze_tests(self, project_id: str, user_query: str) -> str:
        """Analyze tests using local model"""
        # Get test data
        test_summary = self.fern_client.get_test_summary(project_id)
        
        # Create prompt for the model
        prompt = f"""
        Test Intelligence Analysis Request:
        
        User Query: {user_query}
        
        Available Test Data:
        {test_summary}
        
        Please provide analysis and recommendations:
        """
        
        # Generate response
        response = self.generator(
            prompt,
            max_length=len(prompt.split()) + 200,
            num_return_sequences=1,
            temperature=0.7,
            do_sample=True,
            pad_token_id=self.tokenizer.eos_token_id
        )
        
        # Extract generated text
        generated_text = response[0]["generated_text"]
        analysis = generated_text[len(prompt):].strip()
        
        return analysis

# Usage
agent = HuggingFaceTestAgent("microsoft/DialoGPT-medium")
result = agent.analyze_tests("web-app", "What are the most critical test failures?")
print(result)
```

#### Hugging Face Inference API

```python
import requests
from typing import Dict, Any

class HuggingFaceAPIAgent:
    def __init__(self, api_token: str, model_id: str = "microsoft/DialoGPT-large"):
        self.api_token = api_token
        self.model_id = model_id
        self.api_url = f"https://api-inference.huggingface.co/models/{model_id}"
        self.fern_client = FernMyceliumClient()
        
        self.headers = {
            "Authorization": f"Bearer {api_token}",
            "Content-Type": "application/json"
        }
    
    def query_model(self, prompt: str) -> str:
        """Query Hugging Face model via API"""
        payload = {
            "inputs": prompt,
            "parameters": {
                "max_new_tokens": 200,
                "temperature": 0.7,
                "do_sample": True
            }
        }
        
        response = requests.post(
            self.api_url,
            headers=self.headers,
            json=payload
        )
        
        if response.status_code == 200:
            result = response.json()
            if isinstance(result, list) and len(result) > 0:
                return result[0].get("generated_text", "")
        
        return f"Error: {response.status_code} - {response.text}"
    
    def analyze_test_trends(self, project_id: str) -> str:
        """Analyze test trends using HF model"""
        test_data = self.fern_client.get_test_summary(project_id)
        
        prompt = f"""Analyze these test results and provide recommendations:

{test_data}

Analysis:"""
        
        return self.query_model(prompt)

# Usage
agent = HuggingFaceAPIAgent("your-hf-token")
analysis = agent.analyze_test_trends("mobile-app")
print(analysis)
```

### 3. Ollama Integration

For running local models with Ollama.

```python
import requests
import json
from typing import Dict, Any, Generator

class OllamaTestAgent:
    def __init__(self, base_url: str = "http://localhost:11434", model: str = "llama2"):
        self.base_url = base_url
        self.model = model
        self.fern_client = FernMyceliumClient()
    
    def generate_stream(self, prompt: str) -> Generator[str, None, None]:
        """Generate streaming response from Ollama"""
        payload = {
            "model": self.model,
            "prompt": prompt,
            "stream": True
        }
        
        response = requests.post(
            f"{self.base_url}/api/generate",
            json=payload,
            stream=True
        )
        
        for line in response.iter_lines():
            if line:
                try:
                    data = json.loads(line)
                    if "response" in data:
                        yield data["response"]
                except json.JSONDecodeError:
                    continue
    
    def generate(self, prompt: str) -> str:
        """Generate complete response from Ollama"""
        payload = {
            "model": self.model,
            "prompt": prompt,
            "stream": False
        }
        
        response = requests.post(
            f"{self.base_url}/api/generate",
            json=payload
        )
        
        if response.status_code == 200:
            return response.json().get("response", "")
        return f"Error: {response.status_code}"
    
    def analyze_tests_interactive(self, project_id: str, user_query: str) -> str:
        """Interactive test analysis with streaming"""
        test_data = self.fern_client.get_test_summary(project_id)
        
        prompt = f"""You are a test intelligence assistant. Analyze the following test data and answer the user's question.

Test Data:
{test_data}

User Question: {user_query}

Provide specific, actionable insights based on the test data. Focus on:
1. Identifying critical issues
2. Prioritizing fixes
3. Suggesting concrete improvements

Analysis:"""
        
        print("Generating analysis...")
        result = ""
        for chunk in self.generate_stream(prompt):
            print(chunk, end="", flush=True)
            result += chunk
        print("\n")
        
        return result
    
    def batch_analyze_projects(self, project_ids: List[str]) -> Dict[str, str]:
        """Analyze multiple projects"""
        results = {}
        
        for project_id in project_ids:
            print(f"Analyzing project: {project_id}")
            test_data = self.fern_client.get_test_summary(project_id)
            
            prompt = f"""Provide a brief health assessment for this project's tests:

{test_data}

Assessment (2-3 sentences):"""
            
            results[project_id] = self.generate(prompt)
        
        return results

# Usage
agent = OllamaTestAgent(model="llama2:13b")
analysis = agent.analyze_tests_interactive("web-app", "What tests need immediate attention?")
```

### 4. Azure OpenAI Integration

For enterprise Azure OpenAI deployments.

```python
import openai
from typing import Dict, Any, List

class AzureOpenAITestAgent:
    def __init__(self, endpoint: str, api_key: str, deployment_name: str, api_version: str = "2023-12-01-preview"):
        # Configure Azure OpenAI
        openai.api_type = "azure"
        openai.api_base = endpoint
        openai.api_key = api_key
        openai.api_version = api_version
        
        self.deployment_name = deployment_name
        self.fern_client = FernMyceliumClient()
    
    def analyze_with_functions(self, user_query: str) -> str:
        """Analyze using Azure OpenAI with function calling"""
        functions = [
            {
                "name": "get_test_intelligence",
                "description": "Get test intelligence data for a project",
                "parameters": {
                    "type": "object",
                    "properties": {
                        "project_id": {
                            "type": "string",
                            "description": "The project ID to analyze"
                        },
                        "limit": {
                            "type": "integer",
                            "description": "Maximum number of tests to analyze",
                            "default": 20
                        }
                    },
                    "required": ["project_id"]
                }
            }
        ]
        
        messages = [
            {
                "role": "system",
                "content": "You are a test intelligence assistant. Use the get_test_intelligence function to fetch data and provide insights."
            },
            {
                "role": "user",
                "content": user_query
            }
        ]
        
        response = openai.ChatCompletion.create(
            engine=self.deployment_name,
            messages=messages,
            functions=functions,
            function_call="auto",
            temperature=0.3
        )
        
        message = response.choices[0].message
        
        if message.get("function_call"):
            # Handle function call
            function_name = message["function_call"]["name"]
            function_args = json.loads(message["function_call"]["arguments"])
            
            if function_name == "get_test_intelligence":
                test_data = self.fern_client.query_flaky_tests(**function_args)
                
                # Add function result to conversation
                messages.append(message)
                messages.append({
                    "role": "function",
                    "name": function_name,
                    "content": json.dumps(test_data)
                })
                
                # Get final response
                final_response = openai.ChatCompletion.create(
                    engine=self.deployment_name,
                    messages=messages,
                    temperature=0.3
                )
                
                return final_response.choices[0].message.content
        
        return message.content

# Usage
agent = AzureOpenAITestAgent(
    endpoint="https://your-resource.openai.azure.com/",
    api_key="your-api-key",
    deployment_name="gpt-4"
)

result = agent.analyze_with_functions("Analyze flaky tests for project 'payment-service'")
print(result)
```

### 5. Google Vertex AI Integration

For Google Cloud AI Platform integration.

```python
from google.cloud import aiplatform
from google.oauth2 import service_account
import json
from typing import Dict, Any

class VertexAITestAgent:
    def __init__(self, project_id: str, location: str, credentials_path: str = None):
        if credentials_path:
            credentials = service_account.Credentials.from_service_account_file(credentials_path)
            aiplatform.init(project=project_id, location=location, credentials=credentials)
        else:
            aiplatform.init(project=project_id, location=location)
        
        self.project_id = project_id
        self.location = location
        self.fern_client = FernMyceliumClient()
    
    def analyze_with_palm(self, project_id: str, user_query: str, model_name: str = "text-bison") -> str:
        """Analyze tests using PaLM model"""
        # Get test data
        test_data = self.fern_client.get_test_summary(project_id)
        
        # Create prompt
        prompt = f"""
        You are a test intelligence expert. Analyze the following test data and answer the user's question.
        
        Test Intelligence Data:
        {test_data}
        
        User Question: {user_query}
        
        Provide specific recommendations based on the test failure patterns and rates.
        
        Analysis:
        """
        
        # Generate with Vertex AI
        model = aiplatform.gapic.PredictionServiceClient(
            client_options={"api_endpoint": f"{self.location}-aiplatform.googleapis.com"}
        )
        
        endpoint = f"projects/{self.project_id}/locations/{self.location}/publishers/google/models/{model_name}"
        
        instance = {
            "prompt": prompt,
            "max_output_tokens": 512,
            "temperature": 0.3,
            "top_p": 0.8,
            "top_k": 40
        }
        
        response = model.predict(
            endpoint=endpoint,
            instances=[instance]
        )
        
        if response.predictions:
            return response.predictions[0]["content"]
        
        return "No response generated"
    
    def batch_analyze_with_gemini(self, project_ids: List[str]) -> Dict[str, str]:
        """Batch analyze using Gemini Pro"""
        results = {}
        
        # Collect all test data
        all_data = {}
        for project_id in project_ids:
            all_data[project_id] = self.fern_client.get_test_summary(project_id)
        
        # Create comprehensive prompt
        prompt = f"""
        Analyze test intelligence data for multiple projects and provide:
        1. Overall health assessment
        2. Priority ranking of issues
        3. Resource allocation recommendations
        
        Project Data:
        {json.dumps(all_data, indent=2)}
        
        Provide analysis for each project and overall recommendations:
        """
        
        # Use Gemini Pro model
        model = aiplatform.gapic.PredictionServiceClient(
            client_options={"api_endpoint": f"{self.location}-aiplatform.googleapis.com"}
        )
        
        endpoint = f"projects/{self.project_id}/locations/{self.location}/publishers/google/models/gemini-pro"
        
        instance = {
            "prompt": prompt,
            "max_output_tokens": 1024,
            "temperature": 0.2
        }
        
        response = model.predict(
            endpoint=endpoint,
            instances=[instance]
        )
        
        if response.predictions:
            return {"analysis": response.predictions[0]["content"]}
        
        return {"error": "No analysis generated"}

# Usage
agent = VertexAITestAgent(
    project_id="your-gcp-project",
    location="us-central1",
    credentials_path="/path/to/service-account.json"
)

analysis = agent.analyze_with_palm("web-app", "What are the riskiest tests?")
print(analysis)
```

### 6. Custom Agent Framework

Build a generic framework that works with any LLM.

```python
from abc import ABC, abstractmethod
from typing import Dict, Any, List, Optional
import logging

class LLMProvider(ABC):
    """Abstract base class for LLM providers"""
    
    @abstractmethod
    def generate(self, prompt: str, **kwargs) -> str:
        """Generate text using the LLM"""
        pass
    
    @abstractmethod
    def supports_functions(self) -> bool:
        """Check if the LLM supports function calling"""
        pass

class TestIntelligenceAgent:
    """Generic test intelligence agent that works with any LLM"""
    
    def __init__(self, llm_provider: LLMProvider, fern_url: str = "http://localhost:8081"):
        self.llm = llm_provider
        self.fern_client = FernMyceliumClient(fern_url)
        self.logger = logging.getLogger(__name__)
    
    def analyze_project(self, project_id: str, focus_area: str = "general") -> Dict[str, Any]:
        """Analyze a project with specific focus"""
        # Get test data
        test_data = self.fern_client.get_test_summary(project_id)
        
        # Create focused prompt based on area
        prompts = {
            "general": f"Provide a general health assessment of these test results:\n{test_data}",
            "critical": f"Identify the most critical test failures that need immediate attention:\n{test_data}",
            "trends": f"Analyze trends and patterns in these test results:\n{test_data}",
            "resources": f"Recommend resource allocation based on these test patterns:\n{test_data}"
        }
        
        prompt = prompts.get(focus_area, prompts["general"])
        
        try:
            analysis = self.llm.generate(prompt, temperature=0.3, max_tokens=500)
            
            return {
                "project_id": project_id,
                "focus_area": focus_area,
                "analysis": analysis,
                "raw_data": test_data,
                "status": "success"
            }
        except Exception as e:
            self.logger.error(f"Analysis failed for {project_id}: {e}")
            return {
                "project_id": project_id,
                "status": "error",
                "error": str(e)
            }
    
    def compare_projects(self, project_ids: List[str]) -> str:
        """Compare test health across multiple projects"""
        project_data = {}
        
        for project_id in project_ids:
            project_data[project_id] = self.fern_client.get_test_summary(project_id)
        
        comparison_prompt = f"""
        Compare the test health across these projects and rank them by priority:
        
        {json.dumps(project_data, indent=2)}
        
        Provide:
        1. Ranking from most to least problematic
        2. Key differences between projects
        3. Resource allocation recommendations
        """
        
        return self.llm.generate(comparison_prompt, temperature=0.2, max_tokens=800)
    
    def generate_report(self, project_id: str, report_type: str = "executive") -> str:
        """Generate different types of reports"""
        test_data = self.fern_client.get_test_summary(project_id)
        
        report_prompts = {
            "executive": f"""
            Create an executive summary of test health for project '{project_id}':
            
            {test_data}
            
            Format as a brief executive report with key metrics and recommendations.
            """,
            "technical": f"""
            Create a detailed technical report for project '{project_id}':
            
            {test_data}
            
            Include specific test names, failure rates, and technical recommendations.
            """,
            "dashboard": f"""
            Create dashboard-style metrics for project '{project_id}':
            
            {test_data}
            
            Format as key metrics suitable for a monitoring dashboard.
            """
        }
        
        prompt = report_prompts.get(report_type, report_prompts["executive"])
        return self.llm.generate(prompt, temperature=0.1, max_tokens=1000)

# Example LLM Provider implementations
class OpenAIProvider(LLMProvider):
    def __init__(self, api_key: str, model: str = "gpt-3.5-turbo"):
        import openai
        self.client = openai.OpenAI(api_key=api_key)
        self.model = model
    
    def generate(self, prompt: str, **kwargs) -> str:
        response = self.client.chat.completions.create(
            model=self.model,
            messages=[{"role": "user", "content": prompt}],
            **kwargs
        )
        return response.choices[0].message.content
    
    def supports_functions(self) -> bool:
        return True

class OllamaProvider(LLMProvider):
    def __init__(self, model: str = "llama2", base_url: str = "http://localhost:11434"):
        self.model = model
        self.base_url = base_url
    
    def generate(self, prompt: str, **kwargs) -> str:
        import requests
        payload = {
            "model": self.model,
            "prompt": prompt,
            "stream": False,
            **kwargs
        }
        
        response = requests.post(f"{self.base_url}/api/generate", json=payload)
        return response.json().get("response", "")
    
    def supports_functions(self) -> bool:
        return False

# Usage examples
def main():
    # Using with OpenAI
    openai_provider = OpenAIProvider("your-api-key", "gpt-4")
    agent = TestIntelligenceAgent(openai_provider)
    
    analysis = agent.analyze_project("web-app", "critical")
    print(analysis)
    
    # Using with Ollama
    ollama_provider = OllamaProvider("llama2:13b")
    local_agent = TestIntelligenceAgent(ollama_provider)
    
    report = local_agent.generate_report("mobile-app", "technical")
    print(report)
    
    # Comparing projects
    comparison = agent.compare_projects(["web-app", "mobile-app", "api-service"])
    print(comparison)

if __name__ == "__main__":
    main()
```

## WebSocket Integration

For real-time test intelligence streaming:

```python
import asyncio
import websockets
import json
from typing import AsyncGenerator

class WebSocketTestAgent:
    def __init__(self, fern_ws_url: str = "ws://localhost:8081/ws"):
        self.fern_ws_url = fern_ws_url
        self.llm_provider = None  # Your LLM provider
    
    async def real_time_analysis(self) -> AsyncGenerator[str, None]:
        """Stream real-time test analysis"""
        async with websockets.connect(self.fern_ws_url) as websocket:
            async for message in websocket:
                try:
                    data = json.loads(message)
                    
                    if data.get("type") == "test_result":
                        # Analyze new test result in real-time
                        analysis = await self.analyze_test_event(data)
                        yield analysis
                        
                except json.JSONDecodeError:
                    continue
    
    async def analyze_test_event(self, event_data: Dict[str, Any]) -> str:
        """Analyze a single test event"""
        if not event_data.get("test_failure"):
            return ""
        
        prompt = f"""
        A test just failed. Analyze this failure:
        
        Test: {event_data.get('test_name')}
        Project: {event_data.get('project_id')}
        Failure Type: {event_data.get('failure_type')}
        Error Message: {event_data.get('error_message')}
        
        Provide immediate recommendations:
        """
        
        # Use your LLM to analyze
        return await self.llm_provider.generate_async(prompt)

# Usage
async def main():
    agent = WebSocketTestAgent()
    
    async for analysis in agent.real_time_analysis():
        print(f"Real-time analysis: {analysis}")

asyncio.run(main())
```

## Deployment Patterns

### Container-based Deployment

```dockerfile
# Dockerfile for custom LLM agent
FROM python:3.11-slim

WORKDIR /app

COPY requirements.txt .
RUN pip install -r requirements.txt

COPY agent.py .
COPY config/ ./config/

ENV FERN_MYCELIUM_URL=http://fern-mycelium:8081
ENV LLM_PROVIDER=ollama
ENV MODEL_NAME=llama2:7b

CMD ["python", "agent.py"]
```

### Kubernetes Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-intelligence-agent
spec:
  replicas: 1
  selector:
    matchLabels:
      app: test-intelligence-agent
  template:
    metadata:
      labels:
        app: test-intelligence-agent
    spec:
      containers:
      - name: agent
        image: your-registry/test-intelligence-agent:latest
        env:
        - name: FERN_MYCELIUM_URL
          value: "http://fern-mycelium-service:8081"
        - name: LLM_API_KEY
          valueFrom:
            secretKeyRef:
              name: llm-secrets
              key: api-key
        ports:
        - containerPort: 8080
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 30
```

## Best Practices

### 1. Error Handling and Resilience

```python
import backoff
from typing import Optional

class ResilientTestAgent:
    def __init__(self, llm_provider: LLMProvider):
        self.llm = llm_provider
        self.fern_client = FernMyceliumClient()
    
    @backoff.on_exception(
        backoff.expo,
        requests.exceptions.RequestException,
        max_tries=3
    )
    def fetch_with_retry(self, project_id: str) -> Optional[Dict[str, Any]]:
        """Fetch test data with exponential backoff"""
        return self.fern_client.query_flaky_tests(project_id)
    
    def safe_analyze(self, project_id: str) -> Dict[str, Any]:
        """Safely analyze with comprehensive error handling"""
        try:
            # Try to fetch data
            test_data = self.fetch_with_retry(project_id)
            if not test_data:
                return {"error": "No data available"}
            
            # Try to analyze
            analysis = self.llm.generate(f"Analyze: {test_data}")
            
            return {
                "status": "success",
                "project_id": project_id,
                "analysis": analysis
            }
            
        except Exception as e:
            return {
                "status": "error",
                "project_id": project_id,
                "error": str(e),
                "fallback": "Manual review recommended"
            }
```

### 2. Performance Optimization

```python
import asyncio
from concurrent.futures import ThreadPoolExecutor
from typing import List

class OptimizedTestAgent:
    def __init__(self, llm_provider: LLMProvider, max_workers: int = 4):
        self.llm = llm_provider
        self.fern_client = FernMyceliumClient()
        self.executor = ThreadPoolExecutor(max_workers=max_workers)
    
    async def parallel_analysis(self, project_ids: List[str]) -> List[Dict[str, Any]]:
        """Analyze multiple projects in parallel"""
        loop = asyncio.get_event_loop()
        
        tasks = [
            loop.run_in_executor(
                self.executor,
                self.analyze_single_project,
                project_id
            )
            for project_id in project_ids
        ]
        
        return await asyncio.gather(*tasks)
    
    def analyze_single_project(self, project_id: str) -> Dict[str, Any]:
        """Analyze a single project"""
        test_data = self.fern_client.get_test_summary(project_id)
        analysis = self.llm.generate(f"Analyze: {test_data}")
        
        return {
            "project_id": project_id,
            "analysis": analysis
        }
```

### 3. Monitoring and Metrics

```python
import time
import logging
from prometheus_client import Counter, Histogram, start_http_server

class MonitoredTestAgent:
    def __init__(self, llm_provider: LLMProvider):
        self.llm = llm_provider
        self.fern_client = FernMyceliumClient()
        
        # Prometheus metrics
        self.request_count = Counter('test_agent_requests_total', 'Total requests')
        self.request_duration = Histogram('test_agent_request_duration_seconds', 'Request duration')
        self.error_count = Counter('test_agent_errors_total', 'Total errors')
        
        # Start metrics server
        start_http_server(8000)
    
    def analyze_with_monitoring(self, project_id: str) -> Dict[str, Any]:
        """Analyze with full monitoring"""
        self.request_count.inc()
        
        start_time = time.time()
        
        try:
            result = self.analyze_project(project_id)
            return result
            
        except Exception as e:
            self.error_count.inc()
            logging.error(f"Analysis failed for {project_id}: {e}")
            raise
            
        finally:
            duration = time.time() - start_time
            self.request_duration.observe(duration)
```

## Troubleshooting

### Common Issues

1. **Model context limits**: Break large test datasets into smaller chunks
2. **Rate limiting**: Implement proper backoff strategies
3. **Memory usage**: Use streaming for large analyses
4. **Network timeouts**: Configure appropriate timeout values

### Debug Mode

```python
class DebugTestAgent:
    def __init__(self, llm_provider: LLMProvider, debug: bool = False):
        self.llm = llm_provider
        self.debug = debug
        
    def analyze_with_debug(self, project_id: str) -> Dict[str, Any]:
        """Analyze with debug information"""
        if self.debug:
            print(f"Fetching data for project: {project_id}")
        
        test_data = self.fern_client.get_test_summary(project_id)
        
        if self.debug:
            print(f"Test data length: {len(str(test_data))}")
            print(f"First 200 chars: {str(test_data)[:200]}...")
        
        analysis = self.llm.generate(f"Analyze: {test_data}")
        
        if self.debug:
            print(f"Analysis length: {len(analysis)}")
        
        return {"analysis": analysis, "debug_info": {"data_length": len(str(test_data))}}
```

For additional examples and patterns, see the [Usage Examples](USAGE_EXAMPLES.md) documentation.