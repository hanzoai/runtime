/*
 * Copyright 2025 Hanzo Industries Inc.
 * SPDX-License-Identifier: AGPL-3.0
 */

import { Organization } from '@hanzo/api-client'
import { createContext } from 'react'

export interface IOrganizationsContext {
  organizations: Organization[]
  refreshOrganizations: () => Promise<Organization[]>
}

export const OrganizationsContext = createContext<IOrganizationsContext | undefined>(undefined)
