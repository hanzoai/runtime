/*
 * Copyright 2025 Daytona Platforms Inc.
 * SPDX-License-Identifier: AGPL-3.0
 */

import { InjectKV } from '../../common/kv.module'
import { Injectable } from '@nestjs/common'
import { KV } from '@hanzo/kv'

type Acquired = boolean

@Injectable()
export class KVLockProvider {
  constructor(@InjectKV() private readonly kv: KV) {}

  async lock(key: string, ttl: number): Promise<Acquired> {
    const acquired = await this.kv.set(key, '1', 'EX', ttl, 'NX')
    return !!acquired
  }

  async unlock(key: string): Promise<void> {
    await this.kv.del(key)
  }

  async waitForLock(key: string, ttl: number): Promise<void> {
    while (await this.kv.get(key)) {
      await new Promise((resolve) => setTimeout(resolve, 50))
    }

    await this.kv.setex(key, ttl, '1')
  }
}
