// Copyright 2025 Daytona Platforms Inc.
// SPDX-License-Identifier: AGPL-3.0

package sandbox

import (
	"github.com/hanzoai/runtime/cli/internal"
	"github.com/spf13/cobra"
)

var SandboxCmd = &cobra.Command{
	Use:     "sandbox",
	Short:   "Manage Runtime sandboxes",
	Long:    "Commands for managing Runtime sandboxes",
	Aliases: []string{"sandboxes"},
	GroupID: internal.SANDBOX_GROUP,
}

func init() {
	SandboxCmd.AddCommand(ListCmd)
	SandboxCmd.AddCommand(CreateCmd)
	SandboxCmd.AddCommand(InfoCmd)
	SandboxCmd.AddCommand(DeleteCmd)
	SandboxCmd.AddCommand(StartCmd)
	SandboxCmd.AddCommand(StopCmd)
}
