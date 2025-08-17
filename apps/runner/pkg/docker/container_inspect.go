// Copyright 2025 Hanzo Industries Inc.
// SPDX-License-Identifier: AGPL-3.0

package docker

import (
	"context"

	"github.com/docker/docker/api/types"
)

func (d *DockerClient) ContainerInspect(ctx context.Context, containerId string) (types.ContainerJSON, error) {
	return d.apiClient.ContainerInspect(ctx, containerId)
}
