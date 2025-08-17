// Copyright 2025 Hanzo Industries Inc.
// SPDX-License-Identifier: AGPL-3.0

package organization

import (
	"errors"

	"github.com/hanzoai/runtime/cli/config"
	"github.com/hanzoai/runtime/cli/internal"
	"github.com/spf13/cobra"
)

var OrganizationCmd = &cobra.Command{
	Use:     "organization",
	Short:   "Manage Runtime organizations",
	Long:    "Commands for managing Runtime organizations",
	Aliases: []string{"organizations", "org", "orgs"},
	GroupID: internal.USER_GROUP,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if config.IsApiKeyAuth() {
			return errors.New("organization commands are not available when using API key authentication - run `runtime login` to reauthenticate with browser")
		}

		return nil
	},
}

func init() {
	OrganizationCmd.AddCommand(ListCmd)
	OrganizationCmd.AddCommand(CreateCmd)
	OrganizationCmd.AddCommand(UseCmd)
	OrganizationCmd.AddCommand(DeleteCmd)
}
