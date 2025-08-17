/*
 * Copyright 2025 Hanzo Industries Inc.
 * SPDX-License-Identifier: AGPL-3.0
 */

import { Reflector } from '@nestjs/core'
import { OrganizationMemberRole } from '../enums/organization-member-role.enum'

export const RequiredOrganizationMemberRole = Reflector.createDecorator<OrganizationMemberRole>()
