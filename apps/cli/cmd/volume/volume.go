// Copyright 2025 Hanzo Industries Inc.
// SPDX-License-Identifier: AGPL-3.0

package volume

import (
	"github.com/hanzoai/runtime/cli/internal"
	"github.com/spf13/cobra"
)

var VolumeCmd = &cobra.Command{
	Use:     "volume",
	Short:   "Manage Runtime volumes",
	Long:    "Commands for managing Runtime volumes",
	Aliases: []string{"volumes"},
	GroupID: internal.SANDBOX_GROUP,
}

func init() {
	VolumeCmd.AddCommand(ListCmd)
	VolumeCmd.AddCommand(CreateCmd)
	VolumeCmd.AddCommand(GetCmd)
	VolumeCmd.AddCommand(DeleteCmd)
}
