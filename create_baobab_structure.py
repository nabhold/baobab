#!/usr/bin/env python3
"""
create_baobab_structure.py

Creates the Baobab enterprise monorepo directory structure.

Python: 3.14+
"""

from pathlib import Path

ROOT = Path("baobab")

DIRECTORIES = [
    ".devcontainer/docker",

    ".github/ISSUE_TEMPLATE",
    ".github/workflows",

    "apps/admin-portal",
    "apps/customer-portal",
    "apps/mobile",
    "apps/analytics",

    "services/backend/platform/authentication",
    "services/backend/platform/tenancy",
    "services/backend/platform/permissions",
    "services/backend/platform/audit",
    "services/backend/platform/notifications",

    "services/backend/products/crm",
    "services/backend/products/procurement",
    "services/backend/products/inventory",
    "services/backend/products/finance",
    "services/backend/products/projects",
    "services/backend/products/hr",
    "services/backend/products/manufacturing",

    "services/backend/shared",

    "services/ai/agents",
    "services/ai/orchestrators",
    "services/ai/memory",
    "services/ai/rag",
    "services/ai/providers",
    "services/ai/prompts",
    "services/ai/evaluations",
    "services/ai/tools",

    "services/worker/queues",
    "services/worker/schedules",
    "services/worker/monitoring",
    "services/worker/tasks",

    "sdk/python",
    "sdk/typescript",
    "sdk/dart",
    "sdk/java",

    "packages/ui",
    "packages/utilities",
    "packages/python",
    "packages/typescript",

    "plugins/connectors",
    "plugins/marketplace",
    "plugins/reports",
    "plugins/customer_extensions",

    "shared/api",
    "shared/contracts",
    "shared/schemas",
    "shared/events",
    "shared/protobuf",
    "shared/configuration",
    "shared/constants",
    "shared/locales",
    "shared/types",
    "shared/utilities",

    "infrastructure/aws",
    "infrastructure/kubernetes",
    "infrastructure/docker",
    "infrastructure/monitoring",
    "infrastructure/terraform",
    "infrastructure/networking",

    "deployment/compose",
    "deployment/docker",
    "deployment/helm",
    "deployment/kubernetes",

    "tools/bootstrap",
    "tools/codegen",
    "tools/release",
    "tools/migration",
    "tools/documentation",
    "tools/utilities",

    "resources/branding",
    "resources/templates",
    "resources/sample-data",
    "resources/fixtures",
    "resources/certificates",

    "design/figma",
    "design/wireframes",
    "design/branding",
    "design/icons",

    "docs/architecture",
    "docs/governance",
    "docs/specifications",
    "docs/adr",
    "docs/runbooks",
    "docs/diagrams",

    "tests/unit",
    "tests/integration",
    "tests/api",
    "tests/performance",
    "tests/security",
    "tests/e2e",
]

FILES = [
    ".devcontainer/.env.example",
    ".devcontainer/devcontainer.json",
    ".devcontainer/README.md",

    ".github/CODEOWNERS",
    ".github/PULL_REQUEST_TEMPLATE.md",

    ".editorconfig",
    ".gitignore",
    "compose.yaml",
    "Dockerfile",
    "Makefile",
    "pyproject.toml",
    "README.md",
    "LICENSE",
]


def create_directories():
    """Create all directories."""
    for directory in DIRECTORIES:
        path = ROOT / directory
        path.mkdir(parents=True, exist_ok=True)
        print(f"[DIR ] {path}")


def create_files():
    """Create placeholder files if they do not already exist."""
    for file in FILES:
        path = ROOT / file
        path.parent.mkdir(parents=True, exist_ok=True)

        if not path.exists():
            path.touch()
            print(f"[FILE] {path}")
        else:
            print(f"[SKIP] {path} already exists")


def create_gitkeep_files():
    """
    Add .gitkeep files to every empty directory so Git can track them.
    Existing non-empty directories are left untouched.
    """
    for directory in ROOT.rglob("*"):
        if directory.is_dir():
            if not any(directory.iterdir()):
                gitkeep = directory / ".gitkeep"
                gitkeep.touch(exist_ok=True)
                print(f"[KEEP] {gitkeep}")


def main():
    print("\nCreating Baobab Enterprise Repository Structure...\n")

    ROOT.mkdir(exist_ok=True)

    create_directories()
    create_files()
    create_gitkeep_files()

    print("\nDone!")
    print(f"Repository created at: {ROOT.resolve()}")


if __name__ == "__main__":
    main()
