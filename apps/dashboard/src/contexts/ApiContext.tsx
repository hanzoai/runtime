/*
 * Copyright 2025 Hanzo Industries Inc.
 * SPDX-License-Identifier: AGPL-3.0
 */

import { ApiClient } from '@/api/apiClient'
import { createContext } from 'react'

export const ApiContext = createContext<ApiClient | null>(null)
