/*
 * Copyright 2025 Hanzo Industries Inc.
 * SPDX-License-Identifier: AGPL-3.0
 */

import { Sandbox } from '../entities/sandbox.entity'

export class SandboxBackupCreatedEvent {
  constructor(public readonly sandbox: Sandbox) {}
}
