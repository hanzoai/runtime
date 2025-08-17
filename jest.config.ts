/*
 * Copyright 2025 Hanzo Industries Inc.
 * SPDX-License-Identifier: AGPL-3.0
 */

import { getJestProjectsAsync } from '@nx/jest'

export default async () => ({
  projects: await getJestProjectsAsync(),
})
