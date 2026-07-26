/*
 * Copyright 2025 Hanzo Industries Inc.
 * SPDX-License-Identifier: AGPL-3.0
 */

import { Module } from '@nestjs/common'
import { TypeOrmModule } from '@nestjs/typeorm'
import { SandboxUsagePeriod } from './entities/sandbox-usage-period.entity'
import { UsageService } from './services/usage.service'
import { KVLockProvider } from '../sandbox/common/kv-lock.provider'
import { Sandbox } from '../sandbox/entities/sandbox.entity'

@Module({
  imports: [TypeOrmModule.forFeature([SandboxUsagePeriod, Sandbox])],
  providers: [UsageService, KVLockProvider],
  exports: [UsageService],
})
export class UsageModule {}
