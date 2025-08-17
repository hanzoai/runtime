/*
 * Copyright 2025 Hanzo Industries Inc.
 * SPDX-License-Identifier: AGPL-3.0
 */

import { User } from './user.entity'

describe('User', () => {
  it('should be defined', () => {
    expect(new User()).toBeDefined()
  })
})
