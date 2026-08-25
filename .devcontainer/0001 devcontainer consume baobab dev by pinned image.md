diff --git a/.devcontainer/.env.example b/.devcontainer/.env.example
deleted file mode 100644
index 6d2d8ef..0000000
--- a/.devcontainer/.env.example
+++ /dev/null
@@ -1,30 +0,0 @@
-###############################################################################
-# BAOBAB Enterprise Platform
-#
-# Development Environment
-###############################################################################
-
-# -----------------------------------------------------------------------------
-# Environment
-# -----------------------------------------------------------------------------
-
-BAOBAB_ENV=development
-BAOBAB_NAME=baobab
-
-# -----------------------------------------------------------------------------
-# Workspace
-# -----------------------------------------------------------------------------
-
-BAOBAB_HOME=/workspaces/baobab
-
-# -----------------------------------------------------------------------------
-# Logging
-# -----------------------------------------------------------------------------
-
-BAOBAB_LOG_LEVEL=INFO
-
-# -----------------------------------------------------------------------------
-# Verification
-# -----------------------------------------------------------------------------
-
-BAOBAB_VERIFY=true
diff --git a/.devcontainer/README.md b/.devcontainer/README.md
index 2e98094..38f822a 100644
--- a/.devcontainer/README.md
+++ b/.devcontainer/README.md
@@ -2,9 +2,23 @@

 ## Purpose

-The BAOBAB Development Container provides a consistent, reproducible, and cloud-ready development environment for all contributors.
+The BAOBAB Development Container provides a consistent, reproducible, cloud-ready development environment for every contributor. It is used by GitHub Codespaces and Visual Studio Code Dev Containers to ensure every developer works in an identical environment with minimal setup.

-It is used by GitHub Codespaces and Visual Studio Code Dev Containers to ensure every developer works in an identical environment with minimal setup.
+---
+
+# Architecture
+
+Unlike a repository that builds its own development image, BAOBAB consumes a **published, versioned image** maintained in a dedicated repository:
+
+| Component | Location |
+| --- | --- |
+| Image source | <https://github.com/nabhold/baobab-dev> |
+| Published image | `ghcr.io/nabhold/baobab-dev` |
+| Current pin | `1.0.0` (see `devcontainer.json`) |
+
+`baobab-dev` is a deterministic, multi-stage, checksum-verified image build with its own CI pipeline, version-resolution lockfile, and documentation. It already bakes in the full toolchain — Python 3.14, `uv`, PostgreSQL/Redis client tools, Docker CLI, GitHub CLI, Node.js LTS, and Flutter/Dart — at **image build time**, deterministically, once. This repository's `.devcontainer/` therefore does **not** build a local Dockerfile or reinstall language runtimes; it only references the published image and layers on the small amount of setup that is specific to *this* repository.
+
+Reference documentation for the image itself (toolchain contents, helper commands, verification, version policy) lives in the `baobab-dev` repository's own `docs/` tree — see in particular `docs/reference/helper-commands.md` and `docs/introduction/included.md`.

 ---

@@ -12,17 +26,9 @@ It is used by GitHub Codespaces and Visual Studio Code Dev Containers to ensure

 ```text
 .devcontainer/
-├── .env.example
 ├── devcontainer.json
-├── README.md
-└── docker/
-    ├── Dockerfile
-    └── scripts/
-        ├── README.md
-        ├── bootstrap.sh
-        ├── post-create.sh
-        ├── run.sh
-        └── verify.sh
+├── post-create.sh
+└── README.md
 ```

 ---
@@ -31,127 +37,62 @@ It is used by GitHub Codespaces and Visual Studio Code Dev Containers to ensure

 ## `devcontainer.json`

-Defines how GitHub Codespaces and Visual Studio Code provision the BAOBAB development environment.
+Defines how GitHub Codespaces and Visual Studio Code provision the BAOBAB development environment: which published image to use, Dev Container Features, lifecycle hooks, VS Code settings/extensions, and forwarded ports.

-Responsibilities include:
+## `post-create.sh`

-* Development container configuration
-* Dev Container Features
-* Visual Studio Code settings
-* Recommended extensions
-* Forwarded development ports
-* Workspace behaviour
-
----
+A small, repository-owned script that runs *this repository's* setup on top of the image: creating `.env` from `.env.example`, and installing the root `uv` workspace's Python dependencies.

-## `docker/Dockerfile`
-
-Defines the development container image.
-
-Its responsibilities are limited to:
-
-* Selecting the base image
-* Defining image metadata
-* Configuring the container environment
-* Providing the development workspace
-
-It does **not** install application services or define runtime orchestration.
-
-Application runtime is managed separately through Docker Compose.
+This is deliberately **not** delegated to the image's own generic `baobab-post-create` dependency-install logic: that logic installs Python dependencies with **Poetry** whenever it finds a root `pyproject.toml`. BAOBAB standardises on **`uv`** (see `uv.lock`, `Makefile`, and the "UV WORKSPACE READINESS ROADMAP" in `pyproject.toml`) — Poetry is not used in this repository. Everything else the image's `baobab-post-create` provides (git identity priming, toolchain verification) is generic and safe, and is still used directly via `onCreateCommand`.

 ---

-## `docker/scripts/`
-
-Contains shell scripts used to provision and configure the development environment.
-
-Typical responsibilities include:
+# Lifecycle

-* Installing development tooling
-* Configuring shell environments
-* Bootstrapping the workspace
-* Performing post-build configuration
-* Automating repetitive setup tasks
+| Hook | Runs | Does |
+| --- | --- | --- |
+| `onCreateCommand` | Once, at creation | `baobab-post-create --stage=on-create` (image-provided: git identity priming, `baobab-verify`) |
+| `updateContentCommand` | Every content update — the stage GitHub Codespaces **prebuilds** snapshot | `post-create.sh` (`.env` bootstrap, `uv sync`) |
+| `postCreateCommand` | Once, after workspace mount | `post-create.sh` again — a safety net in case a prebuild was stale |
+| `postStartCommand` | Every container start | `baobab-summary` (image-provided environment banner) |

-As the platform evolves, this directory allows provisioning logic to remain modular instead of accumulating in the Dockerfile.
+All steps are idempotent and safe to re-run.

 ---

-# Development Container Features
+# Helper Commands

-BAOBAB uses Dev Container Features wherever practical to install and maintain common development tooling.
+The image installs the following commands directly onto `PATH`; they are available in every container without any additional setup:

-Current Features include:
+| Command | Purpose |
+| --- | --- |
+| `baobab-verify` | Validate the toolchain is installed and correctly configured |
+| `baobab-summary` | Print installed tool versions and environment metadata |
+| `baobab-post-create` | Image-provided lifecycle hook (`--stage=on-create` / `post-create` / `update-content`) |
+| `baobab-bootstrap` | Full provisioning entry point for non-Codespaces use (bare `docker run`, self-hosted runners) |

-| Feature    | Purpose                                                    |
-| ---------- | ---------------------------------------------------------- |
-| Python     | Backend development using Django, FastAPI and `uv`         |
-| Node.js    | Frontend development with Next.js                          |
-| Java       | Android tooling required for Flutter development           |
-| Docker CLI | Build and manage local containers                          |
-| GitHub CLI | Interact with GitHub from within the development container |
-
-Additional Features may be introduced in future sprints as new platform capabilities are implemented.
-
----
-
-# Relationship to the Platform
-
-```text
-GitHub Codespaces
-        │
-        ▼
-Development Container
-        │
-        ▼
-Docker Compose
-        │
-        ▼
-Platform Services
-        ├── Backend (Django)
-        ├── AI (FastAPI)
-        ├── Frontend (Next.js)
-        ├── Mobile (Flutter)
-        ├── Worker
-        ├── PostgreSQL
-        ├── Redis
-        └── Mailpit
-```
-
-The Development Container provides the development environment.
-
-Docker Compose orchestrates the local platform.
-
-The platform services provide BAOBAB's application functionality.
+Run `baobab-verify` any time you want to confirm the environment is healthy, and `baobab-summary` when reporting an issue.

 ---

-# Rebuilding the Development Container
-
-Whenever `devcontainer.json`, `docker/Dockerfile`, or scripts under `docker/scripts/` are modified, rebuild the development container.
+# Rebuilding / Updating the Development Container

-### Visual Studio Code
+Because the container now references a published image rather than building one locally, "rebuilding" means either:

-1. Open the Command Palette.
-2. Select **Dev Containers: Rebuild Container**.
+1. **Picking up the current pin** — Command Palette → **Dev Containers: Rebuild Container** (or recreate the Codespace) re-pulls `ghcr.io/nabhold/baobab-dev:1.0.0` if it has changed.
+2. **Moving to a new image release** — bump the tag in `devcontainer.json`'s `image` field via a reviewed pull request once a new `baobab-dev` release is published, then rebuild as above.

-### GitHub Codespaces
-
-Rebuild or recreate the Codespace from the repository.
+The image tag is deliberately pinned (never `:latest`) so every contributor and every Codespace resolve to byte-identical tooling.

 ---

 # Engineering Principles

-The BAOBAB Development Container follows these principles:
-
 * Reproducible development environments
