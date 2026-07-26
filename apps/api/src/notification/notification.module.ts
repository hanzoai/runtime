/*
 * Copyright 2025 Hanzo Industries Inc.
 * SPDX-License-Identifier: AGPL-3.0
 */

import { Module } from '@nestjs/common'
import { NotificationService } from './services/notification.service'
import { NotificationGateway } from './gateways/notification.gateway'
import { OrganizationModule } from '../organization/organization.module'
import { SandboxModule } from '../sandbox/sandbox.module'
import { KVModule } from '../common/kv.module'
import { AuthModule } from '../auth/auth.module'

@Module({
  imports: [OrganizationModule, SandboxModule, KVModule, AuthModule],
  providers: [NotificationService, NotificationGateway],
  exports: [NotificationService],
})
export class NotificationModule {}
