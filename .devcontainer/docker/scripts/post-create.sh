#!/usr/bin/env bash

# ==============================================================================
# BAOBAB Enterprise Platform
#
# Script  : .devcontainer/post-create.sh
# Purpose : Repository-specific setup that runs on top of the published
#           ghcr.io/nabhold/baobab-dev image.
#
# Why this exists instead of using the image's own baobab-post-create hook
# for dependency installation:
#
#   The baobab-dev image's generic post-create logic installs Python
#   dependencies with Poetry whenever it finds a root pyproject.toml. This
#   repository standardises on `uv` (see uv.lock, Makefile, pyproject.toml
#   "UV WORKSPACE READINESS ROADMAP") — Poetry is not used here at all.
#   Delegating dependency installation to the image's generic hook would
#   silently invoke the wrong package manager, so this repository owns that
#   one step directly instead.
#
# Invoked by devcontainer.json as both `updateContentCommand` (the stage
# GitHub Codespaces prebuilds snapshot) and `postCreateCommand` (a safety
# net re-run on actual container creation). Safe to run repeatedly — every
# step here is idempotent.
#
# Note on the image's own baobab-post-create / baobab-bootstrap commands:
# both unconditionally require a workspace-local config/versions.yaml and
# an executable config/resolve.sh (they call `die` if either is missing),
# a convention this repository has not adopted. They are therefore not
# used here at all — this script deliberately re-implements only the small,
# safe subset of that behaviour (git identity priming, toolchain
# verification) that has no such precondition. baobab-verify and
# baobab-summary have no config/ dependency (they fall back to the image's
# own baked-in configuration) and are used directly below.
#
# Author  : BAOBAB Contributors
# License : Apache License 2.0
# ==============================================================================

set -Eeuo pipefail

WORKSPACE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${WORKSPACE_DIR}"

log() {
    printf '\n\033[1;36m[baobab:post-create]\033[0m %s\n' "$1"
}

warn() {
    printf '\n\033[1;33m[baobab:post-create] WARNING:\033[0m %s\n' "$1"
}

# ------------------------------------------------------------------------------
# Git identity
# ------------------------------------------------------------------------------
# Placeholder identity only if none is already configured (e.g. mounted from
# the host, or set automatically by GitHub Codespaces). safe.directory avoids
# "detected dubious ownership" errors when the container's UID does not
# match the mounted workspace's owning UID.

log "Priming Git configuration"

git config --global --get user.name  >/dev/null 2>&1 || git config --global user.name  "BAOBAB Developer"
git config --global --get user.email >/dev/null 2>&1 || git config --global user.email "dev@example.com"

git config --global --get-all safe.directory 2>/dev/null | grep -Fxq "${WORKSPACE_DIR}" \
    || git config --global --add safe.directory "${WORKSPACE_DIR}"

# ------------------------------------------------------------------------------
# Environment file
# ------------------------------------------------------------------------------

if [[ -f ".env.example" && ! -f ".env" ]]; then
    log "Creating .env from .env.example"
    cp .env.example .env
else
    log ".env already present or no .env.example found."
fi

# ------------------------------------------------------------------------------
# Python workspace dependencies (uv)
# ------------------------------------------------------------------------------
# Installs only what the root uv workspace currently manages (repository
# tooling: lint, typecheck, test, docs — see [dependency-groups] in
# pyproject.toml). services/backend, services/ai, services/worker, and the
# other prospective workspace members are intentionally NOT yet uv workspace
# members (see the "UV WORKSPACE READINESS ROADMAP" in pyproject.toml), so
# this deliberately does not attempt to sync them until that gate is lifted.

if [[ -f "pyproject.toml" ]]; then

    if command -v uv >/dev/null 2>&1; then
        log "Installing root workspace dependencies with uv"
        uv sync --all-groups
    else
        warn "uv not found on PATH. Skipping Python dependency installation."
    fi

else
    log "No root pyproject.toml found. Skipping Python dependency installation."
fi

# ------------------------------------------------------------------------------
# Toolchain verification
# ------------------------------------------------------------------------------
# Non-fatal: a verification issue should surface loudly, not block the
# developer from getting a working shell.

if command -v baobab-verify >/dev/null 2>&1; then
    baobab-verify --quiet || warn "Toolchain verification reported issues. Run 'baobab-verify' for details."
fi

# ------------------------------------------------------------------------------
# Complete
# ------------------------------------------------------------------------------

log "Repository setup complete."

exit 0
