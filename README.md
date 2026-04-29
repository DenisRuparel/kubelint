# KubeLint

## Kubernetes Package Manager with CUE Validation

KubeLint is a production-grade Kubernetes package manager and validation CLI inspired by Timoni and powered by CUE.

It helps teams follow the workflow:

```text
Write Once → Validate Strongly → Deploy Everywhere
```

Instead of maintaining separate Kubernetes manifests for every environment, KubeLint allows you to define reusable templates, environment-specific values, and strong validation rules using CUE schemas.

---

# Why KubeLint?

Managing Kubernetes manifests across multiple environments becomes difficult because:

* YAML duplication grows quickly
* Environment-specific changes create drift
* Validation is weak with plain YAML
* Production mistakes are expensive
* Security issues are often missed before deployment

KubeLint solves this by providing:

* reusable Kubernetes templates
* environment-based values management
* CUE-powered schema validation
* security and policy checks
* diff preview before deployment
* safe deployment workflow
* CI/CD integration support
* production-grade release system

---

# Core Idea

## Project Example

```text
templates/
  deployment.yaml
  service.yaml

values/
  default.yaml
  dev.yaml
  staging.yaml
  prod.yaml

schemas/
  deployment.cue
  service.cue
```

Run:

```bash
kubelint build --env prod
kubelint validate --env prod
kubelint deploy --env prod
```

---

# Supported Commands

```bash
kubelint init
kubelint build
kubelint validate
kubelint diff
kubelint deploy
kubelint doctor
kubelint version
```

Future Commands:

```bash
kubelint rollback
kubelint package
kubelint drift-detect
kubelint test
```

---

# Final Project Structure

```text
kubelint/
│
├── cmd/
│   ├── root.go
│   ├── init.go
│   ├── build.go
│   ├── validate.go
│   ├── diff.go
│   ├── deploy.go
│   ├── doctor.go
│   └── version.go
│
├── internal/
│   ├── config/
│   ├── parser/
│   ├── templater/
│   ├── validator/
│   ├── cueengine/
│   ├── reporter/
│   ├── scorer/
│   ├── kubeclient/
│   ├── deployer/
│   ├── differ/
│   └── utils/
│
├── templates/
├── values/
├── schemas/
├── output/
├── examples/
├── testdata/
├── docs/
├── scripts/
│   └── install.sh
│
├── .github/
│   └── workflows/
│
├── main.go
├── go.mod
├── go.sum
├── README.md
├── CHANGELOG.md
├── LICENSE
└── kubelint.yaml
```

---

# Release Strategy (Version by Version)

We release KubeLint phase by phase like real production tools.

This allows users to download any version they want and helps maintain a professional open-source workflow.

---

## v0.1.0 — CLI Foundation

### Features

* kubelint init
* kubelint version
* kubelint doctor
* Cobra CLI foundation
* project scaffold generation

---

## v0.2.0 — Template Rendering Engine

### Features

* kubelint build
* template rendering
* multi-template support
* multi-document YAML generation

---

## v0.3.0 — Values Engine

### Features

* environment values support
* default/dev/staging/prod values
* CLI overrides using --set
* merge engine

---

## v0.4.0 — CUE Validation Engine

### Features

* kubelint validate
* CUE schema validation
* validation reporting
* required fields and type checks

This is the major milestone release.

---

## v0.5.0 — Security + Policy Engine

### Features

* security rule engine
* severity levels
* scoring engine
* Kubernetes best-practice validation

---

## v0.6.0 — Diff Engine

### Features

* kubelint diff
* cluster state comparison
* change preview before deployment

---

## v0.7.0 — Deployment Engine

### Features

* kubelint deploy
* dry-run support
* namespace creation
* deployment safety checks

---

## v1.0.0 — Stable Production Release

### Features

* CI/CD integration
* GitHub Releases
* install.sh
* tests
* release automation
* documentation
* production-grade readiness

---

# Installation

Users can install either:

## Option 1 — Latest Release

OR

## Option 2 — Specific Version

Both are supported.

---

# Option 1 — Download Latest Release

## Step 1

Open:

```text
GitHub Repository → Releases → Latest Release
```

---

## Step 2

Choose your platform:

### Linux

```text
kubelint-linux-amd64
```

### macOS

```text
kubelint-darwin-amd64
```

### Windows

```text
kubelint-windows-amd64.exe
```

---

## Step 3

Download the correct binary for your OS.

---

## Step 4 — Make Executable (Linux/macOS)

```bash
chmod +x kubelint-linux-amd64
sudo mv kubelint-linux-amd64 /usr/local/bin/kubelint
```

---

## Step 5 — Verify Installation

```bash
kubelint version
```

Expected output:

```text
KubeLint version v1.0.0
```

---

# Option 2 — Download Specific Version

Example: Install v0.4.0 only

---

## Step 1

Open:

```text
GitHub Repository → Releases
```

---

## Step 2

Select the version you want:

```text
v0.1.0
v0.2.0
v0.3.0
v0.4.0
v0.5.0
...
```

---

## Step 3

Choose your platform binary:

### Linux

```text
kubelint-linux-amd64
```

### macOS

```text
kubelint-darwin-amd64
```

### Windows

```text
kubelint-windows-amd64.exe
```

---

## Step 4

Install using the same method as latest release.

---

# Quick Install Script

For latest stable release:

```bash
curl -fsSL <install-script-url> | bash
```

This will:

* detect OS automatically
* download latest binary
* install KubeLint
* verify installation

---

# Example Workflow

```bash
kubelint init my-app
cd my-app

kubelint build --env prod
kubelint validate --env prod
kubelint diff --env prod
kubelint deploy --env prod
```

---

# CI/CD Example

```yaml
- name: Validate Kubernetes Manifests
  run: kubelint validate --env prod

- name: Deploy to Kubernetes
  run: kubelint deploy --env prod
```

Perfect for GitHub Actions pipelines.

---

# Development Roadmap

```text
CLI
→ Templates
→ Values
→ CUE Validation
→ Security Rules
→ Diff
→ Deploy
→ Release
```

Stable architecture + incremental releases.

---

# Why This Project Is Strong

KubeLint is not just another YAML linter.

It combines:

* Timoni-style package management
* CUE-based validation
* security and policy enforcement
* safe deployment workflows
* CI/CD readiness
* versioned open-source releases

This makes it a strong DevOps, Platform Engineering, and Final Year Project with real production value.

---

# Contributing

Contributions, suggestions, and improvements are welcome.

Please open an issue or submit a pull request.

---

# License

MIT License
