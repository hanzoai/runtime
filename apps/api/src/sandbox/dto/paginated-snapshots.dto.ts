/*
 * Copyright 2025 Hanzo Industries Inc.
 * SPDX-License-Identifier: AGPL-3.0
 */

import { ApiProperty } from '@nestjs/swagger'
import { SnapshotDto } from './snapshot.dto'

export class PaginatedSnapshotsDto {
  @ApiProperty({ type: [SnapshotDto] })
  items: SnapshotDto[]

  @ApiProperty()
  total: number

  @ApiProperty()
  page: number

  @ApiProperty()
  totalPages: number
}