-* Infrastructure as Code
-* Minimal manual configuration
-* Cross-platform compatibility
-* Security by default
-* Incremental evolution
-* Separation of concerns
+* Single, centrally governed source of truth for the toolchain
+* Minimal repository-local provisioning logic
+* Separation of concerns between the platform image and this repository's own setup
+* Security by default (the image's own build pipeline owns checksum verification and CVE tracking for the toolchain)

 ---

@@ -162,11 +103,10 @@ The BAOBAB Development Container follows these principles:
 * `docs/governance/engineering-handbook.md`
 * `docs/governance/repository-governance.md`
 * `docs/governance/decision-record-process.md`
+* [`nabhold/baobab-dev`](https://github.com/nabhold/baobab-dev) — image source, version policy, and full toolchain reference

 ---

 # Ownership

-This configuration is maintained by the BAOBAB Engineering Team.
-
-Changes to the development environment should be reviewed through the project's pull request process and, where appropriate, documented through an Architecture Decision Record (ADR).
+This configuration is maintained by the BAOBAB Engineering Team. Changes to the development environment should be reviewed through the project's pull request process and, where appropriate, documented through an Architecture Decision Record (ADR).
diff --git a/.devcontainer/devcontainer.json b/.devcontainer/devcontainer.json
index 1c6fd1d..596071c 100644
--- a/.devcontainer/devcontainer.json
+++ b/.devcontainer/devcontainer.json
@@ -1,19 +1,22 @@
 {
   "name": "BAOBAB Development Environment",

-  "build": {
-    "context": "..",
-    "dockerfile": "docker/Dockerfile",
-    "target": "development",
-    "args": {
-      "UBUNTU_VERSION": "26.04",
-      "BAOBAB_ENV": "development"
-    }
-  },
+  // -----------------------------------------------------------------------
+  // Image
+  // -----------------------------------------------------------------------
+  // Consumes the published, versioned BAOBAB development image directly
+  // rather than building a local Dockerfile on top of it. The image
+  // (github.com/nabhold/baobab-dev) already bakes in the full toolchain
+  // deterministically at build time — Python 3.14, uv, PostgreSQL/Redis
+  // clients, Docker CLI, GitHub CLI, Node.js LTS, Flutter/Dart — so there
+  // is currently nothing for a repo-local build layer to add. Pinned to an
+  // explicit release tag (never :latest) so every contributor and every
+  // Codespace resolves to byte-identical tooling; bump deliberately via PR
+  // when baobab-dev cuts a new release.
+  "image": "ghcr.io/nabhold/baobab-dev:1.0.0",

   "features": {
-    "ghcr.io/devcontainers/features/docker-outside-of-docker:1": {},
-    "ghcr.io/devcontainers/features/github-cli:1": {}
+    "ghcr.io/devcontainers/features/docker-outside-of-docker:1": {}
   },

   "remoteUser": "vscode",
@@ -29,22 +32,44 @@
   ],

   "mounts": [
-    "source=baobab-pip-cache,target=/home/vscode/.cache/pip,type=volume",
     "source=baobab-uv-cache,target=/home/vscode/.cache/uv,type=volume"
   ],

   "containerEnv": {
     "BAOBAB_ENV": "development",
-    "PYTHONUNBUFFERED": "1",
-    "PIP_DISABLE_PIP_VERSION_CHECK": "1",
     "UV_LINK_MODE": "copy"
   },

-  "postCreateCommand": "bash .devcontainer/docker/scripts/run.sh bootstrap",
-
-  "updateContentCommand": "bash .devcontainer/docker/scripts/run.sh post-create",
-
-  "postStartCommand": "bash .devcontainer/docker/scripts/summary.sh",
+  // -----------------------------------------------------------------------
+  // Lifecycle
+  // -----------------------------------------------------------------------
+  // updateContentCommand : runs on every content update, and is the stage
+  //                     GitHub Codespaces prebuilds actually snapshot.
+  //                     Runs a repo-owned script (git identity priming,
+  //                     .env bootstrap, `uv sync`, toolchain verification).
+  //                     Deliberately NOT the image's own generic
+  //                     baobab-post-create hook: that command requires a
+  //                     workspace-local config/versions.yaml +
+  //                     config/resolve.sh (a convention this repo hasn't
+  //                     adopted) and, separately, assumes Poetry for any
+  //                     repo with a root pyproject.toml — this repo uses
+  //                     uv. See .devcontainer/post-create.sh for details.
+  // postCreateCommand : safety-net re-run of the same repo-owned script on
+  //                     actual container/Codespace creation, in case a
+  //                     prebuild was stale relative to the checked-out
+  //                     commit. Every step in it is idempotent.
+  // postStartCommand  : lightweight environment banner on every start.
+  //                     baobab-summary has no workspace config/
+  //                     dependency — it falls back to the image's own
+  //                     baked-in configuration — so it is safe to call
+  //                     directly.
+  // -----------------------------------------------------------------------
+
+  "updateContentCommand": "bash .devcontainer/post-create.sh",
+
+  "postCreateCommand": "bash .devcontainer/post-create.sh",
+
+  "postStartCommand": "baobab-summary",

   "customizations": {
     "vscode": {
@@ -69,7 +94,9 @@

         "terminal.integrated.defaultProfile.linux": "bash",

+        "python.defaultInterpreterPath": ".venv/bin/python",
         "python.analysis.typeCheckingMode": "basic",
+        "python.terminal.activateEnvironment": true,

         "git.autofetch": true,

@@ -80,6 +107,7 @@
         "ms-python.python",
         "ms-python.vscode-pylance",
         "charliermarsh.ruff",
+        "tamasfe.even-better-toml",

         "dbaeumer.vscode-eslint",
         "esbenp.prettier-vscode",
@@ -109,7 +137,7 @@

   "portsAttributes": {
     "3000": {
-      "label": "Next.js Frontend",
+      "label": "Frontend (Next.js, planned)",
       "onAutoForward": "notify"
     },

diff --git a/.devcontainer/docker/Dockerfile b/.devcontainer/docker/Dockerfile
deleted file mode 100644
index 3b54f13..0000000
--- a/.devcontainer/docker/Dockerfile
+++ /dev/null
@@ -1,144 +0,0 @@
-# ==============================================================================
-# BAOBAB Enterprise Platform
-#
-# Dockerfile : Development Container
-#
-# Purpose
-# -------
-# Base image for GitHub Codespaces and VS Code Dev Containers.
-#
-# Responsibilities
-# ----------------
-# • Provide a stable Ubuntu development environment.
-# • Install common operating system packages.
-# • Configure the development container environment.
-# • Prepare the workspace.
-#
-# This image intentionally DOES NOT install language runtimes or project
-# dependencies. Those are provisioned by the BAOBAB Provisioning Framework
-# under:
-#
-#     .devcontainer/docker/scripts/
-#
-# The provisioning sequence is:
-#
-#     run.sh
-#         ├── bootstrap.sh
-#         ├── environment.sh
-#         ├── install.sh
-#         ├── verify.sh
-#         └── summary.sh
-#
-# This image is NOT intended for production deployments.
-#
-# Author  : BAOBAB Contributors
-# License : Apache License 2.0
-# ==============================================================================
-
-FROM ghcr.io/nabhold/baobab-dev:latest AS development
-
-# ------------------------------------------------------------------------------
-# Build Metadata
-# ------------------------------------------------------------------------------
-
-ARG VERSION="0.1.0-dev"
-ARG BUILD_DATE="unknown"
-ARG VCS_REF="unknown"
-ARG BAOBAB_PROVISIONER_VERSION="1.0.0"
-
-LABEL org.opencontainers.image.title="BAOBAB Development Environment" \
-      org.opencontainers.image.description="Development container for the BAOBAB Enterprise Platform." \
-      org.opencontainers.image.vendor="NabHold Group Africa" \
-      org.opencontainers.image.authors="BAOBAB Engineering Team" \
-      org.opencontainers.image.licenses="Apache-2.0" \
-      org.opencontainers.image.source="https://github.com/bjnabs/baobab" \
-      org.opencontainers.image.version="${VERSION}" \
-      org.opencontainers.image.created="${BUILD_DATE}" \
-      org.opencontainers.image.revision="${VCS_REF}"
-
-
-LABEL io.baobab.provisioner.version="${BAOBAB_PROVISIONER_VERSION}"
-
-
-
-# ------------------------------------------------------------------------------
-# Environment
-# ------------------------------------------------------------------------------
-
-ENV DEBIAN_FRONTEND=noninteractive
-ENV BAOBAB_PROVISIONER_VERSION="${BAOBAB_PROVISIONER_VERSION}"
-
-ENV LANG=en_US.UTF-8 \
-    LANGUAGE=en_US:en \
-    LC_ALL=en_US.UTF-8 \
-    TZ=UTC
-
-# ------------------------------------------------------------------------------
-# Default Shell
-# ------------------------------------------------------------------------------
-
-SHELL ["/bin/bash", "-o", "pipefail", "-c"]
-
-# ------------------------------------------------------------------------------
-# Install Base Operating System Packages
-# ------------------------------------------------------------------------------
-
-RUN apt-get update && \
-    apt-get install -y --no-install-recommends \
-        bash \
-        build-essential \
-        ca-certificates \
-        curl \
-        file \
-        git \
-        gnupg \
-        jq \
-        less \
-        locales \
-        lsb-release \
-        nano \
-        openssh-client \
-        pkg-config \
-        procps \
-        rsync \
-        software-properties-common \
-        sudo \
-        tar \
-        tree \
-        unzip \
-        vim \
-        wget \
-        xz-utils \
-        zip && \
-    locale-gen en_US.UTF-8 && \
-    apt-get clean && \
-    rm -rf /var/lib/apt/lists/*
-
-# ------------------------------------------------------------------------------
-# Create Common Development Directories
-# ------------------------------------------------------------------------------
-
-RUN mkdir -p \
-    /opt/baobab \
-    /opt/baobab/cache \
-    /opt/baobab/tools \
-    /var/log/baobab
-
-# ------------------------------------------------------------------------------
-# Workspace
-# ------------------------------------------------------------------------------
-
-WORKDIR /workspaces/baobab
-
-# ------------------------------------------------------------------------------
-# Container Environment
-# ------------------------------------------------------------------------------
-
-ENV BAOBAB_HOME=/workspaces/baobab \
-    BAOBAB_ENV=development
-
-# ------------------------------------------------------------------------------
-# Default Command
-# ------------------------------------------------------------------------------
-
-CMD ["sleep", "infinity"]
diff --git a/.devcontainer/docker/scripts/README.md b/.devcontainer/docker/scripts/README.md
deleted file mode 100644
index 6f5d191..0000000
--- a/.devcontainer/docker/scripts/README.md
+++ /dev/null
@@ -1,108 +0,0 @@
-# Development Container Scripts
-
-## Purpose
-
-This directory contains scripts used to provision, configure, and validate the BAOBAB Development Container.
-
-These scripts support GitHub Codespaces and Visual Studio Code Dev Containers by encapsulating development environment logic outside the Dockerfile.
-
-Keeping provisioning logic in dedicated scripts improves readability, maintainability, and reuse while allowing the Dockerfile to remain focused on defining the container image.
-
----
-
-# Directory Structure
-
-```text
-scripts/
-├── bootstrap.sh
-├── post-create.sh
-├── verify.sh
-└── README.md
-```
-
----
-
-# File Responsibilities
-
-## `bootstrap.sh`
-
-Initialises the development environment before project-specific configuration is applied.
-
-Typical responsibilities include:
-
-* Preparing environment variables
-* Configuring common shell settings
-* Executing shared helper functions
-* Performing lightweight prerequisite checks
-
-This script should remain generic and reusable.
-
----
-
-## `post-create.sh`
-
-Executed after the Development Container has been created.
-
-Typical responsibilities include:
-
-* Verifying installed tooling
-* Initialising repository-specific configuration
-* Preparing the development workspace
-* Running non-destructive setup tasks
-* Displaying useful onboarding information
-
-This script should be safe to execute multiple times.
-
----
-
-## `verify.sh`
-
-Validates that the Development Container has been configured correctly.
-
-Typical verification tasks include:
-
-* Python version
-* Node.js version
-* Java version
-* Docker CLI availability
-* GitHub CLI availability
-* Development tool diagnostics
-
-The script should exit with a non-zero status if a critical dependency is unavailable.
-
----
-
-# Engineering Principles
-
-The scripts in this directory should adhere to the following principles:
-
-* Single Responsibility Principle
-* Idempotent execution wherever practical
-* Small, modular, and reusable functions
-* Clear logging and meaningful error messages
-* Fail fast on unrecoverable errors
-* Avoid hard-coded paths and configuration values
-* Use environment variables where appropriate
-
----
-
-# Future Expansion
-
-As BAOBAB evolves, additional scripts may be introduced to support:
-
-* Flutter SDK configuration
-* `uv` environment preparation
-* Git hook installation
-* Local certificate management
-* Repository bootstrap automation
-* Development diagnostics
-
-Scripts should remain focused and composable rather than becoming large, multi-purpose utilities.
-
----
-
-# Ownership
-
-This directory is maintained by the BAOBAB Engineering Team.
-
-Changes to provisioning or verification scripts should be reviewed through the project's standard pull request process and kept aligned with the Development Container documentation.
diff --git a/.devcontainer/docker/scripts/bootstrap.sh b/.devcontainer/docker/scripts/bootstrap.sh
deleted file mode 100755
index 6c79c0b..0000000
--- a/.devcontainer/docker/scripts/bootstrap.sh
+++ /dev/null
@@ -1,150 +0,0 @@
-#!/usr/bin/env bash
-
-# ==============================================================================
-# BAOBAB Enterprise Platform
-#
-# Script      : bootstrap.sh
-# Purpose     : Performs complete provisioning of the BAOBAB Development
-#               Environment.
-#
-# Description :
-#   Orchestrates the complete provisioning lifecycle for the BAOBAB
-#   Development Environment.
-#
-#   Execution Order:
-#
-#       1. Install development toolchains
-#       2. Configure the environment
-#       3. Initialize the workspace
-#       4. Verify the installation
-#       5. Display the environment summary
-#
-# Author      : BAOBAB Contributors
-# License     : Apache License 2.0
-# ==============================================================================
-
-set -Eeuo pipefail
-
-SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
-
-# ------------------------------------------------------------------------------
-# Execute Script
-# ------------------------------------------------------------------------------
-
-run_step() {
-
-    local script="$1"
-
-    if [[ ! -f "$script" ]]; then
-        echo "ERROR: Missing script:"
-        echo "  $script"
-        exit 1
-    fi
-
-    echo
-    echo "=================================================================="
-    echo "Running: $(basename "$script")"
-    echo "=================================================================="
-    echo
-
-    # Execute explicitly with Bash rather than relying on executable bits.
-    bash "$script"
-}
-
-# ------------------------------------------------------------------------------
-# Install
-# ------------------------------------------------------------------------------
-
-echo
-echo "====================== INSTALL ======================"
-
-INSTALL_SCRIPTS=(
-    system.sh
-    python.sh
-    node.sh
-    java.sh
-    database.sh
-    docker.sh
-    flutter.sh
-    android_sdk.sh
-    dependencies.sh
-)
-
-for script in "${INSTALL_SCRIPTS[@]}"; do
-    run_step "${SCRIPT_DIR}/install/${script}"
-done
-
-# ------------------------------------------------------------------------------
-# Configure
-# ------------------------------------------------------------------------------
-
-echo
-echo "===================== CONFIGURE ====================="
-
-CONFIGURE_SCRIPTS=(
-    environment.sh
-    git.sh
-    shell.sh
-    vscode.sh
-    paths.sh
-)
-
-for script in "${CONFIGURE_SCRIPTS[@]}"; do
-    run_step "${SCRIPT_DIR}/configure/${script}"
-done
-
-# ------------------------------------------------------------------------------
-# Workspace
-# ------------------------------------------------------------------------------
-
-echo
-echo "===================== WORKSPACE ====================="
-
-WORKSPACE_SCRIPTS=(
-    initialize.sh
-    directories.sh
-    permissions.sh
-    cleanup.sh
-)
-
-for script in "${WORKSPACE_SCRIPTS[@]}"; do
-    run_step "${SCRIPT_DIR}/workspace/${script}"
-done
-
-# ------------------------------------------------------------------------------
-# Verify
-# ------------------------------------------------------------------------------
-
-echo
-echo "====================== VERIFY ======================="
-
-VERIFY_SCRIPTS=(
-    system.sh
-    python.sh
-    node.sh
-    java.sh
-    flutter.sh
-    docker.sh
-    database.sh
-    workspace.sh
-)
-
-for script in "${VERIFY_SCRIPTS[@]}"; do
-    run_step "${SCRIPT_DIR}/verify/${script}"
-done
-
-# ------------------------------------------------------------------------------
-# Summary
-# ------------------------------------------------------------------------------
-
-echo
-echo "====================== SUMMARY ======================"
-
-run_step "${SCRIPT_DIR}/summary.sh"
-
-echo
-echo "====================================================="
-echo " BAOBAB Development Environment Ready"
-echo "====================================================="
-
-exit 0
diff --git a/.devcontainer/docker/scripts/config/versions.sh b/.devcontainer/docker/scripts/config/versions.sh
deleted file mode 100644
index 9f3b50f..0000000
--- a/.devcontainer/docker/scripts/config/versions.sh
+++ /dev/null
@@ -1,16 +0,0 @@
-# Operating System
-UBUNTU_VERSION="26.04"
-# Python
-PYTHON_VERSION="3.14"
-# Node.js
-NODE_VERSION="24"
-# Java
-JAVA_VERSION="21"
-# Flutter
-FLUTTER_VERSION="stable"
-# PostgreSQL
-POSTGRES_VERSION="17"
-# Redis
-REDIS_VERSION="8"
-# Docker
-DOCKER_CHANNEL="stable"
diff --git a/.devcontainer/docker/scripts/configure/environment.sh b/.devcontainer/docker/scripts/configure/environment.sh
deleted file mode 100755
index 25dda84..0000000
--- a/.devcontainer/docker/scripts/configure/environment.sh
+++ /dev/null
@@ -1,165 +0,0 @@
-#!/usr/bin/env bash
-
-# ==============================================================================
-# BAOBAB Enterprise Platform
-#
-# Script      : environment.sh
-# Purpose     : Configures the development environment variables required by
-#               the BAOBAB Development Environment.
-#
-# Author      : BAOBAB Contributors
-# License     : Apache License 2.0
-# ==============================================================================
-
-set -Eeuo pipefail
-
-SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
-
-# ------------------------------------------------------------------------------
-# Load Utilities
-# ------------------------------------------------------------------------------
-
-# shellcheck source=/dev/null
-source "${SCRIPT_DIR}/../utils/colors.sh"
-source "${SCRIPT_DIR}/../utils/functions.sh"
-source "${SCRIPT_DIR}/../utils/logging.sh"
-
-log_header "Environment Configuration"
-
-###############################################################################
-# Determine Project Root
-###############################################################################
-
-PROJECT_ROOT="$(get_project_root)"
-
-if [[ -z "${PROJECT_ROOT}" ]]; then
-
-    log_error "Unable to determine the BAOBAB project root."
-    exit 1
-
-fi
-
-if [[ ! -d "${PROJECT_ROOT}" ]]; then
-
-    log_error "Invalid project root: ${PROJECT_ROOT}"
-    exit 1
-
-fi
-
-log_info "Project Root : ${PROJECT_ROOT}"
-
-###############################################################################
-# Project Environment
-###############################################################################
-
-log_section "Project Environment"
-
-ENVIRONMENT_FILE="${PROJECT_ROOT}/.env"
-
-if [[ ! -f "${ENVIRONMENT_FILE}" ]]; then
-
-    if [[ -f "${PROJECT_ROOT}/.env.example" ]]; then
-
-        cp "${PROJECT_ROOT}/.env.example" "${ENVIRONMENT_FILE}"
-
-        log_success ".env created from .env.example"
-
-    else
-
-        log_warning ".env.example not found."
-
-    fi
-
-else
-
-    log_info ".env already exists."
-
-fi
-
-###############################################################################
-# Local Binary Path
-###############################################################################
-
-log_section "Local Binary Path"
-
-LOCAL_BIN='export PATH="$HOME/.local/bin:$PATH"'
-
-if ! grep -Fxq "${LOCAL_BIN}" "${HOME}/.bashrc"; then
-
-    printf '\n# BAOBAB Local Tools\n%s\n' "${LOCAL_BIN}" >> "${HOME}/.bashrc"
-
-    log_success "Added ~/.local/bin to PATH."
-
-else
-
-    log_info "PATH already configured."
-
-fi
-
-###############################################################################
-# BAOBAB Environment Variables
-###############################################################################
-
-log_section "BAOBAB Variables"
-
-#
-# Export immediately for the current provisioning process.
-#
-
-export PROJECT_ROOT
-export BAOBAB_HOME="${PROJECT_ROOT}"
-export BAOBAB_ENV="development"
-
-declare -A VARIABLES=(
-    ["PROJECT_ROOT"]="${PROJECT_ROOT}"
-    ["BAOBAB_HOME"]="${BAOBAB_HOME}"
-    ["BAOBAB_ENV"]="${BAOBAB_ENV}"
-)
-
-for VARIABLE in "${!VARIABLES[@]}"; do
-
-    VALUE="${VARIABLES[$VARIABLE]}"
-
-    if grep -q "^export ${VARIABLE}=" "${HOME}/.bashrc"; then
-
-        log_info "${VARIABLE} already configured."
-
-    else
-
-        printf 'export %s="%s"\n' "${VARIABLE}" "${VALUE}" >> "${HOME}/.bashrc"
-
-        log_success "${VARIABLE} configured."
-
-    fi
-
-done
-
-###############################################################################
-# Reload Environment
-###############################################################################
-
-log_section "Reload Environment"
-
-# shellcheck disable=SC1090
-source "${HOME}/.bashrc"
-
-log_success "Shell environment reloaded."
-
-###############################################################################
-# Summary
-###############################################################################
-
-log_section "Environment Summary"
-
-log_info "PROJECT_ROOT : ${PROJECT_ROOT}"
-log_info "BAOBAB_HOME : ${BAOBAB_HOME}"
-log_info "BAOBAB_ENV  : ${BAOBAB_ENV}"
-
-###############################################################################
-# Complete
-###############################################################################
-
-log_blank
-log_success "Environment configuration completed successfully."
-
-exit 0
diff --git a/.devcontainer/docker/scripts/configure/git.sh b/.devcontainer/docker/scripts/configure/git.sh
deleted file mode 100755
index 40e0636..0000000
--- a/.devcontainer/docker/scripts/configure/git.sh
+++ /dev/null
@@ -1,132 +0,0 @@
-#!/usr/bin/env bash
-
-# ==============================================================================
-# BAOBAB Enterprise Platform
-#
-# Script      : git.sh
-# Purpose     : Configures Git for the BAOBAB Development Environment.
-#
-# Author      : BAOBAB Contributors
-# License     : Apache License 2.0
-# ==============================================================================
-
-set -Eeuo pipefail
-
-SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
-
-# ------------------------------------------------------------------------------
-# Load Utilities
-# ------------------------------------------------------------------------------
-
-# shellcheck source=/dev/null
-source "${SCRIPT_DIR}/../utils/colors.sh"
-source "${SCRIPT_DIR}/../utils/functions.sh"
-source "${SCRIPT_DIR}/../utils/logging.sh"
-
-log_header "Git Configuration"
-
-###############################################################################
-# Verify Git
-###############################################################################
-
-log_section "Git Installation"
-
-if ! command_exists git; then
-    log_error "Git is not installed."
-    exit 1
-fi
-
-log_info "$(git --version)"
-
-###############################################################################
-# Default Configuration
-###############################################################################
-
-log_section "Git Defaults"
-
-git config --global init.defaultBranch main
-git config --global pull.rebase false
-git config --global fetch.prune true
-git config --global push.autoSetupRemote true
-git config --global core.editor "code --wait"
-git config --global core.autocrlf input
-git config --global core.filemode false
-git config --global core.ignorecase false
-git config --global color.ui auto
-
-log_success "Git defaults configured."
-
-###############################################################################
-# Git Identity
-###############################################################################
-
-log_section "Git Identity"
-
-NAME="$(git config --global user.name || true)"
-EMAIL="$(git config --global user.email || true)"
-
-if [[ -n "${NAME}" ]]; then
-    log_info "User Name : ${NAME}"
-else
-    log_warning "Git user.name has not been configured."
-fi
-
-if [[ -n "${EMAIL}" ]]; then
-    log_info "User Email: ${EMAIL}"
-else
-    log_warning "Git user.email has not been configured."
-fi
-
-###############################################################################
-# GitHub CLI
-###############################################################################
-
-log_section "GitHub CLI"
-
-if command_exists gh; then
-
-    log_info "$(gh --version | head -n 1)"
-
-    if gh auth status >/dev/null 2>&1; then
-        log_success "GitHub CLI authenticated."
-    else
-        log_warning "GitHub CLI is not authenticated."
-    fi
-
-else
-
-    log_warning "GitHub CLI is not installed."
-
-fi
-
-###############################################################################
-# Repository
-###############################################################################
-
-log_section "Repository"
-
-PROJECT_ROOT="$(get_project_root)"
-
-cd "${PROJECT_ROOT}"
-
-if git rev-parse --git-dir >/dev/null 2>&1; then
-
-    CURRENT_BRANCH="$(git branch --show-current)"
-
-    log_info "Repository detected."
-    log_info "Current Branch: ${CURRENT_BRANCH}"
-
-else
-
-    log_warning "Current directory is not a Git repository."
-
-fi
-
-###############################################################################
-# Complete
-###############################################################################
-
-log_blank
-log_success "Git configuration completed successfully."
-
-exit 0
diff --git a/.devcontainer/docker/scripts/configure/paths.sh b/.devcontainer/docker/scripts/configure/paths.sh
deleted file mode 100755
index 8628c5b..0000000
--- a/.devcontainer/docker/scripts/configure/paths.sh
+++ /dev/null
@@ -1,141 +0,0 @@
-#!/usr/bin/env bash
-
-# ==============================================================================
-# BAOBAB Enterprise Platform
-#
-# Script      : paths.sh
-# Purpose     : Configures common PATH entries and project-specific environment
-#               variables for the BAOBAB Development Environment.
-#
-# Author      : BAOBAB Contributors
-# License     : Apache License 2.0
-# ==============================================================================
-
-set -Eeuo pipefail
-
-SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
-
-# ------------------------------------------------------------------------------
-# Load Utilities
-# ------------------------------------------------------------------------------
-
-# shellcheck source=/dev/null
-source "${SCRIPT_DIR}/../utils/colors.sh"
-source "${SCRIPT_DIR}/../utils/functions.sh"
-source "${SCRIPT_DIR}/../utils/logging.sh"
-
-log_header "PATH Configuration"
-
-PROJECT_ROOT="$(get_project_root)"
-BASHRC="${HOME}/.bashrc"
-
-###############################################################################
-# BAOBAB PATH Configuration
-###############################################################################
-
-log_section "Shell PATH"
-
-START_MARKER="# >>> BAOBAB PATH CONFIGURATION >>>"
-END_MARKER="# <<< BAOBAB PATH CONFIGURATION <<<"
-
-if grep -Fq "${START_MARKER}" "${BASHRC}"; then
-
-    log_info "BAOBAB PATH configuration already exists."
-
-else
-
-cat >> "${BASHRC}" <<EOF
-
-${START_MARKER}
-
-###############################################################################
-# BAOBAB Project
-###############################################################################
-
-export BAOBAB_HOME="${PROJECT_ROOT}"
-
-###############################################################################
-# Python
-###############################################################################
-
-export PATH="\$HOME/.local/bin:\$PATH"
-
-###############################################################################
-# Flutter
-###############################################################################
-
-if [ -d "\$HOME/flutter/bin" ]; then
-    export PATH="\$HOME/flutter/bin:\$PATH"
-fi
-
-###############################################################################
-# Android SDK
-###############################################################################
-
-if [ -n "\${ANDROID_HOME:-}" ]; then
-    export PATH="\$ANDROID_HOME/platform-tools:\$PATH"
-    export PATH="\$ANDROID_HOME/emulator:\$PATH"
-    export PATH="\$ANDROID_HOME/cmdline-tools/latest/bin:\$PATH"
-fi
-
-###############################################################################
-# Java
-###############################################################################
-
-if [ -n "\${JAVA_HOME:-}" ]; then
-    export PATH="\$JAVA_HOME/bin:\$PATH"
-fi
-
-###############################################################################
-# Node.js
-###############################################################################
-
-if [ -d "${PROJECT_ROOT}/services/frontend/node_modules/.bin" ]; then
-    export PATH="${PROJECT_ROOT}/services/frontend/node_modules/.bin:\$PATH"
-fi
-
-###############################################################################
-# Python Virtual Environment
-###############################################################################
-
-if [ -d "${PROJECT_ROOT}/.venv/bin" ]; then
-    export PATH="${PROJECT_ROOT}/.venv/bin:\$PATH"
-fi
-
-${END_MARKER}
-
-EOF
-
-    log_success "BAOBAB PATH configuration added."
-
-fi
-
-###############################################################################
-# Reload Shell
-###############################################################################
-
-log_section "Reload Shell"
-
-# shellcheck disable=SC1090
-source "${BASHRC}"
-
-log_success "Shell configuration reloaded."
-
-###############################################################################
-# Display Important Paths
-###############################################################################
-
-log_section "Configured Paths"
-
-log_info "BAOBAB_HOME : ${PROJECT_ROOT}"
-log_info "HOME         : ${HOME}"
-log_info "PATH         : ${PATH}"
-
-###############################################################################
-# Complete
-###############################################################################
-
-log_blank
-log_success "PATH configuration completed successfully."
-
-exit 0
diff --git a/.devcontainer/docker/scripts/configure/shell.sh b/.devcontainer/docker/scripts/configure/shell.sh
deleted file mode 100755
index 7f9f8d2..0000000
--- a/.devcontainer/docker/scripts/configure/shell.sh
+++ /dev/null
@@ -1,122 +0,0 @@
-#!/usr/bin/env bash
-
-# ==============================================================================
-# BAOBAB Enterprise Platform
-#
-# Script      : shell.sh
-# Purpose     : Configures the Bash shell for the BAOBAB Development
-#               Environment.
-#
-# Author      : BAOBAB Contributors
-# License     : Apache License 2.0
-# ==============================================================================
-
-set -Eeuo pipefail
-
-SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
-
-# ------------------------------------------------------------------------------
-# Load Utilities
-# ------------------------------------------------------------------------------
-
-# shellcheck source=/dev/null
-source "${SCRIPT_DIR}/../utils/colors.sh"
-source "${SCRIPT_DIR}/../utils/functions.sh"
-source "${SCRIPT_DIR}/../utils/logging.sh"
-
-log_header "Shell Configuration"
-
-BASHRC="${HOME}/.bashrc"
-
-###############################################################################
-# BAOBAB Shell Configuration
-###############################################################################
-
-log_section "BAOBAB Shell Configuration"
-
-START_MARKER="# >>> BAOBAB SHELL CONFIGURATION >>>"
-END_MARKER="# <<< BAOBAB SHELL CONFIGURATION <<<"
-
-if grep -Fq "${START_MARKER}" "${BASHRC}"; then
-
-    log_info "BAOBAB shell configuration already exists."
-
-else
-
-cat >> "${BASHRC}" <<'EOF'
-
-# >>> BAOBAB SHELL CONFIGURATION >>>
-
-###############################################################################
-# History
-###############################################################################
-
-export HISTSIZE=10000
-export HISTFILESIZE=20000
-export HISTCONTROL=ignoreboth
-shopt -s histappend
-
-###############################################################################
-# Aliases
-###############################################################################
-
-alias ll='ls -alF'
-alias la='ls -A'
-alias l='ls -CF'
-
-alias cls='clear'
-
-alias gs='git status'
-alias ga='git add'
-alias gc='git commit'
-alias gp='git push'
-alias gl='git pull'
-
-alias dc='docker compose'
-alias dps='docker ps'
-alias di='docker images'
-
-alias python='python3'
-alias pip='pip3'
-
-###############################################################################
-# Prompt
-###############################################################################
-
-export PS1='\u@\h:\w\$ '
-
-###############################################################################
-# Completion
-###############################################################################
-
-if [ -f /etc/bash_completion ]; then
-    source /etc/bash_completion
-fi
-
-# <<< BAOBAB SHELL CONFIGURATION <<<
-
-EOF
-
-    log_success "BAOBAB shell configuration added."
-
-fi
-
-###############################################################################
-# Reload Shell
-###############################################################################
-
-log_section "Reload Shell"
-
-# shellcheck disable=SC1090
-source "${BASHRC}"
-
-log_success "Shell configuration reloaded."
-
-###############################################################################
-# Complete
-###############################################################################
-
-log_blank
-log_success "Shell configuration completed successfully."
-
-exit 0
diff --git a/.devcontainer/docker/scripts/configure/vscode.sh b/.devcontainer/docker/scripts/configure/vscode.sh
deleted file mode 100755
index 5008b38..0000000
--- a/.devcontainer/docker/scripts/configure/vscode.sh
+++ /dev/null
@@ -1,182 +0,0 @@
-#!/usr/bin/env bash
-
-# ==============================================================================
-# BAOBAB Enterprise Platform
-#
-# Script      : vscode.sh
-# Purpose     : Configures Visual Studio Code for the BAOBAB Development
-#               Environment.
-#
-# Author      : BAOBAB Contributors
-# License     : Apache License 2.0
-# ==============================================================================
-
-set -Eeuo pipefail
-
-SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
-
-# ------------------------------------------------------------------------------
-# Load Utilities
-# ------------------------------------------------------------------------------
-
-# shellcheck source=/dev/null
-source "${SCRIPT_DIR}/../utils/colors.sh"
-source "${SCRIPT_DIR}/../utils/functions.sh"
-source "${SCRIPT_DIR}/../utils/logging.sh"
-
-log_header "Visual Studio Code Configuration"
-
-PROJECT_ROOT="$(get_project_root)"
-VSCODE_DIR="${PROJECT_ROOT}/.vscode"
-
-###############################################################################
-# Create Workspace Directory
-###############################################################################
-
-log_section "Workspace"
-
-ensure_directory "${VSCODE_DIR}"
-
-###############################################################################
-# settings.json
-###############################################################################
-
-log_section "settings.json"
-
-SETTINGS_FILE="${VSCODE_DIR}/settings.json"
-
-if [[ ! -f "${SETTINGS_FILE}" ]]; then
-
-cat > "${SETTINGS_FILE}" <<'EOF'
-{
-    "editor.formatOnSave": true,
-    "editor.codeActionsOnSave": {
-        "source.fixAll": "explicit",
-        "source.organizeImports": "explicit"
-    },
-
-    "files.trimTrailingWhitespace": true,
-    "files.insertFinalNewline": true,
-
-    "terminal.integrated.defaultProfile.linux": "bash",
-
-    "python.defaultInterpreterPath": ".venv/bin/python",
-
-    "python.terminal.activateEnvironment": true,
-
-    "typescript.tsdk": "services/frontend/node_modules/typescript/lib",
-
-    "search.exclude": {
-        "**/.git": true,
-        "**/.next": true,
-        "**/__pycache__": true,
-        "**/.pytest_cache": true,
-        "**/.venv": true
-    }
-}
-EOF
-
-    log_success "settings.json created."
-
-else
-
-    log_info "settings.json already exists."
-
-fi
-
-###############################################################################
-# extensions.json
-###############################################################################
-
-log_section "extensions.json"
-
-EXTENSIONS_FILE="${VSCODE_DIR}/extensions.json"
-
-if [[ ! -f "${EXTENSIONS_FILE}" ]]; then
-
-cat > "${EXTENSIONS_FILE}" <<'EOF'
-{
-    "recommendations": [
-        "ms-python.python",
-        "ms-python.vscode-pylance",
-        "charliermarsh.ruff",
-        "ms-azuretools.vscode-docker",
-        "esbenp.prettier-vscode",
-        "dbaeumer.vscode-eslint",
-        "bradlc.vscode-tailwindcss",
-        "ms-vscode.makefile-tools",
-        "eamodio.gitlens",
-        "github.vscode-github-actions",
-        "ms-kubernetes-tools.vscode-kubernetes-tools",
-        "redhat.vscode-yaml",
-        "humao.rest-client"
-    ]
-}
-EOF
-
-    log_success "extensions.json created."
-
-else
-
-    log_info "extensions.json already exists."
-
-fi
-
-###############################################################################
-# launch.json
-###############################################################################
-
-log_section "launch.json"
-
-LAUNCH_FILE="${VSCODE_DIR}/launch.json"
-
-if [[ ! -f "${LAUNCH_FILE}" ]]; then
-
-cat > "${LAUNCH_FILE}" <<'EOF'
-{
-    "version": "0.2.0",
-    "configurations": []
-}
-EOF
-
-    log_success "launch.json created."
-
-else
-
-    log_info "launch.json already exists."
-
-fi
-
-###############################################################################
-# tasks.json
-###############################################################################
-
-log_section "tasks.json"
-
-TASKS_FILE="${VSCODE_DIR}/tasks.json"
-
-if [[ ! -f "${TASKS_FILE}" ]]; then
-
-cat > "${TASKS_FILE}" <<'EOF'
-{
-    "version": "2.0.0",
-    "tasks": []
-}
-EOF
-
-    log_success "tasks.json created."
-
-else
-
-    log_info "tasks.json already exists."
-
-fi
-
-###############################################################################
-# Complete
-###############################################################################
-
-log_blank
-log_success "Visual Studio Code configuration completed successfully."
-
-exit 0
diff --git a/.devcontainer/docker/scripts/install/database.sh b/.devcontainer/docker/scripts/install/database.sh
deleted file mode 100755
index fdbef0e..0000000
--- a/.devcontainer/docker/scripts/install/database.sh
+++ /dev/null
@@ -1,172 +0,0 @@
-#!/usr/bin/env bash
-
-# ==============================================================================
-# BAOBAB Enterprise Platform
-#
-# Script      : database.sh
-# Purpose     : Installs and configures the database client tools required by
-#               the BAOBAB Development Environment.
-#
-# Author      : BAOBAB Contributors
-# License     : Apache License 2.0
-# ==============================================================================
-
-set -Eeuo pipefail
-
-SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
-
-# ------------------------------------------------------------------------------
-# Load Utilities
-# ------------------------------------------------------------------------------
-
-# shellcheck source=/dev/null
-source "${SCRIPT_DIR}/../utils/colors.sh"
-source "${SCRIPT_DIR}/../utils/functions.sh"
-source "${SCRIPT_DIR}/../utils/logging.sh"
-
-log_header "Database Tool Installation"
-
-###############################################################################
-# PostgreSQL Client
-###############################################################################
-
-log_section "PostgreSQL Client"
-
-if command_exists psql; then
-
-    log_info "PostgreSQL client is already installed."
-    log_info "$(psql --version)"
-
-else
-
-    log_info "Installing PostgreSQL client..."
-
-    sudo apt-get update
-
-    sudo apt-get install -y \
-        postgresql-client
-
-    log_success "PostgreSQL client installed."
-
-fi
-
-###############################################################################
-# Redis CLI
-###############################################################################
-
-log_section "Redis CLI"
-
-if command_exists redis-cli; then
-
-    log_info "Redis CLI is already installed."
-    log_info "$(redis-cli --version)"
-
-else
-
-    log_info "Installing Redis tools..."
-
-    sudo apt-get install -y \
-        redis-tools
-
-    log_success "Redis CLI installed."
-
-fi
-
-###############################################################################
-# PostgreSQL Connectivity
-###############################################################################
-
-log_section "PostgreSQL Connectivity"
-
-PROJECT_ROOT="$(get_project_root)"
-COMPOSE_FILE="${PROJECT_ROOT}/compose.yaml"
-
-if file_exists "${COMPOSE_FILE}"; then
-
-    cd "${PROJECT_ROOT}"
-
-    if docker compose ps db >/dev/null 2>&1; then
-
-        log_info "Database service 'db' detected."
-
-        if docker compose exec -T db pg_isready >/dev/null 2>&1; then
-            log_success "PostgreSQL server is accepting connections."
-        else
-            log_warning "PostgreSQL container exists but is not yet ready."
-        fi
-
-    else
-
-        log_warning "Database container is not running."
-
-    fi
-
-else
-
-    log_warning "compose.yaml not found. Skipping database connectivity checks."
-
-fi
-
-###############################################################################
-# Redis Connectivity
-###############################################################################
-
-log_section "Redis Connectivity"
-
-if file_exists "${COMPOSE_FILE}"; then
-
-    cd "${PROJECT_ROOT}"
-
-    if docker compose ps redis >/dev/null 2>&1; then
-
-        log_info "Redis service detected."
-
-        if docker compose exec -T redis redis-cli ping >/dev/null 2>&1; then
-            log_success "Redis server is responding."
-        else
-            log_warning "Redis container exists but is not yet ready."
-        fi
-
-    else
-
-        log_warning "Redis container is not running."
-
-    fi
-
-fi
-
-###############################################################################
-# Database Environment
-###############################################################################
-
-log_section "Database Environment"
-
-if [[ -f "${PROJECT_ROOT}/.env" ]]; then
-
-    log_info "Project environment file detected."
-
-else
-
-    log_warning ".env file not found."
-
-fi
-
-###############################################################################
-# Summary
-###############################################################################
-
-log_blank
-
-log_info "Installed Database Tools"
-
-command_exists psql && log_info "  • PostgreSQL : $(psql --version)"
-command_exists redis-cli && log_info "  • Redis CLI : $(redis-cli --version)"
-
-###############################################################################
-# Complete
-###############################################################################
-
-log_blank
-log_success "Database environment configured successfully."
-
-exit 0
diff --git a/.devcontainer/docker/scripts/install/docker.sh b/.devcontainer/docker/scripts/install/docker.sh
deleted file mode 100755
index c52a367..0000000
--- a/.devcontainer/docker/scripts/install/docker.sh
+++ /dev/null
@@ -1,136 +0,0 @@
-#!/usr/bin/env bash
-
-# ==============================================================================
-# BAOBAB Enterprise Platform
-#
-# Script      : docker.sh
-# Purpose     : Verifies and configures the Docker development tooling required
-#               by the BAOBAB Development Environment.
-#
-# Author      : BAOBAB Contributors
-# License     : Apache License 2.0
-# ==============================================================================
-
-set -Eeuo pipefail
-
-SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
-
-# ------------------------------------------------------------------------------
-# Load Utilities
-# ------------------------------------------------------------------------------
-
-# shellcheck source=/dev/null
-source "${SCRIPT_DIR}/../utils/colors.sh"
-source "${SCRIPT_DIR}/../utils/functions.sh"
-source "${SCRIPT_DIR}/../utils/logging.sh"
-
-log_header "Docker Installation"
-
-###############################################################################
-# Docker CLI
-###############################################################################
-
-log_section "Docker CLI"
-
-if command_exists docker; then
-
-    log_info "Docker CLI is available."
-    log_info "$(docker --version)"
-
-else
-
-    log_error "Docker CLI is not installed."
-    exit 1
-
-fi
-
-###############################################################################
-# Docker Compose
-###############################################################################
-
-log_section "Docker Compose"
-
-if docker compose version >/dev/null 2>&1; then
-
-    log_info "$(docker compose version)"
-
-else
-
-    log_error "Docker Compose is not available."
-    exit 1
-
-fi
-
-###############################################################################
-# Docker Buildx
-###############################################################################
-
-log_section "Docker Buildx"
-
-if docker buildx version >/dev/null 2>&1; then
-
-    log_info "$(docker buildx version | head -n 1)"
-
-else
-
-    log_warning "Docker Buildx is not available."
-
-fi
-
-###############################################################################
-# Docker Context
-###############################################################################
-
-log_section "Docker Context"
-
-CURRENT_CONTEXT="$(docker context show)"
-
-log_info "Current Context: ${CURRENT_CONTEXT}"
-
-###############################################################################
-# Docker Engine
-###############################################################################
-
-log_section "Docker Engine"
-
-if docker info >/dev/null 2>&1; then
-
-    log_success "Docker Engine is accessible."
-
-else
-
-    log_warning "Docker Engine is not currently running or is unavailable."
-
-fi
-
-###############################################################################
-# BAOBAB Compose File
-###############################################################################
-
-log_section "Compose Configuration"
-
-PROJECT_ROOT="$(get_project_root)"
-COMPOSE_FILE="${PROJECT_ROOT}/compose.yaml"
-
-if file_exists "${COMPOSE_FILE}"; then
-
-    cd "${PROJECT_ROOT}"
-
-    docker compose config >/dev/null
-
-    log_success "compose.yaml validation successful."
-
-else
-
-    log_warning "compose.yaml not found."
-
-fi
-
-###############################################################################
-# Complete
-###############################################################################
-
-log_blank
-log_success "Docker environment verification completed."
-
-exit 0
diff --git a/.devcontainer/docker/scripts/install/flutter.sh b/.devcontainer/docker/scripts/install/flutter.sh
deleted file mode 100755
index b64734b..0000000
--- a/.devcontainer/docker/scripts/install/flutter.sh
+++ /dev/null
@@ -1,133 +0,0 @@
-#!/usr/bin/env bash
-
-# ==============================================================================
-# BAOBAB Enterprise Platform
-#
-# Script      : flutter.sh
-# Purpose     : Installs and configures the Flutter and Dart SDKs for the
-#               BAOBAB Development Environment.
-#
-# Author      : BAOBAB Contributors
-# License     : Apache License 2.0
-# ==============================================================================
-
-set -Eeuo pipefail
-
-SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
-
-# ------------------------------------------------------------------------------
-# Load Utilities
-# ------------------------------------------------------------------------------
-
-# shellcheck source=/dev/null
-source "${SCRIPT_DIR}/../utils/colors.sh"
-source "${SCRIPT_DIR}/../utils/functions.sh"
-source "${SCRIPT_DIR}/../utils/logging.sh"
-
-log_header "Flutter Installation"
-
-PROJECT_ROOT="$(get_project_root)"
-FLUTTER_HOME="${HOME}/flutter"
-FLUTTER_BIN="${FLUTTER_HOME}/bin"
-
-###############################################################################
-# Install Flutter SDK
-###############################################################################
-
-log_section "Flutter SDK"
-
-if command_exists flutter; then
-
-    log_info "Flutter is already installed."
-    log_info "$(flutter --version | head -n 1)"
-
-else
-
-    log_info "Installing Flutter SDK..."
-
-    if [[ ! -d "${FLUTTER_HOME}" ]]; then
-
-        git clone https://github.com/flutter/flutter.git \
-            --branch stable \
-            "${FLUTTER_HOME}"
-
-    fi
-
-    export PATH="${FLUTTER_BIN}:${PATH}"
-
-    if ! grep -q "flutter/bin" "${HOME}/.bashrc"; then
-        {
-            echo ""
-            echo "# Flutter"
-            echo "export PATH=\"${FLUTTER_BIN}:\$PATH\""
-        } >> "${HOME}/.bashrc"
-    fi
-
-    log_success "Flutter SDK installed."
-
-fi
-
-###############################################################################
-# Verify Dart
-###############################################################################
-
-log_section "Dart SDK"
-
-if command_exists dart; then
-
-    log_info "$(dart --version 2>&1)"
-
-else
-
-    log_warning "Dart will become available after Flutter is added to PATH."
-
-fi
-
-###############################################################################
-# Flutter Precache
-###############################################################################
-
-log_section "Flutter Cache"
-
-flutter precache
-
-log_success "Flutter cache downloaded."
-
-###############################################################################
-# Flutter Doctor
-###############################################################################
-
-log_section "Flutter Doctor"
-
-flutter doctor
-
-###############################################################################
-# Mobile Dependencies
-###############################################################################
-
-log_section "Mobile Dependencies"
-
-MOBILE_DIR="${PROJECT_ROOT}/services/mobile"
-
-if directory_exists "${MOBILE_DIR}" && file_exists "${MOBILE_DIR}/pubspec.yaml"; then
-
-    cd "${MOBILE_DIR}"
-
-    flutter pub get
-
-    log_success "Flutter project dependencies installed."
-
-else
-
-    log_warning "Mobile project not initialized. Skipping dependency installation."
-
-fi
-
-###############################################################################
-# Complete
-###############################################################################
-
-log_blank
-log_success "Flutter environment configured successfully."
-
-exit 0
diff --git a/.devcontainer/docker/scripts/install/java.sh b/.devcontainer/docker/scripts/install/java.sh
deleted file mode 100755
index c5d3b99..0000000
--- a/.devcontainer/docker/scripts/install/java.sh
+++ /dev/null
@@ -1,116 +0,0 @@
-#!/usr/bin/env bash
-
-# ==============================================================================
-# BAOBAB Enterprise Platform
-#
-# Script      : java.sh
-# Purpose     : Installs and configures the Java Development Kit (JDK) for the
-#               BAOBAB Development Environment.
-#
-# Author      : BAOBAB Contributors
-# License     : Apache License 2.0
-# ==============================================================================
-
-set -Eeuo pipefail
-
-SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
-
-# ------------------------------------------------------------------------------
-# Load Utilities
-# ------------------------------------------------------------------------------
-
-# shellcheck source=/dev/null
-source "${SCRIPT_DIR}/../utils/colors.sh"
-source "${SCRIPT_DIR}/../utils/functions.sh"
-source "${SCRIPT_DIR}/../utils/logging.sh"
-
-log_header "Java Installation"
-
-###############################################################################
-# Verify Java
-###############################################################################
-
-log_section "Java Development Kit"
-
-if command_exists java && command_exists javac; then
-
-    log_info "Java is already installed."
-    log_info "$(java --version | head -n 1)"
-
-else
-
-    log_info "Installing OpenJDK..."
-
-    sudo apt-get update
-
-    sudo apt-get install -y \
-        openjdk-21-jdk
-
-    log_success "OpenJDK installed."
-
-fi
-
-###############################################################################
-# Verify Java Compiler
-###############################################################################
-
-log_section "Java Compiler"
-
-if command_exists javac; then
-
-    log_info "$(javac --version)"
-
-else
-
-    log_error "Java compiler (javac) is not available."
-    exit 1
-
-fi
-
-###############################################################################
-# Configure JAVA_HOME
-###############################################################################
-
-log_section "JAVA_HOME"
-
-JAVA_HOME_PATH="$(dirname "$(dirname "$(readlink -f "$(command -v javac)")")")"
-
-if [[ -n "${JAVA_HOME_PATH}" ]]; then
-
-    export JAVA_HOME="${JAVA_HOME_PATH}"
-
-    if ! grep -q "JAVA_HOME" "${HOME}/.bashrc"; then
-        {
-            echo ""
-            echo "# Java"
-            echo "export JAVA_HOME=${JAVA_HOME_PATH}"
-            echo 'export PATH="$JAVA_HOME/bin:$PATH"'
-        } >> "${HOME}/.bashrc"
-    fi
-
-    log_info "JAVA_HOME=${JAVA_HOME}"
-
-else
-
-    log_error "Unable to determine JAVA_HOME."
-    exit 1
-
-fi
-
-###############################################################################
-# Verify Environment
-###############################################################################
-
-log_section "Verification"
-
-log_info "$(java --version | head -n 1)"
-log_info "$(javac --version)"
-
-###############################################################################
-# Complete
-###############################################################################
-
-log_blank
-log_success "Java environment configured successfully."
-
-exit 0
diff --git a/.devcontainer/docker/scripts/install/node.sh b/.devcontainer/docker/scripts/install/node.sh
deleted file mode 100755
index 3b6140d..0000000
--- a/.devcontainer/docker/scripts/install/node.sh
+++ /dev/null
@@ -1,120 +0,0 @@
-#!/usr/bin/env bash
-
-# ==============================================================================
-# BAOBAB Enterprise Platform
-#
-# Script      : node.sh
-# Purpose     : Installs and configures the Node.js development environment for
-#               the BAOBAB Development Environment.
-#
-# Author      : BAOBAB Contributors
-# License     : Apache License 2.0
-# ==============================================================================
-
-set -Eeuo pipefail
-
-SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
-
-# ------------------------------------------------------------------------------
-# Load Utilities
-# ------------------------------------------------------------------------------
-
-# shellcheck source=/dev/null
-source "${SCRIPT_DIR}/../utils/colors.sh"
-source "${SCRIPT_DIR}/../utils/functions.sh"
-source "${SCRIPT_DIR}/../utils/logging.sh"
-
-log_header "Node.js Installation"
-
-###############################################################################
-# Verify Node.js
-###############################################################################
-
-log_section "Node.js"
-
-if command_exists node; then
-
-    log_info "Node.js is already installed."
-    log_info "$(node --version)"
-
-else
-
-    log_info "Installing Node.js..."
-
-    sudo apt-get update
-
-    sudo apt-get install -y \
-        nodejs \
-        npm
-
-    log_success "Node.js installed."
-
-fi
-
-###############################################################################
-# Verify npm
-###############################################################################
-
-log_section "npm"
-
-if command_exists npm; then
-
-    log_info "npm is available."
-    log_info "Version: $(npm --version)"
-
-else
-
-    log_error "npm installation failed."
-    exit 1
-
-fi
-
-###############################################################################
-# Enable Corepack
-###############################################################################
-
-log_section "Corepack"
-
-if command_exists corepack; then
-
-    corepack enable
-
-    log_success "Corepack enabled."
-
-else
-
-    log_warning "Corepack is not available."
-
-fi
-
-###############################################################################
-# Install Project Dependencies
-###############################################################################
-
-log_section "Frontend Dependencies"
-
-PROJECT_ROOT="$(get_project_root)"
-FRONTEND_DIR="${PROJECT_ROOT}/services/frontend"
-
-if directory_exists "${FRONTEND_DIR}" && file_exists "${FRONTEND_DIR}/package.json"; then
-
-    cd "${FRONTEND_DIR}"
-
-    npm install
-
-    log_success "Frontend dependencies installed."
-
-else
-
-    log_warning "Frontend project not initialized. Skipping dependency installation."
-
-fi
-
-###############################################################################
-# Complete
-###############################################################################
-
-log_blank
-log_success "Node.js environment configured successfully."
-
-exit 0
diff --git a/.devcontainer/docker/scripts/install/python.sh b/.devcontainer/docker/scripts/install/python.sh
deleted file mode 100755
index b4ce1b6..0000000
--- a/.devcontainer/docker/scripts/install/python.sh
+++ /dev/null
@@ -1,230 +0,0 @@
-#!/usr/bin/env bash
-
-# ==============================================================================
-# BAOBAB Enterprise Platform
-#
-# Script      : python.sh
-# Purpose     : Installs and configures the Python development environment for
-#               the BAOBAB Development Environment.
-#
-# Author      : BAOBAB Contributors
-# License     : Apache License 2.0
-# ==============================================================================
-
-set -Eeuo pipefail
-
-SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
-
-# ------------------------------------------------------------------------------
-# Load Utilities
-# ------------------------------------------------------------------------------
-
-# shellcheck source=/dev/null
-source "${SCRIPT_DIR}/../utils/colors.sh"
-source "${SCRIPT_DIR}/../utils/functions.sh"
-source "${SCRIPT_DIR}/../utils/logging.sh"
-source "${SCRIPT_DIR}/../config/versions.sh"
-
-log_header "Python Installation"
-
-###############################################################################
-# Refresh Package Index
-###############################################################################
-
-sudo apt-get update
-
-###############################################################################
-# Configuration
-###############################################################################
-
-PYTHON_BIN="python3"
-
-PYTHON_PACKAGE="python${PYTHON_VERSION}-full"
-PYTHON_DEV_PACKAGE="python${PYTHON_VERSION}-dev"
-PYTHON_VENV_PACKAGE="python${PYTHON_VERSION}-venv"
-
-UV_HOME="${HOME}/.local/bin"
-
-export PATH="${UV_HOME}:${PATH}"
-
-###############################################################################
-# Python
-###############################################################################
-
-log_section "Python"
-
-if ! command_exists "${PYTHON_BIN}"; then
-
-    log_info "Python not found."
-    log_info "Installing Python ${PYTHON_VERSION}..."
-
-    sudo apt-get install -y \
-        "${PYTHON_PACKAGE}" \
-        "${PYTHON_DEV_PACKAGE}"
-
-fi
-
-CURRENT_VERSION="$(
-${PYTHON_BIN} -c '
-import sys
-print(f"{sys.version_info.major}.{sys.version_info.minor}")
-'
-)"
-
-log_info "Target Python : ${PYTHON_VERSION}"
-log_info "Detected      : ${CURRENT_VERSION}"
-
-if [[ "${CURRENT_VERSION}" != "${PYTHON_VERSION}" ]]; then
-
-    log_warning "Expected Python ${PYTHON_VERSION}."
-    log_warning "Detected Python ${CURRENT_VERSION}."
-
-    sudo apt-get install -y \
-        "${PYTHON_PACKAGE}" \
-        "${PYTHON_DEV_PACKAGE}"
-
-fi
-
-###############################################################################
-# pip
-###############################################################################
-
-log_section "pip"
-
-if ! ${PYTHON_BIN} -m ensurepip --version >/dev/null 2>&1; then
-
-    log_info "Bootstrapping pip..."
-
-    sudo apt-get install -y "${PYTHON_PACKAGE}"
-
-fi
-
-if ${PYTHON_BIN} -m pip --version >/dev/null 2>&1; then
-
-    log_info "$(${PYTHON_BIN} -m pip --version)"
-
-else
-
-    log_error "pip installation failed."
-    exit 1
-
-fi
-
-###############################################################################
-# venv
-###############################################################################
-
-log_section "venv"
-
-if ${PYTHON_BIN} -m venv --help >/dev/null 2>&1; then
-
-    log_info "venv support detected."
-
-else
-
-    log_info "Installing ${PYTHON_VENV_PACKAGE}..."
-
-    sudo apt-get install -y "${PYTHON_VENV_PACKAGE}"
-
-fi
-
-###############################################################################
-# Development Headers
-###############################################################################
-
-log_section "Development Headers"
-
-if dpkg -s "${PYTHON_DEV_PACKAGE}" >/dev/null 2>&1; then
-
-    log_info "${PYTHON_DEV_PACKAGE} already installed."
-
-else
-
-    log_info "Installing ${PYTHON_DEV_PACKAGE}..."
-
-    sudo apt-get install -y "${PYTHON_DEV_PACKAGE}"
-
-fi
-
-###############################################################################
-# Upgrade Python Toolchain
-###############################################################################
-
-log_section "Python Toolchain"
-
-${PYTHON_BIN} -m pip install --upgrade \
-    pip \
-    setuptools \
-    wheel
-
-log_success "Python toolchain upgraded."
-
-###############################################################################
-# uv
-###############################################################################
-
-log_section "uv"
-
-if command_exists uv; then
-
-    log_info "uv already installed."
-    log_info "$(uv --version)"
-
-else
-
-    log_info "Installing uv..."
-
-    curl \
-        --fail \
-        --location \
-        --silent \
-        --show-error \
-        https://astral.sh/uv/install.sh | sh
-
-    if ! grep -q 'HOME/.local/bin' "${HOME}/.bashrc"; then
-        echo 'export PATH="$HOME/.local/bin:$PATH"' >> "${HOME}/.bashrc"
-    fi
-
-    export PATH="${UV_HOME}:${PATH}"
-
-    log_success "uv installed."
-
-fi
-
-###############################################################################
-# Python Workspace
-###############################################################################
-
-log_section "Python Workspace"
-
-PROJECT_ROOT="$(get_project_root)"
-
-cd "${PROJECT_ROOT}"
-
-if [[ ! -f pyproject.toml ]]; then
-
-    log_info "Initializing uv workspace..."
-
-    uv init \
-        --bare \
-        --python "${PYTHON_BIN}"
-
-    log_success "uv workspace initialized."
-
-else
-
-    log_info "Existing pyproject.toml found."
-
-fi
-
-log_info "uv interpreter:"
-uv python find
-
-###############################################################################
-# Complete
-###############################################################################
-
-log_blank
-log_success "Python ${CURRENT_VERSION} environment configured successfully."
-
-exit 0
diff --git a/.devcontainer/docker/scripts/install/system.sh b/.devcontainer/docker/scripts/install/system.sh
deleted file mode 100755
index 3344ade..0000000
--- a/.devcontainer/docker/scripts/install/system.sh
+++ /dev/null
@@ -1,166 +0,0 @@
-#!/usr/bin/env bash
-
-# ==============================================================================
-# BAOBAB Enterprise Platform
-#
-# Script      : system.sh
-# Purpose     : Verifies and prepares the operating system for the BAOBAB
-#               Development Environment.
-#
-# Author      : BAOBAB Contributors
-# License     : Apache License 2.0
-# ==============================================================================
-
-set -Eeuo pipefail
-
-SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
-
-# ------------------------------------------------------------------------------
-# Load Utilities
-# ------------------------------------------------------------------------------
-
-# shellcheck source=/dev/null
-source "${SCRIPT_DIR}/../utils/colors.sh"
-source "${SCRIPT_DIR}/../utils/functions.sh"
-source "${SCRIPT_DIR}/../utils/logging.sh"
-source "${SCRIPT_DIR}/../config/versions.sh"
-
-log_header "System Installation"
-
-###############################################################################
-# Helper Functions
-###############################################################################
-
-ensure_package() {
-
-    local package="$1"
-
-    if dpkg -s "${package}" >/dev/null 2>&1; then
-
-        log_info "${package} already installed."
-
-    else
-
-        log_info "Installing ${package}..."
-
-        sudo apt-get install -y "${package}"
-
-        log_success "${package} installed."
-
-    fi
-}
-
-###############################################################################
-# Verify Operating System
-###############################################################################
-
-log_section "Operating System"
-
-CURRENT_UBUNTU="$(lsb_release -rs)"
-
-log_info "Detected Ubuntu ${CURRENT_UBUNTU}"
-
-if [[ "${CURRENT_UBUNTU}" != "${UBUNTU_VERSION}" ]]; then
-    log_warning "Configured target Ubuntu version: ${UBUNTU_VERSION}"
-fi
-
-###############################################################################
-# Update Package Index
-###############################################################################
-
-log_section "Updating Package Index"
-
-sudo apt-get update
-
-###############################################################################
-# Verify Core Packages
-###############################################################################
-
-log_section "Core Packages"
-
-SYSTEM_PACKAGES=(
-    apt-transport-https
-    build-essential
-    ca-certificates
-    curl
-    file
-    git
-    gnupg
-    jq
-    lsb-release
-    software-properties-common
-    tree
-    unzip
-    wget
-    zip
-)
-
-for package in "${SYSTEM_PACKAGES[@]}"; do
-    ensure_package "${package}"
-done
-
-###############################################################################
-# Configure Package Repositories
-###############################################################################
-
-log_section "Package Repositories"
-
-#
-# Deadsnakes is only required for Ubuntu releases that do not provide the
-# desired Python version.
-#
-
-if [[ "${CURRENT_UBUNTU}" < "26.04" ]]; then
-
-    if ! grep -Rqs "deadsnakes" \
-        /etc/apt/sources.list \
-        /etc/apt/sources.list.d \
-        2>/dev/null; then
-
-        log_info "Adding Deadsnakes Python repository..."
-
-        sudo add-apt-repository -y ppa:deadsnakes/ppa
-
-        sudo apt-get update
-
-        log_success "Deadsnakes repository added."
-
-    else
-
-        log_info "Deadsnakes repository already configured."
-
-    fi
-
-else
-
-    log_info "Ubuntu ${CURRENT_UBUNTU} provides modern Python packages."
-    log_info "Skipping Deadsnakes repository."
-
-fi
-
-#
-# Additional repositories are configured by their respective installers:
-#
-#   install/node.sh      → NodeSource
-#   install/docker.sh    → Docker CE
-#   install/database.sh  → PostgreSQL PGDG
-#   install/flutter.sh   → Flutter
-#
-
-###############################################################################
-# Cleanup
-###############################################################################
-
-log_section "Cleaning Package Cache"
-
-sudo apt-get autoremove -y
-sudo apt-get autoclean -y
-
-###############################################################################
-# Complete
-###############################################################################
-
-log_blank
-log_success "System provisioning completed successfully."
-
-exit 0
diff --git a/.devcontainer/docker/scripts/post-create.sh b/.devcontainer/docker/scripts/post-create.sh
deleted file mode 100755
index f36ea3f..0000000
--- a/.devcontainer/docker/scripts/post-create.sh
+++ /dev/null
@@ -1,111 +0,0 @@
-#!/usr/bin/env bash
-
-# ==============================================================================
-# BAOBAB Enterprise Platform
-#
-# Script      : post-create.sh
-# Purpose     : Performs post-creation configuration for GitHub Codespaces
-#               and VS Code Dev Containers.
-#
-# Author      : BAOBAB Contributors
-# License     : Apache License 2.0
-# ==============================================================================
-
-set -Eeuo pipefail
-
-SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
-
-###############################################################################
-# Helper Functions
-###############################################################################
-
-print_banner() {
-
-    local title="$1"
-
-    echo
-    echo "=================================================================="
-    printf " %-64s\n" "${title}"
-    echo "=================================================================="
-}
-
-run_step() {
-
-    local script="$1"
-
-    if [[ ! -f "${script}" ]]; then
-        echo "ERROR: Missing script:"
-        echo "  ${script}"
-        exit 1
-    fi
-
-    if [[ ! -x "${script}" ]]; then
-        chmod +x "${script}"
-    fi
-
-    echo
-    echo ">>> Running $(basename "${script}")"
-
-    "${script}"
-}
-
-###############################################################################
-# Configure Development Environment
-###############################################################################
-
-print_banner "CONFIGURE"
-
-CONFIGURE_SCRIPTS=(
-    environment.sh
-    git.sh
-    shell.sh
-    vscode.sh
-    paths.sh
-)
-
-for script in "${CONFIGURE_SCRIPTS[@]}"; do
-    run_step "${SCRIPT_DIR}/configure/${script}"
-done
-
-###############################################################################
-# Prepare Workspace
-###############################################################################
-
-print_banner "WORKSPACE"
-
-WORKSPACE_SCRIPTS=(
-    initialize.sh
-    directories.sh
-    permissions.sh
-)
-
-for script in "${WORKSPACE_SCRIPTS[@]}"; do
-    run_step "${SCRIPT_DIR}/workspace/${script}"
-done
-
-###############################################################################
-# Verify Development Environment
-###############################################################################
-
-print_banner "VERIFY"
-
-run_step "${SCRIPT_DIR}/verify.sh"
-
-###############################################################################
-# Display Summary
-###############################################################################
-
-print_banner "SUMMARY"
-
-run_step "${SCRIPT_DIR}/summary.sh"
-
-###############################################################################
-# Complete
-###############################################################################
-
-print_banner "POST-CREATE COMPLETE"
-
-echo "✓ BAOBAB Dev Container is ready."
-echo
-
-exit 0
diff --git a/.devcontainer/docker/scripts/run.sh b/.devcontainer/docker/scripts/run.sh
deleted file mode 100755
index adab40c..0000000
--- a/.devcontainer/docker/scripts/run.sh
+++ /dev/null
@@ -1,113 +0,0 @@
-#!/usr/bin/env bash
-
-# ==============================================================================
-# BAOBAB Enterprise Platform
-#
-# Script      : run.sh
-# Purpose     : Entry point for the BAOBAB Dev Container provisioning framework.
-#
-# Description :
-#   Dispatches high-level provisioning commands to the appropriate orchestration
-#   scripts. This script intentionally contains no installation or provisioning
-#   logic and serves only as the single entry point into the BAOBAB development
-#   environment.
-#
-# Author      : BAOBAB Contributors
-# License     : Apache License 2.0
-# ==============================================================================
-
-set -Eeuo pipefail
-
-SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
-
-# ------------------------------------------------------------------------------
-# Display Help
-# ------------------------------------------------------------------------------
-
-show_help() {
-
-cat <<EOF
-
-BAOBAB Dev Container
-
-Usage
-
-    run.sh [command]
-
-Commands
-
-    bootstrap      Complete machine provisioning
-
-    post-create    Configure the workspace after the container starts
-                   (Default)
-
-    verify         Execute every verification module
-
-    summary        Display an environment summary
-
-    help           Show this help
-
-EOF
-
-}
-
-# ------------------------------------------------------------------------------
-# Execute Script
-# ------------------------------------------------------------------------------
-
-run_script() {
-
-    local script="$1"
-
-    if [[ ! -f "$script" ]]; then
-        echo "ERROR: Script not found:"
-        echo "  $script"
-        exit 1
-    fi
-
-    echo
-    echo "=============================================================="
-    echo "Running: $(basename "$script")"
-    echo "=============================================================="
-    echo
-
-    # Execute explicitly with Bash to avoid dependency on executable bits.
-    bash "$script"
-}
-
-# ------------------------------------------------------------------------------
-# Command Dispatcher
-# ------------------------------------------------------------------------------
-
-case "${1:-post-create}" in
-
-    bootstrap)
-        run_script "${SCRIPT_DIR}/bootstrap.sh"
-        ;;
-
-    post-create)
-        run_script "${SCRIPT_DIR}/post-create.sh"
-        ;;
-
-    verify)
-        run_script "${SCRIPT_DIR}/verify.sh"
-        ;;
-
-    summary)
-        run_script "${SCRIPT_DIR}/summary.sh"
-        ;;
-
-    help|-h|--help)
-        show_help
-        ;;
-
-    *)
-        echo "ERROR: Unknown command: ${1}"
-        echo
-        show_help
-        exit 1
-        ;;
-
-esac
-
-exit 0
diff --git a/.devcontainer/docker/scripts/summary.sh b/.devcontainer/docker/scripts/summary.sh
deleted file mode 100755
index bdca125..0000000
--- a/.devcontainer/docker/scripts/summary.sh
+++ /dev/null
@@ -1,158 +0,0 @@
-#!/usr/bin/env bash
-
-# ==============================================================================
-# BAOBAB Enterprise Platform
-#
-# Script      : summary.sh
-# Purpose     : Displays a summary of the BAOBAB Development Environment.
-#
-# Author      : BAOBAB Contributors
-# License     : Apache License 2.0
-# ==============================================================================
-
-set -Eeuo pipefail
-
-SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
-
-# ------------------------------------------------------------------------------
-# Load Utilities
-# ------------------------------------------------------------------------------
-
-# shellcheck source=/dev/null
-source "${SCRIPT_DIR}/utils/colors.sh"
-source "${SCRIPT_DIR}/utils/functions.sh"
-source "${SCRIPT_DIR}/utils/logging.sh"
-
-PROJECT_ROOT="$(get_project_root)"
-
-log_header "BAOBAB Development Environment Summary"
-
-###############################################################################
-# Operating System
-###############################################################################
-
-log_section "Operating System"
-
-log_info "Kernel       : $(uname -s)"
-log_info "Release      : $(uname -r)"
-log_info "Architecture : $(uname -m)"
-
-###############################################################################
-# Development Tools
-###############################################################################
-
-log_section "Development Tools"
-
-tool_version() {
-    local cmd="$1"
-
-    if command_exists "${cmd}"; then
-        case "${cmd}" in
-            python3) python3 --version ;;
-            node) node --version ;;
-            npm) npm --version ;;
-            java) java --version | head -n 1 ;;
-            javac) javac --version ;;
-            flutter) flutter --version | head -n 1 ;;
-            dart) dart --version 2>&1 ;;
-            docker) docker --version ;;
-            git) git --version ;;
-            uv) uv --version ;;
-            *)
-                "${cmd}" --version 2>/dev/null | head -n 1 || true
-                ;;
-        esac
-    else
-        echo "Not Installed"
-    fi
-}
-
-TOOLS=(
-    git
-    python3
-    uv
-    node
-    npm
-    java
-    javac
-    flutter
-    dart
-    docker
-)
-
-for tool in "${TOOLS[@]}"; do
-    printf "%-12s %s\n" "${tool}" "$(tool_version "${tool}")"
-done
-
-###############################################################################
-# Repository
-###############################################################################
-
-log_section "Repository"
-
-log_info "Project Root : ${PROJECT_ROOT}"
-
-cd "${PROJECT_ROOT}"
-
-if git rev-parse --git-dir >/dev/null 2>&1; then
-
-    log_info "Branch       : $(git branch --show-current)"
-
-    if git diff --quiet; then
-        log_info "Status       : Clean"
-    else
-        log_info "Status       : Modified"
-    fi
-
-fi
-
-###############################################################################
-# Docker
-###############################################################################
-
-log_section "Docker Services"
-
-if command_exists docker && file_exists "${PROJECT_ROOT}/compose.yaml"; then
-
-    docker compose ps || true
-
-else
-
-    log_info "Docker Compose not available."
-
-fi
-
-###############################################################################
-# Workspace
-###############################################################################
-
-log_section "Workspace"
-
-FILE_COUNT="$(find "${PROJECT_ROOT}" -type f | wc -l | tr -d ' ')"
-DIRECTORY_COUNT="$(find "${PROJECT_ROOT}" -type d | wc -l | tr -d ' ')"
-
-log_info "Files        : ${FILE_COUNT}"
-log_info "Directories  : ${DIRECTORY_COUNT}"
-
-###############################################################################
-# Environment Variables
-###############################################################################
-
-log_section "Environment"
-
-for var in BAOBAB_HOME JAVA_HOME ANDROID_HOME; do
-
-    if [[ -n "${!var:-}" ]]; then
-        log_info "${var}=${!var}"
-    fi
-
-done
-
-###############################################################################
-# Finished
-###############################################################################
-
-log_blank
-log_success "BAOBAB Development Environment is ready."
-
-exit 0```````````
diff --git a/.devcontainer/docker/scripts/utils/checks.sh b/.devcontainer/docker/scripts/utils/checks.sh
deleted file mode 100755
index a52051d..0000000
--- a/.devcontainer/docker/scripts/utils/checks.sh
+++ /dev/null
@@ -1,189 +0,0 @@
-#!/usr/bin/env bash
-
-# ==============================================================================
-# BAOBAB Enterprise Platform
-#
-# Script      : checks.sh
-# Purpose     : Provides reusable validation and verification functions for the
-#               BAOBAB Development Environment Provisioning Framework.
-#
-# Author      : BAOBAB Contributors
-# License     : Apache License 2.0
-#
-# Notes
-# -----
-# This script contains reusable validation functions only.
-# It performs no installation or configuration.
-# ==============================================================================
-#
-# shellcheck disable=SC1091
-
-# Prevent multiple sourcing
-[[ -n "${BAOBAB_CHECKS_LOADED:-}" ]] && return
-readonly BAOBAB_CHECKS_LOADED=1
-
-SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
-
-# Load dependencies
-source "${SCRIPT_DIR}/colors.sh"
-source "${SCRIPT_DIR}/functions.sh"
-source "${SCRIPT_DIR}/logging.sh"
-
-###############################################################################
-# Commands
-###############################################################################
-
-check_command() {
-    local command="$1"
-    local version_command="${2:-}"
-
-    if command_exists "${command}"; then
-        if [[ -n "${version_command}" ]]; then
-            local version
-            version=$(eval "${version_command}" 2>/dev/null | head -n 1)
-            log_pass "${command} (${version})"
-        else
-            log_pass "${command}"
-        fi
-        return 0
-    fi
-
-    log_fail "${command}"
-    return 1
-}
-
-###############################################################################
-# Files
-###############################################################################
-
-check_file() {
-    local file="$1"
-
-    if file_exists "${file}"; then
-        log_pass "${file}"
-        return 0
-    fi
-
-    log_fail "${file}"
-    return 1
-}
-
-###############################################################################
-# Directories
-###############################################################################
-
-check_directory() {
-    local directory="$1"
-
-    if directory_exists "${directory}"; then
-        log_pass "${directory}"
-        return 0
-    fi
-
-    log_fail "${directory}"
-    return 1
-}
-
-###############################################################################
-# Environment Variables
-###############################################################################
-
-check_environment_variable() {
-    local variable="$1"
-
-    if [[ -n "${!variable:-}" ]]; then
-        log_pass "${variable}"
-        return 0
-    fi
-
-    log_fail "${variable}"
-    return 1
-}
-
-###############################################################################
-# Ports
-###############################################################################
-
-check_port() {
-    local port="$1"
-
-    if command_exists ss; then
-        if ss -ltn | awk '{print $4}' | grep -q ":${port}$"; then
-            log_pass "Port ${port}"
-            return 0
-        fi
-    elif command_exists netstat; then
-        if netstat -ltn 2>/dev/null | awk '{print $4}' | grep -q ":${port}$"; then
-            log_pass "Port ${port}"
-            return 0
-        fi
-    fi
-
-    log_fail "Port ${port}"
-    return 1
-}
-
-###############################################################################
-# Docker
-###############################################################################
-
-check_container() {
-    local container="$1"
-
-    if command_exists docker &&
-       docker ps --format '{{.Names}}' | grep -Fxq "${container}"; then
-        log_pass "Container ${container}"
-        return 0
-    fi
-
-    log_fail "Container ${container}"
-    return 1
-}
-
-check_network() {
-    local network="$1"
-
-    if command_exists docker &&
-       docker network ls --format '{{.Name}}' | grep -Fxq "${network}"; then
-        log_pass "Network ${network}"
-        return 0
-    fi
-
-    log_fail "Network ${network}"
-    return 1
-}
-
-###############################################################################
-# Summary Counters
-###############################################################################
-
-BAOBAB_CHECKS_PASSED=0
-BAOBAB_CHECKS_FAILED=0
-
-record_check_result() {
-    local status="$1"
-
-    if [[ "${status}" -eq 0 ]]; then
-        ((BAOBAB_CHECKS_PASSED++))
-    else
-        ((BAOBAB_CHECKS_FAILED++))
-    fi
-}
-
-print_check_summary() {
-    log_blank
-    log_section "Verification Summary"
-
-    printf "Passed : %d\n" "${BAOBAB_CHECKS_PASSED}"
-    printf "Failed : %d\n" "${BAOBAB_CHECKS_FAILED}"
-
-    log_blank
-
-    if [[ "${BAOBAB_CHECKS_FAILED}" -eq 0 ]]; then
-        log_success "All verification checks passed."
-        return 0
-    fi
-
-    log_error "One or more verification checks failed."
-    return 1
-}
diff --git a/.devcontainer/docker/scripts/utils/colors.sh b/.devcontainer/docker/scripts/utils/colors.sh
deleted file mode 100755
index 617510d..0000000
--- a/.devcontainer/docker/scripts/utils/colors.sh
+++ /dev/null
@@ -1,91 +0,0 @@
-#!/usr/bin/env bash
-
-# ==============================================================================
-# BAOBAB Enterprise Platform
-#
-# Script      : colors.sh
-# Purpose     : Defines ANSI colour and text formatting constants used by the
-#               BAOBAB Development Environment Provisioning Framework.
-#
-# Author      : BAOBAB Contributors
-# License     : Apache License 2.0
-#
-# Notes
-# -----
-# This script only defines colours and formatting.
-# It intentionally contains no functions or business logic.
-# ==============================================================================
-
-# Prevent multiple sourcing
-[[ -n "${BAOBAB_COLORS_LOADED:-}" ]] && return
-readonly BAOBAB_COLORS_LOADED=1
-
-###############################################################################
-# Colour Support
-###############################################################################
-
-if [[ -t 1 ]] && command -v tput >/dev/null 2>&1; then
-    readonly COLOR_SUPPORTED=true
-else
-    readonly COLOR_SUPPORTED=false
-fi
-
-###############################################################################
-# ANSI Colours
-###############################################################################
-
-if [[ "${COLOR_SUPPORTED}" == true ]]; then
-
-    readonly RESET="\033[0m"
-
-    readonly BLACK="\033[30m"
-    readonly RED="\033[31m"
-    readonly GREEN="\033[32m"
-    readonly YELLOW="\033[33m"
-    readonly BLUE="\033[34m"
-    readonly MAGENTA="\033[35m"
-    readonly CYAN="\033[36m"
-    readonly WHITE="\033[37m"
-
-    readonly BRIGHT_BLACK="\033[90m"
-    readonly BRIGHT_RED="\033[91m"
-    readonly BRIGHT_GREEN="\033[92m"
-    readonly BRIGHT_YELLOW="\033[93m"
-    readonly BRIGHT_BLUE="\033[94m"
-    readonly BRIGHT_MAGENTA="\033[95m"
-    readonly BRIGHT_CYAN="\033[96m"
-    readonly BRIGHT_WHITE="\033[97m"
-
-    readonly BOLD="\033[1m"
-    readonly DIM="\033[2m"
-    readonly ITALIC="\033[3m"
-    readonly UNDERLINE="\033[4m"
-
-else
-
-    readonly RESET=""
-
-    readonly BLACK=""
-    readonly RED=""
-    readonly GREEN=""
-    readonly YELLOW=""
-    readonly BLUE=""
-    readonly MAGENTA=""
-    readonly CYAN=""
-    readonly WHITE=""
-
-    readonly BRIGHT_BLACK=""
-    readonly BRIGHT_RED=""
-    readonly BRIGHT_GREEN=""
-    readonly BRIGHT_YELLOW=""
-    readonly BRIGHT_BLUE=""
-    readonly BRIGHT_MAGENTA=""
-    readonly BRIGHT_CYAN=""
-    readonly BRIGHT_WHITE=""
-
-    readonly BOLD=""
-    readonly DIM=""
-    readonly ITALIC=""
-    readonly UNDERLINE=""
-
-fi
diff --git a/.devcontainer/docker/scripts/utils/functions.sh b/.devcontainer/docker/scripts/utils/functions.sh
deleted file mode 100755
index 2eae05f..0000000
--- a/.devcontainer/docker/scripts/utils/functions.sh
+++ /dev/null
@@ -1,248 +0,0 @@
-#!/usr/bin/env bash
-
-# ==============================================================================
-# BAOBAB Enterprise Platform
-#
-# Script      : functions.sh
-# Purpose     : Provides reusable helper functions for the BAOBAB Development
-#               Environment Provisioning Framework.
-#
-# Author      : BAOBAB Contributors
-# License     : Apache License 2.0
-#
-# Notes
-# -----
-# This script contains generic helper functions only.
-# It intentionally contains no logging, installation, verification,
-# or provisioning logic.
-# ==============================================================================
-
-# Prevent multiple sourcing
-[[ -n "${BAOBAB_FUNCTIONS_LOADED:-}" ]] && return
-readonly BAOBAB_FUNCTIONS_LOADED=1
-
-###############################################################################
-# Paths
-###############################################################################
-
-get_script_directory() {
-
-    local source="${BASH_SOURCE[1]}"
-
-    while [[ -h "${source}" ]]; do
-
-        local dir
-
-        dir="$(cd -P "$(dirname "${source}")" >/dev/null 2>&1 && pwd)"
-
-        source="$(readlink "${source}")"
-
-        [[ "${source}" != /* ]] && source="${dir}/${source}"
-
-    done
-
-    cd -P "$(dirname "${source}")" >/dev/null 2>&1 && pwd
-
-}
-
-get_project_root() {
-
-    local script_dir
-    local project_root
-
-    script_dir="$(get_script_directory)"
-
-    # --------------------------------------------------------------------------
-    # Preferred Method
-    # --------------------------------------------------------------------------
-    # Ask Git for the repository root. This is reliable regardless of the
-    # current working directory.
-    # --------------------------------------------------------------------------
-
-    if command -v git >/dev/null 2>&1; then
-
-        if project_root="$(git -C "${script_dir}" rev-parse --show-toplevel 2>/dev/null)"; then
-
-            printf '%s\n' "${project_root}"
-            return 0
-
-        fi
-
-    fi
-
-    # --------------------------------------------------------------------------
-    # Fallback Method
-    # --------------------------------------------------------------------------
-    #
-    # Repository Layout
-    #
-    # baobab/
-    # └── .devcontainer/
-    #     └── docker/
-    #         └── scripts/
-    #             ├── configure/
-    #             ├── installers/
-    #             ├── utilities/
-    #             └── utils/
-    #                 └── functions.sh
-    #
-    # utils -> scripts -> docker -> .devcontainer -> baobab
-    # --------------------------------------------------------------------------
-
-    project_root="$(
-        cd "${script_dir}/../../../.." >/dev/null 2>&1
-        pwd
-    )"
-
-    # --------------------------------------------------------------------------
-    # Validation
-    # --------------------------------------------------------------------------
-
-    if [[ ! -d "${project_root}/.devcontainer" ]]; then
-
-        printf 'ERROR: Unable to determine the BAOBAB project root.\n' >&2
-        printf 'Resolved path: %s\n' "${project_root}" >&2
-
-        return 1
-
-    fi
-
-    printf '%s\n' "${project_root}"
-
-}
-
-###############################################################################
-# Time
-###############################################################################
-
-timestamp() {
-
-    date +"%Y-%m-%d %H:%M:%S"
-
-}
-
-###############################################################################
-# File System
-###############################################################################
-
-command_exists() {
-
-    command -v "$1" >/dev/null 2>&1
-
-}
-
-file_exists() {
-
-    [[ -f "$1" ]]
-
-}
-
-directory_exists() {
-
-    [[ -d "$1" ]]
-
-}
-
-ensure_directory() {
-
-    mkdir -p "$1"
-
-}
-
-ensure_executable() {
-
-    chmod +x "$1"
-
-}
-
-###############################################################################
-# String Utilities
-###############################################################################
-
-trim() {
-
-    local value="$*"
-
-    value="${value#"${value%%[![:space:]]*}"}"
-    value="${value%"${value##*[![:space:]]}"}"
-
-    printf '%s\n' "$value"
-
-}
-
-join_by() {
-
-    local delimiter="$1"
-
-    shift
-
-    local first=1
-
-    for item in "$@"; do
-
-        if [[ ${first} -eq 1 ]]; then
-
-            printf "%s" "$item"
-
-            first=0
-
-        else
-
-            printf "%s%s" "$delimiter" "$item"
-
-        fi
-
-    done
-
-    printf "\n"
-
-}
-
-###############################################################################
-# Environment
-###############################################################################
-
-is_linux() {
-
-    [[ "$(uname -s)" == "Linux" ]]
-
-}
-
-is_macos() {
-
-    [[ "$(uname -s)" == "Darwin" ]]
-
-}
-
-is_root() {
-
-    [[ "${EUID}" -eq 0 ]]
-
-}
-
-###############################################################################
-# Execution
-###############################################################################
-
-run_command() {
-
-    "$@"
-
-}
-
-safe_source() {
-
-    local script="$1"
-
-    if [[ -f "${script}" ]]; then
-
-        # shellcheck source=/dev/null
-        source "${script}"
-
-    else
-
-        return 1
-
-    fi
-
-}
diff --git a/.devcontainer/docker/scripts/utils/logging.sh b/.devcontainer/docker/scripts/utils/logging.sh
deleted file mode 100755
index 4581d6a..0000000
--- a/.devcontainer/docker/scripts/utils/logging.sh
+++ /dev/null
@@ -1,114 +0,0 @@
-#!/usr/bin/env bash
-
-# ==============================================================================
-# BAOBAB Enterprise Platform
-#
-# Script      : logging.sh
-# Purpose     : Provides standardized logging functions for the BAOBAB
-#               Development Environment Provisioning Framework.
-#
-# Author      : BAOBAB Contributors
-# License     : Apache License 2.0
-#
-# Notes
-# -----
-# This script provides logging utilities only.
-# ==============================================================================
-#
-# shellcheck disable=SC1091
-
-# Prevent multiple sourcing
-[[ -n "${BAOBAB_LOGGING_LOADED:-}" ]] && return
-readonly BAOBAB_LOGGING_LOADED=1
-
-SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
-
-# Load dependencies
-source "${SCRIPT_DIR}/colors.sh"
-source "${SCRIPT_DIR}/functions.sh"
-
-###############################################################################
-# Logging
-###############################################################################
-
-log() {
-    local level="$1"
-    local colour="$2"
-    shift 2
-
-    printf "%b[%s] [%s]%b %s\n" \
-        "${colour}" \
-        "$(timestamp)" \
-        "${level}" \
-        "${RESET}" \
-        "$*"
-}
-
-###############################################################################
-# Standard Messages
-###############################################################################
-
-log_info() {
-    log "INFO" "${CYAN}" "$@"
-}
-
-log_success() {
-    log "SUCCESS" "${GREEN}" "$@"
-}
-
-log_warning() {
-    log "WARNING" "${YELLOW}" "$@"
-}
-
-log_error() {
-    log "ERROR" "${RED}" "$@"
-}
-
-log_debug() {
-    log "DEBUG" "${MAGENTA}" "$@"
-}
-
-###############################################################################
-# Headings
-###############################################################################
-
-log_header() {
-    echo
-    printf "%b============================================================%b\n" "${BOLD}${BLUE}" "${RESET}"
-    printf "%b%s%b\n" "${BOLD}${BLUE}" "$*" "${RESET}"
-    printf "%b============================================================%b\n" "${BOLD}${BLUE}" "${RESET}"
-    echo
-}
-
-log_section() {
-    echo
-    printf "%b---- %s ----%b\n" "${BOLD}${CYAN}" "$*" "${RESET}"
-}
-
-log_subsection() {
-    printf "%b• %s%b\n" "${BOLD}" "$*" "${RESET}"
-}
-
-###############################################################################
-# Status
-###############################################################################
-
-log_pass() {
-    printf "%b✔%b %s\n" "${GREEN}" "${RESET}" "$*"
-}
-
-log_fail() {
-    printf "%b✘%b %s\n" "${RED}" "${RESET}" "$*"
-}
-
-log_skip() {
-    printf "%b➜%b %s\n" "${YELLOW}" "${RESET}" "$*"
-}
-
-###############################################################################
-# Blank Line
-###############################################################################
-
-log_blank() {
-    echo
-}
diff --git a/.devcontainer/docker/scripts/verify.sh b/.devcontainer/docker/scripts/verify.sh
deleted file mode 100755
index c6f77c7..0000000
--- a/.devcontainer/docker/scripts/verify.sh
+++ /dev/null
@@ -1,92 +0,0 @@
-#!/usr/bin/env bash
-
-# ==============================================================================
-# BAOBAB Enterprise Platform
-#
-# Script      : verify.sh
-# Purpose     : Executes all BAOBAB verification modules.
-#
-# Author      : BAOBAB Contributors
-# License     : Apache License 2.0
-# ==============================================================================
-
-set -Eeuo pipefail
-
-SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
-
-###############################################################################
-# Helper
-###############################################################################
-
-run_step() {
-
-    local script="$1"
-
-    if [[ ! -f "${script}" ]]; then
-        echo "Missing verification script: ${script}"
-        exit 1
-    fi
-
-    chmod +x "${script}"
-
-    echo
-    echo "------------------------------------------------------------"
-    echo "Running $(basename "${script}")"
-    echo "------------------------------------------------------------"
-
-    "${script}"
-}
-
-###############################################################################
-# Verification Modules
-###############################################################################
-
-VERIFY_SCRIPTS=(
-    system.sh
-    python.sh
-    node.sh
-    java.sh
-    flutter.sh
-    docker.sh
-    database.sh
-    workspace.sh
-)
-
-echo
-echo "============================================================"
-echo " BAOBAB Development Environment Verification"
-echo "============================================================"
-
-FAILED=0
-
-for script in "${VERIFY_SCRIPTS[@]}"; do
-
-    if ! run_step "${SCRIPT_DIR}/verify/${script}"; then
-        FAILED=1
-    fi
-
-done
-
-###############################################################################
-# Final Result
-###############################################################################
-
-echo
-
-if [[ "${FAILED}" -eq 0 ]]; then
-
-    echo "============================================================"
-    echo " All verification modules completed successfully."
-    echo "============================================================"
-
-else
-
-    echo "============================================================"
-    echo " One or more verification modules failed."
-    echo "============================================================"
-
-    exit 1
-
-fi
-
-exit 0
diff --git a/.devcontainer/docker/scripts/verify/database.sh b/.devcontainer/docker/scripts/verify/database.sh
deleted file mode 100755
index d5fe338..0000000
--- a/.devcontainer/docker/scripts/verify/database.sh
+++ /dev/null
@@ -1,200 +0,0 @@
-#!/usr/bin/env bash
-
-# ==============================================================================
-# BAOBAB Enterprise Platform
-#
-# Script      : database.sh
-# Purpose     : Verifies the BAOBAB data services.
-#
-# Author      : BAOBAB Contributors
-# License     : Apache License 2.0
-# ==============================================================================
-
-set -Eeuo pipefail
-
-SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
-
-# ------------------------------------------------------------------------------
-# Load Utilities
-# ------------------------------------------------------------------------------
-
-# shellcheck source=/dev/null
-source "${SCRIPT_DIR}/../utils/colors.sh"
-source "${SCRIPT_DIR}/../utils/functions.sh"
-source "${SCRIPT_DIR}/../utils/logging.sh"
-source "${SCRIPT_DIR}/../utils/checks.sh"
-
-log_header "Database Verification"
-
-PROJECT_ROOT="$(get_project_root)"
-
-FAILED=0
-
-###############################################################################
-# PostgreSQL Client
-###############################################################################
-
-log_section "PostgreSQL"
-
-if command_exists psql; then
-
-    log_pass "psql"
-    log_info "$(psql --version)"
-
-else
-
-    log_warning "PostgreSQL client not installed."
-
-fi
-
-###############################################################################
-# Redis Client
-###############################################################################
-
-log_section "Redis"
-
-if command_exists redis-cli; then
-
-    log_pass "redis-cli"
-    log_info "$(redis-cli --version)"
-
-else
-
-    log_warning "Redis client not installed."
-
-fi
-
-###############################################################################
-# Docker Compose Services
-###############################################################################
-
-log_section "Docker Services"
-
-cd "${PROJECT_ROOT}"
-
-if ! file_exists "compose.yaml"; then
-
-    log_fail "compose.yaml not found."
-    exit 1
-
-fi
-
-SERVICES=(
-    db
-    redis
-    minio
-)
-
-for service in "${SERVICES[@]}"; do
-
-    if docker compose ps --services | grep -qx "${service}"; then
-
-        if [[ "$(docker compose ps -q "${service}")" != "" ]]; then
-
-            STATUS="$(docker inspect \
-                --format='{{.State.Status}}' \
-                "$(docker compose ps -q "${service}")")"
-
-            if [[ "${STATUS}" == "running" ]]; then
-
-                log_pass "${service}"
-
-            else
-
-                log_warning "${service} (${STATUS})"
-
-            fi
-
-        else
-
-            log_warning "${service} not created."
-
-        fi
-
-    else
-
-        log_info "${service} not defined."
-
-    fi
-
-done
-
-###############################################################################
-# PostgreSQL Health
-###############################################################################
-
-log_section "PostgreSQL Health"
-
-if docker compose ps --services | grep -qx "db"; then
-
-    if docker compose exec -T db pg_isready >/dev/null 2>&1; then
-
-        log_pass "PostgreSQL accepting connections."
-
-    else
-
-        log_warning "PostgreSQL not ready."
-
-    fi
-
-fi
-
-###############################################################################
-# Redis Health
-###############################################################################
-
-log_section "Redis Health"
-
-if docker compose ps --services | grep -qx "redis"; then
-
-    if docker compose exec -T redis redis-cli ping >/dev/null 2>&1; then
-
-        log_pass "Redis responding."
-
-    else
-
-        log_warning "Redis not responding."
-
-    fi
-
-fi
-
-###############################################################################
-# MinIO Health
-###############################################################################
-
-log_section "MinIO"
-
-if docker compose ps --services | grep -qx "minio"; then
-
-    if docker compose exec -T minio sh -c "test -d /data" >/dev/null 2>&1; then
-
-        log_pass "MinIO data volume."
-
-    else
-
-        log_warning "MinIO volume unavailable."
-
-    fi
-
-fi
-
-###############################################################################
-# Summary
-###############################################################################
-
-if [[ "${FAILED}" -eq 0 ]]; then
-
-    log_blank
-    log_success "Database verification completed successfully."
-
-else
-
-    log_blank
-    log_error "Database verification failed."
-
-    exit 1
-
-fi
-
-exit 0
diff --git a/.devcontainer/docker/scripts/verify/docker.sh b/.devcontainer/docker/scripts/verify/docker.sh
deleted file mode 100755
index ee8f56d..0000000
--- a/.devcontainer/docker/scripts/verify/docker.sh
+++ /dev/null
@@ -1,195 +0,0 @@
-#!/usr/bin/env bash
-
-# ==============================================================================
-# BAOBAB Enterprise Platform
-#
-# Script      : docker.sh
-# Purpose     : Verifies the Docker development environment.
-#
-# Author      : BAOBAB Contributors
-# License     : Apache License 2.0
-# ==============================================================================
-
-set -Eeuo pipefail
-
-SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
-
-# ------------------------------------------------------------------------------
-# Load Utilities
-# ------------------------------------------------------------------------------
-
-# shellcheck source=/dev/null
-source "${SCRIPT_DIR}/../utils/colors.sh"
-source "${SCRIPT_DIR}/../utils/functions.sh"
-source "${SCRIPT_DIR}/../utils/logging.sh"
-source "${SCRIPT_DIR}/../utils/checks.sh"
-
-log_header "Docker Verification"
-
-PROJECT_ROOT="$(get_project_root)"
-
-FAILED=0
-
-###############################################################################
-# Docker CLI
-###############################################################################
-
-log_section "Docker CLI"
-
-if command_exists docker; then
-
-    log_pass "Docker"
-    log_info "$(docker --version)"
-
-else
-
-    log_fail "Docker"
-    FAILED=1
-
-fi
-
-###############################################################################
-# Docker Compose
-###############################################################################
-
-log_section "Docker Compose"
-
-if docker compose version >/dev/null 2>&1; then
-
-    log_pass "Docker Compose"
-    log_info "$(docker compose version)"
-
-else
-
-    log_fail "Docker Compose"
-    FAILED=1
-
-fi
-
-###############################################################################
-# Docker Buildx
-###############################################################################
-
-log_section "Docker Buildx"
-
-if docker buildx version >/dev/null 2>&1; then
-
-    log_pass "Buildx"
-    log_info "$(docker buildx version | head -n 1)"
-
-else
-
-    log_warning "Docker Buildx not available."
-
-fi
-
-###############################################################################
-# Docker Engine
-###############################################################################
-
-log_section "Docker Engine"
-
-if docker info >/dev/null 2>&1; then
-
-    log_pass "Docker Engine"
-
-else
-
-    log_fail "Docker Engine"
-    FAILED=1
-
-fi
-
-###############################################################################
-# Docker Context
-###############################################################################
-
-log_section "Docker Context"
-
-CURRENT_CONTEXT="$(docker context show)"
-
-log_info "Current Context : ${CURRENT_CONTEXT}"
-
-###############################################################################
-# Docker Compose Configuration
-###############################################################################
-
-log_section "compose.yaml"
-
-COMPOSE_FILE="${PROJECT_ROOT}/compose.yaml"
-
-if file_exists "${COMPOSE_FILE}"; then
-
-    cd "${PROJECT_ROOT}"
-
-    if docker compose config >/dev/null 2>&1; then
-
-        log_pass "compose.yaml"
-
-    else
-
-        log_fail "compose.yaml"
-
-        FAILED=1
-
-    fi
-
-else
-
-    log_fail "compose.yaml"
-
-    FAILED=1
-
-fi
-
-###############################################################################
-# Running Services
-###############################################################################
-
-log_section "Containers"
-
-if docker compose ps >/dev/null 2>&1; then
-
-    docker compose ps
-
-else
-
-    log_warning "No running compose project."
-
-fi
-
-###############################################################################
-# Docker Images
-###############################################################################
-
-log_section "Images"
-
-docker images --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}" || true
-
-###############################################################################
-# Docker Networks
-###############################################################################
-
-log_section "Networks"
-
-docker network ls || true
-
-###############################################################################
-# Summary
-###############################################################################
-
-if [[ "${FAILED}" -eq 0 ]]; then
-
-    log_blank
-    log_success "Docker verification completed successfully."
-
-else
-
-    log_blank
-    log_error "Docker verification failed."
-
-    exit 1
-
-fi
-
-exit 0
diff --git a/.devcontainer/docker/scripts/verify/flutter.sh b/.devcontainer/docker/scripts/verify/flutter.sh
deleted file mode 100755
index edffeed..0000000
--- a/.devcontainer/docker/scripts/verify/flutter.sh
+++ /dev/null
@@ -1,186 +0,0 @@
-#!/usr/bin/env bash
-
-# ==============================================================================
-# BAOBAB Enterprise Platform
-#
-# Script      : flutter.sh
-# Purpose     : Verifies the Flutter and Dart development environment.
-#
-# Author      : BAOBAB Contributors
-# License     : Apache License 2.0
-# ==============================================================================
-
-set -Eeuo pipefail
-
-SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
-
-# ------------------------------------------------------------------------------
-# Load Utilities
-# ------------------------------------------------------------------------------
-
-# shellcheck source=/dev/null
-source "${SCRIPT_DIR}/../utils/colors.sh"
-source "${SCRIPT_DIR}/../utils/functions.sh"
-source "${SCRIPT_DIR}/../utils/logging.sh"
-source "${SCRIPT_DIR}/../utils/checks.sh"
-
-log_header "Flutter Verification"
-
-PROJECT_ROOT="$(get_project_root)"
-MOBILE_DIR="${PROJECT_ROOT}/services/mobile"
-
-FAILED=0
-
-###############################################################################
-# Flutter SDK
-###############################################################################
-
-log_section "Flutter SDK"
-
-if command_exists flutter; then
-
-    log_pass "Flutter"
-    log_info "$(flutter --version | head -n 1)"
-
-else
-
-    log_fail "Flutter"
-    FAILED=1
-
-fi
-
-###############################################################################
-# Dart SDK
-###############################################################################
-
-log_section "Dart SDK"
-
-if command_exists dart; then
-
-    log_pass "Dart"
-    log_info "$(dart --version 2>&1)"
-
-else
-
-    log_fail "Dart"
-    FAILED=1
-
-fi
-
-###############################################################################
-# Mobile Project
-###############################################################################
-
-log_section "Mobile Project"
-
-if directory_exists "${MOBILE_DIR}"; then
-
-    log_pass "services/mobile"
-
-    if file_exists "${MOBILE_DIR}/pubspec.yaml"; then
-        log_pass "pubspec.yaml"
-    else
-        log_warning "pubspec.yaml not found."
-    fi
-
-else
-
-    log_warning "Mobile project has not been initialized."
-
-fi
-
-###############################################################################
-# Flutter Doctor
-###############################################################################
-
-log_section "Flutter Doctor"
-
-if command_exists flutter; then
-
-    if flutter doctor >/dev/null 2>&1; then
-
-        log_pass "Flutter Doctor"
-
-    else
-
-        log_warning "Flutter Doctor reported issues."
-        flutter doctor
-
-    fi
-
-fi
-
-###############################################################################
-# Android SDK
-###############################################################################
-
-log_section "Android SDK"
-
-if [[ -n "${ANDROID_HOME:-}" ]]; then
-
-    log_pass "ANDROID_HOME"
-    log_info "${ANDROID_HOME}"
-
-    if [[ -d "${ANDROID_HOME}" ]]; then
-        log_pass "Android SDK"
-    else
-        log_warning "ANDROID_HOME directory does not exist."
-    fi
-
-else
-
-    log_info "ANDROID_HOME not configured."
-
-fi
-
-###############################################################################
-# Android Platform Tools
-###############################################################################
-
-log_section "Platform Tools"
-
-if command_exists adb; then
-
-    log_pass "adb"
-    log_info "$(adb version | head -n 1)"
-
-else
-
-    log_info "adb not installed."
-
-fi
-
-###############################################################################
-# Emulator
-###############################################################################
-
-log_section "Android Emulator"
-
-if command_exists emulator; then
-
-    log_pass "Android Emulator"
-
-else
-
-    log_info "Android Emulator not installed."
-
-fi
-
-###############################################################################
-# Summary
-###############################################################################
-
-if [[ "${FAILED}" -eq 0 ]]; then
-
-    log_blank
-    log_success "Flutter verification completed successfully."
-
-else
-
-    log_blank
-    log_error "Flutter verification failed."
-    exit 1
-
-fi
-
-exit 0
diff --git a/.devcontainer/docker/scripts/verify/java.sh b/.devcontainer/docker/scripts/verify/java.sh
deleted file mode 100755
index 2b7b7a3..0000000
--- a/.devcontainer/docker/scripts/verify/java.sh
+++ /dev/null
@@ -1,159 +0,0 @@
-#!/usr/bin/env bash
-
-# ==============================================================================
-# BAOBAB Enterprise Platform
-#
-# Script      : java.sh
-# Purpose     : Verifies the Java Development Kit (JDK) installation and
-#               configuration for the BAOBAB Development Environment.
-#
-# Author      : BAOBAB Contributors
-# License     : Apache License 2.0
-# ==============================================================================
-
-set -Eeuo pipefail
-
-SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
-
-# ------------------------------------------------------------------------------
-# Load Utilities
-# ------------------------------------------------------------------------------
-
-# shellcheck source=/dev/null
-source "${SCRIPT_DIR}/../utils/colors.sh"
-source "${SCRIPT_DIR}/../utils/functions.sh"
-source "${SCRIPT_DIR}/../utils/logging.sh"
-source "${SCRIPT_DIR}/../utils/checks.sh"
-
-log_header "Java Verification"
-
-FAILED=0
-
-###############################################################################
-# Java Runtime
-###############################################################################
-
-log_section "Java Runtime"
-
-if command_exists java; then
-
-    log_pass "Java Runtime"
-    log_info "$(java --version | head -n 1)"
-
-else
-
-    log_fail "Java Runtime"
-    FAILED=1
-
-fi
-
-###############################################################################
-# Java Compiler
-###############################################################################
-
-log_section "Java Compiler"
-
-if command_exists javac; then
-
-    log_pass "Java Compiler"
-    log_info "$(javac --version)"
-
-else
-
-    log_fail "Java Compiler"
-    FAILED=1
-
-fi
-
-###############################################################################
-# JAVA_HOME
-###############################################################################
-
-log_section "JAVA_HOME"
-
-if [[ -n "${JAVA_HOME:-}" ]]; then
-
-    log_pass "JAVA_HOME"
-    log_info "${JAVA_HOME}"
-
-    if [[ -d "${JAVA_HOME}" ]]; then
-        log_pass "JAVA_HOME directory"
-    else
-        log_warning "JAVA_HOME points to a non-existent directory."
-    fi
-
-else
-
-    log_warning "JAVA_HOME is not configured."
-
-fi
-
-###############################################################################
-# PATH Verification
-###############################################################################
-
-log_section "PATH"
-
-JAVA_BINARY="$(command -v java || true)"
-
-if [[ -n "${JAVA_BINARY}" ]]; then
-
-    log_info "Java executable: ${JAVA_BINARY}"
-
-else
-
-    log_warning "Java executable not found on PATH."
-
-fi
-
-###############################################################################
-# JDK Verification
-###############################################################################
-
-log_section "JDK"
-
-if command_exists jar; then
-
-    log_pass "jar"
-    log_info "$(jar --version)"
-
-else
-
-    log_warning "jar utility not available."
-
-fi
-
-###############################################################################
-# Android Compatibility
-###############################################################################
-
-log_section "Android Toolchain"
-
-if [[ -n "${ANDROID_HOME:-}" ]]; then
-
-    log_info "ANDROID_HOME=${ANDROID_HOME}"
-
-else
-
-    log_info "ANDROID_HOME is not configured (expected until Android SDK is installed)."
-
-fi
-
-###############################################################################
-# Summary
-###############################################################################
-
-if [[ "${FAILED}" -eq 0 ]]; then
-
-    log_blank
-    log_success "Java verification completed successfully."
-
-else
-
-    log_blank
-    log_error "Java verification failed."
-    exit 1
-
-fi
-
-exit 0
diff --git a/.devcontainer/docker/scripts/verify/node.sh b/.devcontainer/docker/scripts/verify/node.sh
deleted file mode 100755
index fb0c7f5..0000000
--- a/.devcontainer/docker/scripts/verify/node.sh
+++ /dev/null
@@ -1,139 +0,0 @@
-#!/usr/bin/env bash
-
-# ==============================================================================
-# BAOBAB Enterprise Platform
-#
-# Script      : node.sh
-# Purpose     : Verifies the Node.js development environment.
-#
-# Author      : BAOBAB Contributors
-# License     : Apache License 2.0
-# ==============================================================================
-
-set -Eeuo pipefail
-
-SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
-
-# ------------------------------------------------------------------------------
-# Load Utilities
-# ------------------------------------------------------------------------------
-
-# shellcheck source=/dev/null
-source "${SCRIPT_DIR}/../utils/colors.sh"
-source "${SCRIPT_DIR}/../utils/functions.sh"
-source "${SCRIPT_DIR}/../utils/logging.sh"
-source "${SCRIPT_DIR}/../utils/checks.sh"
-
-log_header "Node.js Verification"
-
-PROJECT_ROOT="$(get_project_root)"
-
-FAILED=0
-FRONTEND_DIR="${PROJECT_ROOT}/services/frontend"
-
-###############################################################################
-# Node.js
-###############################################################################
-
-log_section "Node.js"
-
-if command_exists node; then
-    log_pass "Node.js"
-    log_info "$(node --version)"
-else
-    log_fail "Node.js"
-    FAILED=1
-fi
-
-###############################################################################
-# npm
-###############################################################################
-
-log_section "npm"
-
-if command_exists npm; then
-    log_pass "npm"
-    log_info "$(npm --version)"
-else
-    log_fail "npm"
-    FAILED=1
-fi
-
-###############################################################################
-# Corepack
-###############################################################################
-
-log_section "Corepack"
-
-if command_exists corepack; then
-    log_pass "Corepack"
-    log_info "$(corepack --version)"
-else
-    log_warning "Corepack not installed."
-fi
-
-###############################################################################
-# Frontend Project
-###############################################################################
-
-log_section "Frontend Project"
-
-if directory_exists "${FRONTEND_DIR}"; then
-
-    log_pass "services/frontend"
-
-    if file_exists "${FRONTEND_DIR}/package.json"; then
-        log_pass "package.json"
-    else
-        log_warning "package.json not found."
-    fi
-
-    if file_exists "${FRONTEND_DIR}/package-lock.json"; then
-        log_pass "package-lock.json"
-    else
-        log_info "package-lock.json not found."
-    fi
-
-    if directory_exists "${FRONTEND_DIR}/node_modules"; then
-        log_pass "node_modules"
-    else
-        log_warning "Dependencies have not been installed."
-    fi
-
-else
-
-    log_warning "Frontend project has not been initialized."
-
-fi
-
-###############################################################################
-# TypeScript
-###############################################################################
-
-log_section "TypeScript"
-
-if [[ -x "${FRONTEND_DIR}/node_modules/.bin/tsc" ]]; then
-
-    log_pass "TypeScript"
-    log_info "$("${FRONTEND_DIR}/node_modules/.bin/tsc" --version)"
-
-else
-
-    log_warning "TypeScript compiler not available."
-
-fi
-
-###############################################################################
-# Summary
-###############################################################################
-
-if [[ "${FAILED}" -eq 0 ]]; then
-    log_blank
-    log_success "Node.js verification completed successfully."
-else
-    log_blank
-    log_error "Node.js verification failed."
-    exit 1
-fi
-
-exit 0
diff --git a/.devcontainer/docker/scripts/verify/python.sh b/.devcontainer/docker/scripts/verify/python.sh
deleted file mode 100755
index 5220501..0000000
--- a/.devcontainer/docker/scripts/verify/python.sh
+++ /dev/null
@@ -1,150 +0,0 @@
-#!/usr/bin/env bash
-
-# ==============================================================================
-# BAOBAB Enterprise Platform
-#
-# Script      : python.sh
-# Purpose     : Verifies the Python development environment.
-#
-# Author      : BAOBAB Contributors
-# License     : Apache License 2.0
-# ==============================================================================
-
-set -Eeuo pipefail
-
-SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
-
-# ------------------------------------------------------------------------------
-# Load Utilities
-# ------------------------------------------------------------------------------
-
-# shellcheck source=/dev/null
-source "${SCRIPT_DIR}/../utils/colors.sh"
-source "${SCRIPT_DIR}/../utils/functions.sh"
-source "${SCRIPT_DIR}/../utils/logging.sh"
-source "${SCRIPT_DIR}/../utils/checks.sh"
-
-log_header "Python Verification"
-
-PROJECT_ROOT="$(get_project_root)"
-
-FAILED=0
-
-###############################################################################
-# Python
-###############################################################################
-
-log_section "Python"
-
-if command_exists python3; then
-    log_pass "Python3"
-    log_info "$(python3 --version)"
-else
-    log_fail "Python3"
-    FAILED=1
-fi
-
-###############################################################################
-# pip
-###############################################################################
-
-log_section "pip"
-
-if command_exists pip3; then
-    log_pass "pip3"
-    log_info "$(pip3 --version)"
-else
-    log_fail "pip3"
-    FAILED=1
-fi
-
-###############################################################################
-# uv
-###############################################################################
-
-log_section "uv"
-
-if command_exists uv; then
-    log_pass "uv"
-    log_info "$(uv --version)"
-else
-    log_fail "uv"
-    FAILED=1
-fi
-
-###############################################################################
-# Virtual Environment
-###############################################################################
-
-log_section "Virtual Environment"
-
-if [[ -d "${PROJECT_ROOT}/.venv" ]]; then
-    log_pass ".venv"
-else
-    log_warning ".venv has not been created yet."
-fi
-
-###############################################################################
-# Project Configuration
-###############################################################################
-
-log_section "Project Configuration"
-
-if file_exists "${PROJECT_ROOT}/pyproject.toml"; then
-    log_pass "pyproject.toml"
-else
-    log_warning "pyproject.toml not found."
-fi
-
-if file_exists "${PROJECT_ROOT}/uv.lock"; then
-    log_pass "uv.lock"
-else
-    log_warning "uv.lock not found."
-fi
-
-###############################################################################
-# Python Services
-###############################################################################
-
-log_section "Python Services"
-
-SERVICES=(
-    "services/backend"
-    "services/ai"
-    "services/worker"
-)
-
-for service in "${SERVICES[@]}"; do
-
-    if directory_exists "${PROJECT_ROOT}/${service}"; then
-
-        log_info "${service}"
-
-        if file_exists "${PROJECT_ROOT}/${service}/pyproject.toml"; then
-            log_pass "pyproject.toml"
-        else
-            log_warning "No pyproject.toml"
-        fi
-
-    else
-
-        log_warning "${service} not yet initialized."
-
-    fi
-
-done
-
-###############################################################################
-# Summary
-###############################################################################
-
-if [[ "${FAILED}" -eq 0 ]]; then
-    log_blank
-    log_success "Python verification completed successfully."
-else
-    log_blank
-    log_error "Python verification failed."
-    exit 1
-fi
-
-exit 0
diff --git a/.devcontainer/docker/scripts/verify/system.sh b/.devcontainer/docker/scripts/verify/system.sh
deleted file mode 100755
index 9a83389..0000000
--- a/.devcontainer/docker/scripts/verify/system.sh
+++ /dev/null
@@ -1,98 +0,0 @@
-#!/usr/bin/env bash
-
-# ==============================================================================
-# BAOBAB Enterprise Platform
-#
-# Script      : system.sh
-# Purpose     : Verifies the operating system and core development tools.
-#
-# Author      : BAOBAB Contributors
-# License     : Apache License 2.0
-# ==============================================================================
-
-set -Eeuo pipefail
-
-SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
-
-# ------------------------------------------------------------------------------
-# Load Utilities
-# ------------------------------------------------------------------------------
-
-# shellcheck source=/dev/null
-source "${SCRIPT_DIR}/../utils/colors.sh"
-source "${SCRIPT_DIR}/../utils/functions.sh"
-source "${SCRIPT_DIR}/../utils/logging.sh"
-source "${SCRIPT_DIR}/../utils/checks.sh"
-
-log_header "System Verification"
-
-###############################################################################
-# Operating System
-###############################################################################
-
-log_section "Operating System"
-
-log_info "Kernel       : $(uname -s)"
-log_info "Release      : $(uname -r)"
-log_info "Architecture : $(uname -m)"
-
-###############################################################################
-# Required Commands
-###############################################################################
-
-log_section "Core Tools"
-
-TOOLS=(
-    bash
-    curl
-    wget
-    git
-    jq
-    tree
-    unzip
-    zip
-)
-
-FAILED=0
-
-for tool in "${TOOLS[@]}"; do
-
-    if command_exists "${tool}"; then
-        log_pass "${tool}"
-    else
-        log_fail "${tool}"
-        FAILED=1
-    fi
-
-done
-
-###############################################################################
-# Available Disk Space
-###############################################################################
-
-log_section "Disk Space"
-
-df -h .
-
-###############################################################################
-# Memory
-###############################################################################
-
-log_section "Memory"
-
-free -h || true
-
-###############################################################################
-# Complete
-###############################################################################
-
-if [[ "${FAILED}" -eq 0 ]]; then
-    log_blank
-    log_success "System verification passed."
-else
-    log_blank
-    log_error "System verification failed."
-    exit 1
-fi
-
-exit 0
diff --git a/.devcontainer/docker/scripts/verify/workspace.sh b/.devcontainer/docker/scripts/verify/workspace.sh
deleted file mode 100755
index 1c71d55..0000000
--- a/.devcontainer/docker/scripts/verify/workspace.sh
+++ /dev/null
@@ -1,234 +0,0 @@
-#!/usr/bin/env bash
-
-# ==============================================================================
-# BAOBAB Enterprise Platform
-#
-# Script      : workspace.sh
-# Purpose     : Verifies the BAOBAB workspace structure and repository health.
-#
-# Author      : BAOBAB Contributors
-# License     : Apache License 2.0
-# ==============================================================================
-
-set -Eeuo pipefail
-
-SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
-
-# ------------------------------------------------------------------------------
-# Load Utilities
-# ------------------------------------------------------------------------------
-
-# shellcheck source=/dev/null
-source "${SCRIPT_DIR}/../utils/colors.sh"
-source "${SCRIPT_DIR}/../utils/functions.sh"
-source "${SCRIPT_DIR}/../utils/logging.sh"
-source "${SCRIPT_DIR}/../utils/checks.sh"
-
-log_header "Workspace Verification"
-
-PROJECT_ROOT="$(get_project_root)"
-
-FAILED=0
-
-###############################################################################
-# Repository Root
-###############################################################################
-
-log_section "Repository"
-
-ROOT_FILES=(
-    ".gitignore"
-    ".editorconfig"
-    ".env.example"
-    "README.md"
-    "LICENSE"
-    "compose.yaml"
-)
-
-for file in "${ROOT_FILES[@]}"; do
-
-    if file_exists "${PROJECT_ROOT}/${file}"; then
-
-        log_pass "${file}"
-
-    else
-
-        log_fail "${file}"
-        FAILED=1
-
-    fi
-
-done
-
-###############################################################################
-# Project Directories
-###############################################################################
-
-log_section "Project Structure"
-
-DIRECTORIES=(
-    ".devcontainer"
-    ".github"
-
-    "docs"
-
-    "services"
-
-    "packages"
-
-    "shared"
-
-    "resources"
-
-    "tests"
-
-    "infrastructure"
-)
-
-for directory in "${DIRECTORIES[@]}"; do
-
-    if directory_exists "${PROJECT_ROOT}/${directory}"; then
-
-        log_pass "${directory}"
-
-    else
-
-        log_fail "${directory}"
-        FAILED=1
-
-    fi
-
-done
-
-###############################################################################
-# Service Directories
-###############################################################################
-
-log_section "Services"
-
-SERVICES=(
-    "backend"
-    "frontend"
-    "mobile"
-    "ai"
-    "worker"
-)
-
-for service in "${SERVICES[@]}"; do
-
-    if directory_exists "${PROJECT_ROOT}/services/${service}"; then
-
-        log_pass "${service}"
-
-    else
-
-        log_warning "${service} not initialized."
-
-    fi
-
-done
-
-###############################################################################
-# Environment Files
-###############################################################################
-
-log_section "Environment"
-
-ENVIRONMENT_ROOT="${PROJECT_ROOT}/resources/environments"
-
-if directory_exists "${ENVIRONMENT_ROOT}"; then
-
-    log_pass "resources/environments"
-
-    ENVIRONMENTS=(
-        development
-        testing
-        staging
-        production
-    )
-
-    for env in "${ENVIRONMENTS[@]}"; do
-
-        if directory_exists "${ENVIRONMENT_ROOT}/${env}"; then
-
-            log_pass "${env}"
-
-        else
-
-            log_warning "${env}"
-
-        fi
-
-    done
-
-else
-
-    log_warning "resources/environments not found."
-
-fi
-
-###############################################################################
-# Git Repository
-###############################################################################
-
-log_section "Git"
-
-cd "${PROJECT_ROOT}"
-
-if git rev-parse --git-dir >/dev/null 2>&1; then
-
-    log_pass "Git repository"
-
-    BRANCH="$(git branch --show-current)"
-
-    log_info "Branch : ${BRANCH}"
-
-    if git diff --quiet; then
-
-        log_pass "Working tree clean"
-
-    else
-
-        log_warning "Working tree contains changes."
-
-    fi
-
-else
-
-    log_fail "Git repository"
-    FAILED=1
-
-fi
-
-###############################################################################
-# Workspace Statistics
-###############################################################################
-
-log_section "Workspace Statistics"
-
-FILE_COUNT="$(find "${PROJECT_ROOT}" -type f | wc -l | tr -d ' ')"
-
-DIRECTORY_COUNT="$(find "${PROJECT_ROOT}" -type d | wc -l | tr -d ' ')"
-
-log_info "Files       : ${FILE_COUNT}"
-log_info "Directories : ${DIRECTORY_COUNT}"
-
-###############################################################################
-# Summary
-###############################################################################
-
-if [[ "${FAILED}" -eq 0 ]]; then
-
-    log_blank
-    log_success "Workspace verification completed successfully."
-
-else
-
-    log_blank
-    log_error "Workspace verification failed."
-
-    exit 1
-
-fi
-
-exit 0
diff --git a/.devcontainer/docker/scripts/workspace/cleanup.sh b/.devcontainer/docker/scripts/workspace/cleanup.sh
deleted file mode 100755
index d2a20c4..0000000
--- a/.devcontainer/docker/scripts/workspace/cleanup.sh
+++ /dev/null
@@ -1,156 +0,0 @@
-#!/usr/bin/env bash
-
-# ==============================================================================
-# BAOBAB Enterprise Platform
-#
-# Script      : cleanup.sh
-# Purpose     : Cleans temporary files, caches, and provisioning artifacts from
-#               the BAOBAB Development Environment.
-#
-# Author      : BAOBAB Contributors
-# License     : Apache License 2.0
-# ==============================================================================
-
-set -Eeuo pipefail
-
-SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
-
-# ------------------------------------------------------------------------------
-# Load Utilities
-# ------------------------------------------------------------------------------
-
-# shellcheck source=/dev/null
-source "${SCRIPT_DIR}/../utils/colors.sh"
-source "${SCRIPT_DIR}/../utils/functions.sh"
-source "${SCRIPT_DIR}/../utils/logging.sh"
-
-log_header "Workspace Cleanup"
-
-PROJECT_ROOT="$(get_project_root)"
-
-###############################################################################
-# Remove Temporary Workspace Files
-###############################################################################
-
-log_section "Workspace"
-
-TEMP_DIRECTORIES=(
-    ".cache"
-    ".logs"
-    ".tmp"
-)
-
-for directory in "${TEMP_DIRECTORIES[@]}"; do
-
-    TARGET="${PROJECT_ROOT}/${directory}"
-
-    if [[ -d "${TARGET}" ]]; then
-        rm -rf "${TARGET:?}/"*
-        log_success "Cleaned ${directory}"
-    else
-        log_info "${directory} not found."
-    fi
-
-done
-
-###############################################################################
-# Python
-###############################################################################
-
-log_section "Python"
-
-find "${PROJECT_ROOT}" \
-    -type d \
-    -name "__pycache__" \
-    -prune \
-    -exec rm -rf {} +
-
-find "${PROJECT_ROOT}" \
-    -type f \
-    -name "*.pyc" \
-    -delete
-
-find "${PROJECT_ROOT}" \
-    -type f \
-    -name "*.pyo" \
-    -delete
-
-find "${PROJECT_ROOT}" \
-    -type f \
-    -name "*.pyd" \
-    -delete
-
-log_success "Python cache cleaned."
-
-###############################################################################
-# Node.js
-###############################################################################
-
-log_section "Node.js"
-
-find "${PROJECT_ROOT}" \
-    -type d \
-    -name ".next" \
-    -prune \
-    -exec rm -rf {} +
-
-find "${PROJECT_ROOT}" \
-    -type d \
-    -name ".turbo" \
-    -prune \
-    -exec rm -rf {} +
-
-log_success "Frontend cache cleaned."
-
-###############################################################################
-# Flutter
-###############################################################################
-
-log_section "Flutter"
-
-MOBILE_DIR="${PROJECT_ROOT}/services/mobile"
-
-if [[ -d "${MOBILE_DIR}" ]]; then
-
-    (
-        cd "${MOBILE_DIR}"
-
-        if command -v flutter >/dev/null 2>&1; then
-            flutter clean >/dev/null 2>&1 || true
-        fi
-    )
-
-    log_success "Flutter build artifacts cleaned."
-
-else
-
-    log_info "Mobile project not initialized."
-
-fi
-
-###############################################################################
-# Docker
-###############################################################################
-
-log_section "Docker"
-
-if command -v docker >/dev/null 2>&1; then
-
-    docker builder prune -f >/dev/null 2>&1 || true
-
-    log_success "Docker builder cache cleaned."
-
-else
-
-    log_info "Docker not available."
-
-fi
-
-###############################################################################
-# Summary
-###############################################################################
-
-log_blank
-log_success "Workspace cleanup completed successfully."
-
-exit 0
diff --git a/.devcontainer/docker/scripts/workspace/directories.sh b/.devcontainer/docker/scripts/workspace/directories.sh
deleted file mode 100755
index c8d56b6..0000000
--- a/.devcontainer/docker/scripts/workspace/directories.sh
+++ /dev/null
@@ -1,131 +0,0 @@
-#!/usr/bin/env bash
-
-# ==============================================================================
-# BAOBAB Enterprise Platform
-#
-# Script      : directories.sh
-# Purpose     : Creates and verifies the BAOBAB workspace directory structure.
-#
-# Author      : BAOBAB Contributors
-# License     : Apache License 2.0
-# ==============================================================================
-
-set -Eeuo pipefail
-
-SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
-
-# ------------------------------------------------------------------------------
-# Load Utilities
-# ------------------------------------------------------------------------------
-
-# shellcheck source=/dev/null
-source "${SCRIPT_DIR}/../utils/colors.sh"
-source "${SCRIPT_DIR}/../utils/functions.sh"
-source "${SCRIPT_DIR}/../utils/logging.sh"
-
-log_header "Workspace Directory Structure"
-
-PROJECT_ROOT="$(get_project_root)"
-
-###############################################################################
-# Directory Structure
-###############################################################################
-
-log_section "Creating Directory Structure"
-
-DIRECTORIES=(
-
-    # -------------------------------------------------------------------------
-    # Root
-    # -------------------------------------------------------------------------
-
-    ".cache"
-    ".logs"
-    ".tmp"
-
-    # -------------------------------------------------------------------------
-    # Resources
-    # -------------------------------------------------------------------------
-
-    "resources/backups"
-    "resources/sample-data"
-    "resources/seeds"
-    "resources/templates"
-
-    # -------------------------------------------------------------------------
-    # Backend
-    # -------------------------------------------------------------------------
-
-    "services/backend/media"
-    "services/backend/static"
-    "services/backend/logs"
-
-    # -------------------------------------------------------------------------
-    # Frontend
-    # -------------------------------------------------------------------------
-
-    "services/frontend/public"
-    "services/frontend/tests"
-
-    # -------------------------------------------------------------------------
-    # Mobile
-    # -------------------------------------------------------------------------
-
-    "services/mobile/test"
-
-    # -------------------------------------------------------------------------
-    # AI
-    # -------------------------------------------------------------------------
-
-    "services/ai/tests"
-
-    # -------------------------------------------------------------------------
-    # Worker
-    # -------------------------------------------------------------------------
-
-    "services/worker/tasks"
-    "services/worker/schedules"
-    "services/worker/monitoring"
-
-    # -------------------------------------------------------------------------
-    # Shared
-    # -------------------------------------------------------------------------
-
-    "shared/documentation"
-
-    # -------------------------------------------------------------------------
-    # Tests
-    # -------------------------------------------------------------------------
-
-    "tests/fixtures"
-    "tests/unit"
-    "tests/integration"
-    "tests/api"
-    "tests/security"
-)
-
-for directory in "${DIRECTORIES[@]}"; do
-
-    ensure_directory "${PROJECT_ROOT}/${directory}"
-
-done
-
-###############################################################################
-# Directory Report
-###############################################################################
-
-log_section "Workspace"
-
-DIRECTORY_COUNT="${#DIRECTORIES[@]}"
-
-log_info "Directories managed : ${DIRECTORY_COUNT}"
-
-###############################################################################
-# Complete
-###############################################################################
-
-log_blank
-log_success "Workspace directory structure verified."
-
-exit 0
-
diff --git a/.devcontainer/docker/scripts/workspace/initialize.sh b/.devcontainer/docker/scripts/workspace/initialize.sh
deleted file mode 100755
index 57384b4..0000000
--- a/.devcontainer/docker/scripts/workspace/initialize.sh
+++ /dev/null
@@ -1,145 +0,0 @@
-#!/usr/bin/env bash
-
-# ==============================================================================
-# BAOBAB Enterprise Platform
-#
-# Script      : initialize.sh
-# Purpose     : Initializes the BAOBAB workspace after the development
-#               environment has been provisioned.
-#
-# Author      : BAOBAB Contributors
-# License     : Apache License 2.0
-# ==============================================================================
-
-set -Eeuo pipefail
-
-SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
-
-# ------------------------------------------------------------------------------
-# Load Utilities
-# ------------------------------------------------------------------------------
-
-# shellcheck source=/dev/null
-source "${SCRIPT_DIR}/../utils/colors.sh"
-source "${SCRIPT_DIR}/../utils/functions.sh"
-source "${SCRIPT_DIR}/../utils/logging.sh"
-
-log_header "Workspace Initialization"
-
-PROJECT_ROOT="$(get_project_root)"
-
-###############################################################################
-# Verify Workspace
-###############################################################################
-
-log_section "Workspace"
-
-if [[ ! -d "${PROJECT_ROOT}" ]]; then
-    log_error "Workspace not found."
-    exit 1
-fi
-
-log_success "Workspace located."
-log_info "Project Root: ${PROJECT_ROOT}"
-
-###############################################################################
-# Required Directories
-###############################################################################
-
-log_section "Required Directories"
-
-DIRECTORIES=(
-    ".cache"
-    ".logs"
-    ".tmp"
-
-    "resources/backups"
-
-    "tests/fixtures"
-
-    "services/backend/media"
-
-    "services/backend/static"
-
-    "services/frontend/public"
-
-    "services/ai/tests"
-
-    "services/worker/tasks"
-)
-
-for directory in "${DIRECTORIES[@]}"; do
-
-    ensure_directory "${PROJECT_ROOT}/${directory}"
-
-done
-
-log_success "Workspace directories verified."
-
-###############################################################################
-# Required Files
-###############################################################################
-
-log_section "Repository Files"
-
-FILES=(
-    ".env"
-    ".env.example"
-    ".gitignore"
-    "README.md"
-    "LICENSE"
-    "compose.yaml"
-)
-
-for file in "${FILES[@]}"; do
-
-    if file_exists "${PROJECT_ROOT}/${file}"; then
-        log_pass "${file}"
-    else
-        log_warning "${file} is missing."
-    fi
-
-done
-
-###############################################################################
-# Git Repository
-###############################################################################
-
-log_section "Git Repository"
-
-cd "${PROJECT_ROOT}"
-
-if git rev-parse --git-dir >/dev/null 2>&1; then
-
-    log_success "Git repository initialized."
-
-else
-
-    log_warning "Git repository not initialized."
-
-fi
-
-###############################################################################
-# Python Workspace
-###############################################################################
-
-log_section "Python"
-
-if file_exists "${PROJECT_ROOT}/pyproject.toml"; then
-
-    log_success "pyproject.toml found."
-
-else
-
-    log_warning "pyproject.toml not found."
-
-fi
-
-###############################################################################
-# Complete
-###############################################################################
-
-log_blank
-log_success "Workspace initialization completed successfully."
-
-exit 0
diff --git a/.devcontainer/docker/scripts/workspace/permissions.sh b/.devcontainer/docker/scripts/workspace/permissions.sh
deleted file mode 100755
index 55d5131..0000000
--- a/.devcontainer/docker/scripts/workspace/permissions.sh
+++ /dev/null
@@ -1,126 +0,0 @@
-#!/usr/bin/env bash
-
-# ==============================================================================
-# BAOBAB Enterprise Platform
-#
-# Script      : permissions.sh
-# Purpose     : Applies standard file and directory permissions to the BAOBAB
-#               workspace.
-#
-# Author      : BAOBAB Contributors
-# License     : Apache License 2.0
-# ==============================================================================
-
-set -Eeuo pipefail
-
-SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
-
-# ------------------------------------------------------------------------------
-# Load Utilities
-# ------------------------------------------------------------------------------
-
-# shellcheck source=/dev/null
-source "${SCRIPT_DIR}/../utils/colors.sh"
-source "${SCRIPT_DIR}/../utils/functions.sh"
-source "${SCRIPT_DIR}/../utils/logging.sh"
-
-log_header "Workspace Permissions"
-
-PROJECT_ROOT="$(get_project_root)"
-
-###############################################################################
-# Validate Workspace
-###############################################################################
-
-log_section "Workspace"
-
-if [[ ! -d "${PROJECT_ROOT}" ]]; then
-    log_error "Workspace not found."
-    exit 1
-fi
-
-###############################################################################
-# Make Shell Scripts Executable
-###############################################################################
-
-log_section "Shell Scripts"
-
-SCRIPT_ROOT="${PROJECT_ROOT}/.devcontainer/docker/scripts"
-
-if [[ -d "${SCRIPT_ROOT}" ]]; then
-
-    find "${SCRIPT_ROOT}" -type f -name "*.sh" -exec chmod +x {} \;
-
-    log_success "Shell scripts are executable."
-
-else
-
-    log_warning "Scripts directory not found."
-
-fi
-
-###############################################################################
-# Make Utility Scripts Executable
-###############################################################################
-
-log_section "Repository Scripts"
-
-REPOSITORY_SCRIPTS="${PROJECT_ROOT}/infrastructure/scripts"
-
-if [[ -d "${REPOSITORY_SCRIPTS}" ]]; then
-
-    find "${REPOSITORY_SCRIPTS}" -type f -name "*.sh" -exec chmod +x {} \;
-
-    log_success "Infrastructure scripts are executable."
-
-else
-
-    log_info "No infrastructure scripts found."
-
-fi
-
-###############################################################################
-# Directory Permissions
-###############################################################################
-
-log_section "Directory Permissions"
-
-find "${PROJECT_ROOT}" \
-    -type d \
-    -exec chmod 755 {} \;
-
-log_success "Directory permissions updated."
-
-###############################################################################
-# File Permissions
-###############################################################################
-
-log_section "File Permissions"
-
-find "${PROJECT_ROOT}" \
-    -type f \
-    ! -name "*.sh" \
-    -exec chmod 644 {} \;
-
-log_success "File permissions updated."
-
-###############################################################################
-# Git Executable Bit
-###############################################################################
-
-log_section "Git Index"
-
-cd "${PROJECT_ROOT}"
-
-git update-index --refresh >/dev/null 2>&1 || true
-
-log_success "Git index refreshed."
-
-###############################################################################
-# Summary
-###############################################################################
-
-log_blank
-log_success "Workspace permissions configured successfully."
-
-exit 0
diff --git a/.devcontainer/post-create.sh b/.devcontainer/post-create.sh
new file mode 100755
index 0000000..41339c0
--- /dev/null
+++ b/.devcontainer/post-create.sh
@@ -0,0 +1,119 @@
+#!/usr/bin/env bash
+
+# ==============================================================================
+# BAOBAB Enterprise Platform
+#
+# Script  : .devcontainer/post-create.sh
+# Purpose : Repository-specific setup that runs on top of the published
+#           ghcr.io/nabhold/baobab-dev image.
+#
+# Why this exists instead of using the image's own baobab-post-create hook
+# for dependency installation:
+#
+#   The baobab-dev image's generic post-create logic installs Python
+#   dependencies with Poetry whenever it finds a root pyproject.toml. This
+#   repository standardises on `uv` (see uv.lock, Makefile, pyproject.toml
+#   "UV WORKSPACE READINESS ROADMAP") — Poetry is not used here at all.
+#   Delegating dependency installation to the image's generic hook would
+#   silently invoke the wrong package manager, so this repository owns that
+#   one step directly instead.
+#
+# Invoked by devcontainer.json as both `updateContentCommand` (the stage
+# GitHub Codespaces prebuilds snapshot) and `postCreateCommand` (a safety
+# net re-run on actual container creation). Safe to run repeatedly — every
+# step here is idempotent.
+#
+# Note on the image's own baobab-post-create / baobab-bootstrap commands:
+# both unconditionally require a workspace-local config/versions.yaml and
+# an executable config/resolve.sh (they call `die` if either is missing),
+# a convention this repository has not adopted. They are therefore not
+# used here at all — this script deliberately re-implements only the small,
+# safe subset of that behaviour (git identity priming, toolchain
+# verification) that has no such precondition. baobab-verify and
+# baobab-summary have no config/ dependency (they fall back to the image's
+# own baked-in configuration) and are used directly below.
+#
+# Author  : BAOBAB Contributors
+# License : Apache License 2.0
+# ==============================================================================
+
+set -Eeuo pipefail
+
+WORKSPACE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
+cd "${WORKSPACE_DIR}"
+
+log() {
+    printf '\n\033[1;36m[baobab:post-create]\033[0m %s\n' "$1"
+}
+
+warn() {
+    printf '\n\033[1;33m[baobab:post-create] WARNING:\033[0m %s\n' "$1"
+}
+
+# ------------------------------------------------------------------------------
+# Git identity
+# ------------------------------------------------------------------------------
+# Placeholder identity only if none is already configured (e.g. mounted from
+# the host, or set automatically by GitHub Codespaces). safe.directory avoids
+# "detected dubious ownership" errors when the container's UID does not
+# match the mounted workspace's owning UID.
+
+log "Priming Git configuration"
+
+git config --global --get user.name  >/dev/null 2>&1 || git config --global user.name  "BAOBAB Developer"
+git config --global --get user.email >/dev/null 2>&1 || git config --global user.email "dev@example.com"
+
+git config --global --get-all safe.directory 2>/dev/null | grep -Fxq "${WORKSPACE_DIR}" \
+    || git config --global --add safe.directory "${WORKSPACE_DIR}"
+
+# ------------------------------------------------------------------------------
+# Environment file
+# ------------------------------------------------------------------------------
+
+if [[ -f ".env.example" && ! -f ".env" ]]; then
+    log "Creating .env from .env.example"
+    cp .env.example .env
+else
+    log ".env already present or no .env.example found."
+fi
+
+# ------------------------------------------------------------------------------
+# Python workspace dependencies (uv)
+# ------------------------------------------------------------------------------
+# Installs only what the root uv workspace currently manages (repository
+# tooling: lint, typecheck, test, docs — see [dependency-groups] in
+# pyproject.toml). services/backend, services/ai, services/worker, and the
+# other prospective workspace members are intentionally NOT yet uv workspace
+# members (see the "UV WORKSPACE READINESS ROADMAP" in pyproject.toml), so
+# this deliberately does not attempt to sync them until that gate is lifted.
+
+if [[ -f "pyproject.toml" ]]; then
+
+    if command -v uv >/dev/null 2>&1; then
+        log "Installing root workspace dependencies with uv"
+        uv sync --all-groups
+    else
+        warn "uv not found on PATH. Skipping Python dependency installation."
+    fi
+
+else
+    log "No root pyproject.toml found. Skipping Python dependency installation."
+fi
+
+# ------------------------------------------------------------------------------
+# Toolchain verification
+# ------------------------------------------------------------------------------
+# Non-fatal: a verification issue should surface loudly, not block the
+# developer from getting a working shell.
+
+if command -v baobab-verify >/dev/null 2>&1; then
+    baobab-verify --quiet || warn "Toolchain verification reported issues. Run 'baobab-verify' for details."
+fi
+
+# ------------------------------------------------------------------------------
+# Complete
+# ------------------------------------------------------------------------------
+
+log "Repository setup complete."
+
+exit 0
diff --git a/.env.example b/.env.example
index 8b13789..c7eec3a 100644
--- a/.env.example
+++ b/.env.example
@@ -1 +1,41 @@
+###############################################################################
+# BAOBAB Enterprise Platform
+# Local Development Environment
+#
+# Copy this file to .env and adjust values as needed:
+#
+#     cp .env.example .env
+#
+# .env is git-ignored and must never be committed. Values below are safe,
+# non-secret local-development defaults only — they are not suitable for any
+# shared, staging, or production environment.
+###############################################################################

+# -----------------------------------------------------------------------------
+# PostgreSQL
+# -----------------------------------------------------------------------------
+# Required by compose.yaml's "db" service (no built-in fallback).
+
+POSTGRES_DB=baobab
+POSTGRES_USER=baobab
+POSTGRES_PASSWORD=baobab
+POSTGRES_PORT=5432
+
+# -----------------------------------------------------------------------------
+# Redis
+# -----------------------------------------------------------------------------
+
+REDIS_PORT=6379
+
+# -----------------------------------------------------------------------------
+# NGINX
+# -----------------------------------------------------------------------------
+
+NGINX_HTTP_PORT=80
+
+# -----------------------------------------------------------------------------
+# Mailpit (local SMTP capture)
+# -----------------------------------------------------------------------------
+
+SMTP_PORT=1025
+MAILPIT_UI_PORT=8025
diff --git a/docs/adr/0001-consume-baobab-dev-image-by-pinned-reference.md b/docs/adr/0001-consume-baobab-dev-image-by-pinned-reference.md
new file mode 100644
index 0000000..7892110
--- /dev/null
+++ b/docs/adr/0001-consume-baobab-dev-image-by-pinned-reference.md
@@ -0,0 +1,81 @@
+# 0001. Consume `baobab-dev` by Pinned Image Reference
+
+## Status
+
+Accepted
+
+## Date
+
+2026-08-23
+
+## Context
+
+The BAOBAB monorepo (`nabhold/baobab`) requires a Dev Container / GitHub Codespaces environment providing a consistent, reproducible toolchain (Python, Node.js, Flutter/Dart, PostgreSQL/Redis client tools, Docker CLI, GitHub CLI) for every contributor.
+
+A companion repository, [`nabhold/baobab-dev`](https://github.com/nabhold/baobab-dev), publishes a deterministic, checksum-verified, multi-stage container image to `ghcr.io/nabhold/baobab-dev`, with its own version-resolution pipeline (`config/versions.yaml` → `config/resolve.sh` → `config/versions.lock` → Dockerfile), CI-driven publish workflow, and semantic version tags (current release: `1.0.0`).
+
+An initial `.devcontainer/` configuration existed in `nabhold/baobab` prior to this ADR. On review it was found to:
+
+- Build a local Dockerfile `FROM ghcr.io/nabhold/baobab-dev:latest`, re-implementing (with several defects) provisioning logic — installing Python, Node.js, Java, Flutter, PostgreSQL/Redis client tools, Docker CLI — that the `baobab-dev` image already provides, deterministically, at image-build time.
+- Reference the base image by `:latest` rather than a pinned release, undermining reproducibility.
+- Contain a hard syntax error (`summary.sh`) and references to two non-existent install scripts (`install/android_sdk.sh`, `install/dependencies.sh`) that would abort container provisioning on every fresh Codespace.
+- Invert the Dev Container lifecycle hook order (a lightweight `updateContentCommand` ran before the heavy `postCreateCommand` it depended on).
+- Assume a `services/frontend` / `services/mobile` layout that does not match the repository's actual `apps/*` structure.
+
+## Problem Statement
+
+How should `nabhold/baobab`'s Dev Container / Codespaces configuration obtain its development toolchain, in a way that is reproducible, low-maintenance, and correctly reflects this repository's actual dependency-management tooling (`uv`, not Poetry)?
+
+## Decision Drivers
+
+- Reproducibility: every contributor and every Codespace must resolve to byte-identical tooling.
+- Minimal repository-local provisioning logic — avoid duplicating what `baobab-dev` already owns and tests.
+- Correctness: the solution must not silently invoke the wrong Python package manager (this repository uses `uv`, not Poetry).
+- Fast Codespaces startup, including effective use of Codespaces prebuilds.
+- Governance: `baobab-dev`'s publish pipeline is owned and versioned independently, in its own repository.
+
+## Considered Options
+
+1. **Build a local Dockerfile `FROM ghcr.io/nabhold/baobab-dev:<tag>`**, layering repository-specific OS packages or configuration on top.
+2. **Reference `ghcr.io/nabhold/baobab-dev` directly via `devcontainer.json`'s `image` field**, with no local build layer, and handle repository-specific setup (dependency install, `.env` bootstrap) via lifecycle hook scripts.
+
+## Decision
+
+**Option 2 — reference the image directly, pinned to an explicit release tag (`ghcr.io/nabhold/baobab-dev:1.0.0`), with no local Dockerfile build layer.**
+
+At the time of this decision, no repository-specific OS package or build customization is required beyond what the published image already provides. A local build layer would only add container build/pull latency without adding capability.
+
+Repository-specific setup is handled by a small, repository-owned script (`.devcontainer/post-create.sh`), invoked from `updateContentCommand` and `postCreateCommand`, rather than delegating to the image's own generic `baobab-post-create` command. This is necessary because `baobab-post-create`:
+
+- Installs Python dependencies via Poetry whenever it finds a root `pyproject.toml`, which is incorrect for this repository (`uv` is the standardized tool — see `uv.lock`, `Makefile`, and the "UV Workspace Readiness Roadmap" in `pyproject.toml`).
+- Unconditionally requires a workspace-local `config/versions.yaml` and executable `config/resolve.sh`, a convention this repository has not adopted (and does not need to, since version resolution is already owned by `baobab-dev`).
+
+The image's `baobab-verify` and `baobab-summary` commands have no such workspace-local dependency (they fall back to the image's own baked-in configuration) and are used directly.
+
+## Consequences
+
+**Positive**
+
+- Codespace/container creation is fast and deterministic — the toolchain is already built into the pulled image; nothing is installed at container-creation time.
+- Eliminates an entire class of repository-local provisioning bugs (the prior configuration's syntax error and missing-script references are removed along with the tree that contained them).
+- `uv sync` is guaranteed to invoke the correct package manager for this repository's dependencies.
+- The image tag is the single, explicit, reviewable point of upgrade — bumping the toolchain is a one-line pull request, not a script change.
+
+**Negative / Trade-offs**
+
+- Any repository-specific OS package or system-level customization (should one become necessary) requires either a new `baobab-dev` release or reintroducing a local build layer — a deliberate trade-off, not an oversight, and one this ADR can be revisited for if it occurs.
+- The two repositories (`baobab`, `baobab-dev`) must be upgraded in a coordinated but decoupled fashion; a `baobab-dev` release does not automatically propagate to `baobab` — the image tag must be bumped explicitly here.
+
+## Alternatives Rejected
+
+- **Continue building a local Dockerfile on top of `baobab-dev`** — rejected: no current requirement justifies the added build/pull latency, and it was the direct cause of the defects found in the prior configuration (drift between the local Dockerfile/scripts and what the base image already provided).
+- **Delegate dependency installation to the image's own `baobab-post-create`** — rejected: assumes Poetry, and requires a workspace-local `config/` directory this repository does not maintain.
+- **Pin to `:latest`** — rejected: undermines reproducibility, the primary decision driver.
+
+## References
+
+- [`nabhold/baobab-dev`](https://github.com/nabhold/baobab-dev) — image source repository
+- `ghcr.io/nabhold/baobab-dev` — published image
+- `.devcontainer/devcontainer.json`, `.devcontainer/post-create.sh`, `.devcontainer/README.md` (this repository)
+- `pyproject.toml` — "UV Workspace Readiness Roadmap"
+- `docs/governance/decision-record-process.md`
