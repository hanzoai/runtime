/*
 * Copyright 2025 Hanzo Industries Inc.
 * SPDX-License-Identifier: AGPL-3.0
 */

import { Injectable, NestMiddleware } from '@nestjs/common'
import { Request, Response, NextFunction } from 'express'
// import { version } from '../../../package.json'

@Injectable()
export class VersionHeaderMiddleware implements NestMiddleware {
  use(req: Request, res: Response, next: NextFunction) {
    // TODO: Fetch version from package.json
    // res.setHeader('X-Runtime-Api-Version', `v${version}`)
    next()
  }
}
