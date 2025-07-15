/*
 * Copyright 2025 Hanzo Industries Inc.
 * SPDX-License-Identifier: Apache-2.0
 */

/**
 * @module Errors
 */

/**
 * Base error for HanzoRuntime SDK.
 */
export class HanzoRuntimeError extends Error {}

export class HanzoRuntimeNotFoundError extends HanzoRuntimeError {}
