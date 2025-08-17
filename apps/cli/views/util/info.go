// Copyright 2025 Hanzo Industries Inc.
// SPDX-License-Identifier: AGPL-3.0

package util

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/hanzoai/runtime/cli/views/common"
)

const PropertyNameWidth = 16

var PropertyNameStyle = lipgloss.NewStyle().
	Foreground(common.LightGray)

var PropertyValueStyle = lipgloss.NewStyle().
	Foreground(common.Light).
	Bold(true)
