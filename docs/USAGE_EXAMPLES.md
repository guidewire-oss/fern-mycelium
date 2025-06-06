# Fern-Mycelium Usage Examples and Tutorials

## Overview

This document provides practical examples and step-by-step tutorials for using the fern-mycelium MCP server in various scenarios. From basic queries to advanced AI agent integrations, these examples will help you get the most out of your test intelligence system.

## Getting Started

### 1. Basic Health Check

First, verify your server is running correctly:

```bash
# Check server health
curl http://localhost:8081/healthz

# Expected response:
# {"status":"ok","message":"fern-mycelium is healthy 🍄"}
```

### 2. Simple GraphQL Query

Test the basic GraphQL functionality:

```bash
curl -X POST http://localhost:8081/query \
  -H "Content-Type: application/json" \
  -d '{
    "query": "{ flakyTests(limit: 3, projectID: \"demo\") { testName passRate failureRate } }"
  }'
```

### 3. Using the Test Client

Run the provided test client to see the system in action:

```bash
# Run the Golang test client
go run test_mcp_client.go

# Expected output:
# 🧪 Testing fern-mycelium MCP Server
# ==================================================
# 
# 1. Testing Health Check...
# ✅ Health: ok - fern-mycelium is healthy 🍄
# 
# 2. Testing Flaky Test Detection...
# ✅ Found 3 tests with flaky behavior:
#    🔴 FLAKY TestFlakyDBConnection
#       - Pass Rate: 50.0%
#       - Failure Rate: 50.0%
#       - Total Runs: 4
```

## Common Use Cases

### 1. Daily Test Health Monitoring

Create a script to monitor test health across multiple projects:

```python
#!/usr/bin/env python3
"""
Daily test health monitoring script
"""
import requests
import json
from datetime import datetime
import smtplib
from email.mime.text import MimeText

class DailyTestMonitor:
    def __init__(self, server_url="http://localhost:8081"):
        self.server_url = server_url
        self.projects = ["web-app", "mobile-app", "api-service", "auth-service"]
    
    def get_project_health(self, project_id):
        """Get health metrics for a project"""
        query = f"""
        {{
            flakyTests(limit: 100, projectID: "{project_id}") {{
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
            json={"query": query}
        )
        
        if response.status_code == 200:
            data = response.json().get("data", {}).get("flakyTests", [])
            return self.analyze_health(data)
        return None
    
    def analyze_health(self, tests):
        """Analyze test health metrics"""
        if not tests:
            return {"status": "no_data", "total": 0}
        
        total_tests = len(tests)
        high_risk = len([t for t in tests if t["failureRate"] > 0.3])
        moderate_risk = len([t for t in tests if 0 < t["failureRate"] <= 0.3])
        stable = len([t for t in tests if t["failureRate"] == 0])
        
        # Calculate overall health score (0-100)
        health_score = ((stable * 100) + (moderate_risk * 50) + (high_risk * 0)) / total_tests
        
        return {
            "status": "healthy" if health_score > 80 else "warning" if health_score > 60 else "critical",
            "health_score": round(health_score, 1),
            "total": total_tests,
            "high_risk": high_risk,
            "moderate_risk": moderate_risk,
            "stable": stable,
            "most_problematic": sorted(tests, key=lambda x: x["failureRate"], reverse=True)[:3]
        }
    
    def generate_report(self):
        """Generate daily health report"""
        report = f"# Daily Test Health Report - {datetime.now().strftime('%Y-%m-%d')}\n\n"
        
        overall_health = []
        
        for project in self.projects:
            health = self.get_project_health(project)
            if health:
                overall_health.append(health["health_score"])
                
                status_emoji = {
                    "healthy": "🟢",
                    "warning": "🟡", 
                    "critical": "🔴",
                    "no_data": "⚫"
                }
                
                report += f"## {project} {status_emoji[health['status']]}\n"
                report += f"- **Health Score**: {health['health_score']}%\n"
                report += f"- **Total Tests**: {health['total']}\n"
                report += f"- **High Risk**: {health['high_risk']}\n"
                report += f"- **Moderate Risk**: {health['moderate_risk']}\n"
                report += f"- **Stable**: {health['stable']}\n"
                
                if health['most_problematic']:
                    report += f"- **Most Problematic**:\n"
                    for test in health['most_problematic']:
                        if test['failureRate'] > 0:
                            report += f"  - {test['testName']}: {test['failureRate']:.1%} failure rate\n"
                
                report += "\n"
        
        # Overall summary
        if overall_health:
            avg_health = sum(overall_health) / len(overall_health)
            report += f"## Overall System Health: {avg_health:.1f}%\n\n"
            
            if avg_health < 70:
                report += "⚠️ **ATTENTION REQUIRED**: Multiple projects have concerning test stability.\n"
            elif avg_health < 85:
                report += "👀 **MONITOR CLOSELY**: Some projects need attention.\n"
            else:
                report += "✅ **HEALTHY**: Test suite is in good condition.\n"
        
        return report
    
    def send_email_report(self, report, recipients):
        """Send report via email"""
        msg = MimeText(report)
        msg['Subject'] = f"Daily Test Health Report - {datetime.now().strftime('%Y-%m-%d')}"
        msg['From'] = "test-intelligence@yourcompany.com"
        msg['To'] = ", ".join(recipients)
        
        # Configure your SMTP server
        with smtplib.SMTP('localhost') as server:
            server.send_message(msg)

# Usage
if __name__ == "__main__":
    monitor = DailyTestMonitor()
    report = monitor.generate_report()
    print(report)
    
    # Optionally send via email
    # monitor.send_email_report(report, ["team@yourcompany.com"])
```

