/*
 * Copyright 2025 Hanzo Industries Inc.
 * SPDX-License-Identifier: AGPL-3.0
 */

// Isolation names the boundary a sandbox runs behind. An organization carries
// the one its sandboxes get; each sandbox records what it was created with, so
// changing the organization's choice moves new sandboxes and leaves running
// ones where they are.
export enum Isolation {
  // Syscalls are served by a userspace kernel rather than the host's.
  GVISOR = 'gvisor',
  // A microVM behind a hardware boundary. Requires KVM on the runner's node.
  FIRECRACKER = 'firecracker',
  // Namespaces and cgroups only, with no boundary against code that attacks
  // the kernel. Runners accept it only from the organizations they name.
  RUNC = 'runc',
}

// What an organization that has configured nothing gets.
export const ISOLATION_DEFAULT = Isolation.GVISOR
