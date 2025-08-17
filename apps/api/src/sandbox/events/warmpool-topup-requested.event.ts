/*
 * Copyright 2025 Hanzo Industries Inc.
 * SPDX-License-Identifier: AGPL-3.0
 */

import { WarmPool } from '../entities/warm-pool.entity'

export class WarmPoolTopUpRequested {
  constructor(public readonly warmPool: WarmPool) {}
}