### 2. CI/CD Integration

Integrate test intelligence into your CI/CD pipeline:

```yaml
# .github/workflows/test-intelligence.yml
name: Test Intelligence Analysis

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  test-analysis:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    
    - name: Setup Go
      uses: actions/setup-go@v3
      with:
        go-version: '1.21'
    
    - name: Run Tests
      run: go test ./... -json > test-results.json
    
    - name: Analyze Test Intelligence
      run: |
        curl -X POST ${{ secrets.FERN_MYCELIUM_URL }}/query \
          -H "Content-Type: application/json" \
          -d '{
            "query": "{ flakyTests(limit: 10, projectID: \"${{ github.repository }}\") { testName failureRate runCount } }"
          }' > flaky-tests.json
    
    - name: Generate Intelligence Report
      run: |
        python3 << 'EOF'
        import json
        import os
        
        # Load flaky test data
        with open('flaky-tests.json') as f:
            data = json.load(f)
        
        tests = data.get('data', {}).get('flakyTests', [])
        high_risk = [t for t in tests if t['failureRate'] > 0.3]
        
        if high_risk:
            print("🚨 High-risk flaky tests detected:")
            for test in high_risk:
                print(f"- {test['testName']}: {test['failureRate']:.1%} failure rate")
            
            # Set output for GitHub Actions
            with open(os.environ['GITHUB_OUTPUT'], 'a') as f:
                f.write(f"high_risk_count={len(high_risk)}\n")
        else:
            print("✅ No high-risk flaky tests detected")
            with open(os.environ['GITHUB_OUTPUT'], 'a') as f:
                f.write("high_risk_count=0\n")
        EOF
      id: analysis
    
    - name: Comment on PR
      if: github.event_name == 'pull_request' && steps.analysis.outputs.high_risk_count > 0
      uses: actions/github-script@v6
      with:
        script: |
          github.rest.issues.createComment({
            issue_number: context.issue.number,
            owner: context.repo.owner,
            repo: context.repo.repo,
            body: '🚨 **Test Intelligence Alert**: This PR affects a repository with ${{ steps.analysis.outputs.high_risk_count }} high-risk flaky tests. Please review test stability before merging.'
          })
```

### 3. Slack Integration

Create a Slack bot that provides test intelligence on demand:

```python
from slack_bolt import App
from slack_bolt.adapter.socket_mode import SocketModeHandler
import requests
import json

app = App(token="your-slack-bot-token")

class TestIntelligenceBot:
    def __init__(self, fern_url="http://localhost:8081"):
        self.fern_url = fern_url
    
    def get_project_summary(self, project_id):
        """Get formatted project summary for Slack"""
        query = f"""
        {{
            flakyTests(limit: 20, projectID: "{project_id}") {{
                testName
                passRate
                failureRate
                runCount
                lastFailure
            }}
        }}
        """
        
        response = requests.post(
            f"{self.fern_url}/query",
            json={"query": query}
        )
        
        if response.status_code != 200:
            return "❌ Unable to fetch test data"
        
        data = response.json().get("data", {}).get("flakyTests", [])
        
        if not data:
            return f"📊 No test data found for project `{project_id}`"
        
        # Analyze data
        total = len(data)
        high_risk = [t for t in data if t["failureRate"] > 0.3]
        moderate_risk = [t for t in data if 0 < t["failureRate"] <= 0.3]
        stable = [t for t in data if t["failureRate"] == 0]
        
        # Format for Slack
        blocks = [
            {
                "type": "header",
                "text": {
                    "type": "plain_text",
                    "text": f"📊 Test Intelligence: {project_id}"
                }
            },
            {
                "type": "section",
                "fields": [
                    {"type": "mrkdwn", "text": f"*Total Tests:* {total}"},
                    {"type": "mrkdwn", "text": f"*Health Score:* {((len(stable) / total) * 100):.1f}%"},
                    {"type": "mrkdwn", "text": f"*🔴 High Risk:* {len(high_risk)}"},
                    {"type": "mrkdwn", "text": f"*🟡 Moderate Risk:* {len(moderate_risk)}"},
                    {"type": "mrkdwn", "text": f"*🟢 Stable:* {len(stable)}"}
                ]
            }
        ]
        
        if high_risk:
            critical_list = "\n".join([
                f"• {test['testName']}: {test['failureRate']:.1%} failure rate"
                for test in high_risk[:5]
            ])
            
            blocks.append({
                "type": "section",
                "text": {
                    "type": "mrkdwn",
                    "text": f"*🚨 Critical Issues:*\n{critical_list}"
                }
            })
        
        return {"blocks": blocks}

bot = TestIntelligenceBot()

@app.command("/test-health")
def test_health_command(ack, respond, command):
    """Handle /test-health slash command"""
    ack()
    
    project_id = command['text'].strip() or "web-app"
    
    try:
        result = bot.get_project_summary(project_id)
        
        if isinstance(result, dict):
            respond(result)
        else:
            respond(result)
    except Exception as e:
        respond(f"❌ Error: {str(e)}")

@app.message("test intelligence")
def handle_test_intelligence(message, say):
    """Handle mentions of 'test intelligence'"""
    say("👋 I can help with test intelligence! Use `/test-health [project-name]` to get started.")

# Start the app
if __name__ == "__main__":
    handler = SocketModeHandler(app, "your-app-token")
    handler.start()
```

