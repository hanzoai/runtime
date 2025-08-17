// Copyright 2025 Hanzo Industries Inc.
// SPDX-License-Identifier: AGPL-3.0

package volume

import (
	"context"
	"fmt"

	"github.com/hanzoai/runtime/cli/apiclient"
	"github.com/hanzoai/runtime/cli/cmd/common"
	view_common "github.com/hanzoai/runtime/cli/views/common"
	"github.com/spf13/cobra"
)

var DeleteCmd = &cobra.Command{
	Use:     "delete [VOLUME_ID]",
	Short:   "Delete a volume",
	Args:    cobra.ExactArgs(1),
	Aliases: common.GetAliases("delete"),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		apiClient, err := apiclient.GetApiClient(nil, nil)
		if err != nil {
			return err
		}

		res, err := apiClient.VolumesAPI.DeleteVolume(ctx, args[0]).Execute()
		if err != nil {
			return apiclient.HandleErrorResponse(res, err)
		}

		view_common.RenderInfoMessageBold(fmt.Sprintf("Volume %s deleted", args[0]))
		return nil
	},
}
