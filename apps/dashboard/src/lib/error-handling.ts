/*
 * Copyright 2025 Daytona Platforms Inc.
 * SPDX-License-Identifier: AGPL-3.0
 */

import { Action, toast } from 'sonner'
import { RuntimeError } from '@/api/errors'

export function handleApiError(error: unknown, message: string, toastAction?: React.ReactNode | Action) {
  const isRuntimeError = error instanceof RuntimeError

  toast.error(message, {
    description: isRuntimeError ? error.message : 'Please try again or check the console for more details',
    action: toastAction,
  })

  if (!isRuntimeError) {
    console.error(message, error)
  }
}
