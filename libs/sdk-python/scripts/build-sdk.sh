#!/usr/bin/env bash
# Copyright 2025 Hanzo Industries Inc.
# SPDX-License-Identifier: Apache-2.0

set -e

echo "→ build-sdk"

if [ -n "$PYPI_PKG_VERSION" ] || [ -n "$DEFAULT_PACKAGE_VERSION" ]; then
  VER="${PYPI_PKG_VERSION:-$DEFAULT_PACKAGE_VERSION}"
  poetry version "$VER"
else
  echo "Using version from pyproject.toml"
fi

poetry build

mv src/hanzo_runtime src/hanzo_runtime_sdk
sed -i 's/^name = "[^"]*"/name = "hanzo_runtime_sdk"/' pyproject.toml
poetry build
mv src/hanzo_runtime_sdk src/hanzo_runtime
sed -i 's/^name = "[^"]*"/name = "hanzo_runtime"/' pyproject.toml
