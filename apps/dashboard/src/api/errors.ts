/*
 * Copyright 2025 Hanzo Industries Inc.
 * SPDX-License-Identifier: AGPL-3.0
 */

export class RuntimeError extends Error {
  public static fromError(error: Error): RuntimeError {
    if (String(error).includes('Organization is suspended')) {
      return new OrganizationSuspendedError(error.message)
    }

    return new RuntimeError(error.message)
  }

  public static fromString(error: string): RuntimeError {
    return RuntimeError.fromError(new Error(error))
  }
}

export class OrganizationSuspendedError extends RuntimeError {}
