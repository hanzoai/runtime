// Copyright 2025 Hanzo Industries Inc.
// SPDX-License-Identifier: AGPL-3.0

package docker

import (
	"context"
	"testing"

	"github.com/hanzoai/runtime/apps/runner/pkg/api/dto"
	"github.com/docker/docker/api/types/system"
	"github.com/docker/docker/client"
)

// everyRuntime is a node that can provide every boundary.
var everyRuntime = map[string]bool{"runc": true, "runsc": true, "kata-fc": true}

// fakeDocker answers the one daemon question the host config asks. Any other
// call is a test reaching further than it means to, and panics.
type fakeDocker struct{ client.APIClient }

func (fakeDocker) Info(context.Context) (system.Info, error) {
	return system.Info{DriverStatus: [][2]string{{"Backing Filesystem", "extfs"}}}, nil
}

func TestResolveIsolation(t *testing.T) {
	cases := []struct {
		name     string
		runtimes map[string]bool
		runcOrgs map[string]bool
		asked    dto.Isolation
		org      string
		want     string
		wantErr  bool
	}{{
		name:     "saying nothing is isolated",
		runtimes: everyRuntime,
		asked:    "",
		org:      "acme",
		want:     "runsc",
	}, {
		name:     "gvisor resolves to runsc",
		runtimes: everyRuntime,
		asked:    dto.IsolationGvisor,
		org:      "acme",
		want:     "runsc",
	}, {
		name:     "firecracker resolves to the kata shim",
		runtimes: everyRuntime,
		asked:    dto.IsolationFirecracker,
		org:      "acme",
		want:     "kata-fc",
	}, {
		name:     "runc is refused to an org this runner does not name",
		runtimes: everyRuntime,
		runcOrgs: map[string]bool{"hanzo-internal": true},
		asked:    dto.IsolationRunc,
		org:      "acme",
		wantErr:  true,
	}, {
		name:     "runc is allowed to an org this runner names",
		runtimes: everyRuntime,
		runcOrgs: map[string]bool{"hanzo-internal": true},
		asked:    dto.IsolationRunc,
		org:      "hanzo-internal",
		want:     "runc",
	}, {
		name:     "an empty runc list names nobody",
		runtimes: everyRuntime,
		runcOrgs: map[string]bool{},
		asked:    dto.IsolationRunc,
		org:      "hanzo-internal",
		wantErr:  true,
	}, {
		name:     "a nameless caller does not reach runc",
		runtimes: everyRuntime,
		runcOrgs: map[string]bool{"hanzo-internal": true},
		asked:    dto.IsolationRunc,
		org:      "",
		wantErr:  true,
	}, {
		name:     "an unknown isolation is refused rather than guessed",
		runtimes: everyRuntime,
		asked:    dto.Isolation("none"),
		org:      "acme",
		wantErr:  true,
	}, {
		name:     "a node without runsc refuses rather than falling back to runc",
		runtimes: map[string]bool{"runc": true},
		runcOrgs: map[string]bool{"acme": true},
		asked:    dto.IsolationGvisor,
		org:      "acme",
		wantErr:  true,
	}, {
		name:     "a node without the kata shim refuses firecracker",
		runtimes: map[string]bool{"runc": true, "runsc": true},
		asked:    dto.IsolationFirecracker,
		org:      "acme",
		wantErr:  true,
	}}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			d := &DockerClient{runtimes: c.runtimes, runcOrgs: c.runcOrgs}

			iso, runtime, err := d.resolveIsolation(c.asked, c.org)
			if c.wantErr {
				if err == nil {
					t.Fatalf("resolved to %q, wanted a refusal", runtime)
				}
				return
			}
			if err != nil {
				t.Fatalf("refused: %v", err)
			}
			if runtime != c.want {
				t.Fatalf("runtime = %q, want %q", runtime, c.want)
			}
			if runtimeFor[iso] != runtime {
				t.Fatalf("reported isolation %q does not name runtime %q", iso, runtime)
			}
		})
	}
}

// A sandbox that asks for nothing must come out isolated, and it must never
// come out with the host capabilities the runner used to hand every sandbox.
func TestHostConfigDropsPrivilege(t *testing.T) {
	d := &DockerClient{runtimes: everyRuntime, daemonPath: "/usr/local/bin/runtime", apiClient: fakeDocker{}}

	_, runtime, err := d.resolveIsolation("", "acme")
	if err != nil {
		t.Fatalf("refused: %v", err)
	}

	hostConfig, err := d.getContainerHostConfig(t.Context(), dto.CreateSandboxDTO{
		Id:          "sandbox",
		OrgId:       "acme",
		CpuQuota:    2,
		MemoryQuota: 4,
	}, nil, runtime)
	if err != nil {
		t.Fatalf("host config: %v", err)
	}

	if hostConfig.Privileged {
		t.Fatal("sandbox is privileged")
	}
	if hostConfig.Runtime != "runsc" {
		t.Fatalf("runtime = %q, want runsc", hostConfig.Runtime)
	}
	if len(hostConfig.CapDrop) != 1 || hostConfig.CapDrop[0] != "ALL" {
		t.Fatalf("CapDrop = %v, want [ALL]", hostConfig.CapDrop)
	}

	for _, unwanted := range []string{"SYS_ADMIN", "NET_RAW", "MKNOD", "SYS_CHROOT", "SETFCAP", "SETPCAP", "NET_BIND_SERVICE"} {
		for _, granted := range hostConfig.CapAdd {
			if granted == unwanted {
				t.Fatalf("capability %s is granted", unwanted)
			}
		}
	}
}

// The label a container carries is what lets a running sandbox be attributed to
// an organization and to the boundary it runs behind.
func TestContainerLabelsRecordOrgAndIsolation(t *testing.T) {
	d := &DockerClient{runtimes: everyRuntime}

	config := d.getContainerCreateConfig(dto.CreateSandboxDTO{
		Id:       "sandbox",
		OrgId:    "acme",
		Snapshot: "python:3.12",
	}, dto.IsolationGvisor)

	if config.Labels["hanzo.ai/org"] != "acme" {
		t.Fatalf("org label = %q", config.Labels["hanzo.ai/org"])
	}
	if config.Labels["hanzo.ai/isolation"] != string(dto.IsolationGvisor) {
		t.Fatalf("isolation label = %q", config.Labels["hanzo.ai/isolation"])
	}
}
