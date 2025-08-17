/*
 * Copyright 2025 Hanzo Industries Inc.
 * SPDX-License-Identifier: AGPL-3.0
 */

module.exports = {
  apps: [
    {
      name: 'runtime',
      script: './dist/apps/api/main.js',
      instances: 4,
      exec_mode: 'cluster',
      watch: false,
      env: {
        NODE_ENV: 'production',
        PM2_CLUSTER: 'true',
      },
      wait_ready: true,
      kill_timeout: 3000,
      listen_timeout: 10000,
    },
  ],
}
