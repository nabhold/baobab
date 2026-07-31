#
# Baobab Enterprise Platform
#--------------------------------
#
# File        : Makefile
# Repository  : https://github.com/nabhold/baobab
# Organization: NABHOLD Group Africa
#
#------------------------------------------------------------------------------
# PURPOSE
#------------------------------------------------------------------------------
#
# The Makefile provides a single, consistent command-line interface for
# developers, CI/CD pipelines, and automation.
#
# It orchestrates:
#
#   • uv
#   • Docker
#   • MkDocs
#   • Testing
#   • Linting
#   • Formatting
#   • Security Scanning
#   • Documentation
#
# Usage
# -----
#
#   make help
#
# Example
# -------
#
#   make install
#   make verify
#   make docs
#   make serve
#   make test
#   make lint
#   make clean
#
#
.DEFAULT_GOAL := help

SHELL := /bin/bash

PROJECT := baobab

PYTHON := python3.14

UV := uv

MKDOCS := $(UV) run mkdocs

PYTEST := $(UV) run pytest

RUFF := $(UV) run ruff

BLACK := $(UV) run black

PIP_AUDIT := $(UV) run pip-audit

BANDIT := $(UV) run bandit

# Help
#------------
.PHONY: help

help:
	@echo ""
	@echo "==============================================================="
	@echo "                 Baobab Enterprise Platform"
	@echo "==============================================================="
	@echo ""
	@echo "Environment"
	@echo "  install        Install all dependency groups"
	@echo "  sync           Synchronize dependencies"
	@echo "  upgrade        Upgrade dependency lock file"
	@echo ""
	@echo "Development"
	@echo "  format         Format source code"
	@echo "  lint           Lint source code"
	@echo "  typecheck      Run static type checking"
	@echo "  test           Run unit tests"
	@echo "  security       Run security scans"
	@echo ""
	@echo "Documentation"
	@echo "  docs           Build documentation"
	@echo "  serve          Start documentation server"
	@echo ""
	@echo "Containers"
	@echo "  up             Start Docker services"
	@echo "  down           Stop Docker services"
	@echo "  logs           View Docker logs"
	@echo ""
	@echo "Repository"
	@echo "  verify         Run all verification checks"
	@echo "  clean          Remove generated artifacts"
	@echo ""

# Environment
#------------
.PHONY: install

install:
	$(UV) sync --all-groups

.PHONY: sync

sync:
	$(UV) sync

.PHONY: upgrade

upgrade:
	$(UV) lock --upgrade

# Formatting
#---------------
.PHONY: format

format:
	$(RUFF) format .

# Linting
#------------
.PHONY: lint

lint:
	$(RUFF) check .

# Static Analysis
#------------
.PHONY: typecheck

typecheck:
	$(UV) run pyright

# Testing
#----------------

.PHONY: test

test:
	$(PYTEST)

#  Security
#------------------
.PHONY: security

security:
	$(BANDIT) -r services
	$(PIP_AUDIT)

# Documentation
#------------------------

.PHONY: docs

docs:
	$(MKDOCS) build --strict

.PHONY: serve

serve:
	$(MKDOCS) serve

#  Docker
#----------------

.PHONY: up

up:
	docker compose up -d

.PHONY: down

down:
	docker compose down

.PHONY: logs

logs:
	docker compose logs -f

#  Verification
#-----------------
.PHONY: verify

verify: lint typecheck test docs security

#
# Cleanup
#------------------

.PHONY: clean

clean:
	find . -type d -name "__pycache__" -exec rm -rf {} +
	find . -type d -name ".pytest_cache" -exec rm -rf {} +
	find . -type d -name ".ruff_cache" -exec rm -rf {} +
	find . -type d -name ".mypy_cache" -exec rm -rf {} +
	find . -type d -name "site" -exec rm -rf {} +
	find . -type d -name "dist" -exec rm -rf {} +
	find . -type d -name "build" -exec rm -rf {} +
	find . -type f -name "*.pyc" -delete

#-------------------
# End of File

