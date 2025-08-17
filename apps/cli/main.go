// Copyright 2025 Hanzo Industries Inc.
// SPDX-License-Identifier: AGPL-3.0

package main

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/hanzoai/runtime/cli/cmd"
	"github.com/hanzoai/runtime/cli/cmd/auth"
	"github.com/hanzoai/runtime/cli/cmd/mcp"
	"github.com/hanzoai/runtime/cli/cmd/organization"
	"github.com/hanzoai/runtime/cli/cmd/sandbox"
	"github.com/hanzoai/runtime/cli/cmd/snapshot"
	"github.com/hanzoai/runtime/cli/cmd/volume"
	"github.com/hanzoai/runtime/cli/internal"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:               "runtime",
	Short:             "Runtime CLI",
	Long:              "Command line interface for Runtime Sandboxes",
	DisableAutoGenTag: true,
	SilenceUsage:      true,
	SilenceErrors:     true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	rootCmd.AddGroup(&cobra.Group{ID: internal.USER_GROUP, Title: "User"})
	rootCmd.AddGroup(&cobra.Group{ID: internal.SANDBOX_GROUP, Title: "Sandbox"})

	rootCmd.AddCommand(auth.LoginCmd)
	rootCmd.AddCommand(auth.LogoutCmd)
	rootCmd.AddCommand(sandbox.SandboxCmd)
	rootCmd.AddCommand(snapshot.SnapshotsCmd)
	rootCmd.AddCommand(volume.VolumeCmd)
	rootCmd.AddCommand(organization.OrganizationCmd)
	rootCmd.AddCommand(mcp.MCPCmd)
	rootCmd.AddCommand(cmd.DocsCmd)
	rootCmd.AddCommand(cmd.AutoCompleteCmd)
	rootCmd.AddCommand(cmd.GenerateDocsCmd)
	rootCmd.AddCommand(cmd.VersionCmd)

	rootCmd.CompletionOptions.HiddenDefaultCmd = true
	rootCmd.PersistentFlags().BoolP("help", "", false, "help for runtime")
	rootCmd.Flags().BoolP("version", "v", false, "Display the version of Runtime")

	rootCmd.PreRun = func(command *cobra.Command, args []string) {
		versionFlag, _ := command.Flags().GetBool("version")
		if versionFlag {
			err := cmd.VersionCmd.RunE(command, []string{})
			if err != nil {
				log.Fatal(err)
			}
			os.Exit(0)
		}
	}
}

func main() {
	_ = godotenv.Load()

	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
