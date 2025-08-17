// Copyright 2025 Hanzo Industries Inc.
// SPDX-License-Identifier: AGPL-3.0

package daemon

import (
	"embed"
)

//go:embed static/*
var static embed.FS
