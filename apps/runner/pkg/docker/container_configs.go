// Copyright 2025 Hanzo Industries Inc.
// SPDX-License-Identifier: AGPL-3.0

package docker

import (
	"context"
	"errors"
	"fmt"

	"github.com/hanzoai/runtime/apps/runner/cmd/runner/config"
	"github.com/hanzoai/runtime/apps/runner/pkg/api/dto"
	"github.com/docker/docker/api/types/network"

	"github.com/docker/docker/api/types/container"
)

// sandboxCaps is what the sandbox keeps once every capability is dropped. Each
// one is here because something in the daemon or the image exercises it:
//
//	CHOWN, DAC_OVERRIDE, FOWNER, FSETID  the toolbox file API (chmod, chown) and
//	                                     package installs into system paths
//	SETUID, SETGID                       sudo and su, which the sandbox image
//	                                     grants the sandbox user
//	KILL                                 process and session management
//	AUDIT_WRITE                          sudo writes an audit record before it
//	                                     runs; without this it still works but
//	                                     prints a permission error every time
//
// Docker's default set also carries MKNOD, NET_RAW, NET_BIND_SERVICE,
// SYS_CHROOT, SETPCAP and SETFCAP. Nothing in the runner, the daemon or the
// computer-use plugin calls for them, so they stay dropped.
var sandboxCaps = []string{
	"CHOWN",
	"DAC_OVERRIDE",
	"FOWNER",
	"FSETID",
	"SETUID",
	"SETGID",
	"KILL",
	"AUDIT_WRITE",
}

func (d *DockerClient) getContainerConfigs(ctx context.Context, sandboxDto dto.CreateSandboxDTO, volumeMountPathBinds []string, isolation dto.Isolation, runtime string) (*container.Config, *container.HostConfig, *network.NetworkingConfig, error) {
	containerConfig := d.getContainerCreateConfig(sandboxDto, isolation)

	hostConfig, err := d.getContainerHostConfig(ctx, sandboxDto, volumeMountPathBinds, runtime)
	if err != nil {
		return nil, nil, nil, err
	}

	networkingConfig := d.getContainerNetworkingConfig(ctx)
	return containerConfig, hostConfig, networkingConfig, nil
}

func (d *DockerClient) getContainerCreateConfig(sandboxDto dto.CreateSandboxDTO, isolation dto.Isolation) *container.Config {
	envVars := []string{
		"RUNTIME_SANDBOX_ID=" + sandboxDto.Id,
		"RUNTIME_SANDBOX_SNAPSHOT=" + sandboxDto.Snapshot,
		"RUNTIME_SANDBOX_USER=" + sandboxDto.OsUser,
	}

	for key, value := range sandboxDto.Env {
		envVars = append(envVars, fmt.Sprintf("%s=%s", key, value))
	}

	return &container.Config{
		Hostname: sandboxDto.Id,
		Image:    sandboxDto.Snapshot,
		// User:         sandboxDto.OsUser,
		Env:          envVars,
		Entrypoint:   sandboxDto.Entrypoint,
		AttachStdout: true,
		AttachStderr: true,
		// Whose the sandbox is and what it runs behind, readable from the
		// container itself rather than only from the sandbox id.
		Labels: map[string]string{
			"hanzo.ai/org":       sandboxDto.OrgId,
			"hanzo.ai/isolation": string(isolation),
		},
	}
}

func (d *DockerClient) getContainerHostConfig(ctx context.Context, sandboxDto dto.CreateSandboxDTO, volumeMountPathBinds []string, runtime string) (*container.HostConfig, error) {
	var binds []string

	binds = append(binds, fmt.Sprintf("%s:/usr/local/bin/runtime:ro", d.daemonPath))

	// Mount the plugin if available
	if d.computerUsePluginPath != "" {
		binds = append(binds, fmt.Sprintf("%s:/usr/local/lib/runtime-computer-use:ro", d.computerUsePluginPath))
	}

	if len(volumeMountPathBinds) > 0 {
		binds = append(binds, volumeMountPathBinds...)
	}

	hostConfig := &container.HostConfig{
		Runtime:    runtime,
		CapDrop:    []string{"ALL"},
		CapAdd:     sandboxCaps,
		ExtraHosts: []string{"host.docker.internal:host-gateway"},
		Resources: container.Resources{
			CPUPeriod:  100000,
			CPUQuota:   sandboxDto.CpuQuota * 100000,
			Memory:     sandboxDto.MemoryQuota * 1024 * 1024 * 1024,
			MemorySwap: sandboxDto.MemoryQuota * 1024 * 1024 * 1024,
		},
		Binds: binds,
	}

	filesystem, err := d.getFilesystem(ctx)
	if err != nil {
		return nil, err
	}

	if filesystem == "xfs" {
		hostConfig.StorageOpt = map[string]string{
			"size": fmt.Sprintf("%dG", sandboxDto.StorageQuota),
		}
	}

	return hostConfig, nil
}

func (d *DockerClient) getContainerNetworkingConfig(_ context.Context) *network.NetworkingConfig {
	containerNetwork := config.GetContainerNetwork()
	if containerNetwork != "" {
		return &network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{
				containerNetwork: {},
			},
		}
	}
	return nil
}

func (d *DockerClient) getFilesystem(ctx context.Context) (string, error) {
	info, err := d.apiClient.Info(ctx)
	if err != nil {
		return "", err
	}

	for _, driver := range info.DriverStatus {
		if driver[0] == "Backing Filesystem" {
			return driver[1], nil
		}
	}

	return "", errors.New("filesystem not found")
}
