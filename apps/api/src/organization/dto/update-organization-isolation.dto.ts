/*
 * Copyright 2025 Hanzo Industries Inc.
 * SPDX-License-Identifier: AGPL-3.0
 */

import { ApiProperty, ApiSchema } from '@nestjs/swagger'
import { IsEnum, IsString, MinLength } from 'class-validator'
import { Isolation } from '../../sandbox/enums/isolation.enum'

@ApiSchema({ name: 'UpdateOrganizationIsolation' })
export class UpdateOrganizationIsolationDto {
  @ApiProperty({ enum: Isolation, enumName: 'Isolation' })
  @IsEnum(Isolation)
  isolation: Isolation

  // Every organization that runs without a boundary is a deliberate exception,
  // so the reason travels with the change rather than living in someone's
  // memory of why it was made.
  @ApiProperty({ description: 'Why this organization runs behind this boundary' })
  @IsString()
  @MinLength(1)
  reason: string
}
