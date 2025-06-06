# Claude Desktop Integration with Fern-Mycelium

## Overview

This guide provides step-by-step instructions for integrating the fern-mycelium MCP server with Claude Desktop, enabling you to have intelligent conversations about your test data directly within Claude.

## Prerequisites

- Claude Desktop application installed
- Fern-mycelium MCP server deployed and running
- PostgreSQL database with test data
- Go runtime environment (if running locally)

## Setup Instructions

### 1. Locate Claude Desktop Configuration

Find your Claude Desktop configuration file:

**macOS:**
```bash
~/Library/Application Support/Claude/claude_desktop_config.json
```

**Windows:**
```bash
%APPDATA%\Claude\claude_desktop_config.json
```

**Linux:**
```bash
~/.config/Claude/claude_desktop_config.json
```

### 2. Create Configuration File

If the configuration file doesn't exist, create it:

```bash
# macOS/Linux
mkdir -p "$(dirname "$HOME/Library/Application Support/Claude/claude_desktop_config.json")"
touch "$HOME/Library/Application Support/Claude/claude_desktop_config.json"
```

### 3. Configure MCP Server

Add the fern-mycelium server to your configuration:

#### Option A: Local Development Setup

```json
{
  "mcpServers": {
    "fern-mycelium": {
      "command": "go",
      "args": ["run", "/path/to/fern-mycelium/main.go", "serve", "--port", "8081"],
      "env": {
        "DB_URL": "postgres://user:password@localhost:5432/fern_mycelium",
        "LOG_LEVEL": "info"
      }
    }
  }
}
```

#### Option B: Compiled Binary Setup

```json
{
  "mcpServers": {
    "fern-mycelium": {
      "command": "/path/to/fern-mycelium-binary",
      "args": ["serve", "--port", "8081"],
      "env": {
        "DB_URL": "postgres://user:password@localhost:5432/fern_mycelium",
        "LOG_LEVEL": "info"
      }
    }
  }
}
```

#### Option C: Docker Setup

```json
{
  "mcpServers": {
    "fern-mycelium": {
      "command": "docker",
      "args": [
        "run", 
        "--rm", 
        "-p", "8081:8081",
        "-e", "DB_URL=postgres://user:password@host.docker.internal:5432/fern_mycelium",
        "anoop2811/fern-mycelium:latest",
        "serve"
      ]
    }
  }
}
```

#### Option C: Remote Server Setup

```json
{
  "mcpServers": {
    "fern-mycelium": {
      "command": "curl",
      "args": ["-X", "POST", "https://your-server.com:8081/query"],
      "env": {
        "MCP_SERVER_URL": "https://your-server.com:8081"
      }
    }
  }
}
```

### 4. Restart Claude Desktop

After updating the configuration:

1. Completely quit Claude Desktop
2. Restart the application
3. Look for the fern-mycelium connection indicator in the interface

## Verification

### 1. Check Server Status

Verify your MCP server is running:

```bash
curl http://localhost:8081/healthz
```

Expected response:
```json
{"status":"ok","message":"fern-mycelium is healthy 🍄"}
```

### 2. Test Claude Integration

Open Claude Desktop and try these commands:

```
"Can you check if the fern-mycelium test intelligence system is available?"

"Show me the current flaky tests for project 'web-app'"

"Analyze test stability patterns and give me recommendations"
```

### 3. Verify MCP Connection

Claude Desktop should show an MCP connection indicator when fern-mycelium is properly connected.

## Usage Examples

### Basic Test Analysis

```
User: "What are the most problematic tests in our codebase?"

Claude: I'll query the fern-mycelium test intelligence system to analyze your flaky tests.

[Claude queries the MCP server and provides analysis based on the data]
```

### Project-Specific Analysis

```
User: "Show me test reliability metrics for the 'authentication' project"

Claude: Let me fetch the test data for the authentication project and analyze the patterns.

[Claude provides detailed analysis of authentication module tests]
```

### Trend Analysis

```
User: "Which tests have been failing most frequently recently?"

Claude: I'll analyze the recent test execution data to identify the most problematic tests.

[Claude provides prioritized list with failure rates and recommendations]
```

### Automated Reporting

```
User: "Generate a weekly test health report"

Claude: I'll create a comprehensive test health report based on the latest data from fern-mycelium.

[Claude generates formatted report with charts, trends, and action items]
```

## Advanced Configuration

### Environment Variables

Configure additional environment variables for enhanced functionality:

```json
{
  "mcpServers": {
    "fern-mycelium": {
      "command": "go",
      "args": ["run", "/path/to/fern-mycelium/main.go", "serve"],
      "env": {
        "DB_URL": "postgres://user:password@localhost:5432/fern_mycelium",
        "LOG_LEVEL": "debug",
        "MAX_QUERY_LIMIT": "100",
        "CACHE_TTL": "300",
        "METRICS_ENABLED": "true",
        "CORS_ENABLED": "true"
      }
    }
  }
}
```

### Multiple Projects

Configure multiple project environments:

```json
{
  "mcpServers": {
    "fern-mycelium-prod": {
      "command": "go",
      "args": ["run", "/path/to/fern-mycelium/main.go", "serve", "--port", "8081"],
      "env": {
        "DB_URL": "postgres://user:password@prod-db:5432/fern_mycelium",
        "ENVIRONMENT": "production"
      }
    },
    "fern-mycelium-staging": {
      "command": "go",
      "args": ["run", "/path/to/fern-mycelium/main.go", "serve", "--port", "8082"],
      "env": {
        "DB_URL": "postgres://user:password@staging-db:5432/fern_mycelium",
        "ENVIRONMENT": "staging"
      }
    }
  }
}
```

### Security Configuration

For production environments with authentication:

```json
{
  "mcpServers": {
    "fern-mycelium": {
      "command": "go",
      "args": ["run", "/path/to/fern-mycelium/main.go", "serve"],
      "env": {
        "DB_URL": "postgres://user:password@localhost:5432/fern_mycelium",
        "AUTH_ENABLED": "true",
        "API_KEY": "your-secure-api-key",
        "TLS_ENABLED": "true",
        "CERT_PATH": "/path/to/cert.pem",
        "KEY_PATH": "/path/to/key.pem"
      }
    }
  }
}
```

## Troubleshooting

### Common Issues

#### 1. Server Not Starting

**Symptoms:** Claude can't connect to fern-mycelium

**Solutions:**
```bash
# Check if port is available
netstat -an | grep 8081

# Verify database connection
psql "$DB_URL" -c "SELECT 1"

# Check logs
go run main.go serve --port 8081 --log-level debug
```

#### 2. Database Connection Issues

**Symptoms:** Server starts but queries fail

**Solutions:**
```bash
# Test database connectivity
curl http://localhost:8081/query -X POST \
  -H "Content-Type: application/json" \
  -d '{"query": "{ __schema { types { name } } }"}'

# Check database schema
psql "$DB_URL" -c "\dt"
```

#### 3. Permission Issues

**Symptoms:** Claude Desktop can't execute commands

**Solutions:**
```bash
# Make binary executable
chmod +x /path/to/fern-mycelium-binary

# Check Go installation
go version

# Verify path permissions
ls -la /path/to/fern-mycelium/
```

#### 4. MCP Connection Problems

**Symptoms:** Claude shows no MCP servers connected

**Solutions:**
1. Verify JSON syntax in configuration file
2. Restart Claude Desktop completely
3. Check Claude Desktop logs (if available)
4. Test server manually before configuring Claude

### Debug Mode

Enable debug mode for detailed logging:

```json
{
  "mcpServers": {
    "fern-mycelium": {
      "command": "go",
      "args": ["run", "/path/to/fern-mycelium/main.go", "serve", "--debug"],
      "env": {
        "DB_URL": "postgres://user:password@localhost:5432/fern_mycelium",
        "LOG_LEVEL": "debug",
        "MCP_DEBUG": "true"
      }
    }
  }
}
```

### Logs and Monitoring

Monitor your MCP server:

```bash
# Check server logs
tail -f /tmp/fern-mycelium.log

# Monitor database connections
psql "$DB_URL" -c "SELECT * FROM pg_stat_activity WHERE application_name LIKE '%fern%'"

# Check system resources
top -p $(pgrep -f fern-mycelium)
```

## Best Practices

### 1. Resource Management

- Set appropriate memory limits for the Go process
- Configure database connection pooling
- Use caching for frequently accessed data

### 2. Security

- Use environment variables for sensitive data
- Enable TLS for production deployments
- Implement rate limiting and authentication

### 3. Performance

- Optimize database queries for large datasets
- Use connection pooling
- Implement proper indexing on database tables

### 4. Monitoring

- Set up health checks
- Monitor query performance
- Track MCP server uptime and connectivity

## Example Conversations

### Test Health Check

```
User: "Is our test suite healthy?"

Claude: Let me check your test intelligence data from fern-mycelium.

Based on the current data, I found:
- 85% of tests are stable (0% failure rate)
- 12% have moderate flakiness (1-30% failure rate)  
- 3% are highly problematic (>30% failure rate)

Recommendations:
1. Focus on the 3 tests with >30% failure rate
2. Monitor the 8 tests in the moderate risk category
3. Your test suite is generally healthy but needs attention on critical failures
```

### Detailed Analysis

```
User: "Give me a detailed analysis of our authentication tests"

Claude: I'll analyze the authentication module tests from fern-mycelium.

Authentication Test Analysis:
- Total tests: 45
- Pass rate: 94.2%
- Most problematic: "TestOAuthFlow" (45% failure rate)
- Recently stable: "TestBasicAuth" (100% pass rate, 120 runs)

Critical Issues:
1. TestOAuthFlow failing due to token expiration timing
2. TestSSOIntegration intermittent network issues
3. TestPasswordReset needs database state cleanup

Recommendations:
1. Add retry logic to OAuth tests
2. Mock external SSO service calls
3. Implement proper test isolation for password reset
```

## Support

For additional help:

- Check the main [MCP Integration Guide](MCP_INTEGRATION.md)
- Review [Usage Examples](USAGE_EXAMPLES.md)  
- Report issues on GitHub
- Join the community discussions