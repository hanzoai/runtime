// Copyright 2025 Hanzo Industries Inc.
// SPDX-License-Identifier: AGPL-3.0

package docker

import (
	"context"

	"github.com/hanzoai/common-go/pkg/timer"
	"github.com/docker/docker/api/types/container"

	log "github.com/sirupsen/logrus"
)

func (d *DockerClient) startRuntimeDaemon(ctx context.Context, containerId string) error {
	defer timer.Timer()()

	execOptions := container.ExecOptions{
		Cmd:          []string{"sh", "-c", "/usr/local/bin/runtime"},
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
	}

	execStartOptions := container.ExecStartOptions{
		Detach: false,
	}

	result, err := d.execSync(ctx, containerId, execOptions, execStartOptions)
	if err != nil {
		log.Errorf("Error starting Runtime daemon: %s", err.Error())
		return nil
	}

	if result.ExitCode != 0 && result.StdErr != "" {
		log.Errorf("Error starting Runtime daemon: %s", string(result.StdErr))
		return nil
	}

	return nil
}