### 4. Development Workflow Integration

Integrate test intelligence into your development workflow:

```python
#!/usr/bin/env python3
"""
Pre-commit hook for test intelligence
"""
import sys
import subprocess
import requests
import json

def get_affected_projects():
    """Get projects affected by current changes"""
    try:
        # Get changed files
        result = subprocess.run(
            ["git", "diff", "--cached", "--name-only"],
            capture_output=True,
            text=True
        )
        
        changed_files = result.stdout.strip().split("\n")
        
        # Map files to projects (customize based on your structure)
        projects = set()
        for file in changed_files:
            if file.startswith("web/"):
                projects.add("web-app")
            elif file.startswith("mobile/"):
                projects.add("mobile-app")
            elif file.startswith("api/"):
                projects.add("api-service")
        
        return list(projects)
    except:
        return []

def check_test_health(project_id):
    """Check test health for a project"""
    query = f"""
    {{
        flakyTests(limit: 50, projectID: "{project_id}") {{
            testName
            failureRate
            runCount
        }}
    }}
    """
    
    try:
        response = requests.post(
            "http://localhost:8081/query",
            json={"query": query},
            timeout=10
        )
        
        if response.status_code == 200:
            data = response.json().get("data", {}).get("flakyTests", [])
            critical_tests = [t for t in data if t["failureRate"] > 0.5]
            return critical_tests
        
    except requests.exceptions.ConnectionError:
        print("⚠️  Warning: Test intelligence server not available")
        return []
    except:
        return []

def main():
    """Main pre-commit hook"""
    affected_projects = get_affected_projects()
    
    if not affected_projects:
        print("✅ No projects affected, skipping test intelligence check")
        return 0
    
    print(f"🔍 Checking test intelligence for: {', '.join(affected_projects)}")
    
    critical_issues = {}
    
    for project in affected_projects:
        critical_tests = check_test_health(project)
        if critical_tests:
            critical_issues[project] = critical_tests
    
    if critical_issues:
        print("\n🚨 CRITICAL TEST STABILITY ISSUES DETECTED:")
        print("=" * 50)
        
        for project, tests in critical_issues.items():
            print(f"\n📊 {project}:")
            for test in tests[:3]:  # Show top 3 most critical
                print(f"  ❌ {test['testName']}: {test['failureRate']:.1%} failure rate")
        
        print("\n⚠️  These projects have severely unstable tests.")
        print("Consider fixing critical test issues before committing.")
        
        # Ask for confirmation
        response = input("\nContinue with commit anyway? (y/N): ")
        if response.lower() != 'y':
            print("Commit aborted.")
            return 1
    
    print("✅ Test intelligence check passed")
    return 0

if __name__ == "__main__":
    sys.exit(main())
```

### 5. Performance Testing Integration

Use test intelligence for performance test analysis:

