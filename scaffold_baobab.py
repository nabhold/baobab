#!/usr/bin/env python3
"""
Baobab Enterprise Platform Monorepo Scaffold & Reconciler

Compares the current workspace or target directory with the official
Baobab platform monorepo directory tree, creating missing directories and
placeholder files without touching or overwriting existing assets.

Usage:
    python scaffold_baobab.py [--target-dir PATH] [--dry-run] [--verbose]
"""

import argparse
from pathlib import Path
import sys
from typing import List, Tuple

# =============================================================================
# Target Monorepo Specification
# =============================================================================

TARGET_DIRECTORIES: List[str] = [
    ".devcontainer",
    ".devcontainer/docker",
    ".devcontainer/docker/scripts",
    ".devcontainer/docker/scripts/config",
    ".devcontainer/docker/scripts/configure",
    ".devcontainer/docker/scripts/install",
    ".devcontainer/docker/scripts/utils",
    ".devcontainer/docker/scripts/verify",
    ".devcontainer/docker/scripts/workspace",
    ".github",
    ".github/ISSUE_TEMPLATE",
    ".github/workflows",
    ".vscode",
    "apps",
    "apps/admin-portal",
    "apps/admin-portal/src",
    "apps/admin-portal/public",
    "apps/admin-portal/tests",
    "apps/admin-portal/docs",
    "apps/customer-portal",
    "apps/customer-portal/src",
    "apps/customer-portal/public",
    "apps/customer-portal/tests",
    "apps/customer-portal/docs",
    "apps/analytics",
    "apps/analytics/src",
    "apps/analytics/public",
    "apps/analytics/tests",
    "apps/analytics/docs",
    "apps/mobile",
    "apps/mobile/lib",
    "apps/mobile/test",
    "apps/mobile/integration_test",
    "apps/mobile/assets",
    "apps/mobile/docs",
    "services",
    "services/backend",
    "services/backend/src",
    "services/backend/src/platform",
    "services/backend/src/platform/authentication",
    "services/backend/src/platform/tenancy",
    "services/backend/src/platform/permissions",
    "services/backend/src/platform/audit",
    "services/backend/src/platform/notifications",
    "services/backend/src/products",
    "services/backend/src/products/crm",
    "services/backend/src/products/procurement",
    "services/backend/src/products/inventory",
    "services/backend/src/products/finance",
    "services/backend/src/products/projects",
    "services/backend/src/products/hr",
    "services/backend/src/products/manufacturing",
    "services/backend/src/shared",
    "services/backend/tests",
    "services/backend/migrations",
    "services/backend/scripts",
    "services/backend/docs",
    "services/ai",
    "services/ai/src",
    "services/ai/src/agents",
    "services/ai/src/orchestrators",
    "services/ai/src/memory",
    "services/ai/src/rag",
    "services/ai/src/providers",
    "services/ai/src/prompts",
    "services/ai/src/evaluations",
    "services/ai/src/tools",
    "services/ai/tests",
    "services/ai/notebooks",
    "services/ai/scripts",
    "services/ai/docs",
    "services/worker",
    "services/worker/src",
    "services/worker/src/queues",
    "services/worker/src/schedules",
    "services/worker/src/monitoring",
    "services/worker/src/tasks",
    "services/worker/tests",
    "services/worker/scripts",
    "services/worker/docs",
    "sdk",
    "sdk/python",
    "sdk/python/src",
    "sdk/python/tests",
    "sdk/python/examples",
    "sdk/typescript",
    "sdk/typescript/src",
    "sdk/typescript/tests",
    "sdk/dart",
    "sdk/dart/lib",
    "sdk/dart/test",
    "sdk/java",
    "sdk/java/src",
    "sdk/java/test",
    "packages",
    "packages/python",
    "packages/python/src",
    "packages/python/tests",
    "packages/typescript",
    "packages/typescript/src",
    "packages/typescript/tests",
    "packages/ui",
    "packages/utilities",
    "shared",
    "shared/python",
    "shared/python/src",
    "shared/python/tests",
    "shared/typescript",
    "shared/dart",
    "shared/java",
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
    "plugins",
    "plugins/connectors",
    "plugins/marketplace",
    "plugins/reports",
    "plugins/customer_extensions",
    "infrastructure",
    "infrastructure/aws",
    "infrastructure/kubernetes",
    "infrastructure/docker",
    "infrastructure/monitoring",
    "infrastructure/terraform",
    "infrastructure/networking",
    "deployment",
    "deployment/compose",
    "deployment/docker",
    "deployment/helm",
    "deployment/kubernetes",
    "tools",
    "tools/bootstrap",
    "tools/codegen",
    "tools/documentation",
    "tools/migration",
    "tools/release",
    "tools/utilities",
    "resources",
    "resources/branding",
    "resources/certificates",
    "resources/fixtures",
    "resources/sample-data",
    "resources/templates",
    "design",
    "design/branding",
    "design/figma",
    "design/icons",
    "design/wireframes",
    "docs",
    "docs/introduction",
    "docs/architecture",
    "docs/adr",
    "docs/api",
    "docs/ai",
    "docs/data",
    "docs/development",
    "docs/devops",
    "docs/diagrams",
    "docs/domain",
    "docs/governance",
    "docs/runbooks",
    "docs/security",
    "docs/specifications",
    "docs/testing",
    "tests",
    "tests/unit",
    "tests/integration",
    "tests/e2e",
    "tests/performance",
    "tests/security",
    "tests/fixtures",
]

