/*
 * Copyright 2025 Hanzo Industries Inc.
 * SPDX-License-Identifier: AGPL-3.0
 */

import { Volume } from '../entities/volume.entity'

export class VolumeLastUsedAtUpdatedEvent {
  constructor(public readonly volume: Volume) {}
}