```python
import requests
import json
import statistics
from datetime import datetime, timedelta

class PerformanceTestAnalyzer:
    def __init__(self, fern_url="http://localhost:8081"):
        self.fern_url = fern_url
    
    def analyze_performance_trends(self, project_id):
        """Analyze performance test trends"""
        # This assumes you have performance data in your test results
        query = f"""
        {{
            flakyTests(limit: 100, projectID: "{project_id}") {{
                testName
                passRate
                failureRate
                runCount
                lastFailure
            }}
        }}
        """
        
        response = requests.post(
            f"{self.fern_url}/query",
            json={"query": query}
        )
        
        if response.status_code != 200:
            return {"error": "Unable to fetch data"}
        
        tests = response.json().get("data", {}).get("flakyTests", [])
        
        # Focus on performance-related tests
        perf_tests = [t for t in tests if "performance" in t["testName"].lower() 
                     or "load" in t["testName"].lower() 
                     or "stress" in t["testName"].lower()]
        
        if not perf_tests:
            return {"message": "No performance tests found"}
        
        # Analyze stability of performance tests
        unstable_perf = [t for t in perf_tests if t["failureRate"] > 0.1]
        
        analysis = {
            "total_performance_tests": len(perf_tests),
            "unstable_performance_tests": len(unstable_perf),
            "stability_score": ((len(perf_tests) - len(unstable_perf)) / len(perf_tests)) * 100,
            "recommendations": []
        }
        
        if unstable_perf:
            analysis["recommendations"].append(
                "Review unstable performance tests - they may indicate infrastructure issues"
            )
            analysis["critical_tests"] = [
                f"{test['testName']} ({test['failureRate']:.1%} failure rate)"
                for test in unstable_perf
            ]
        
        if analysis["stability_score"] < 80:
            analysis["recommendations"].append(
                "Performance test suite needs attention - consider reviewing test environments"
            )
        
        return analysis
    
    def generate_performance_report(self, projects):
        """Generate comprehensive performance report"""
        report = f"# Performance Test Intelligence Report\n"
        report += f"Generated: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}\n\n"
        
        for project in projects:
            analysis = self.analyze_performance_trends(project)
            
            report += f"## {project}\n"
            
            if "error" in analysis:
                report += f"❌ {analysis['error']}\n\n"
                continue
            
            if "message" in analysis:
                report += f"ℹ️  {analysis['message']}\n\n"
                continue
            
            stability_emoji = "🟢" if analysis['stability_score'] > 90 else "🟡" if analysis['stability_score'] > 70 else "🔴"
            
            report += f"- **Stability Score**: {stability_emoji} {analysis['stability_score']:.1f}%\n"
            report += f"- **Total Performance Tests**: {analysis['total_performance_tests']}\n"
            report += f"- **Unstable Tests**: {analysis['unstable_performance_tests']}\n"
            
            if analysis.get('critical_tests'):
                report += f"- **Critical Issues**:\n"
                for test in analysis['critical_tests']:
                    report += f"  - {test}\n"
            
            if analysis.get('recommendations'):
                report += f"- **Recommendations**:\n"
                for rec in analysis['recommendations']:
                    report += f"  - {rec}\n"
            
            report += "\n"
        
        return report

# Usage
analyzer = PerformanceTestAnalyzer()
projects = ["api-service", "web-app", "mobile-app"]
report = analyzer.generate_performance_report(projects)
print(report)
```

### 6. AI-Powered Test Recommendations

Create an AI system that provides intelligent test recommendations:

