
# BAOBAB

> **Strong Roots. Inspired Growth.**

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
![Python](https://img.shields.io/badge/Python-3.14-blue)
![Django](https://img.shields.io/badge/Django-6.0-092E20)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-17-336791)
![Docker](https://img.shields.io/badge/Docker-Ready-2496ED)
![GitHub Actions](https://img.shields.io/badge/CI-GitHub%20Actions-2088FF)
![Dev Container](https://img.shields.io/badge/DevContainer-Supported-0db7ed)
![GitHub Codespaces](https://img.shields.io/badge/GitHub-Codespaces-181717)
![Status](https://img.shields.io/badge/Status-Active-success)

**An Enterprise Software Platform for Building Secure, Scalable, Multi-Tenant and AI-Enabled Digital Solutions.**

---

# Executive Summary

BAOBAB is an enterprise-grade, open-source software platform engineered to accelerate the development and deployment of modern digital solutions. It provides a secure, scalable, and maintainable foundation for organisations seeking to build intelligent, cloud-native applications capable of supporting complex business operations.

Unlike traditional monolithic applications, BAOBAB embraces a service-oriented, polyglot architecture where each service is implemented using the technology best suited to its purpose. The platform integrates backend services, web and mobile applications, artificial intelligence capabilities, background workers, and shared platform components into a cohesive ecosystem that promotes flexibility, scalability, and long-term maintainability.

Built upon modern software engineering principles, BAOBAB incorporates containerised development environments, Infrastructure as Code, continuous integration and delivery pipelines, enterprise security practices, and comprehensive documentation. These capabilities enable development teams to collaborate effectively while delivering reliable, high-quality software.

Whether deployed for a startup, government institution, non-profit organisation, or multinational enterprise, BAOBAB is designed to adapt, scale, and evolve with changing business requirements.

---

# Project Status

> **Current Phase:** Foundation Development

BAOBAB is currently under active development as an enterprise software platform. The project is focused on establishing a robust architectural foundation, professional engineering standards, and automated development workflows before the implementation of business capabilities.

Current priorities include:

- Enterprise platform architecture
- Multi-service repository structure
- Development standards and governance
- Continuous Integration and Continuous Deployment (CI/CD)
- Security-first engineering practices
- Documentation and contributor experience

The project roadmap will continue to evolve as new capabilities and services are introduced.

---

# Vision

To become a trusted enterprise software platform that empowers organisations to design, develop, deploy, and scale secure, intelligent, and resilient digital solutions through modern engineering practices, open collaboration, and responsible innovation.

---

# Mission

To provide developers, businesses, and institutions with a comprehensive platform that combines robust architecture, enterprise-grade security, modern development workflows, artificial intelligence capabilities, and cloud-native technologies into a unified ecosystem that accelerates digital transformation while maintaining exceptional quality, reliability, and maintainability.

---

# Guiding Principles

Every architectural decision and engineering practice within BAOBAB is guided by the following principles.

## Security by Design

Security is incorporated into every layer of the platform rather than added as an afterthought. Authentication, authorisation, secure coding practices, dependency management, secrets management, and infrastructure security form an integral part of the development lifecycle.

## Modularity First

Each service should have a clearly defined responsibility and remain independently maintainable. Loose coupling and high cohesion promote scalability, extensibility, and easier long-term maintenance.

## Polyglot by Design

Different problems require different technologies. BAOBAB embraces a polyglot architecture, allowing each service to be implemented using the language, framework, and tooling most appropriate to its domain while maintaining consistent engineering standards across the platform.

## Developer Experience

Excellent software begins with an excellent development experience. BAOBAB provides reproducible development environments, automated tooling, clear documentation, and streamlined workflows that enable developers to focus on delivering value.

## Documentation as Code

Documentation is treated as a first-class deliverable. It evolves alongside the source code, ensuring that architectural decisions, implementation details, and operational procedures remain accurate, discoverable, and maintainable.

## Automation Wherever Practical

Repetitive activities should be automated wherever possible. Continuous Integration, automated testing, dependency management, security scanning, and deployment pipelines help improve consistency, reliability, and development velocity.

## Scalability from Day One

Architectural decisions are made with future growth in mind. The platform is designed to support increasing workloads, additional services, new deployment environments, and expanding development teams without requiring fundamental redesign.

## Quality Before Speed

Every contribution should satisfy established coding, testing, documentation, and security standards before being integrated into the platform. Long-term maintainability always takes precedence over short-term convenience.

---

# Key Features

BAOBAB is designed as a comprehensive enterprise software platform with capabilities that support modern software development across multiple domains.

### Enterprise Architecture

- Service-oriented architecture
- Polyglot technology stack
- Multi-tenant platform design
- Domain-driven organisation
- Modular and extensible services

### Development Experience

- Docker-based development environments
- GitHub Codespaces support
- Dev Container configuration
- Automated development workflows
- Professional engineering standards
- Infrastructure as Code

### Security

- Security-first development practices
- Automated dependency scanning
- Secrets management
- Secure authentication and authorisation
- Continuous security analysis
- Secure deployment pipelines

### Artificial Intelligence

- AI service architecture
- Agent-based workflows
- Retrieval-Augmented Generation (RAG)
- Model provider abstraction
- Prompt management
- Evaluation framework

### Cloud & DevOps

- GitHub Actions CI/CD
- Containerised deployments
- AWS-ready infrastructure
- Kubernetes-ready architecture
- Automated releases
- Environment configuration management

### Platform Engineering

- Shared APIs and contracts
- Event-driven communication
- Shared schemas and utilities
- Plugin architecture
- Monitoring and observability
- Enterprise documentation

---

# Technology Stack

BAOBAB adopts a polyglot technology strategy that allows each service to leverage the technologies best suited to its responsibilities while maintaining consistency through shared engineering standards.

| Layer | Technology |
|--------|------------|
| Backend Services | Python 3.14 · Django 6.0 |
| Frontend Services | React · Next.js · TypeScript *(Planned)* |
| Mobile Applications | Flutter *(Planned)* |
| AI Services | Python · LangChain · Model Context Protocol (MCP) · OpenAI-Compatible Providers *(Planned)* |
| Background Workers | Python · Celery *(Planned)* |
| Database | PostgreSQL 17 |
| Caching | Redis *(Planned)* |
| Object Storage | Amazon S3-Compatible Storage |
| Containers | Docker |
| Container Orchestration | Kubernetes *(Future)* |
| Infrastructure | AWS |
| Infrastructure as Code | Terraform |
| Reverse Proxy | NGINX |
| Monitoring | Prometheus · Grafana · Loki · OpenTelemetry |
| Continuous Integration | GitHub Actions |
| Development Environment | Docker Dev Containers · GitHub Codespaces |
| Version Control | Git · GitHub |

---

---

# Enterprise Architecture

BAOBAB is designed as a **service-oriented enterprise software platform** that embraces a **polyglot architecture**. Rather than organising the repository around a single framework or programming language, the platform is composed of independent services, each implemented using the technologies most appropriate to its domain.

This architectural approach promotes scalability, maintainability, and flexibility while allowing development teams to evolve individual services without unnecessarily impacting the wider platform.

The platform is organised around several major architectural layers:

- **Platform Services** – Independent business and infrastructure services.
- **Shared Platform Components** – APIs, contracts, schemas, utilities, and common resources shared across services.
- **Infrastructure** – Deployment, networking, monitoring, and cloud resources.
- **Documentation** – Enterprise architecture, design decisions, operational procedures, and technical specifications.
- **Resources** – Shared assets, templates, fixtures, branding, and environment configuration.

This structure supports future expansion into additional business domains, cloud environments, mobile applications, AI services, and customer-specific extensions.

---

# High-Level Platform Architecture

```text
                                   BAOBAB Platform

 ┌─────────────────────────────────────────────────────────────────────────────┐
 │                             Client Applications                            │
 │                                                                             │
 │    Web Applications        Mobile Applications        External Systems       │
 └──────────────────────────────────┬──────────────────────────────────────────┘
                                    │
                                    ▼
 ┌─────────────────────────────────────────────────────────────────────────────┐
 │                            Platform Services                               │
 │                                                                             │
 │  Backend Service   Frontend Service   AI Service   Worker Service           │
 └──────────────────────────────────┬──────────────────────────────────────────┘
                                    │
                                    ▼
 ┌─────────────────────────────────────────────────────────────────────────────┐
 │                     Shared Platform Components                              │
 │                                                                             │
 │ APIs • Contracts • Events • Schemas • Utilities • Localisation             │
 └──────────────────────────────────┬──────────────────────────────────────────┘
                                    │
                                    ▼
 ┌─────────────────────────────────────────────────────────────────────────────┐
 │                        Platform Infrastructure                              │
 │                                                                             │
 │ PostgreSQL • Redis • Object Storage • Monitoring • AWS • Kubernetes        │
 └─────────────────────────────────────────────────────────────────────────────┘
```

---

# Repository Structure

The repository is organised around **services**, **shared platform assets**, **enterprise documentation**, and **Infrastructure as Code**. This organisation reflects BAOBAB's service-oriented architecture and supports the independent evolution of each platform capability.

```text
Baobab/
│
├── .devcontainer/
│
├── .github/
│   └── workflows/
│
├── docs/
│
├── infrastructure/
│
├── resources/
│
├── services/
│   ├── backend/
│   ├── frontend/
│   ├── mobile/
│   ├── ai/
│   └── worker/
│
├── shared/
│
├── plugins/
│
├── tests/
│
├── compose.yaml
├── Dockerfile
├── Makefile
├── pyproject.toml
├── README.md
├── LICENSE
├── CHANGELOG.md
├── .env.example
└── uv.lock
```

> **Note:** The complete repository structure, including all nested directories, is documented within the Enterprise Architecture documentation under `docs/01-enterprise-architecture/`.

---

# Repository Organisation

The repository is organised into logical areas, each serving a distinct responsibility within the platform.

| Directory | Purpose |
|-----------|---------|
| `.devcontainer/` | Docker Dev Container configuration for consistent development environments. |
| `.github/` | GitHub workflows, automation, issue templates, repository configuration, and CI/CD pipelines. |
| `docs/` | Enterprise architecture, design decisions, development guides, API documentation, operational runbooks, and technical specifications. |
| `infrastructure/` | Infrastructure as Code, cloud resources, networking, monitoring, deployment, and platform automation. |
| `resources/` | Shared project resources including branding, templates, fixtures, sample data, certificates, backups, and environment assets. |
| `services/` | Independent platform services implementing business capabilities using the most appropriate technologies. |
| `shared/` | APIs, contracts, schemas, shared libraries, localisation resources, events, utilities, and common platform assets. |
| `plugins/` | Optional platform extensions, customer integrations, reporting modules, and experimental features. |
| `tests/` | Cross-platform testing organised by testing discipline rather than by service implementation. |

---

# Platform Services

Each platform service is developed independently while conforming to common engineering standards, security requirements, and operational practices.

## Backend Service

The Backend Service provides the core business capabilities of the platform. It is responsible for domain logic, authentication, authorisation, APIs, tenant management, data persistence, and administrative functionality.

**Primary Technologies**

- Python 3.14
- Django 6.0
- PostgreSQL 17
- Wagtail CMS

---

## Frontend Service

The Frontend Service delivers modern web user interfaces and client-side experiences. It consumes platform APIs and presents responsive, accessible, and intuitive interfaces for end users.

**Primary Technologies**

- React *(Planned)*
- Next.js *(Planned)*
- TypeScript *(Planned)*

---

## Mobile Service

The Mobile Service provides native cross-platform applications for Android, iOS, desktop, and web clients, enabling users to interact with BAOBAB from a wide range of devices.

**Primary Technologies**

- Flutter *(Planned)*

---

## AI Service

The AI Service encapsulates artificial intelligence capabilities, including large language model integration, Retrieval-Augmented Generation (RAG), agent orchestration, prompt management, embeddings, and evaluation pipelines.

The service is designed to remain provider-agnostic, enabling organisations to integrate commercial or open-source models according to their operational and regulatory requirements.

---

## Worker Service

The Worker Service executes asynchronous and scheduled workloads, including background processing, event handling, notifications, long-running tasks, and platform automation.

---

# Shared Platform Components

Rather than duplicating functionality across services, BAOBAB centralises reusable assets within the `shared/` directory.

Shared components include:

- API specifications
- OpenAPI definitions
- GraphQL schemas
- Shared request and response contracts
- Event definitions
- Validation schemas
- Shared data types
- Platform constants
- Common utilities
- Localisation resources
- Documentation assets

This approach promotes consistency, reduces duplication, and enables independent services to evolve while maintaining well-defined contracts.

---

# Extensibility

BAOBAB is designed as an extensible platform.

The `plugins/` directory provides a structured mechanism for integrating customer-specific extensions, third-party connectors, reporting modules, and experimental functionality without modifying the platform's core services.

This approach enables organisations to customise the platform while preserving a clean separation between core functionality and optional extensions.
---

# Getting Started

This section provides the information required to set up a local development environment and begin contributing to BAOBAB.

The recommended development experience uses **Docker Dev Containers** or **GitHub Codespaces**, ensuring every contributor works in a consistent, reproducible environment.

---

# Prerequisites

Before working with BAOBAB, ensure the following software is available.

| Software | Recommended Version |
|-----------|--------------------|
| Git | Latest Stable |
| Docker Desktop | Latest Stable |
| Visual Studio Code | Latest Stable |
| Dev Containers Extension | Latest Stable |
| GitHub CLI *(Optional)* | Latest Stable |

For contributors working directly on the Backend Service:

| Software | Version |
|----------|---------|
| Python | 3.14 |
| uv | Latest |
| PostgreSQL | 17 |

---

# Clone the Repository

Clone the repository using Git.

```bash
git clone https://github.com/<organisation>/baobab.git

cd baobab
```

---

# Development Environment

BAOBAB supports two recommended development environments.

## Option 1 — Docker Dev Containers (Recommended)

Open the repository using Visual Studio Code.

When prompted:

```
Reopen in Container
```

The development container automatically provisions:

- Python
- PostgreSQL client tools
- Development dependencies
- Linters
- Formatters
- Testing tools
- Git utilities

No manual configuration should be required.

---

## Option 2 — GitHub Codespaces

Open the repository on GitHub.

Select:

```
Code

↓

Codespaces

↓

Create Codespace
```

The environment is automatically configured using the same Dev Container specification used for local development.

---

# Python Dependencies

BAOBAB uses **uv** for dependency management.

Install project dependencies:

```bash
uv sync
```

To update dependencies:

```bash
uv lock
```

---

# Environment Configuration

Copy the example environment file.

```bash
cp .env.example .env
```

Update the environment variables according to your local development environment.

---

# Start Platform Services

Start supporting services using Docker Compose.

```bash
docker compose up
```

or

```bash
docker compose up --build
```

---

# Backend Service

Navigate to the Backend Service.

```bash
cd services/backend
```

Run database migrations.

```bash
python manage.py migrate
```

Create an administrator account.

```bash
python manage.py createsuperuser
```

Start the development server.

```bash
python manage.py runserver
```

The backend will normally be available at:

```
http://127.0.0.1:8000
```

---

# Frontend Service

The Frontend Service will eventually run independently.

Example:

```bash
cd services/frontend

npm install

npm run dev
```

---

# AI Service

Example local development.

```bash
cd services/ai

uv sync

python app/main.py
```

---

# Running Tests

Execute all automated tests.

```bash
pytest
```

Run only unit tests.

```bash
pytest tests/unit
```

Run integration tests.

```bash
pytest tests/integration
```

Run security tests.

```bash
pytest tests/security
```

Run performance tests.

```bash
pytest tests/performance
```

---

# Code Quality

Before submitting any Pull Request, run the project quality checks.

Example:

```bash
ruff check

ruff format

pytest
```

Additional quality checks may execute automatically through GitHub Actions.

---

# Development Workflow

Every change to BAOBAB follows a structured engineering workflow designed to maintain quality, traceability, and collaboration.

```text
Issue

   │

   ▼

Create Feature Branch

   │

   ▼

Implement Feature

   │

   ▼

Run Local Tests

   │

   ▼

Update Documentation

   │

   ▼

Commit Changes

   │

   ▼

Open Pull Request

   │

   ▼

Automated CI Validation

   │

   ▼

Code Review

   │

   ▼

Merge into Main
```

---

# Branching Strategy

BAOBAB follows a feature branch workflow.

| Branch | Purpose |
|---------|----------|
| `main` | Production-ready code |
| `develop` *(Optional)* | Integration branch for ongoing development |
| `feature/*` | New functionality |
| `bugfix/*` | Bug fixes |
| `hotfix/*` | Critical production fixes |
| `release/*` | Release preparation |

---

# Commit Message Convention

Meaningful commit messages improve traceability and project history.

Recommended format:

```text
type(scope): short description
```

Examples:

```text
feat(auth): add tenant authentication

fix(api): resolve pagination issue

docs(readme): update installation guide

test(ai): improve evaluation coverage

refactor(worker): simplify scheduling logic
```

---

# Continuous Integration

Every Pull Request triggers automated quality checks, including:

- Code formatting
- Static analysis
- Unit testing
- Integration testing
- Security scanning
- Dependency auditing
- Build validation

A Pull Request should not be merged until all required checks have completed successfully.

---

# Contributor Expectations

Every contributor is expected to:

- Follow the Coding Standards.
- Follow the Documentation Standards.
- Follow the Testing Standards.
- Keep documentation up to date.
- Write meaningful commit messages.
- Maintain high test coverage.
- Submit focused, reviewable Pull Requests.
- Respect the Code of Conduct.

These expectations are described in greater detail within the project's governance documentation.

---

# Documentation

BAOBAB is developed using a **documentation-first** approach. Documentation is maintained alongside the source code to ensure architectural decisions, engineering practices, and operational procedures remain accurate, discoverable, and maintainable.

The `/docs` directory serves as the project's primary technical knowledge base.

## Documentation Hub

| Document | Description |
|----------|-------------|
| `README.md` | Project overview and getting started guide. |
| `CONTRIBUTING.md` | Contribution process, pull requests, branching strategy, and development workflow. |
| `CODE_OF_CONDUCT.md` | Community standards and expected behaviour. |
| `SECURITY.md` | Responsible vulnerability disclosure process and security policy. |
| `CHANGELOG.md` | Project release history following *Keep a Changelog* principles. |
| `docs/Coding-Standards.md` | Coding conventions and engineering practices. |
| `docs/Testing-Standards.md` | Testing philosophy, quality expectations, and coverage requirements. |
| `docs/Documentation-Standards.md` | Documentation style guide and documentation governance. |

## Enterprise Documentation

The `docs/` directory is organised into specialised areas covering the entire platform lifecycle.

```text
docs/
├── 00-introduction/
├── 01-enterprise-architecture/
├── 02-domain-model/
├── 03-development/
├── 04-devops/
├── 05-data/
├── 06-ai/
├── 07-security/
├── 08-api/
├── 09-testing/
├── diagrams/
├── specifications/
├── adr/
└── runbooks/
```

Each section contains documentation specific to a functional area of the platform, making it easier for contributors, architects, operators, and stakeholders to locate relevant information.

---

# Roadmap

BAOBAB follows an incremental roadmap focused on delivering a secure, scalable, and extensible enterprise platform.

## Phase 1 — Foundation

- Repository governance
- Enterprise architecture
- Development standards
- CI/CD pipelines
- Infrastructure as Code
- Documentation framework

## Phase 2 — Core Platform

- Multi-tenant framework
- Identity and access management
- Administrative portal
- Shared platform services
- Initial API implementation

## Phase 3 — Digital Services

- Frontend application
- Mobile application
- Background worker service
- Notifications
- Reporting

## Phase 4 — Artificial Intelligence

- AI platform services
- Agent orchestration
- Retrieval-Augmented Generation (RAG)
- Prompt management
- Evaluation framework
- Intelligent automation

## Phase 5 — Enterprise Operations

- Monitoring and observability
- Performance optimisation
- Disaster recovery
- High availability
- Multi-region deployment
- Enterprise integrations

The roadmap will evolve as the project matures and community contributions shape future priorities.

---

# Security

Security is a foundational principle of BAOBAB and is integrated throughout the software development lifecycle.

If you discover a security vulnerability, **please do not disclose it publicly**.

Instead, follow the responsible disclosure process described in `SECURITY.md`.

Security practices include:

- Secure software development lifecycle (SSDLC)
- Automated dependency auditing
- Static Application Security Testing (SAST)
- Secrets management
- Infrastructure security
- Continuous vulnerability assessment
- Secure deployment pipelines

Security is a shared responsibility, and every contributor is expected to follow the project's security policies and engineering standards.

---

# Contributing

Community contributions are welcome and encouraged.

Whether you are fixing bugs, improving documentation, developing new platform capabilities, or proposing architectural improvements, your contributions help strengthen BAOBAB.

Before contributing, please read:

- `CONTRIBUTING.md`
- `CODE_OF_CONDUCT.md`
- `docs/Coding-Standards.md`
- `docs/Testing-Standards.md`
- `docs/Documentation-Standards.md`

All Pull Requests are reviewed before being merged into the project.

We encourage contributions that are:

- Well documented
- Well tested
- Focused and easy to review
- Consistent with the project's architectural principles
- Respectful of established engineering standards

---

# Community

BAOBAB aims to cultivate a professional, respectful, and inclusive engineering community.

All contributors are expected to adhere to the project's `CODE_OF_CONDUCT.md`.

We believe successful open-source software is built not only through excellent engineering, but also through respectful collaboration, constructive communication, and shared responsibility.

---

# License

Copyright © 2026 BAOBAB Contributors.

Licensed under the **Apache License, Version 2.0** (the "License").

You may not use this project except in compliance with the License.

You may obtain a copy of the License at:

```
http://www.apache.org/licenses/LICENSE-2.0
```

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an **"AS IS" BASIS**, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.

See the `LICENSE` file for the complete license text.

---

# Acknowledgements

BAOBAB draws inspiration from the African baobab tree—an enduring symbol of resilience, wisdom, community, and sustainable growth. These qualities reflect the platform's commitment to building software that is robust, adaptable, and designed for the long term.

The project is made possible by the global open-source community, whose collective contributions to programming languages, frameworks, tools, and standards continue to advance modern software engineering.

We also recognise the importance of open collaboration, responsible innovation, and knowledge sharing in shaping technology that delivers meaningful impact for organisations and communities alike.

---

# Support

If you have questions, suggestions, or wish to report issues, please use the project's GitHub Issues and Discussions once they become available.

Feature requests, bug reports, and architectural proposals are always welcome.

---

# Maintainers

BAOBAB is maintained by its core maintainers together with the broader community of contributors.

Project governance and contribution processes are documented in the repository's governance documentation.

---

<div align="center">

## Strong Roots. Inspired Growth.

**Building the next generation of enterprise software—securely, collaboratively, and sustainably.**

</div>
