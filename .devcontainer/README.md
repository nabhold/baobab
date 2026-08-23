# BAOBAB Development Container

## Purpose

The BAOBAB Development Container provides a consistent, reproducible, cloud-ready development environment for every contributor. It is used by GitHub Codespaces and Visual Studio Code Dev Containers to ensure every developer works in an identical environment with minimal setup.

---

# Architecture

Unlike a repository that builds its own development image, BAOBAB consumes a **published, versioned image** maintained in a dedicated repository:

| Component | Location |
| --- | --- |
| Image source | <https://github.com/nabhold/baobab-dev> |
| Published image | `ghcr.io/nabhold/baobab-dev` |
| Current pin | `1.0.0` (see `devcontainer.json`) |

`baobab-dev` is a deterministic, multi-stage, checksum-verified image build with its own CI pipeline, version-resolution lockfile, and documentation. It already bakes in the full toolchain — Python 3.14, `uv`, PostgreSQL/Redis client tools, Docker CLI, GitHub CLI, Node.js LTS, and Flutter/Dart — at **image build time**, deterministically, once. This repository's `.devcontainer/` therefore does **not** build a local Dockerfile or reinstall language runtimes; it only references the published image and layers on the small amount of setup that is specific to *this* repository.

Reference documentation for the image itself (toolchain contents, helper commands, verification, version policy) lives in the `baobab-dev` repository's own `docs/` tree — see in particular `docs/reference/helper-commands.md` and `docs/introduction/included.md`.

---

# Directory Structure

```text
.devcontainer/
├── devcontainer.json
├── post-create.sh
└── README.md
```

---

# Component Responsibilities

## `devcontainer.json`

Defines how GitHub Codespaces and Visual Studio Code provision the BAOBAB development environment: which published image to use, Dev Container Features, lifecycle hooks, VS Code settings/extensions, and forwarded ports.

## `post-create.sh`

A small, repository-owned script that runs *this repository's* setup on top of the image: creating `.env` from `.env.example`, and installing the root `uv` workspace's Python dependencies.

This is deliberately **not** delegated to the image's own generic `baobab-post-create` dependency-install logic: that logic installs Python dependencies with **Poetry** whenever it finds a root `pyproject.toml`. BAOBAB standardises on **`uv`** (see `uv.lock`, `Makefile`, and the "UV WORKSPACE READINESS ROADMAP" in `pyproject.toml`) — Poetry is not used in this repository. Everything else the image's `baobab-post-create` provides (git identity priming, toolchain verification) is generic and safe, and is still used directly via `onCreateCommand`.

---

# Lifecycle

| Hook | Runs | Does |
| --- | --- | --- |
| `onCreateCommand` | Once, at creation | `baobab-post-create --stage=on-create` (image-provided: git identity priming, `baobab-verify`) |
| `updateContentCommand` | Every content update — the stage GitHub Codespaces **prebuilds** snapshot | `post-create.sh` (`.env` bootstrap, `uv sync`) |
| `postCreateCommand` | Once, after workspace mount | `post-create.sh` again — a safety net in case a prebuild was stale |
| `postStartCommand` | Every container start | `baobab-summary` (image-provided environment banner) |

All steps are idempotent and safe to re-run.

---

# Helper Commands

The image installs the following commands directly onto `PATH`; they are available in every container without any additional setup:

| Command | Purpose |
| --- | --- |
| `baobab-verify` | Validate the toolchain is installed and correctly configured |
| `baobab-summary` | Print installed tool versions and environment metadata |
| `baobab-post-create` | Image-provided lifecycle hook (`--stage=on-create` / `post-create` / `update-content`) |
| `baobab-bootstrap` | Full provisioning entry point for non-Codespaces use (bare `docker run`, self-hosted runners) |

Run `baobab-verify` any time you want to confirm the environment is healthy, and `baobab-summary` when reporting an issue.

---

# Rebuilding / Updating the Development Container

Because the container now references a published image rather than building one locally, "rebuilding" means either:

1. **Picking up the current pin** — Command Palette → **Dev Containers: Rebuild Container** (or recreate the Codespace) re-pulls `ghcr.io/nabhold/baobab-dev:1.0.0` if it has changed.
2. **Moving to a new image release** — bump the tag in `devcontainer.json`'s `image` field via a reviewed pull request once a new `baobab-dev` release is published, then rebuild as above.

The image tag is deliberately pinned (never `:latest`) so every contributor and every Codespace resolve to byte-identical tooling.

---

# Engineering Principles

* Reproducible development environments
* Single, centrally governed source of truth for the toolchain
* Minimal repository-local provisioning logic
* Separation of concerns between the platform image and this repository's own setup
* Security by default (the image's own build pipeline owns checksum verification and CVE tracking for the toolchain)

---

# Related Documentation

* `docs/03-development/`
* `docs/04-devops/`
* `docs/governance/engineering-handbook.md`
* `docs/governance/repository-governance.md`
* `docs/governance/decision-record-process.md`
* [`nabhold/baobab-dev`](https://github.com/nabhold/baobab-dev) — image source, version policy, and full toolchain reference

---

# Ownership

This configuration is maintained by the BAOBAB Engineering Team. Changes to the development environment should be reviewed through the project's pull request process and, where appropriate, documented through an Architecture Decision Record (ADR).
