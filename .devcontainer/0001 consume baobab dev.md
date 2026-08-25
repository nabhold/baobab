# 0001. Consume `baobab-dev` by Pinned Image Reference

## Status

Accepted

## Date

2026-08-23

## Context

The BAOBAB monorepo (`nabhold/baobab`) requires a Dev Container / GitHub Codespaces environment providing a consistent, reproducible toolchain (Python, Node.js, Flutter/Dart, PostgreSQL/Redis client tools, Docker CLI, GitHub CLI) for every contributor.

A companion repository, [`nabhold/baobab-dev`](https://github.com/nabhold/baobab-dev), publishes a deterministic, checksum-verified, multi-stage container image to `ghcr.io/nabhold/baobab-dev`, with its own version-resolution pipeline (`config/versions.yaml` → `config/resolve.sh` → `config/versions.lock` → Dockerfile), CI-driven publish workflow, and semantic version tags (current release: `1.0.0`).

An initial `.devcontainer/` configuration existed in `nabhold/baobab` prior to this ADR. On review it was found to:

- Build a local Dockerfile `FROM ghcr.io/nabhold/baobab-dev:latest`, re-implementing (with several defects) provisioning logic — installing Python, Node.js, Java, Flutter, PostgreSQL/Redis client tools, Docker CLI — that the `baobab-dev` image already provides, deterministically, at image-build time.
- Reference the base image by `:latest` rather than a pinned release, undermining reproducibility.
- Contain a hard syntax error (`summary.sh`) and references to two non-existent install scripts (`install/android_sdk.sh`, `install/dependencies.sh`) that would abort container provisioning on every fresh Codespace.
- Invert the Dev Container lifecycle hook order (a lightweight `updateContentCommand` ran before the heavy `postCreateCommand` it depended on).
- Assume a `services/frontend` / `services/mobile` layout that does not match the repository's actual `apps/*` structure.

## Problem Statement

How should `nabhold/baobab`'s Dev Container / Codespaces configuration obtain its development toolchain, in a way that is reproducible, low-maintenance, and correctly reflects this repository's actual dependency-management tooling (`uv`, not Poetry)?

## Decision Drivers

- Reproducibility: every contributor and every Codespace must resolve to byte-identical tooling.
- Minimal repository-local provisioning logic — avoid duplicating what `baobab-dev` already owns and tests.
- Correctness: the solution must not silently invoke the wrong Python package manager (this repository uses `uv`, not Poetry).
- Fast Codespaces startup, including effective use of Codespaces prebuilds.
- Governance: `baobab-dev`'s publish pipeline is owned and versioned independently, in its own repository.

## Considered Options

1. **Build a local Dockerfile `FROM ghcr.io/nabhold/baobab-dev:<tag>`**, layering repository-specific OS packages or configuration on top.
2. **Reference `ghcr.io/nabhold/baobab-dev` directly via `devcontainer.json`'s `image` field**, with no local build layer, and handle repository-specific setup (dependency install, `.env` bootstrap) via lifecycle hook scripts.

## Decision

**Option 2 — reference the image directly, pinned to an explicit release tag (`ghcr.io/nabhold/baobab-dev:1.0.0`), with no local Dockerfile build layer.**

At the time of this decision, no repository-specific OS package or build customization is required beyond what the published image already provides. A local build layer would only add container build/pull latency without adding capability.

Repository-specific setup is handled by a small, repository-owned script (`.devcontainer/post-create.sh`), invoked from `updateContentCommand` and `postCreateCommand`, rather than delegating to the image's own generic `baobab-post-create` command. This is necessary because `baobab-post-create`:

- Installs Python dependencies via Poetry whenever it finds a root `pyproject.toml`, which is incorrect for this repository (`uv` is the standardized tool — see `uv.lock`, `Makefile`, and the "UV Workspace Readiness Roadmap" in `pyproject.toml`).
- Unconditionally requires a workspace-local `config/versions.yaml` and executable `config/resolve.sh`, a convention this repository has not adopted (and does not need to, since version resolution is already owned by `baobab-dev`).

The image's `baobab-verify` and `baobab-summary` commands have no such workspace-local dependency (they fall back to the image's own baked-in configuration) and are used directly.

## Consequences

**Positive**

- Codespace/container creation is fast and deterministic — the toolchain is already built into the pulled image; nothing is installed at container-creation time.
- Eliminates an entire class of repository-local provisioning bugs (the prior configuration's syntax error and missing-script references are removed along with the tree that contained them).
- `uv sync` is guaranteed to invoke the correct package manager for this repository's dependencies.
- The image tag is the single, explicit, reviewable point of upgrade — bumping the toolchain is a one-line pull request, not a script change.

**Negative / Trade-offs**

- Any repository-specific OS package or system-level customization (should one become necessary) requires either a new `baobab-dev` release or reintroducing a local build layer — a deliberate trade-off, not an oversight, and one this ADR can be revisited for if it occurs.
- The two repositories (`baobab`, `baobab-dev`) must be upgraded in a coordinated but decoupled fashion; a `baobab-dev` release does not automatically propagate to `baobab` — the image tag must be bumped explicitly here.

## Alternatives Rejected

- **Continue building a local Dockerfile on top of `baobab-dev`** — rejected: no current requirement justifies the added build/pull latency, and it was the direct cause of the defects found in the prior configuration (drift between the local Dockerfile/scripts and what the base image already provided).
- **Delegate dependency installation to the image's own `baobab-post-create`** — rejected: assumes Poetry, and requires a workspace-local `config/` directory this repository does not maintain.
- **Pin to `:latest`** — rejected: undermines reproducibility, the primary decision driver.

## References

- [`nabhold/baobab-dev`](https://github.com/nabhold/baobab-dev) — image source repository
- `ghcr.io/nabhold/baobab-dev` — published image
- `.devcontainer/devcontainer.json`, `.devcontainer/post-create.sh`, `.devcontainer/README.md` (this repository)
- `pyproject.toml` — "UV Workspace Readiness Roadmap"
- `docs/governance/decision-record-process.md`
