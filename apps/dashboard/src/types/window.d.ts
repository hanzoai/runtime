/*
 * Copyright 2025 Hanzo Industries Inc.
 * SPDX-License-Identifier: AGPL-3.0
 */

interface Window {
  pylon?: {
    chat_settings: {
      app_id: string
      email: string
      name: string
      avatar_url?: string
      email_hash?: string
    }
  }
}