```python
import openai
import requests
import json
from typing import List, Dict, Any

class AITestRecommendationEngine:
    def __init__(self, openai_api_key: str, fern_url: str = "http://localhost:8081"):
        self.openai_client = openai.OpenAI(api_key=openai_api_key)
        self.fern_url = fern_url
    
    def get_comprehensive_test_data(self, project_id: str) -> Dict[str, Any]:
        """Get comprehensive test data for analysis"""
        query = f"""
        {{
            flakyTests(limit: 100, projectID: "{project_id}") {{
                testID
                testName
                passRate
                failureRate
                runCount
                lastFailure
            }}
        }}
        """
        
        response = requests.post(
            f"{self.fern_url}/query",
            json={"query": query}
        )
        
        if response.status_code == 200:
            return response.json().get("data", {}).get("flakyTests", [])
        return []
    
    def generate_ai_recommendations(self, project_id: str) -> str:
        """Generate AI-powered recommendations"""
        test_data = self.get_comprehensive_test_data(project_id)
        
        if not test_data:
            return "No test data available for analysis."
        
        # Prepare data for AI analysis
        data_summary = {
            "project": project_id,
            "total_tests": len(test_data),
            "high_risk_tests": [t for t in test_data if t["failureRate"] > 0.3],
            "moderate_risk_tests": [t for t in test_data if 0 < t["failureRate"] <= 0.3],
            "stable_tests": [t for t in test_data if t["failureRate"] == 0],
            "most_problematic": sorted(test_data, key=lambda x: x["failureRate"], reverse=True)[:10]
        }
        
        prompt = f"""
        You are a senior test engineer analyzing test intelligence data. Provide specific, actionable recommendations based on this test analysis:

        Project: {project_id}
        Total Tests: {data_summary['total_tests']}
        High Risk Tests (>30% failure): {len(data_summary['high_risk_tests'])}
        Moderate Risk Tests (1-30% failure): {len(data_summary['moderate_risk_tests'])}
        Stable Tests (0% failure): {len(data_summary['stable_tests'])}

        Most Problematic Tests:
        {json.dumps([{{'name': t['testName'], 'failure_rate': t['failureRate'], 'run_count': t['runCount']}} for t in data_summary['most_problematic']], indent=2)}

        Provide:
        1. Priority ranking of issues (Critical/High/Medium/Low)
        2. Specific recommendations for each high-risk test
        3. Resource allocation suggestions
        4. Technical strategies for improving test stability
        5. Estimated effort for each recommendation

        Format your response as a structured action plan.
        """
        
        response = self.openai_client.chat.completions.create(
            model="gpt-4",
            messages=[{"role": "user", "content": prompt}],
            temperature=0.3,
            max_tokens=1500
        )
        
        return response.choices[0].message.content
    
    def compare_projects_and_recommend(self, project_ids: List[str]) -> str:
        """Compare multiple projects and provide strategic recommendations"""
        all_data = {}
        
        for project_id in project_ids:
            all_data[project_id] = self.get_comprehensive_test_data(project_id)
        
        # Calculate metrics for each project
        project_metrics = {}
        for project_id, tests in all_data.items():
            if tests:
                total = len(tests)
                high_risk = len([t for t in tests if t["failureRate"] > 0.3])
                health_score = ((total - high_risk) / total) * 100 if total > 0 else 0
                
                project_metrics[project_id] = {
                    "total_tests": total,
                    "high_risk_count": high_risk,
                    "health_score": health_score,
                    "most_critical": sorted(tests, key=lambda x: x["failureRate"], reverse=True)[:3]
                }
        
        prompt = f"""
        You are a test engineering director analyzing test health across multiple projects. Provide strategic recommendations for resource allocation and improvement priorities.

        Project Metrics:
        {json.dumps(project_metrics, indent=2)}

        Provide:
        1. Project ranking by urgency (most to least critical)
        2. Strategic resource allocation recommendations
        3. Cross-project patterns and systemic issues
        4. Timeline for addressing critical issues
        5. ROI analysis for test stability improvements
        6. Team/skill requirements for each improvement area

        Focus on organizational impact and strategic decision-making.
        """
        
        response = self.openai_client.chat.completions.create(
            model="gpt-4",
            messages=[{"role": "user", "content": prompt}],
            temperature=0.2,
            max_tokens=2000
        )
        
        return response.choices[0].message.content
    
    def generate_weekly_intelligence_report(self, project_ids: List[str]) -> str:
        """Generate comprehensive weekly intelligence report"""
        project_comparison = self.compare_projects_and_recommend(project_ids)
        
        individual_recommendations = {}
        for project_id in project_ids:
            individual_recommendations[project_id] = self.generate_ai_recommendations(project_id)
        
        # Create comprehensive report
        report = f"""# Weekly Test Intelligence Report
Generated: {datetime.now().strftime('%Y-%m-%d')}

## Executive Summary
{project_comparison}

## Individual Project Recommendations
"""
        
        for project_id, recommendations in individual_recommendations.items():
            report += f"\n### {project_id}\n{recommendations}\n"
        
        return report

# Usage examples
def main():
    engine = AITestRecommendationEngine("your-openai-api-key")
    
    # Single project analysis
    recommendations = engine.generate_ai_recommendations("web-app")
    print("=== Single Project Recommendations ===")
    print(recommendations)
    
    # Multi-project strategic analysis
    projects = ["web-app", "mobile-app", "api-service"]
    strategic_analysis = engine.compare_projects_and_recommend(projects)
    print("\n=== Strategic Analysis ===")
    print(strategic_analysis)
    
    # Weekly report
    weekly_report = engine.generate_weekly_intelligence_report(projects)
    print("\n=== Weekly Report ===")
    print(weekly_report)

if __name__ == "__main__":
    main()
```

## Advanced Patterns

### 1. Real-time Test Intelligence Dashboard

Create a real-time dashboard using WebSockets:

```html
<!DOCTYPE html>
<html>
<head>
    <title>Test Intelligence Dashboard</title>
    <script src="https://cdn.jsdelivr.net/npm/chart.js"></script>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .metric-card { 
            display: inline-block; 
            margin: 10px; 
            padding: 20px; 
            border: 1px solid #ddd; 
            border-radius: 8px; 
            background: #f9f9f9; 
        }
        .critical { border-left: 5px solid #ff4444; }
        .warning { border-left: 5px solid #ffaa00; }
        .healthy { border-left: 5px solid #00aa44; }
        #chart-container { width: 80%; margin: 20px auto; }
    </style>
</head>
<body>
    <h1>🧪 Test Intelligence Dashboard</h1>
    
    <div id="metrics-container">
        <!-- Metrics will be populated here -->
    </div>
    
    <div id="chart-container">
        <canvas id="healthChart"></canvas>
    </div>
    
    <div id="alerts-container">
        <h2>🚨 Active Alerts</h2>
        <div id="alerts-list"></div>
    </div>

    <script>
        class TestIntelligenceDashboard {
            constructor() {
                this.fern_url = 'http://localhost:8081';
                this.projects = ['web-app', 'mobile-app', 'api-service', 'auth-service'];
                this.chart = null;
                this.init();
            }
            
            async init() {
                await this.updateDashboard();
                setInterval(() => this.updateDashboard(), 30000); // Update every 30 seconds
                this.initChart();
            }
            
            async fetchProjectData(projectId) {
                const query = `{
                    flakyTests(limit: 50, projectID: "${projectId}") {
                        testName
                        passRate
                        failureRate
                        runCount
                        lastFailure
                    }
                }`;
                
                try {
                    const response = await fetch(`${this.fern_url}/query`, {
                        method: 'POST',
                        headers: { 'Content-Type': 'application/json' },
                        body: JSON.stringify({ query })
                    });
                    
                    const data = await response.json();
                    return data.data?.flakyTests || [];
                } catch (error) {
                    console.error(`Error fetching data for ${projectId}:`, error);
                    return [];
                }
            }
            
            analyzeProjectHealth(tests) {
                if (!tests.length) return { status: 'no_data', score: 0, total: 0 };
                
                const total = tests.length;
                const highRisk = tests.filter(t => t.failureRate > 0.3).length;
                const moderateRisk = tests.filter(t => t.failureRate > 0 && t.failureRate <= 0.3).length;
                const stable = tests.filter(t => t.failureRate === 0).length;
                
                const score = ((stable * 100) + (moderateRisk * 50)) / total;
                
                return {
                    status: score > 80 ? 'healthy' : score > 60 ? 'warning' : 'critical',
                    score: Math.round(score),
                    total,
                    highRisk,
                    moderateRisk,
                    stable,
                    mostProblematic: tests
                        .filter(t => t.failureRate > 0)
                        .sort((a, b) => b.failureRate - a.failureRate)
                        .slice(0, 3)
                };
            }
            
            async updateDashboard() {
                const projectHealthData = {};
                
                // Fetch data for all projects
                for (const project of this.projects) {
                    const tests = await this.fetchProjectData(project);
                    projectHealthData[project] = this.analyzeProjectHealth(tests);
                }
                
                this.updateMetrics(projectHealthData);
                this.updateChart(projectHealthData);
                this.updateAlerts(projectHealthData);
            }
            
            updateMetrics(healthData) {
                const container = document.getElementById('metrics-container');
                container.innerHTML = '';
                
                for (const [project, health] of Object.entries(healthData)) {
                    const card = document.createElement('div');
                    card.className = `metric-card ${health.status}`;
                    card.innerHTML = `
                        <h3>${project}</h3>
                        <div style="font-size: 24px; font-weight: bold;">${health.score}%</div>
                        <div>Health Score</div>
                        <div style="margin-top: 10px; font-size: 12px;">
                            Total: ${health.total} | 
                            High Risk: ${health.highRisk} | 
                            Stable: ${health.stable}
                        </div>
                    `;
                    container.appendChild(card);
                }
            }
            
            initChart() {
                const ctx = document.getElementById('healthChart').getContext('2d');
                this.chart = new Chart(ctx, {
                    type: 'bar',
                    data: {
                        labels: this.projects,
                        datasets: [{
                            label: 'Health Score',
                            data: [],
                            backgroundColor: 'rgba(54, 162, 235, 0.6)',
                            borderColor: 'rgba(54, 162, 235, 1)',
                            borderWidth: 1
                        }]
                    },
                    options: {
                        responsive: true,
                        scales: {
                            y: {
                                beginAtZero: true,
                                max: 100
                            }
                        }
                    }
                });
            }
            
            updateChart(healthData) {
                if (!this.chart) return;
                
                const scores = this.projects.map(project => healthData[project]?.score || 0);
                const colors = scores.map(score => 
                    score > 80 ? 'rgba(76, 175, 80, 0.6)' :
                    score > 60 ? 'rgba(255, 193, 7, 0.6)' :
                    'rgba(244, 67, 54, 0.6)'
                );
                
                this.chart.data.datasets[0].data = scores;
                this.chart.data.datasets[0].backgroundColor = colors;
                this.chart.update();
            }
            
            updateAlerts(healthData) {
                const container = document.getElementById('alerts-list');
                container.innerHTML = '';
                
                const alerts = [];
                
                for (const [project, health] of Object.entries(healthData)) {
                    if (health.status === 'critical') {
                        alerts.push({
                            level: 'critical',
                            message: `${project}: Critical test stability issues (${health.score}% health score)`,
                            details: health.mostProblematic
                        });
                    } else if (health.highRisk > 0) {
                        alerts.push({
                            level: 'warning',
                            message: `${project}: ${health.highRisk} high-risk tests detected`,
                            details: health.mostProblematic
                        });
                    }
                }
                
                if (alerts.length === 0) {
                    container.innerHTML = '<div style="color: green;">✅ No active alerts</div>';
                    return;
                }
                
                alerts.forEach(alert => {
                    const alertDiv = document.createElement('div');
                    alertDiv.style.margin = '10px 0';
                    alertDiv.style.padding = '10px';
                    alertDiv.style.border = `1px solid ${alert.level === 'critical' ? '#ff4444' : '#ffaa00'}`;
                    alertDiv.style.borderRadius = '4px';
                    alertDiv.style.backgroundColor = alert.level === 'critical' ? '#ffe6e6' : '#fff3cd';
                    
                    const icon = alert.level === 'critical' ? '🚨' : '⚠️';
                    alertDiv.innerHTML = `
                        <div><strong>${icon} ${alert.message}</strong></div>
                        ${alert.details.map(test => 
                            `<div style="margin-left: 20px; font-size: 12px;">
                                • ${test.testName}: ${(test.failureRate * 100).toFixed(1)}% failure rate
                            </div>`
                        ).join('')}
                    `;
                    
                    container.appendChild(alertDiv);
                });
            }
        }
        
        // Initialize dashboard
        new TestIntelligenceDashboard();
    </script>
</body>
</html>
```

