# ==============================================================================
# BAOBAB Enterprise Platform
#
# Dockerfile : Development Container
#
# Purpose
# -------
# Base image for GitHub Codespaces and VS Code Dev Containers.
#
# Responsibilities
# ----------------
# • Provide a stable Ubuntu development environment.
# • Install common operating system packages.
# • Configure the development container environment.
# • Prepare the workspace.
#
# This image intentionally DOES NOT install language runtimes or project
# dependencies. Those are provisioned by the BAOBAB Provisioning Framework
# under:
#
#     .devcontainer/docker/scripts/
#
# The provisioning sequence is:
#
#     run.sh
#         ├── bootstrap.sh
#         ├── environment.sh
#         ├── install.sh
#         ├── verify.sh
#         └── summary.sh
#
# This image is NOT intended for production deployments.
#
# Author  : BAOBAB Contributors
# License : Apache License 2.0
# ==============================================================================

FROM ghcr.io/nabhold/baobab-dev:1.0.0 AS development

# ------------------------------------------------------------------------------
# Build Metadata
# ------------------------------------------------------------------------------

ARG VERSION="0.1.0-dev"
ARG BUILD_DATE="unknown"
ARG VCS_REF="unknown"
ARG BAOBAB_PROVISIONER_VERSION="1.0.0"

LABEL org.opencontainers.image.title="BAOBAB Development Environment" \
      org.opencontainers.image.description="Development container for the BAOBAB Enterprise Platform." \
      org.opencontainers.image.vendor="NabHold Group Africa" \
      org.opencontainers.image.authors="BAOBAB Engineering Team" \
      org.opencontainers.image.licenses="Apache-2.0" \
      org.opencontainers.image.source="https://github.com/bjnabs/baobab" \
      org.opencontainers.image.version="${VERSION}" \
      org.opencontainers.image.created="${BUILD_DATE}" \
      org.opencontainers.image.revision="${VCS_REF}"


LABEL io.baobab.provisioner.version="${BAOBAB_PROVISIONER_VERSION}"



# ------------------------------------------------------------------------------
# Environment
# ------------------------------------------------------------------------------

ENV DEBIAN_FRONTEND=noninteractive
ENV BAOBAB_PROVISIONER_VERSION="${BAOBAB_PROVISIONER_VERSION}"

ENV LANG=en_US.UTF-8 \
    LANGUAGE=en_US:en \
    LC_ALL=en_US.UTF-8 \
    TZ=UTC

# ------------------------------------------------------------------------------
# Default Shell
# ------------------------------------------------------------------------------

SHELL ["/bin/bash", "-o", "pipefail", "-c"]

# ------------------------------------------------------------------------------
# Install Base Operating System Packages
# ------------------------------------------------------------------------------

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        bash \
        build-essential \
        ca-certificates \
        curl \
        file \
        git \
        gnupg \
        jq \
        less \
        locales \
        lsb-release \
        nano \
        openssh-client \
        pkg-config \
        procps \
        rsync \
        software-properties-common \
        sudo \
        tar \
        tree \
        unzip \
        vim \
        wget \
        xz-utils \
        zip && \
    locale-gen en_US.UTF-8 && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# ------------------------------------------------------------------------------
# Create Common Development Directories
# ------------------------------------------------------------------------------

RUN mkdir -p \
    /opt/baobab \
    /opt/baobab/cache \
    /opt/baobab/tools \
    /var/log/baobab

# ------------------------------------------------------------------------------
# Workspace
# ------------------------------------------------------------------------------

WORKDIR /workspaces/baobab

# ------------------------------------------------------------------------------
# Container Environment
# ------------------------------------------------------------------------------

ENV BAOBAB_HOME=/workspaces/baobab \
    BAOBAB_ENV=development

# ------------------------------------------------------------------------------
# Default Command
# ------------------------------------------------------------------------------

CMD ["sleep", "infinity"]
