/*
 * Copyright 2025 Hanzo Industries Inc.
 * SPDX-License-Identifier: AGPL-3.0
 */

import { DocumentBuilder } from '@nestjs/swagger'

const getOpenApiConfig = (oidcIssuer: string) =>
  new DocumentBuilder()
    .setTitle('Runtime')
    .addServer('http://localhost:3000')
    .setDescription('Runtime AI platform API Docs')
    .setContact('Hanzo Industries Inc.', 'https://www.hanzo.ai', 'support@runtime.com')
    .setVersion('1.0')
    .addBearerAuth({
      type: 'http',
      scheme: 'bearer',
      description: 'API Key access',
    })
    .addOAuth2({
      type: 'openIdConnect',
      flows: undefined,
      openIdConnectUrl: `${oidcIssuer}/.well-known/openid-configuration`,
    })
    .build()

export { getOpenApiConfig }