### 2. Automated Test Triage System

```python
import requests
import json
from dataclasses import dataclass
from typing import List, Dict, Any
from enum import Enum

class Priority(Enum):
    CRITICAL = "critical"
    HIGH = "high"
    MEDIUM = "medium"
    LOW = "low"

@dataclass
class TestIssue:
    test_name: str
    project_id: str
    failure_rate: float
    run_count: int
    last_failure: str
    priority: Priority
    estimated_effort: str
    recommended_action: str

class AutomatedTestTriage:
    def __init__(self, fern_url="http://localhost:8081"):
        self.fern_url = fern_url
        
        # Triage rules
        self.priority_rules = [
            (lambda t: t['failureRate'] > 0.7, Priority.CRITICAL, "immediate", "Block production deployments"),
            (lambda t: t['failureRate'] > 0.5, Priority.CRITICAL, "1-2 days", "Fix before next release"),
            (lambda t: t['failureRate'] > 0.3, Priority.HIGH, "3-5 days", "Investigate root cause"),
            (lambda t: t['failureRate'] > 0.1, Priority.MEDIUM, "1-2 weeks", "Monitor and improve"),
            (lambda t: t['failureRate'] > 0, Priority.LOW, "monthly", "Review during maintenance")
        ]
    
    def fetch_all_test_data(self, projects: List[str]) -> Dict[str, List[Dict]]:
        """Fetch test data for all projects"""
        all_data = {}
        
        for project in projects:
            query = f"""
            {{
                flakyTests(limit: 100, projectID: "{project}") {{
                    testName
                    passRate
                    failureRate
                    runCount
                    lastFailure
                }}
            }}
            """
            
            try:
                response = requests.post(
                    f"{self.fern_url}/query",
                    json={"query": query},
                    timeout=30
                )
                
                if response.status_code == 200:
                    data = response.json().get("data", {}).get("flakyTests", [])
                    all_data[project] = data
                else:
                    all_data[project] = []
                    
            except Exception as e:
                print(f"Error fetching data for {project}: {e}")
                all_data[project] = []
        
        return all_data
    
    def triage_test(self, test_data: Dict[str, Any], project_id: str) -> TestIssue:
        """Triage a single test"""
        # Apply priority rules
        priority = Priority.LOW
        effort = "unknown"
        action = "review"
        
        for rule_func, rule_priority, rule_effort, rule_action in self.priority_rules:
            if rule_func(test_data):
                priority = rule_priority
                effort = rule_effort
                action = rule_action
                break
        
        # Additional heuristics
        if "auth" in test_data['testName'].lower() and test_data['failureRate'] > 0.2:
            priority = Priority.CRITICAL
            action = "Security-critical test - immediate review required"
        
        if "payment" in test_data['testName'].lower() and test_data['failureRate'] > 0.1:
            priority = Priority.HIGH
            action = "Business-critical test - expedited fix needed"
        
        return TestIssue(
            test_name=test_data['testName'],
            project_id=project_id,
            failure_rate=test_data['failureRate'],
            run_count=test_data['runCount'],
            last_failure=test_data.get('lastFailure', 'unknown'),
            priority=priority,
            estimated_effort=effort,
            recommended_action=action
        )
    
    def generate_triage_report(self, projects: List[str]) -> str:
        """Generate comprehensive triage report"""
        all_data = self.fetch_all_test_data(projects)
        
        # Triage all tests
        all_issues = []
        for project, tests in all_data.items():
            for test in tests:
                if test['failureRate'] > 0:  # Only triage failing tests
                    issue = self.triage_test(test, project)
                    all_issues.append(issue)
        
        # Sort by priority and failure rate
        priority_order = {Priority.CRITICAL: 0, Priority.HIGH: 1, Priority.MEDIUM: 2, Priority.LOW: 3}
        all_issues.sort(key=lambda x: (priority_order[x.priority], -x.failure_rate))
        
        # Generate report
        report = f"# Automated Test Triage Report\n"
        report += f"Generated: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}\n\n"
        
        # Summary
        summary = {p: len([i for i in all_issues if i.priority == p]) for p in Priority}
        report += f"## Summary\n"
        report += f"- 🚨 Critical: {summary[Priority.CRITICAL]}\n"
        report += f"- 🔴 High: {summary[Priority.HIGH]}\n"
        report += f"- 🟡 Medium: {summary[Priority.MEDIUM]}\n"
        report += f"- 🟢 Low: {summary[Priority.LOW]}\n\n"
        
        # Detailed breakdown by priority
        for priority in [Priority.CRITICAL, Priority.HIGH, Priority.MEDIUM, Priority.LOW]:
            priority_issues = [i for i in all_issues if i.priority == priority]
            
            if not priority_issues:
                continue
                
            emoji = {"critical": "🚨", "high": "🔴", "medium": "🟡", "low": "🟢"}[priority.value]
            
            report += f"## {emoji} {priority.value.title()} Priority Issues\n\n"
            
            for issue in priority_issues:
                report += f"### {issue.test_name} ({issue.project_id})\n"
                report += f"- **Failure Rate**: {issue.failure_rate:.1%}\n"
                report += f"- **Run Count**: {issue.run_count}\n"
                report += f"- **Estimated Effort**: {issue.estimated_effort}\n"
                report += f"- **Recommended Action**: {issue.recommended_action}\n"
                report += f"- **Last Failure**: {issue.last_failure}\n\n"
        
        # Action items
        critical_issues = [i for i in all_issues if i.priority == Priority.CRITICAL]
        if critical_issues:
            report += f"## 🚨 Immediate Action Required\n\n"
            report += f"The following {len(critical_issues)} tests require immediate attention:\n\n"
            
            for issue in critical_issues:
                report += f"1. **{issue.test_name}** ({issue.project_id})\n"
                report += f"   - {issue.failure_rate:.1%} failure rate\n"
                report += f"   - Action: {issue.recommended_action}\n\n"
        
        return report
    
    def generate_tickets(self, projects: List[str]) -> List[Dict[str, Any]]:
        """Generate JIRA/GitHub tickets for high-priority issues"""
        all_data = self.fetch_all_test_data(projects)
        tickets = []
        
        for project, tests in all_data.items():
            high_priority_tests = []
            
            for test in tests:
                if test['failureRate'] > 0.3:  # High priority threshold
                    issue = self.triage_test(test, project)
                    high_priority_tests.append(issue)
            
            if high_priority_tests:
                # Group by project
                ticket = {
                    "title": f"Fix flaky tests in {project}",
                    "description": self._generate_ticket_description(high_priority_tests),
                    "labels": ["flaky-tests", "testing", project],
                    "priority": "high" if any(i.priority == Priority.CRITICAL for i in high_priority_tests) else "medium",
                    "assignee": "test-team",
                    "project": project
                }
                tickets.append(ticket)
        
        return tickets
    
    def _generate_ticket_description(self, issues: List[TestIssue]) -> str:
        """Generate ticket description from issues"""
        desc = f"## Flaky Test Issues\n\n"
        desc += f"This ticket tracks {len(issues)} flaky tests that need attention.\n\n"
        
        desc += f"### Issues:\n\n"
        for issue in issues:
            desc += f"- **{issue.test_name}**\n"
            desc += f"  - Failure Rate: {issue.failure_rate:.1%}\n"
            desc += f"  - Priority: {issue.priority.value}\n"
            desc += f"  - Recommended Action: {issue.recommended_action}\n\n"
        
        desc += f"### Acceptance Criteria:\n\n"
        desc += f"- [ ] All tests have <10% failure rate\n"
        desc += f"- [ ] Root causes identified and documented\n"
        desc += f"- [ ] Fixes implemented and verified\n"
        desc += f"- [ ] Tests run stably for 1 week\n\n"
        
        return desc

# Usage
def main():
    triage = AutomatedTestTriage()
    projects = ["web-app", "mobile-app", "api-service"]
    
    # Generate triage report
    report = triage.generate_triage_report(projects)
    print(report)
    
    # Generate tickets for high-priority issues
    tickets = triage.generate_tickets(projects)
    for ticket in tickets:
        print(f"Ticket: {ticket['title']}")
        print(f"Priority: {ticket['priority']}")
        print("---")

if __name__ == "__main__":
    main()
```

## Best Practices

### 1. Data Collection Strategy

- **Comprehensive Coverage**: Ensure all critical test suites report to fern-mycelium
- **Historical Data**: Maintain sufficient historical data for trend analysis
- **Real-time Updates**: Implement real-time test result streaming for immediate insights

### 2. Alert Management

- **Graduated Alerting**: Use different alert thresholds for different types of tests
- **Alert Fatigue Prevention**: Implement intelligent alert suppression and grouping
- **Escalation Policies**: Define clear escalation paths for different severity levels

### 3. Integration Patterns

- **API-First Design**: Build integrations that can work with any LLM or AI system
- **Graceful Degradation**: Ensure systems work even when test intelligence is unavailable
- **Security**: Implement proper authentication and rate limiting for production use

These examples demonstrate the versatility and power of fern-mycelium for improving test intelligence and development workflows. Adapt them to your specific needs and infrastructure.