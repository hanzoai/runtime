// Copyright 2025 Hanzo Industries Inc.
// SPDX-License-Identifier: AGPL-3.0

package models

import (
	"time"

	"github.com/hanzoai/runtime/apps/runner/pkg/models/enums"
)

type CacheData struct {
	SandboxState    enums.SandboxState
	BackupState     enums.BackupState
	DestructionTime *time.Time
}