TARGET_FILES: List[str] = [
    # Top-level config & metadata files
    ".editorconfig",
    ".env.example",
    ".gitattributes",
    ".gitignore",
    ".markdownlint.jsonc",
    ".pre-commit-config.yaml",
    ".yamllint.yml",
    "CHANGELOG.md",
    "CODE_OF_CONDUCT.md",
    "CONTRIBUTING.md",
    "compose.override.yaml",
    "compose.yaml",
    "Dockerfile",
    "LICENSE",
    "Makefile",
    "mkdocs.yml",
    "pyproject.toml",
    "pytest.ini",
    "README.md",
    "ruff.toml",
    "SECURITY.md",
    "taplo.toml",
    "uv.lock",
    # Sub-component configuration & doc files
    ".devcontainer/.env.example",
    ".devcontainer/README.md",
    ".devcontainer/devcontainer.json",
    ".devcontainer/docker/Dockerfile",
    ".devcontainer/docker/scripts/bootstrap.sh",
    ".devcontainer/docker/scripts/run.sh",
    ".devcontainer/docker/scripts/post-create.sh",
    ".devcontainer/docker/scripts/verify.sh",
    ".devcontainer/docker/scripts/summary.sh",
    ".github/CODEOWNERS",
    ".github/PULL_REQUEST_TEMPLATE.md",
    "apps/admin-portal/package.json",
    "apps/admin-portal/README.md",
    "apps/customer-portal/package.json",
    "apps/customer-portal/README.md",
    "apps/analytics/package.json",
    "apps/analytics/README.md",
    "apps/mobile/pubspec.yaml",
    "apps/mobile/README.md",
    "services/backend/pyproject.toml",
    "services/backend/README.md",
    "services/backend/manage.py",
    "services/ai/pyproject.toml",
    "services/ai/README.md",
    "services/worker/pyproject.toml",
    "services/worker/README.md",
    "sdk/python/pyproject.toml",
    "sdk/python/README.md",
    "sdk/typescript/package.json",
    "sdk/typescript/README.md",
    "sdk/dart/pubspec.yaml",
    "sdk/dart/README.md",
    "sdk/java/pom.xml",
    "sdk/java/README.md",
    "packages/python/pyproject.toml",
    "packages/python/README.md",
    "packages/typescript/package.json",
    "packages/typescript/README.md",
    "shared/python/pyproject.toml",
    "shared/python/README.md",
]


# =============================================================================
# Core Reconciliation Logic
# =============================================================================

def reconcile_tree(target_root: Path, dry_run: bool = False, verbose: bool = False) -> Tuple[int, int, int, int]:
    """
    Compares the target directory against the specification, creating missing items.

    Returns:
        Tuple[dirs_created, dirs_skipped, files_created, files_skipped]
    """
    dirs_created = 0
    dirs_skipped = 0
    files_created = 0
    files_skipped = 0

    print(f"\n🚀 Reconciling directory structure at: {target_root.resolve()}")
    if dry_run:
        print("⚠️  DRY RUN MODE ENABLED — No changes will be written to disk.\n")

    # Ensure target root directory exists
    if not target_root.exists():
        if not dry_run:
            target_root.mkdir(parents=True, exist_ok=True)
        print(f"[CREATED DIR] {target_root}")
        dirs_created += 1

    # 1. Reconcile Directories
    for dir_rel in TARGET_DIRECTORIES:
        dir_path = target_root / dir_rel
        if dir_path.exists():
            dirs_skipped += 1
            if verbose:
                print(f"[EXISTS DIR ] {dir_rel}")
        else:
            dirs_created += 1
            if not dry_run:
                dir_path.mkdir(parents=True, exist_ok=True)
            print(f"[CREATED DIR] {dir_rel}")

    # 2. Reconcile Files
    for file_rel in TARGET_FILES:
        file_path = target_root / file_rel
        if file_path.exists():
            files_skipped += 1
            if verbose:
                print(f"[EXISTS FILE] {file_rel}")
        else:
            files_created += 1
            if not dry_run:
                # Ensure parent directory exists before creating the file
                file_path.parent.mkdir(parents=True, exist_ok=True)
                file_path.touch(exist_ok=True)
            print(f"[CREATED FILE] {file_rel}")

    return dirs_created, dirs_skipped, files_created, files_skipped


# =============================================================================
# CLI Entrypoint
# =============================================================================

def main():
    parser = argparse.ArgumentParser(
        description="Scaffold and reconcile the Baobab Enterprise Platform directory structure."
    )
    parser.add_argument(
        "--target-dir",
        "-t",
        type=str,
        default=".",
        help="Root target directory path (defaults to current directory '.')",
    )
    parser.add_argument(
        "--dry-run",
        "-d",
        action="store_true",
        help="Preview actions without creating directories or files.",
    )
    parser.add_argument(
        "--verbose",
        "-v",
        action="store_true",
        help="Display existing items during comparison.",
    )

    args = parser.parse_args()
    target_root = Path(args.target_dir)

    dirs_created, dirs_skipped, files_created, files_skipped = reconcile_tree(
        target_root=target_root,
        dry_run=args.dry_run,
        verbose=args.verbose,
    )

    # Output Summary Table
    print("\n" + "=" * 60)
    print(" 📊 RECONCILIATION SUMMARY")
    print("=" * 60)
    print(f" Directories Created : {dirs_created}")
    print(f" Directories Skipped : {dirs_skipped}")
    print(f" Files Created       : {files_created}")
    print(f" Files Skipped       : {files_skipped}")
    print(f" Total Items Checked : {len(TARGET_DIRECTORIES) + len(TARGET_FILES)}")
    print("=" * 60)

    if args.dry_run:
        print("💡 To apply these changes, run the script without --dry-run / -d.")
    else:
        print("✅ Baobab Monorepo structure successfully aligned!")


if __name__ == "__main__":
    main()
