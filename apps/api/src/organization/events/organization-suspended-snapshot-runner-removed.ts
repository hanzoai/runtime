/*
 * Copyright 2025 Hanzo Industries Inc.
 * SPDX-License-Identifier: AGPL-3.0
 */

export class OrganizationSuspendedSnapshotRunnerRemovedEvent {
  constructor(public readonly snapshotRunnerId: string) {}
}
