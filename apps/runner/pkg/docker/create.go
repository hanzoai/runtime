// Copyright 2025 Hanzo Industries Inc.
// SPDX-License-Identifier: AGPL-3.0

package docker

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hanzoai/runtime/libs/common-go/pkg/timer"
	"github.com/hanzoai/runtime/apps/runner/internal/constants"
	"github.com/hanzoai/runtime/apps/runner/pkg/api/dto"
	"github.com/hanzoai/runtime/apps/runner/pkg/common"
	"github.com/hanzoai/runtime/apps/runner/pkg/models/enums"
	"github.com/docker/docker/errdefs"

	log "github.com/sirupsen/logrus"
)

func (d *DockerClient) Create(ctx context.Context, sandboxDto dto.CreateSandboxDTO) (string, error) {
	defer timer.Timer()()

	// Settle the boundary before anything expensive happens, so a sandbox this
	// runner cannot isolate is refused rather than refused after an image pull.
	isolation, runtime, err := d.resolveIsolation(sandboxDto.Isolation, sandboxDto.OrgId)
	if err != nil {
		return "", err
	}

	startTime := time.Now()
	defer func() {
		common.ContainerOperationDuration.WithLabelValues("create", string(isolation)).Observe(time.Since(startTime).Seconds())
	}()

	state, err := d.DeduceSandboxState(ctx, sandboxDto.Id)
	if err != nil && state == enums.SandboxStateError {
		return "", err
	}

	if state == enums.SandboxStateStarted || state == enums.SandboxStatePullingSnapshot || state == enums.SandboxStateStarting {
		return sandboxDto.Id, nil
	}

	if state == enums.SandboxStateStopped || state == enums.SandboxStateCreating {
		err = d.Start(ctx, sandboxDto.Id)
		if err != nil {
			return "", err
		}

		return sandboxDto.Id, nil
	}

	d.cache.SetSandboxState(ctx, sandboxDto.Id, enums.SandboxStateCreating)

	ctx = context.WithValue(ctx, constants.ID_KEY, sandboxDto.Id)
	err = d.PullImage(ctx, sandboxDto.Snapshot, sandboxDto.Registry)
	if err != nil {
		return "", err
	}

	d.cache.SetSandboxState(ctx, sandboxDto.Id, enums.SandboxStateCreating)

	err = d.validateImageArchitecture(ctx, sandboxDto.Snapshot)
	if err != nil {
		log.Errorf("ERROR: %s.\n", err.Error())
		return "", err
	}

	volumeMountPathBinds := make([]string, 0)
	if sandboxDto.Volumes != nil {
		volumeMountPathBinds, err = d.getVolumesMountPathBinds(ctx, sandboxDto.Volumes)
		if err != nil {
			return "", err
		}
	}

	containerConfig, hostConfig, networkingConfig, err := d.getContainerConfigs(ctx, sandboxDto, volumeMountPathBinds, isolation, runtime)
	if err != nil {
		return "", err
	}

	c, err := d.apiClient.ContainerCreate(ctx, containerConfig, hostConfig, networkingConfig, nil, sandboxDto.Id)
	if err != nil {
		return "", err
	}

	created := time.Now()

	err = d.Start(ctx, sandboxDto.Id)
	if err != nil {
		return "", err
	}

	// One record per sandbox carrying what it ran behind and what it was, so a
	// comparison across boundaries and across workload shapes is a query.
	log.WithFields(log.Fields{
		"sandbox":   sandboxDto.Id,
		"org":       sandboxDto.OrgId,
		"isolation": isolation,
		"runtime":   runtime,
		"snapshot":  sandboxDto.Snapshot,
		"cpu":       sandboxDto.CpuQuota,
		"memory":    sandboxDto.MemoryQuota,
		"create_ms": created.Sub(startTime).Milliseconds(),
		"start_ms":  time.Since(created).Milliseconds(),
		"total_ms":  time.Since(startTime).Milliseconds(),
	}).Info("sandbox created")

	return c.ID, nil
}

func (p *DockerClient) validateImageArchitecture(ctx context.Context, image string) error {
	defer timer.Timer()()

	inspect, _, err := p.apiClient.ImageInspectWithRaw(ctx, image)
	if err != nil {
		if errdefs.IsNotFound(err) {
			return err
		}
		return fmt.Errorf("failed to inspect image: %w", err)
	}

	arch := strings.ToLower(inspect.Architecture)
	validArchs := []string{"amd64", "x86_64"}

	for _, validArch := range validArchs {
		if arch == validArch {
			return nil
		}
	}

	return common.NewConflictError(fmt.Errorf("image %s architecture (%s) is not x64 compatible", image, inspect.Architecture))
}
