/*
 * Copyright 2025 Hanzo Industries Inc.
 * SPDX-License-Identifier: AGPL-3.0
 */

export class RunnerNotReadyError extends Error {
  constructor(message: string) {
    super(message)
    this.name = 'RunnerNotReadyError'
  }
}
