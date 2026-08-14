// Copyright 2025 Hanzo Industries Inc.
// SPDX-License-Identifier: AGPL-3.0

package docker

import (
	"context"
	"fmt"

	"github.com/hanzoai/runtime/apps/runner/pkg/api/dto"
	"github.com/hanzoai/runtime/apps/runner/pkg/common"
	"github.com/docker/docker/client"
)

// runtimeFor names the Docker runtime that provides each boundary. These are
// the names the tools register themselves under: `runsc install` writes
// "runsc", and Kata's Firecracker shim registers as "kata-fc".
var runtimeFor = map[dto.Isolation]string{
	dto.IsolationGvisor:      "runsc",
	dto.IsolationFirecracker: "kata-fc",
	dto.IsolationRunc:        "runc",
}

// Runtimes reports the runtime names the Docker daemon offers, which is what
// decides whether a boundary can be provided on this node.
func Runtimes(ctx context.Context, apiClient client.APIClient) (map[string]bool, error) {
	info, err := apiClient.Info(ctx)
	if err != nil {
		return nil, err
	}

	names := make(map[string]bool, len(info.Runtimes))
	for name := range info.Runtimes {
		names[name] = true
	}

	return names, nil
}

// IsolationRuntime names the Docker runtime that serves a boundary, for
// reporting which boundaries a node can provide.
func IsolationRuntime(iso dto.Isolation) string {
	return runtimeFor[iso]
}

// resolveIsolation returns the boundary a sandbox will actually run behind and
// the Docker runtime that serves it.
//
// An absent selection resolves to dto.IsolationDefault, so a caller that says
// nothing is isolated. A boundary the node cannot provide is an error rather
// than a substitution: nothing here returns a weaker runtime than the one that
// was asked for.
func (d *DockerClient) resolveIsolation(iso dto.Isolation, orgId string) (dto.Isolation, string, error) {
	if iso == "" {
		iso = dto.IsolationDefault
	}

	runtime, known := runtimeFor[iso]
	if !known {
		return "", "", common.NewBadRequestError(fmt.Errorf("unknown isolation %q", iso))
	}

	// runc has no boundary against code that attacks the kernel, so it is
	// reachable only by the organizations this runner names, and the empty
	// list names nobody.
	if iso == dto.IsolationRunc && !d.runcOrgs[orgId] {
		return "", "", common.NewBadRequestError(fmt.Errorf("isolation %q is not available to organization %q on this runner", iso, orgId))
	}

	if !d.runtimes[runtime] {
		return "", "", common.NewBadRequestError(fmt.Errorf("isolation %q needs the %q runtime, which this runner does not have", iso, runtime))
	}

	return iso, runtime, nil
}
