/*
 * Copyright 2025 Hanzo Industries Inc.
 * SPDX-License-Identifier: AGPL-3.0
 */

import { UserOrganizationInvitationsContext } from '@/contexts/UserOrganizationInvitationsContext'
import { useContext } from 'react'

export function useUserOrganizationInvitations() {
  const context = useContext(UserOrganizationInvitationsContext)

  if (!context) {
    throw new Error('useUserOrganizationInvitations must be used within a UserOrganizationInvitationsProvider')
  }

  return context
}
