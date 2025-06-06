# fern-mycelium

[![CI](https://github.com/guidewire-oss/fern-mycelium/workflows/Fern%20Mycelium%20CI%20Pipeline/badge.svg)](https://github.com/guidewire-oss/fern-mycelium/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/guidewire-oss/fern-mycelium)](https://goreportcard.com/report/github.com/guidewire-oss/fern-mycelium)
[![codecov](https://codecov.io/gh/guidewire-oss/fern-mycelium/branch/main/graph/badge.svg)](https://codecov.io/gh/guidewire-oss/fern-mycelium)
[![License](https://img.shields.io/github/license/guidewire-oss/fern-mycelium.svg)](LICENSE)
[![Release](https://img.shields.io/github/release/guidewire-oss/fern-mycelium.svg)](https://github.com/guidewire-oss/fern-mycelium/releases)
[![OpenSSF Scorecard](https://api.securityscorecards.dev/projects/github.com/guidewire-oss/fern-mycelium/badge)](https://securityscorecards.dev/viewer/?uri=github.com/guidewire-oss/fern-mycelium)

**The intelligent context layer beneath your test ecosystem.**

fern-mycelium is an open-source, extensible context engine that augments your test reporting system with rich insights, test intelligence, and AI agent integration via the **Model Context Protocol (MCP)**.

It sits on top of [fern-reporter](https://github.com/guidewire-oss/fern-reporter) and collects structured, high-fidelity context from test executions—initially from **Ginkgo** (Go) and **JUnit** (Java) via compatible adapters—and serves it through GraphQL and RESTful APIs to power analytics dashboards, developer feedback loops, and autonomous test agents.

---

## 🌱 Why "Mycelium"?

In nature, mycelium is the underground neural network of fungal threads that enables communication and resource sharing between plants. Similarly, **fern-mycelium** is the *substrate of test intelligence* that:

- Connects test executions with context
- Enables analysis, pattern detection, and historical insight
- Powers agents that observe, learn, and assist

---

## 🚀 Project Goals

1. **Expose test execution context via MCP**
   - Normalize data from Ginkgo, JUnit, and other adapters
   - Serve historical and real-time query interfaces

2. **Provide a foundation for intelligent agents**
   - Agents like *Test Coach*, *Postmortem Generator*, and *Flaky Detector*

3. **Drive actionable test analytics**
   - Flakiness scores, latency trends, and test evolution metrics

4. **Keep it open and adaptable**
   - Not bound to Ginkgo alone; built for plugin-style test source integration

---

## 🔹 Initial Capabilities

- [x] Historical context tracking for test runs
- [x] Flaky test identification framework (score calculation)
- [x] Latency and performance metrics
- [x] MCP-compatible GraphQL and REST endpoints
- [ ] Basic agents (feedback suggester, postmortem generator)
- [ ] Dashboard UI (Fern-UI extension planned)

---

## 🔄 Roadmap

| Phase | Focus | Description |
|-------|-------|-------------|
| **1** | Foundation | Schema extensions, MCP APIs, Ginkgo+JUnit adapters |
| **2** | Analytics | Flakiness scores, latency trends, test confidence metrics |
| **3** | Agents | Test Coach, Postmortem Generator, Prioritizer |
| **4** | Dev Experience | Slack/GitHub feedback bots, VSCode plugins |
| **5** | Extensibility | Plugin system, Mycelium SDK, Agent templates |

---

## 🧠 Planned Agent Capabilities

fern-mycelium is designed to power AI agents that help developers and QA teams reason about their tests and systems with minimal manual intervention. These agents are built to be pluggable and context-aware through the Model Context Protocol.

| Agent | Purpose |
|-------|---------|
| **Test Coach** | Reviews historical test data to suggest improvements, isolate brittle specs, and guide refactoring. |
| **Postmortem Generator** | Automatically drafts failure reports and incident summaries based on test runs and historical context. |
| **Predictive Prioritizer** | Reorders test execution based on failure likelihood, impact, and recent code changes. |
| **Flakiness Detector** | Flags intermittent test behavior, surfaces patterns, and scores test reliability. |
| **Feedback Assistant** | Leaves contextual PR comments or Slack messages when tests fail, enriched with test history and runtime conditions. |
| **QA Coach** *(future)* | Tracks test coverage quality trends, team-level reliability, and provides feedback loops for test effectiveness. |

These agents will be accessible through APIs and optionally embedded into tools like GitHub Actions, CI/CD pipelines, or IDE extensions.

---

## 📊 Architecture (High Level)

```text
                    ┌──────────────────────────────┐
                    │     Test Suites (CI/CD)      │
                    │ Ginkgo | JUnit | Pytest (TBD)│
                    └──────────────────────────────┘
                                 ↓
                       [ fern-ginkgo-client ]
                       [ fern-junit-adapter ]
                                 ↓
                    ┌──────────────────────────────┐
                    │      fern-reporter DB        │
                    └──────────────────────────────┘
                                 ↓
                    ┌──────────────────────────────┐
                    │       fern-mycelium API      │
                    │  - GraphQL/REST (MCP)        │
                    └──────────────────────────────┘
                                 ↓
      ┌──────────────────────────────┬──────────────────────────────┐
      │    Fern-UI Dashboards        │      Autonomous Test Agents  │
      │  (Latency, Flake Maps, etc) │  (Test Coach, Postmortem AI) │
      └──────────────────────────────┴──────────────────────────────┘
```

---

## 🚧 Under Construction

This project is in **active early-stage development**. We're currently:

- Validating schema patterns across test types
- Exposing test context over MCP
- Building the first generation of agents and queries

We welcome ideas, feedback, and early contributors—whether you’re working with Ginkgo, JUnit, or other test frameworks.

---

## 🛌 Contributing

- See [CONTRIBUTING.md](./CONTRIBUTING.md) (coming soon)
- Open discussions in Issues
- Suggest test adapters or agents you'd like to build!

---

## 🌐 License

Apache 2.0

---

## 🛍️ Questions?

- Want to integrate a new test adapter?
- Building a custom agent?
- Interested in shaping the Model Context Protocol?

**Start a discussion or open an issue — we’d love to hear from you.**
