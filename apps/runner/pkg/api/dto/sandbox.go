// Copyright 2025 Hanzo Industries Inc.
// SPDX-License-Identifier: AGPL-3.0

package dto

// Isolation names the boundary a sandbox runs behind. It is chosen per
// organization and travels with the create request, so the runner never has to
// ask anyone which boundary a sandbox gets.
//
// The empty value is not a selection. Anything that resolves it treats absent
// as IsolationDefault, so a caller that says nothing is isolated.
type Isolation string

const (
	// IsolationGvisor runs the sandbox under gVisor, which serves the guest's
	// syscalls from a userspace kernel rather than passing them to the host.
	IsolationGvisor Isolation = "gvisor"
	// IsolationFirecracker runs the sandbox in a microVM behind a hardware
	// boundary. It requires KVM on the node.
	IsolationFirecracker Isolation = "firecracker"
	// IsolationRunc runs the sandbox in namespaces and cgroups only, with no
	// boundary against code that attacks the kernel. Only organizations named
	// in the runner's RUNC_ORGS may select it.
	IsolationRunc Isolation = "runc"
)

// IsolationDefault is what an organization gets when it has configured nothing.
const IsolationDefault = IsolationGvisor

type CreateSandboxDTO struct {
	Id           string            `json:"id" validate:"required"`
	FromVolumeId string            `json:"fromVolumeId,omitempty"`
	OrgId        string            `json:"orgId" validate:"required"`
	Isolation    Isolation         `json:"isolation,omitempty"`
	Snapshot     string            `json:"snapshot" validate:"required"`
	OsUser       string            `json:"osUser" validate:"required"`
	CpuQuota     int64             `json:"cpuQuota" validate:"min=1"`
	GpuQuota     int64             `json:"gpuQuota" validate:"min=0"`
	MemoryQuota  int64             `json:"memoryQuota" validate:"min=1"`
	StorageQuota int64             `json:"storageQuota" validate:"min=1"`
	Env          map[string]string `json:"env,omitempty"`
	Registry     *RegistryDTO      `json:"registry,omitempty"`
	Entrypoint   []string          `json:"entrypoint,omitempty"`
	Volumes      []VolumeDTO       `json:"volumes,omitempty"`
} //	@name	CreateSandboxDTO

type ResizeSandboxDTO struct {
	Cpu    int64 `json:"cpu" validate:"min=1"`
	Gpu    int64 `json:"gpu" validate:"min=0"`
	Memory int64 `json:"memory" validate:"min=1"`
} //	@name	ResizeSandboxDTO
